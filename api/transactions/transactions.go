package transactions

import (
	"context"
	"log"
	"net/http"

	"golang-interview-exercise/api/subscribe"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransaction(c echo.Context) error {
   mongodb := c.Get("mongodb").(* mongo.Client)

   address := c.Param("address")
   collectionAddresses := mongodb.Database("mydatabase").Collection("addresses")
   collectionTransaction := mongodb.Database("mydatabase").Collection("transaction")

   filter := bson.M{
    "$or": []bson.M{
        {"from": address},
        {"to": address},
    },
    }

    var searchAddress subscribe.SubscribeRequestBody

   err := collectionAddresses.FindOne(context.Background(), bson.M{"address": address}).Decode(&searchAddress)
   if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
    }

   var result []TransactionsType


   cursor,err := collectionTransaction.Find(context.Background(), filter)
   if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err)
    }

    for cursor.Next(context.Background()) {
        var item TransactionsType
        if err := cursor.Decode(&item); err != nil {
            log.Fatalf("Failed to decode document: %v", err)
        }

        result = append(result, item)
    }

    return c.JSON(http.StatusOK, result)
}