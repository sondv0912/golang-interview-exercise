package block

import (
	"log"
	"math/big"
	"net/http"
	"time"

	"golang-interview-exercise/ethereum"
	"golang-interview-exercise/utils"
	"golang-interview-exercise/utils/context_utils"

	"github.com/labstack/echo/v4"
)

func GetCurrentBlock(c echo.Context) error {
	client, err := ethereum.GetClientEthereum()
	if err != nil {
		log.Printf("Error initializing Ethereum client: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to connect to Ethereum client")
	}

	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve latest block number")
	}

	block, err := utils.GetBlockByBlockNumber(client, big.NewInt(int64(blockNumber)))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve block data")
	}

	transaction := block.Transactions()

	blockData := map[string]interface{}{
		"block_number":      blockNumber,
		"transaction_count": len(transaction),
		"transaction":       transaction,
	}

	return c.JSON(http.StatusOK, blockData)
}
