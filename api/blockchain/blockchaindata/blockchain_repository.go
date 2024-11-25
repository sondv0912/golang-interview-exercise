package blockchaindata

import (
	"context"
	"fmt"

	"golang-interview-exercise/api/blockchain/blockchainmodel"
	"golang-interview-exercise/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

type BlockchainRepository struct {
}

func New() *BlockchainRepository {
	return &BlockchainRepository{}
}

func (r *BlockchainRepository) FindByAddress(ctx context.Context, address string) (*blockchainmodel.SubscribeRequestBody, error) {

	collection, err := mongodb.GetCollection("mydatabase", "addresses")
	if err != nil {
		return nil, fmt.Errorf("error finding transaction: %v", err)
	}

	var result blockchainmodel.SubscribeRequestBody
	err = collection.FindOne(ctx, bson.M{"address": address}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error finding address: %v", err)
	}

	return &result, nil
}
