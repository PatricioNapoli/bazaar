package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	APIKey            string
	InfuraEndpoint    string
	ExcludeDeadTokens bool
	IncludeFees       bool
	TokensFile        string
	PairsFile         string
	OutputFilename    string
	WETHAddr          string
	ReservesAddr      string
	DexSwapGas        float64
	InitialWETH       float64
}

// NewConfig creates a configuration struct from environment vars.
func NewConfig() Config {
	log.Printf("loading config from env vars")

	c := Config{
		APIKey:            "",
		InfuraEndpoint:    "https://mainnet.infura.io/v3/",
		ExcludeDeadTokens: true,
		IncludeFees:       false,
		TokensFile:        "assets/tokens.json",
		PairsFile:         "assets/uni_sushi_paths.json",
		OutputFilename:    "output/output.json",
		WETHAddr:          "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		ReservesAddr:      "0x416355755f32b2710ce38725ed0fa102ce7d07e6",
		DexSwapGas:        150000.0,
		InitialWETH:       1.0,
	}

	if env := os.Getenv("BAZAAR_INFURA_KEY"); envIsValid(env) {
		c.APIKey = env
	}

	if env := os.Getenv("BAZAAR_EXCLUDE_DEAD_TOKENS"); envIsValid(env) {
		c.ExcludeDeadTokens = intIsTrue(env)
	}

	if env := os.Getenv("BAZAAR_INCLUDE_FEES"); envIsValid(env) {
		c.IncludeFees = intIsTrue(env)
	}

	if env := os.Getenv("BAZAAR_OUTPUT_FILENAME"); envIsValid(env) {
		c.OutputFilename = env
	}

	return c
}

func envIsValid(s string) bool {
	return len(s) > 0
}

func intIsTrue(s string) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panicf("env var invalid format, should be 1 or 0: %s", s)
	}

	return i == 1
}
