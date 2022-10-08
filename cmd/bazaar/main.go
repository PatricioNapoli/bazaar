package main

import (
	"fmt"
	"github.com/PatricioNapoli/bazaar/pkg/arbiter"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/console"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"github.com/PatricioNapoli/bazaar/pkg/watchdog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	console.PrintArt()

	log.Printf("launching Bazaar on epoch ms %d", utils.GetMillis())

	t := tokens.GetTokens("./assets/tokens.json")

	log.Printf("attempting arbitrage on entry coin: ")

	j, _ := utils.ToPrettyJSON(t["0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"])
	fmt.Println(string(j))

	paths, swaps := pairs.NewPaths("./assets/uni_sushi_paths.json", t)

	client := chain.NewEthClient("https://mainnet.infura.io/v3/")

	wd := watchdog.NewWatchdog(client, swaps)

	broker := make(chan string)
	go wd.Start(broker)

	arb := arbiter.NewArbiter(paths, client)
	go arb.Start(broker)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-stop

	log.Println("stopping all arbiter operations")

	broker <- "stop"
	<-broker
}
