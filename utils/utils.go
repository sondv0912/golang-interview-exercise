package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"golang-interview-exercise/api/transactions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func getTransactionData(trans *types.Transaction) *transactions.TransactionsType {
	chainId := trans.ChainId()
	var sender common.Address
	var err error
	if trans.Type() == 0 {
		sender, err = types.Sender(types.NewEIP155Signer(chainId), trans)
	} else if trans.Type() == 2 { // EIP-1559 transaction
		sender, err = types.Sender(types.NewLondonSigner(chainId), trans) // Mainnet chain ID = 1
	}
	if err != nil {
		log.Fatal(err, `for getting the sender`)
	}
	to := trans.To()
	if to != nil {
		return &transactions.TransactionsType{
			Hash:     trans.Hash().Hex(),
			Nonce:    trans.Nonce(),
			From:     sender.Hex(),
			To:       trans.To().String(),
		}
	}

	return &transactions.TransactionsType{
		Hash:     trans.Hash().Hex(),
		Nonce:    trans.Nonce(),
		From:     sender.Hex(),
		To:       "",
	}
}

func CheckNewBlock(mongodb *mongo.Client) (*ethclient.Client, error) {
	client, _ := InitEthereumClient()
	var lastBlockNumber *big.Int
	collectionAddresses := mongodb.Database("mydatabase").Collection("addresses")
	collectionTransaction := mongodb.Database("mydatabase").Collection("transaction")

	for {
		currentBlockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Fatalf("Failed to retrieve block number: %v", err)
		}
		if lastBlockNumber == nil || currentBlockNumber > lastBlockNumber.Uint64() {
			var addresses []string
			cursor, err := collectionAddresses.Find(context.Background(), bson.M{})
			if err != nil {
				log.Fatalf("Failed to find documents: %v", err)
			}
			for cursor.Next(context.Background()) {
				var address struct {
					Address string `bson:"address"`
				}
				if err := cursor.Decode(&address); err != nil {
					log.Fatalf("Failed to decode document: %v", err)
				}

				addresses = append(addresses, address.Address)
			}
			set := make(map[string]struct{})
			for _, id := range addresses {
				set[id] = struct{}{}
			}
			lastBlockNumber = big.NewInt(int64(currentBlockNumber))
			block, _ := GetBlockByBlockNumber(client, lastBlockNumber)
			transactions := block.Transactions()
			for _, item := range transactions {
				if item != nil {
					transaction := getTransactionData(item)
					var err error
					if _, exists := set[transaction.From]; exists {
						_,err = collectionTransaction.InsertOne(context.Background(),transaction)
					}
					if _, exists := set[transaction.To]; exists {
						_,err = collectionTransaction.InsertOne(context.Background(),transaction)
					}
					if err != nil {
						log.Fatalf("Failed to decode document: %v", err)
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
