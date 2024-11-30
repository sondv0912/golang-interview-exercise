package main

import (
	"fmt"
	"log"

	"golang-interview-exercise/api"
	"golang-interview-exercise/database/mongodb"
	"golang-interview-exercise/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	_, err := mongodb.GetMongoClient()
	if err != nil {
		log.Fatalf("Error initializing MongoDB client: %v", err)
	}

	if err := api.RegisterAPI(e); err != nil {
		fmt.Println("Error register api: %w", err)
	}

	go utils.CheckNewBlock()

	e.Logger.Fatal(e.Start(":1323"))
}
