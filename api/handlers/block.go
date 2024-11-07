package handlers

import (
	"context"
	"log"
	"math/big"
	"net/http"

	"golang-interview-exercise/utils"

	"github.com/labstack/echo/v4"
)

func GetCurrentBlock(c echo.Context) error {
	client, err := utils.InitEthereumClient()
	if err != nil {
		log.Printf("Error initializing Ethereum client: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to connect to Ethereum client")
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve latest block number")
	}

	block, err := utils.GetBlockByBlockNumber(client, big.NewInt(int64(blockNumber)))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve block data")
	}

	transaction := block.Transactions()
	transactionCount := len(transaction)

	blockData := map[string]interface{}{
		"block_number":      blockNumber,
		"transaction_count": transactionCount,
		"transaction":       transaction,
	}

	return c.JSON(http.StatusOK, blockData)
}
