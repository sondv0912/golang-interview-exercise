package blockchainapi

import (
	"fmt"
	"net/http"
	"time"

	"golang-interview-exercise/utils/context_utils"

	"github.com/labstack/echo/v4"
)

func (api *BlockchainAPI) GetCurrentBlockNumber(context echo.Context) error {
	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	data, err := api.svc.GetCurrentBlock(ctx)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get block number: %v", err),
		})
	}

	return context.JSON(http.StatusOK, data)

}
