package main

import (
	"context"
	"fmt"
	"log"

	"golang-interview-exercise/api"
	"golang-interview-exercise/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	e := echo.New()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

	err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mongodb", client)
			return next(c)
		}
	})


	if err := api.RegisterAPI(e); err != nil {
		fmt.Println("Error register api: %w", err)
	}

	go utils.CheckNewBlock(client)

	e.Logger.Fatal(e.Start(":1323"))
}
