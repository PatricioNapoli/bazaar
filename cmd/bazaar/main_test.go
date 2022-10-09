package main

import (
	"context"
	"github.com/PatricioNapoli/bazaar/pkg/arbiter"
	"github.com/PatricioNapoli/bazaar/pkg/chain"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/PatricioNapoli/bazaar/pkg/pairs"
	"github.com/PatricioNapoli/bazaar/pkg/tokens"
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"github.com/PatricioNapoli/bazaar/pkg/watchdog"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestReduceBigInt(t *testing.T) {
	res := utils.ReduceBigInt(new(big.Int).SetInt64(10000), 4)

	AssertFloat(t, res, 1.0)
}

func TestNewConfig(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = ""

	expected := string(ReadFile(t, "test/expected/config.json"))

	res, err := utils.ToJSON(cfg)
	if err != nil {
		t.Errorf("%v", err)
	}

	AssertString(t, string(res), expected)
}

func TestNewEthClient(t *testing.T) {
	cfg := config.NewConfig()

	eth := chain.NewEthClient(cfg)

	_, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Errorf("failed when fetching suggested gas: %v", err)
	}
}

func TestGetTokens(t *testing.T) {
	cfg := config.NewConfig()
	cfg.TokensFile = "test/fixtures/tokens.json"

	tkns := tokens.GetTokens(cfg)

	AssertExpectedJSON(t, tkns, "test/expected/tokens.json")
}

func TestNewPaths(t *testing.T) {
	cfg := config.NewConfig()
	cfg.TokensFile = "test/fixtures/tokens.json"
	cfg.PairsFile = "test/fixtures/pairs.json"

	tkns := tokens.GetTokens(cfg)

	paths, swaps := pairs.NewPaths(cfg, tkns)

	AssertExpectedJSON(t, paths, "test/expected/paths.json")
	AssertExpectedJSON(t, swaps, "test/expected/swaps.json")
}

func TestNewWatchdog(t *testing.T) {
	cfg := config.NewConfig()
	cfg.TokensFile = "test/fixtures/tokens.json"
	cfg.PairsFile = "test/fixtures/pairs.json"
	cfg.RatePrecision = 5

	tkns := tokens.GetTokens(cfg)
	_, swaps := pairs.NewPaths(cfg, tkns)

	eth := chain.NewEthClient(cfg)

	wd := watchdog.NewWatchdog(cfg, eth, swaps)
	wd.Start()

	AssertExpectedJSON(t, swaps, "test/expected/swaps_wd.json")
}

func TestNewArbiter(t *testing.T) {
	cfg := config.NewConfig()
	cfg.TokensFile = "test/fixtures/tokens.json"
	cfg.PairsFile = "test/fixtures/pairs.json"
	cfg.OutputFilename = "output/output_test.json"
	cfg.PrettyPrintOutput = false
	cfg.RatePrecision = 6
	cfg.IncludeFees = false

	tkns := tokens.GetTokens(cfg)
	paths, swaps := pairs.NewPaths(cfg, tkns)

	eth := chain.NewEthClient(cfg)

	wd := watchdog.NewWatchdog(cfg, eth, swaps)
	wd.Start()

	arb := arbiter.NewArbiter(cfg, eth, paths)
	arb.Start()

	out := ReadFile(t, cfg.OutputFilename)
	exp := ReadFile(t, "test/expected/arbiter.json")

	AssertString(t, string(out), string(exp))
}

func ReadFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJSON
}

func AssertExpectedJSON(t *testing.T, obj interface{}, filename string) {
	expected := string(ReadFile(t, filename))

	res, err := utils.ToJSON(obj)
	if err != nil {
		t.Errorf("%v", err)
	}

	AssertString(t, string(res), expected)
}

func AssertString(t *testing.T, res string, expected string) {
	if res != expected {
		t.Errorf("got: %s - expected: %s", res, expected)
	}
}

func AssertFloat(t *testing.T, res float64, expected float64) {
	if res != expected {
		t.Errorf("got: %f - expected: %f", res, expected)
	}
}
