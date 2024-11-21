package ethereum

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethClient *ethclient.Client
	once      sync.Once
	errInit   error
)

func GetClientEthereum() (*ethclient.Client, error) {
	once.Do(func() {
		ethClient, errInit = ethclient.Dial("https://mainnet.infura.io/v3/afd3d007db6f48fb946468a2877b5151")
		if errInit != nil {
			errInit = fmt.Errorf("failed to connect to Ethereum client: %w", errInit)
		}
	})
	return ethClient, errInit
}
