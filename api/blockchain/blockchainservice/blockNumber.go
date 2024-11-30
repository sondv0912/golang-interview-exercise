package blockchainservice

import (
	"context"
	"fmt"
	"math/big"

	"golang-interview-exercise/ethereum"
	"golang-interview-exercise/utils"
)

func (svc *Service) GetCurrentBlock(ctx context.Context) (map[string]interface{}, error) {
	client, err := ethereum.GetClientEthereum()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to Ethereum client: %w", err)
	}

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve latest block number: %w", err)
	}

	block, err := utils.GetBlockByBlockNumber(client, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve block data: %w", err)
	}

	transaction := block.Transactions()

	blockData := map[string]interface{}{
		"block_number":      blockNumber,
		"transaction_count": len(transaction),
		"transaction":       transaction,
	}

	return blockData, nil
}
