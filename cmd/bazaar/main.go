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

	_, swaps := pairs.NewPaths("./assets/uni_sushi_paths.json", t)

	client := chain.NewEthClient("https://mainnet.infura.io/v3/")

	wd := watchdog.NewWatchdog(client, swaps)

	broker := make(chan string)
	go wd.Start(broker)

	arb := arbiter.NewArbiter()
	go arb.Start(broker)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-stop

	log.Println("stopping all arbiter operations")

	broker <- "stop"
	<-broker
}

// Maintain a Watchdog fetching reserves and updating a map with swap balances. Swaps without balance cannot be used.
// Discard tokens size 2 with each swap size 1, impossible to arbitrage.
// Represent each path swap array as a node graph, should traverse the node forward and back, looking for
// a positive profit at the end of traversal, should be easy after fetching rates
// Nodes can be built TOKATOKB (UNI) -> TOKBTOKC (UNI) -> TOKCTOKA (UNI) and are valid, along with any market combination
// in between (node traversal cannot backtrack, but it will be attempted reversed), single path, either through UNI or SUSHI
// DO NOT attempt coming out through same market, same pair, unless total swaps is greater than 1
