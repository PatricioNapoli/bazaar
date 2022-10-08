package watchdog

import (
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math"
	"math/big"
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
		log.Fatalf("failed when binding to contract: %s", err)
	}

	return Watchdog{
		Net:       client,
		Swaps:     swaps,
		Addresses: addresses,
		Contract:  contract,
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

	reserves, err := wd.Contract.ViewPair(nil, wd.Addresses)
	if err != nil {
		log.Fatalf("failed when fetching contract pair reserves: %s", err)
	}

	for i := 0; i < len(wd.Swaps); i++ {
		t0dec, _ := strconv.Atoi(wd.Swaps[i].Token0.Decimals)
		wd.Swaps[i].Token0Reserve = ReduceBigInt(reserves[i*2], t0dec)

		t1dec, _ := strconv.Atoi(wd.Swaps[i].Token1.Decimals)
		wd.Swaps[i].Token1Reserve = ReduceBigInt(reserves[i*2+1], t1dec)

		if wd.Swaps[i].Token0Reserve < 1 || wd.Swaps[i].Token1Reserve < 1 {
			wd.Swaps[i].HasReserves = false
		} else {
			wd.Swaps[i].Rate0to1 = wd.Swaps[i].Token1Reserve / wd.Swaps[i].Token0Reserve
			wd.Swaps[i].Rate1to0 = wd.Swaps[i].Token0Reserve / wd.Swaps[i].Token1Reserve
			wd.Swaps[i].HasReserves = true
		}
	}

	broker <- "refresh"
}

func ReduceBigInt(n *big.Int, decimals int) float64 {
	mul := math.Pow(10.0, -float64(decimals))

	flt := new(big.Float).SetPrec(128).SetInt(n)
	flt = new(big.Float).Mul(flt, new(big.Float).SetFloat64(mul))

	final, _ := flt.Float64()
	return final
}
