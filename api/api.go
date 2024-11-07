package api

import (
	"golang-interview-exercise/api/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterAPI(e *echo.Echo) error {
	e.GET("/block", handlers.GetCurrentBlock)
	e.GET("/subscribe", handlers.GetSubscribe)
	e.GET("/transaction", handlers.GetTransaction)

	return nil
}