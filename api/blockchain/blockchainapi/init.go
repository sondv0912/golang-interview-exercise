package blockchainapi

import (
	"fmt"

	"golang-interview-exercise/api/blockchain"
	"golang-interview-exercise/api/blockchain/blockchainservice"

	"github.com/labstack/echo/v4"
)

type BlockchainAPI struct {
	svc *blockchainservice.Service
}

func (m *BlockchainAPI) initService() error {
	var err error
	m.svc, err = blockchain.New()
	if err != nil {
		return fmt.Errorf("failed to create market service: %w", err)
	}

	return nil
}

func Init(g *echo.Group) error {
	api := &BlockchainAPI{}

	if err := api.initService(); err != nil {
		return err
	}

	g.GET("/blockNumber", api.GetCurrentBlockNumber)

	return nil
}
