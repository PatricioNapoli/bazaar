package arbiter

import (
	"context"
	"fmt"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"log"
	"math"
	"os"
)

type SwapNode struct {
	Children []*SwapNode `json:"-"`
	Target   tokens.Token
	Swap     *pairs.Swap
}

type Arbitration struct {
	Balance float64
	Path    []*SwapNode
}

type Arbiter struct {
	Net       chain.EthClient `json:"-"`
	RootNodes []SwapNode
	Config    config.Config
}

func NewArbiter(paths []pairs.Path, cfg config.Config, client chain.EthClient) Arbiter {
	rnodes := make([]SwapNode, len(paths))

	for x, p := range paths {
		rn := SwapNode{
			Children: make([]*SwapNode, 0),
			Swap:     nil, // Starting node is root, has no swap
		}

		skip := false
		lastNodes := []*SwapNode{&rn}
		for _, tswp := range p.TokenSwaps {

			if cfg.ExcludeDeadTokens {
				if tswp.Token.Address == "0x6b583cf4aba7bf9d6f8a51b3f1f7c7b2ce59bf15" || tswp.Token.Address == "0xd233d1f6fd11640081abb8db125f722b5dc729dc" {
					skip = true
					break
				}
			}

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

		if !skip {
			rnodes[x] = rn
		}
	}

	return Arbiter{
		Net:       client,
		RootNodes: rnodes,
		Config:    cfg,
	}
}

func (arb *Arbiter) Start() {
	log.Println("searching possible arbitration paths")

	gp, err := arb.Net.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Panicf("failed when fetching suggested gas: %s", err)
	}

	gpgwei := utils.ReduceBigInt(gp, 9)
	totalGas := gpgwei * 160000

	log.Printf("current gas price in gwei: %f", totalGas)

	gpeth := totalGas / math.Pow(10, 9)

	millisBefore := utils.GetMillis()

	balance := 1.0
	highest := Arbitration{}
	pathGroups := make([]Arbitration, 0)
	for _, rn := range arb.RootNodes {
		for _, node := range rn.Children {
			currArbitration := Arbitration{
				Balance: balance,
				Path:    make([]*SwapNode, 0),
			}
			arbitrations := make([]Arbitration, 0)
			arb.graphArbitration(&arbitrations, &currArbitration, node)

			for _, a := range arbitrations {
				pathLength := float64(len(a.Path))

				if arb.Config.IncludeFees {
					a.Balance -= gpeth * pathLength             // network fees
					a.Balance -= a.Balance * 0.003 * pathLength // DEX swap fees
				}

				if a.Balance > balance {
					pathGroups = append(pathGroups, a)

					if a.Balance > highest.Balance {
						highest = a
					}
				}
			}
		}
	}

	log.Printf("found %d routes for market making in %dms", len(pathGroups), utils.GetMillis()-millisBefore)

	writePathsToFile(pathGroups)
}

func (arb *Arbiter) graphArbitration(arbitrations *[]Arbitration, currArb *Arbitration, node *SwapNode) {
	if node.Target.Address == node.Swap.Token0.Address {
		exch := currArb.Balance * node.Swap.Rate1to0
		if node.Swap.Token0Reserve < exch {
			currArb.Balance = 1.0
			return
		}
		currArb.Balance = exch
	} else {
		exch := currArb.Balance * node.Swap.Rate0to1
		if node.Swap.Token1Reserve < exch {
			currArb.Balance = 1.0
			return
		}
		currArb.Balance = exch
	}

	currArb.Path = append(currArb.Path, node)

	balanceBefore := currArb.Balance
	for i, n := range node.Children {
		// Refuse to hop to itself or to same target
		if n.Swap.Address == node.Swap.Address || n.Target.Address == node.Target.Address {
			continue
		}

		if i > 0 {
			newPath := currArb.Path[:len(currArb.Path)-1]
			currArb = &Arbitration{
				Balance: balanceBefore,
				Path:    newPath,
			}
		}
		arb.graphArbitration(arbitrations, currArb, n)
	}

	if len(node.Children) == 0 {
		*arbitrations = append(*arbitrations, *currArb)
	}
}

func writePathsToFile(paths []Arbitration) {
	file, err := os.Create("output.txt")
	if err != nil {
		log.Panicf("unable to create arbitration log file: %v", err)
	}
	defer file.Close()

	bytes, _ := utils.ToPrettyJSON(paths)
	fmt.Fprintf(file, "%s", string(bytes))
}
