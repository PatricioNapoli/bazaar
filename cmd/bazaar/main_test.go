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
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestReduceBigInt(t *testing.T) {
	res := utils.ReduceBigInt(new(big.Int).SetInt64(10000), 4)

	AssertBigInt(t, res, new(big.Int).SetInt64(1))
}

func TestExtendBigInt(t *testing.T) {
	res := utils.ExtendBigInt(new(big.Int).SetInt64(1), 4)

	AssertBigInt(t, res, new(big.Int).SetInt64(10000))
}

func TestExtractFees(t *testing.T) {
	res := utils.ExtractFees(new(big.Int).SetInt64(100), 0.5, 2)

	AssertBigInt(t, res, new(big.Int).SetInt64(50))
}

func TestNewConfig(t *testing.T) {
	cfg := config.New()
	cfg.APIKey = ""

	expected := string(ReadFile(t, "test/expected/config.json"))

	res, err := utils.ToJSON(cfg)
	if err != nil {
		t.Errorf("%v", err)
	}

	AssertString(t, string(res), expected)
}

func TestNewEthClient(t *testing.T) {
	cfg := config.New()

	eth := chain.NewEthClient(cfg)

	_, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Errorf("failed when fetching suggested gas: %v", err)
	}
}

func TestGetTokens(t *testing.T) {
	cfg := config.New()
	cfg.TokensFile = "test/fixtures/tokens.json"

	tkns := tokens.New(cfg)

	AssertExpectedJSON(t, tkns, "test/expected/tokens.json")
}

func TestNewPaths(t *testing.T) {
	cfg := config.New()
	cfg.TokensFile = "test/fixtures/tokens.json"
	cfg.PairsFile = "test/fixtures/pairs.json"

	tkns := tokens.New(cfg)

	paths, swaps := pairs.New(cfg, tkns)

	AssertExpectedJSON(t, paths, "test/expected/paths.json")
	AssertExpectedJSON(t, swaps, "test/expected/swaps.json")
}

func TestNewWatchdog(t *testing.T) {
	cfg := config.New()
	cfg.TokensFile = "test/fixtures/tokens.json"
	cfg.PairsFile = "test/fixtures/pairs.json"

	tkns := tokens.New(cfg)
	_, swaps := pairs.New(cfg, tkns)

	eth := chain.NewEthClient(cfg)

	wd := watchdog.New(cfg, eth, swaps)
	wd.Start()
}

func TestNewArbiter(t *testing.T) {
	cfg := config.New()
	cfg.TokensFile = "test/fixtures/tokens_full.json"
	cfg.PairsFile = "test/fixtures/pairs_full.json"
	cfg.OutputFilename = "output/output_test.json"
	cfg.PrettyPrintOutput = false
	cfg.IncludeGas = true

	tkns := tokens.New(cfg)
	paths, swaps := pairs.New(cfg, tkns)

	eth := chain.NewEthClient(cfg)

	wd := watchdog.New(cfg, eth, swaps)
	wd.Start()

	arb := arbiter.New(cfg, eth, paths)
	arb.Start()

	out := ReadFile(t, cfg.OutputFilename)
	arbs := make([]arbiter.Arbitration, 0)

	err := utils.FromJSON(out, &arbs)
	if err != nil {
		t.Errorf("invalid arbitration json: %v", err)
	}

	if len(arbs) == 0 {
		t.Errorf("failed to build arbitrations")
	}
}

func ReadFile(t *testing.T, inputFile string) []byte {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return input
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

func AssertBigInt(t *testing.T, res *big.Int, expected *big.Int) {
	if res.Cmp(expected) != 0 {
		t.Errorf("got: %s - expected: %s", res, expected)
	}
}
