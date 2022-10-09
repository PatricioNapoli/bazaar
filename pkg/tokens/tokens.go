package tokens

import (
	"github.com/PatricioNapoli/bazaar/pkg/utils"
	"log"
)

type Token struct {
	Address  string
	Symbol   string
	Name     string
	Decimals string
}

type Tokens struct {
	tokens map[string]Token
}

// GetTokens loads a tokens file into an indexed map through token's address.
func GetTokens(file string) map[string]Token {
	log.Printf("loading token info in %s", file)

	f, err := utils.ReadFile(file)
	if err != nil {
		log.Panicf("failed when reading file %s - %v", file, err)
	}

	tokens := map[string]Token{}
	tokenArr := make([]Token, 0)

	err = utils.FromJSON(f, &tokenArr)
	if err != nil {
		log.Panicf("failed deserializing token file to array - %v", err)
	}

	for _, t := range tokenArr {
		tokens[t.Address] = t
	}

	return tokens
}
