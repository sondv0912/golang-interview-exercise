//go:build wireinject
// +build wireinject

package blockchain

import (
	"golang-interview-exercise/api/blockchain/blockchaindata"
	"golang-interview-exercise/api/blockchain/blockchainservice"

	"github.com/google/wire"
)

func New() (*blockchainservice.Service, error) {
	wire.Build(
		blockchaindata.New,
		blockchainservice.New,
	)

	return &blockchainservice.Service{}, nil

}
