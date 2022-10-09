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

	log.Printf("loading config from env vars")

	cfg := config.NewConfig()

	t := tokens.GetTokens("./assets/tokens.json")

	log.Printf("attempting arbitrage on entry coin: ")

	j, _ := utils.ToPrettyJSON(t["0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"])
	fmt.Println(string(j))

	paths, swaps := pairs.NewPaths("./assets/uni_sushi_paths.json", t)

	client := chain.NewEthClient(cfg, "https://mainnet.infura.io/v3/")

	wd := watchdog.NewWatchdog(client, swaps)
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
