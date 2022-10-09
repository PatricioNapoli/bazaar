package pairs

import (
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"log"
	"sort"
)

type Swap struct {
	Address       string
	Token0        tokens.Token
	Token1        tokens.Token
	Token0Reserve float64
	Token1Reserve float64
	Rate0to1      float64
	Rate1to0      float64
}

type TokenSwap struct {
	Token tokens.Token
	Swaps []*Swap
}

type Path struct {
	TokenSwaps []TokenSwap
}

// NewPaths builds a sane struct representation parse from
// the provided swaps file, this improves later processing.
func NewPaths(cfg config.Config, tokenMap map[string]tokens.Token) ([]Path, []*Swap) {
	log.Printf("loading pairs info in %s", cfg.PairsFile)

	f, err := utils.ReadFile(cfg.PairsFile)
	if err != nil {
		log.Panicf("failed when reading file %s - %v", cfg.PairsFile, err)
	}

	paths := make([]Path, 0)
	swaps := make([]*Swap, 0)
	prePaths := make([]interface{}, 0)

	err = utils.FromJSON(f, &prePaths)
	if err != nil {
		log.Panicf("failed deserializing token file to array - %v", err)
	}

	for _, rootPath := range prePaths {
		rp := Path{}

		rootPathC := rootPath.([]interface{})

		for i, tswap := range rootPathC {

			tswapC := tswap.([]interface{})
			ts := TokenSwap{}
			ts.Token = tokenMap[tswapC[0].(string)]

			for _, swap := range tswapC[1].([]interface{}) {

				var t0 tokens.Token
				if i == 0 {
					t0 = getTokenFromPath(rootPathC[len(rootPathC)-1], tokenMap)
				} else {
					t0 = getTokenFromPath(rootPathC[i-1], tokenMap)
				}

				t1 := ts.Token

				setTokenOrder(&t0, &t1)

				swp := Swap{Address: swap.(string), Token0: t0, Token1: t1}

				swaps = append(swaps, &swp)
				ts.Swaps = append(ts.Swaps, &swp)
			}

			rp.TokenSwaps = append(rp.TokenSwaps, ts)
		}

		if len(rp.TokenSwaps) != 0 {
			paths = append(paths, rp)
		}
	}

	return paths, swaps
}

func getTokenFromPath(t interface{}, tokenMap map[string]tokens.Token) tokens.Token {
	return tokenMap[t.([]interface{})[0].(string)]
}

func setTokenOrder(t0 *tokens.Token, t1 *tokens.Token) {
	tkns := []string{t0.Address, t1.Address}
	sort.Strings(tkns)

	if t0.Address != tkns[0] {
		temp := *t1
		*t1 = *t0
		*t0 = temp
	}
}
