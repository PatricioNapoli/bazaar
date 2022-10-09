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

func NewEthClient(cfg config.Config, endpoint string) EthClient {
	log.Printf("connecting to ETH mainnet through %s", endpoint)

	if len(cfg.APIKey) == 0 {
		log.Panic("infura key not set or empty")
	}

	client, err := ethclient.Dial(fmt.Sprintf("%s%s", endpoint, cfg.APIKey))
	if err != nil {
		log.Panicf("error when dialing eth mainnet: %v", err)
	}

	log.Println("connected to ETH mainnet")

	return EthClient{Client: client}
}
