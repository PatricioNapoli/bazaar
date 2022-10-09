package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	APIKey            string
	ExcludeDeadTokens bool
	IncludeFees       bool
}

func NewConfig() Config {
	c := Config{
		APIKey:            "",
		ExcludeDeadTokens: true,
		IncludeFees:       true,
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
