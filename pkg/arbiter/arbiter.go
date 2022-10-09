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

// NewArbiter creates an arbitration node graph structure, that's capable of easily
// traversing it and finding potential market arbitration paths.
func NewArbiter(cfg config.Config, client chain.EthClient, paths []pairs.Path) Arbiter {
	// These are excluded because they do not belong to Uniswap or Sushiswap,
	// furthermore, it appears that they have been disabled and reserves have not
	// been rebalanced, leading to non-valid exchange rates, since they are not
	// consistent with exchange rates seen in Uniswap and Sushiswap contracts.
	dollarProtocolAddr0 := "0x6b583cf4aba7bf9d6f8a51b3f1f7c7b2ce59bf15"
	dollarProtocolAddr1 := "0xd233d1f6fd11640081abb8db125f722b5dc729dc"

	rnodes := make([]SwapNode, len(paths))

	for x, p := range paths {
		rn := SwapNode{
			Children: make([]*SwapNode, 0),
			Swap:     nil,
		}

		skip := false
		lastNodes := []*SwapNode{&rn}
		for _, tswp := range p.TokenSwaps {

			if cfg.ExcludeDeadTokens {
				if tswp.Token.Address == dollarProtocolAddr0 || tswp.Token.Address == dollarProtocolAddr1 {
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

// Start initiates the node graph exploration and writes the output file.
func (arb *Arbiter) Start() {
	log.Println("searching possible arbitration paths")

	gp, err := arb.Net.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Panicf("failed when fetching suggested gas: %v", err)
	}

	gpgwei := utils.ReduceBigInt(gp, 9)
	totalGas := gpgwei * arb.Config.DexSwapGas

	if arb.Config.IncludeFees {
		log.Printf("current gas price in gwei: %f", totalGas)
	}

	gpeth := totalGas / math.Pow(10, 9)

	millisBefore := utils.GetMillis()

	balance := arb.Config.InitialWETH
	pathGroups := make([]Arbitration, 0)

	// Splitting this operation in many go threads yielded bad results
	// Threads benefit when the CPU computation of the unit surpasses
	// that of the thread scheduling
	for _, rn := range arb.RootNodes {

		for _, node := range rn.Children {
			currArbitration := Arbitration{
				Balance: balance,
				Path:    make([]*SwapNode, 0),
			}
			arbs := make([]Arbitration, 0)
			arb.exploreArbitrationGraph(&arbs, &currArbitration, node)

			for _, a := range arbs {
				pathLength := float64(len(a.Path))

				if arb.Config.IncludeFees {
					a.Balance -= gpeth * pathLength             // network fees
					a.Balance -= a.Balance * 0.003 * pathLength // DEX swap fees
				}

				if a.Balance > balance {
					pathGroups = append(pathGroups, a)
				}
			}
		}

	}

	log.Printf("found %d routes for market making in %dms", len(pathGroups), utils.GetMillis()-millisBefore)

	writePathsToFile(arb.Config, pathGroups, arb.Config.OutputFilename)
}

func (arb *Arbiter) exploreArbitrationGraph(arbs *[]Arbitration, currArb *Arbitration, node *SwapNode) {
	if node.Target.Address == node.Swap.Token0.Address {
		exch := currArb.Balance * node.Swap.Rate1to0
		if node.Swap.Token0Reserve < exch {
			*arbs = make([]Arbitration, 0)
			return
		}
		currArb.Balance = exch
	} else {
		exch := currArb.Balance * node.Swap.Rate0to1
		if node.Swap.Token1Reserve < exch {
			*arbs = make([]Arbitration, 0)
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

		if i > 0 && n.Target.Address != arb.Config.WETHAddr {
			newPath := currArb.Path[:len(currArb.Path)-1]
			currArb = &Arbitration{
				Balance: balanceBefore,
				Path:    newPath,
			}
		}
		arb.exploreArbitrationGraph(arbs, currArb, n)
	}

	if len(node.Children) == 0 {
		*arbs = append(*arbs, *currArb)
	}
}

func writePathsToFile(cfg config.Config, paths []Arbitration, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("unable to create arbitration log file: %v", err)
	}
	defer file.Close()

	var bytes []byte

	if cfg.PrettyPrintOutput {
		bytes, _ = utils.ToPrettyJSON(paths)
	} else {
		bytes, _ = utils.ToJSON(paths)
	}

	log.Printf("writing arbitration output to %s", filename)
	fmt.Fprintf(file, "%s", string(bytes))
}
