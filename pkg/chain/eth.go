package chain

import (
	"fmt"
	"github.com/PatricioNapoli/bazaar/pkg/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type EthClient struct {
	Client *ethclient.Client
}

// NewEthClient creates an ethereum client from provided endpoint and config for API key
func NewEthClient(cfg config.Config) EthClient {
	log.Printf("connecting to ETH mainnet through %s", cfg.InfuraEndpoint)

	if len(cfg.APIKey) == 0 {
		log.Panic("infura key not set or empty")
	}

	client, err := ethclient.Dial(fmt.Sprintf("%s%s", cfg.InfuraEndpoint, cfg.APIKey))
	if err != nil {
		log.Panicf("error when dialing eth mainnet: %v", err)
	}

	log.Println("connected to ETH mainnet")

	return EthClient{Client: client}
}
