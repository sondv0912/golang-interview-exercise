package transactions

import (
	"net/http"
	"time"

	"golang-interview-exercise/api/subscribe"
	"golang-interview-exercise/database/mongodb"
	"golang-interview-exercise/utils/context_utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransaction(c echo.Context) error {

	address := c.Param("address")

	collectionAddresses, err := mongodb.GetCollection("mydatabase", "addresses")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting collection: "+err.Error())
	}

	collectionTransaction, err := mongodb.GetCollection("mydatabase", "transaction")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting collection: "+err.Error())
	}

	filter := bson.M{
		"$or": []bson.M{
			{"from": address},
			{"to": address},
		},
	}

	var searchAddress subscribe.SubscribeRequestBody

	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	err = collectionAddresses.FindOne(ctx, bson.M{"address": address}).Decode(&searchAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, "Address not found")
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Error querying address: "+err.Error())
	}

	var result []TransactionsType

	cursor, err := collectionTransaction.Find(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error querying transactions: "+err.Error())
	}

	for cursor.Next(ctx) {
		var item TransactionsType
		if err := cursor.Decode(&item); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode transaction: "+err.Error())
		}

		result = append(result, item)
	}

	return c.JSON(http.StatusOK, result)
}
