package main

import (
	"fmt"
	"github.com/PatricioNapoli/bazaar/pkg/arbiter"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/console"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"github.com/PatricioNapoli/bazaar/pkg/watchdog"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	wethAddr := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	reservesAddr := "0x416355755f32b2710ce38725ed0fa102ce7d07e6"

	tokensFile := "./assets/tokens.json"
	pairsFile := "./assets/uni_sushi_paths.json"

	infuraEndpoint := "https://mainnet.infura.io/v3/"

	console.PrintArt()
	log.Printf("launching Bazaar on epoch ms %d", utils.GetMillis())

	cfg := config.NewConfig()
	client := chain.NewEthClient(cfg, infuraEndpoint)

	tkns := tokens.GetTokens(tokensFile)

	log.Printf("attempting arbitrage on entry coin: ")
	wethJson, _ := utils.ToPrettyJSON(tkns[wethAddr])
	fmt.Println(string(wethJson))

	paths, swaps := pairs.NewPaths(pairsFile, tkns)

	wd := watchdog.NewWatchdog(reservesAddr, client, swaps)
	arb := arbiter.NewArbiter(paths, cfg, client)

	var wg sync.WaitGroup
	keep := true

	wg.Add(1)
	go func() {
		for keep {
			wd.Start()
			arb.Start()
		}

		log.Println("stopping all bazaar operations")
		wg.Done()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	<-stop
	keep = false
	wg.Wait()
}
