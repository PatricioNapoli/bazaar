package arbiter

import (
	"context"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"log"
)

type SwapNode struct {
	Children []*SwapNode
	Target   tokens.Token
	Swap     *pairs.Swap
}

type Arbiter struct {
	Net       chain.EthClient
	RootNodes []SwapNode
}

type Arbitration struct {
	FinalCoins float64
	Path       []*SwapNode
}

func NewArbiter(paths []pairs.Path, client chain.EthClient) Arbiter {
	rnodes := make([]SwapNode, len(paths))

	for x, p := range paths {
		rn := SwapNode{
			Children: make([]*SwapNode, 0),
			Swap:     nil, // Starting node is root, has no swap
		}

		lastNodes := []*SwapNode{&rn}
		for _, tswp := range p.TokenSwaps {

			nodes := make([]*SwapNode, 2)

			for j, swap := range tswp.Swaps {
				node := SwapNode{
					Children: make([]*SwapNode, 0),
					Swap:     swap,
					Target:   tswp.Token,
				}

				nodes[j] = &node

				for _, n := range lastNodes {
					if n != nil {
						n.Children = append(n.Children, &node)
					}
				}
			}

			lastNodes = nodes
		}

		rnodes[x] = rn
	}

	return Arbiter{
		Net:       client,
		RootNodes: rnodes,
	}
}

func (arb *Arbiter) Start(broker chan string) {
	log.Println("starting bazaar arbiter")

	for msg := range broker {
		if msg == "refresh" {
			arb.Search()
		}
	}
}

func (arb *Arbiter) Search() {
	log.Println("searching possible arbitration paths")

	_, err := arb.Net.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed when fetching suggested gas: %s", err)
	}

	// Gas price + 0.3% fee for swapping. Slippage should also be considered if transactions are to be executed.

	for _, rn := range arb.RootNodes {
		for _, node := range rn.Children {
			arb.graphArbitration(node, 1.0)
		}
	}
}

func (arb *Arbiter) graphArbitration(node *SwapNode, coins float64) float64 {
	if node.Target.Address == node.Swap.Token0.Address {
		coins = coins * node.Swap.Rate0to1
	}

	for _, n := range node.Children {
		if n.Swap.Address == node.Swap.Address {
			continue
		}

		arb.graphArbitration(n, coins)
	}

	return coins
}
