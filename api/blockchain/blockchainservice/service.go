package blockchainservice

import (
	"context"

	"golang-interview-exercise/api/blockchain/blockchainmodel"
)

type BlockchainRepository interface {
	GetAddress(ctx context.Context, address string) (*blockchainmodel.SubscribeRequestBody, error)
	AddNewAddress(ctx context.Context, address string) (*blockchainmodel.SubscribeRequestBody, error)
	AddNewTransaction(ctx context.Context, transaction blockchainmodel.TransactionsType) (*blockchainmodel.TransactionsType, error)
	GetTransactions(ctx context.Context, transaction blockchainmodel.TransactionsType) (*[]blockchainmodel.TransactionsType, error)
}

type Service struct {
	blockChainRepo BlockchainRepository
}

func New(blockchainRepository *BlockchainRepository) (*Service, error) {
	svc := Service{
		blockChainRepo: *blockchainRepository,
	}

	return &svc, nil
}
