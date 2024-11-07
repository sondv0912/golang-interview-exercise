package main

import (
	"fmt"

	"golang-interview-exercise/api"
	"golang-interview-exercise/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	if err := api.RegisterAPI(e); err != nil {
		fmt.Println("Error register api: %w", err)
	}

	go utils.CheckNewBlock()

	e.Logger.Fatal(e.Start(":1323"))
}
