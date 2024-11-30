package subscribe

import (
	"net/http"
	"time"

	"golang-interview-exercise/database/mongodb"
	"golang-interview-exercise/utils/context_utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostSubscribe(c echo.Context) error {
	req := &SubscribeRequestBody{}

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address notfound")
	}

	collection, err := mongodb.GetCollection("mydatabase", "addresses")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting collection: "+err.Error())
	}

	ctx, cancel := context_utils.CreateTimeoutContext(5 * time.Second)
	defer cancel()

	var searchAddress SubscribeRequestBody

	err = collection.FindOne(ctx, bson.M{"address": req.Address}).Decode(&searchAddress)
	if err != nil && err != mongo.ErrNoDocuments {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error querying address: "+err.Error())
	}

	if searchAddress.Address != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address has existed")
	}

	_, err = collection.InsertOne(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert address: "+err.Error())
	}

	return c.JSON(http.StatusOK, req)
}
