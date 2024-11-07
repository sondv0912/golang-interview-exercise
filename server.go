package main

import (
	"fmt"

	"golang-interview-exercise/api"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	if err := api.RegisterAPI(e); err != nil {
		fmt.Println("Error register api: %w", err)
	}

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	e.Logger.Fatal(e.Start(":1323"))
}