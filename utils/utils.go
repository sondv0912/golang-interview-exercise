package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ParseBlock(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	// Lấy block thông qua blockNumber
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve block: %v", err)
	}

	// Trả về block đã parse được
	return block, nil
}