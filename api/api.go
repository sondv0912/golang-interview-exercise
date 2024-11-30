package api

import (
	"fmt"

	"golang-interview-exercise/api/block"
	"golang-interview-exercise/api/blockchain/blockchainapi"
	"golang-interview-exercise/api/subscribe"
	"golang-interview-exercise/api/transactions"

	"github.com/labstack/echo/v4"
)

func RegisterAPI(e *echo.Echo) error {
	e.GET("/block", block.GetCurrentBlock)
	e.POST("/subscribe", subscribe.PostSubscribe)
	e.GET("/transaction/:address", transactions.GetTransaction)

	if err := blockchainapi.Init(e.Group("/blockchain")); err != nil {
		return fmt.Errorf("failed to init collections module: %w", err)
	}

	return nil
}
