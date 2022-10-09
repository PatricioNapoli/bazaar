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
	console.PrintArt()
	log.Printf("launching Bazaar on epoch ms %d", utils.GetMillis())

	cfg := config.NewConfig()

	client := chain.NewEthClient(cfg)

	tkns := tokens.GetTokens(cfg)
	paths, swaps := pairs.NewPaths(cfg, tkns)

	log.Printf("attempting arbitrage on entry coin: ")
	wethJson, _ := utils.ToPrettyJSON(tkns[cfg.WETHAddr])
	fmt.Println(string(wethJson))

	wd := watchdog.NewWatchdog(cfg, client, swaps)
	arb := arbiter.NewArbiter(cfg, client, paths)

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
