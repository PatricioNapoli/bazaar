package watchdog

import (
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"strconv"
)

type Watchdog struct {
	Net       chain.EthClient
	Swaps     []*pairs.Swap
	Addresses []common.Address
	Contract  *chain.Chain
}

func NewWatchdog(client chain.EthClient, swaps []*pairs.Swap) Watchdog {
	addresses := make([]common.Address, len(swaps))

	for i, swp := range swaps {
		addresses[i] = common.HexToAddress(swp.Address)
	}

	address := common.HexToAddress("0x416355755f32b2710ce38725ed0fa102ce7d07e6")
	contract, err := chain.NewChain(address, client.Client)
	if err != nil {
		log.Panicf("failed when binding to contract: %s", err)
	}

	return Watchdog{
		Net:       client,
		Swaps:     swaps,
		Addresses: addresses,
		Contract:  contract,
	}
}

func (wd *Watchdog) Start() {
	log.Println("refreshing reserves from DEX")

	reserves, err := wd.Contract.ViewPair(nil, wd.Addresses)
	if err != nil {
		log.Panicf("failed when fetching contract pair reserves: %s", err)
	}

	for i := 0; i < len(wd.Swaps); i++ {
		t0dec, _ := strconv.Atoi(wd.Swaps[i].Token0.Decimals)
		wd.Swaps[i].Token0Reserve = utils.ReduceBigInt(reserves[i*2], t0dec)

		t1dec, _ := strconv.Atoi(wd.Swaps[i].Token1.Decimals)
		wd.Swaps[i].Token1Reserve = utils.ReduceBigInt(reserves[i*2+1], t1dec)

		if wd.Swaps[i].Token0Reserve > 0 && wd.Swaps[i].Token1Reserve > 0 {
			wd.Swaps[i].Rate0to1 = wd.Swaps[i].Token1Reserve / wd.Swaps[i].Token0Reserve
			wd.Swaps[i].Rate1to0 = wd.Swaps[i].Token0Reserve / wd.Swaps[i].Token1Reserve
		}
	}
}
