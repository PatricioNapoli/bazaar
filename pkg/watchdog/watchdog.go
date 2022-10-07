package watchdog

import (
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

type Watchdog struct {
	Net       chain.EthClient
	Swaps     []*pairs.Swap
	Addresses []common.Address
}

func NewWatchdog(client chain.EthClient, swaps []*pairs.Swap) Watchdog {
	addresses := make([]common.Address, len(swaps))

	for i, swp := range swaps {
		addresses[i] = common.HexToAddress(swp.Address)
	}

	return Watchdog{
		Net:       client,
		Swaps:     swaps,
		Addresses: addresses,
	}
}

func (wd *Watchdog) Start(broker chan string) {
	log.Println("starting DEX token reserves watchdog")

	for {
		select {
		case msg, ok := <-broker:
			log.Printf("%s watchdog", msg)
			if ok {
				close(broker)
			} else {
				break
			}
		default:
			wd.RefreshRates(broker)
		}
	}
}

func (wd *Watchdog) RefreshRates(broker chan string) {
	log.Println("refreshing pairs reserves data")

	address := common.HexToAddress("0x416355755f32b2710ce38725ed0fa102ce7d07e6")
	instance, err := chain.NewChain(address, wd.Net.Client)
	if err != nil {
		log.Fatal(err)
	}

	reserves, err := instance.ViewPair(nil, wd.Addresses)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(wd.Swaps); i++ {
		wd.Swaps[i].Token0.Reserve = *reserves[i*2]
		wd.Swaps[i].Token1.Reserve = *reserves[i*2+1]

		/*		rt := *new(big.Int).Div(&wd.Swaps[i].Token1.Reserve, &wd.Swaps[i].Token0.Reserve)

				wd.Swaps[i].Rate = *new(big.Int).Div(&wd.Swaps[i].Token1.Reserve, &wd.Swaps[i].Token0.Reserve)*/
	}

	broker <- "refresh"
}
