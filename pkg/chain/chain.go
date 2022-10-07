package chain

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	Client *ethclient.Client
}

func NewEthClient(endpoint string) EthClient {
	log.Printf("connecting to ETH mainnet through %s", endpoint)

	key := os.Getenv("INFURA_KEY")

	if len(key) == 0 {
		log.Panic("infura key not set or empty")
	}

	client, err := ethclient.Dial(fmt.Sprintf("%s%s", endpoint, key))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to eth mainnet")

	return EthClient{Client: client}
}
