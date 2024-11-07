package handlers

import (
	"context"
	"log"
	"math/big"
	"net/http"

	"golang-interview-exercise/utils"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

func GetCurrentBlock(c echo.Context) error {
    client, err := ethclient.Dial("https://mainnet.infura.io/v3/afd3d007db6f48fb946468a2877b5151")
	if err != nil {
		log.Fatal(err)
	}

	// Lấy block mới nhất
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
        return echo.NewHTTPError(400, err)
    }

    block, err :=  utils.ParseBlock(client, big.NewInt(int64(blockNumber)))
    if err != nil {
        return echo.NewHTTPError(400, err)
    }

    blockData := map[string]interface{}{
        "block_number": blockNumber,
        "block_info": block,
    }

    // Trả về dữ liệu dưới dạng JSON
    return c.JSON(http.StatusOK, blockData)
}