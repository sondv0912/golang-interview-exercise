package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"golang-interview-exercise/api/transactions"
	"golang-interview-exercise/database/mongodb"
	"golang-interview-exercise/ethereum"
	"golang-interview-exercise/utils/context_utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBlockByBlockNumber(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block: %v", err)
	}

	return block, nil
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
			Hash:  trans.Hash().Hex(),
			Nonce: trans.Nonce(),
			From:  sender.Hex(),
			To:    trans.To().String(),
		}
	}

	return &transactions.TransactionsType{
		Hash:  trans.Hash().Hex(),
		Nonce: trans.Nonce(),
		From:  sender.Hex(),
		To:    "",
	}
}

func CheckNewBlock() error {
	client, err := ethereum.GetClientEthereum()
	if err != nil {
		return fmt.Errorf("error initializing Ethereum client: %v", err)
	}

	collectionAddresses, err := mongodb.GetCollection("mydatabase", "addresses")
	if err != nil {
		return fmt.Errorf("error getting collection: %v", err)
	}

	collectionTransaction, err := mongodb.GetCollection("mydatabase", "transaction")
	if err != nil {
		return fmt.Errorf("error getting collection: %v", err)
	}

	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	var lastBlockNumber *big.Int

	for {
		currentBlockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve block number: %v", err)
		}

		if lastBlockNumber == nil || currentBlockNumber > lastBlockNumber.Uint64() {
			addresses, err := getAddresses(ctx, collectionAddresses)
			if err != nil {
				return err
			}

			addressSet := make(map[string]struct{})
			for _, address := range addresses {
				addressSet[address] = struct{}{}
			}

			lastBlockNumber = big.NewInt(int64(currentBlockNumber))

			block, err := GetBlockByBlockNumber(client, lastBlockNumber)
			if err != nil {
				return fmt.Errorf("failed to get block: %v", err)
			}

			err = processTransactions(ctx, block, addressSet, collectionTransaction)
			if err != nil {
				return fmt.Errorf("failed to process transactions: %v", err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func getAddresses(ctx context.Context, collection *mongo.Collection) ([]string, error) {
	var addresses []string
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var address struct {
			Address string `bson:"address"`
		}
		if err := cursor.Decode(&address); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		addresses = append(addresses, address.Address)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return addresses, nil
}

func processTransactions(ctx context.Context, block *types.Block, addressSet map[string]struct{}, collection *mongo.Collection) error {
	transactions := block.Transactions()
	for _, item := range transactions {
		if item == nil {
			continue
		}

		transaction := getTransactionData(item)
		var err error

		if _, exists := addressSet[transaction.From]; exists {
			_, err = collection.InsertOne(ctx, transaction)
		}
		if _, exists := addressSet[transaction.To]; exists {
			_, err = collection.InsertOne(ctx, transaction)
		}

		if err != nil {
			return fmt.Errorf("failed to insert transaction: %v", err)
		}
	}
	return nil
}
