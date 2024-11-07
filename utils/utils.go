package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"golang-interview-exercise/api/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBlockByBlockNumber(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}

	return block, nil
}

func InitEthereumClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/afd3d007db6f48fb946468a2877b5151")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	return client, nil
}

func getTransactionData(trans *types.Transaction) *services.TransactionsType {
	chainId := trans.ChainId()
	var sender common.Address
	var err error
	if trans.Type() == 0 {
		sender, err = types.Sender(types.NewLondonSigner(chainId), trans)
	} else if trans.Type() == 2 { // EIP-1559 transaction
		sender, err = types.Sender(types.NewLondonSigner(chainId), trans) // Mainnet chain ID = 1
	}
	if err != nil {
		log.Fatal(err, `for getting the sender`)
	}

	return &services.TransactionsType{
		Hash:     trans.Hash(),
		Value:    trans.Value().String(),
		Gas:      trans.Gas(),
		GasPrice: trans.GasPrice().String(),
		Nonce:    trans.Nonce(),
		From:     sender.Hex(),
		To:       trans.To().String(),
	}
}

func CheckNewBlock() (*ethclient.Client, error) {
	client, _ := InitEthereumClient()
	var lastBlockNumber *big.Int
	for {
		currentBlockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Fatalf("Failed to retrieve block number: %v", err)
		}

		if lastBlockNumber == nil || currentBlockNumber > lastBlockNumber.Uint64() {
			lastBlockNumber = big.NewInt(int64(currentBlockNumber))
			block, _ := GetBlockByBlockNumber(client, lastBlockNumber)
			transaction := block.Transactions()
			for i := 0; i < transaction.Len(); i++ {
				if transaction[i] != nil {
					fmt.Println(getTransactionData(transaction[i]).From)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
