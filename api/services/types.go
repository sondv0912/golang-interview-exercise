package services

import "github.com/ethereum/go-ethereum/common"

type TransactionsType struct {
	Hash     common.Hash
	Value    string
	Gas      uint64
	GasPrice string
	Nonce    uint64
	Data     string
	From     string
	To       string
}
