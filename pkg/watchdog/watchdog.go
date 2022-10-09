package watchdog

import (
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
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

// NewWatchdog creates a token reserves watchdog that
// maintains rates updated to their latest exchange when started
// using the provided contract address.
func NewWatchdog(cfg config.Config, client chain.EthClient, swaps []*pairs.Swap) Watchdog {
	addresses := make([]common.Address, len(swaps))

	for i, swp := range swaps {
		addresses[i] = common.HexToAddress(swp.Address)
	}

	address := common.HexToAddress(cfg.ReservesAddr)
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

// Start starts the watchdog call to the net, using a special contract pre compiled ABI.
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
