package arbiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"log"
	"math/big"
	"os"
)

type SwapNode struct {
	Children []*SwapNode `json:"-"`
	Target   tokens.Token
	Swap     *pairs.Swap
}

type Arbitration struct {
	Balance *big.Int
	Path    []*SwapNode
}

type Arbiter struct {
	Net       chain.EthClient `json:"-"`
	RootNodes []SwapNode
	Config    config.Config
}

// New creates an arbitration node graph structure, that's capable of easily
// traversing it and finding potential market arbitration paths.
func New(cfg config.Config, client chain.EthClient, paths []pairs.Path) Arbiter {
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

	dexSwapWei := utils.ExtendBigInt(new(big.Int).SetInt64(arb.Config.DexSwapGas), 9)
	totalGas := new(big.Int).Mul(gp, dexSwapWei)

	if arb.Config.IncludeGas {
		log.Printf("current gas price in gwei: %s", gp.String())
	}

	millisBefore := utils.GetMillis()

	initialBalance := utils.ExtendBigInt(new(big.Int).SetInt64(int64(arb.Config.InitialWETH)), 18)
	balance := new(big.Int).Set(initialBalance)
	pathGroups := make([]Arbitration, 0)

	if arb.Config.IncludeFees {
		balance = utils.ExtractFees(balance, arb.Config.DEXFee, 3)
	}

	// Splitting this operation in many go threads yielded bad results
	// Threads benefit when the CPU computation of the unit surpasses
	// that of the thread scheduling load and context switches
	for _, rn := range arb.RootNodes {

		for _, node := range rn.Children {
			currArbitration := Arbitration{
				Balance: balance,
				Path:    make([]*SwapNode, 0),
			}
			arbs := make([]Arbitration, 0)
			arb.findPaths(&arbs, &currArbitration, node)

			for _, a := range arbs {
				pathLength := len(a.Path)

				if arb.Config.IncludeGas {
					totalWei := new(big.Int).Mul(totalGas, new(big.Int).SetInt64(int64(pathLength)))
					a.Balance = new(big.Int).Rem(a.Balance, totalWei)
				}

				if a.Balance.Cmp(initialBalance) == 1 {
					pathGroups = append(pathGroups, a)
				}
			}
		}

	}

	log.Printf("found %d routes for market making in %dms", len(pathGroups), utils.GetMillis()-millisBefore)

	writePathsToFile(arb.Config, pathGroups, arb.Config.OutputFilename)
}

func (arb *Arbiter) findPaths(arbs *[]Arbitration, currArb *Arbitration, node *SwapNode) {
	if !node.Swap.HasReserves {
		currArb.Balance = new(big.Int)
		return
	}

	// Check what token we are swapping to
	if node.Target.Address == node.Swap.Token0.Address {
		out, err := getOutTokens(node.Swap.Token1Reserve, node.Swap.Token0Reserve, node.Swap.K, currArb.Balance)

		if err != nil {
			currArb.Balance = out
			return
		}
		currArb.Balance = out
	} else {
		out, err := getOutTokens(node.Swap.Token0Reserve, node.Swap.Token1Reserve, node.Swap.K, currArb.Balance)

		if err != nil {
			currArb.Balance = out
			return
		}
		currArb.Balance = out
	}

	currArb.Path = append(currArb.Path, node)

	balanceBefore := currArb.Balance
	for i, n := range node.Children {
		// Refuse to hop to itself or to same target token.
		if n.Swap.Address == node.Swap.Address || n.Target.Address == node.Target.Address {
			continue
		}

		// Check for second token swap in same Target, it means it is a different path.
		if i > 0 {
			newPath := currArb.Path[:len(currArb.Path)-1]
			currArb = &Arbitration{
				Balance: balanceBefore,
				Path:    newPath,
			}
		}
		arb.findPaths(arbs, currArb, n)
	}

	// Finish arbitration when no more children found.
	if len(node.Children) == 0 {
		*arbs = append(*arbs, *currArb)
	}
}

func getOutTokens(from *big.Int, to *big.Int, K *big.Int, balance *big.Int) (*big.Int, error) {
	newFrom := new(big.Int).Add(from, balance)

	if newFrom.Cmp(K) == 1 {
		return new(big.Int), errors.New("not enough balance")
	}

	newTo := new(big.Int).Div(K, newFrom)
	out := new(big.Int).Rem(to, newTo)

	return out, nil
}

func writePathsToFile(cfg config.Config, paths []Arbitration, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("unable to create arbitration log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed when closing file %s", filename)
		}
	}(file)

	var bytes []byte

	if cfg.PrettyPrintOutput {
		bytes, _ = utils.ToPrettyJSON(paths)
	} else {
		bytes, _ = utils.ToJSON(paths)
	}

	log.Printf("writing arbitration output to %s", filename)
	_, err = fmt.Fprintf(file, "%s", string(bytes))
	if err != nil {
		log.Printf("failed when writing file %s", filename)
	}
}
