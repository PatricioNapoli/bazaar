package watchdog

import (
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

type Watchdog struct {
	Net       chain.EthClient
	Swaps     []*pairs.Swap
	Addresses []common.Address
	Contract  *chain.Chain
	Config    config.Config
}

// New creates a token reserves watchdog that
// maintains rates updated to their latest exchange when started
// using the provided contract address.
func New(cfg config.Config, client chain.EthClient, swaps []*pairs.Swap) Watchdog {
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
		Config:    cfg,
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
		wd.Swaps[i].Token0Reserve = reserves[i*2]
		wd.Swaps[i].Token1Reserve = reserves[i*2+1]

		if wd.Swaps[i].Token0Reserve.Cmp(new(big.Int)) == 0 || wd.Swaps[i].Token1Reserve.Cmp(new(big.Int)) == 0 {
			wd.Swaps[i].HasReserves = false
		}

		wd.Swaps[i].K = new(big.Int).Mul(wd.Swaps[i].Token0Reserve, wd.Swaps[i].Token1Reserve)
	}
}
