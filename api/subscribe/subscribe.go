package subscribe

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostSubscribe(c echo.Context) error {
	req := new(SubscribeRequestBody)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if req.Address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address notfound")
	}

	mongodb := c.Get("mongodb").(*mongo.Client)

	collection := mongodb.Database("mydatabase").Collection("addresses")

	var searchAddress SubscribeRequestBody
	collection.FindOne(context.Background(), bson.M{"address": req.Address}).Decode(&searchAddress)

	if searchAddress.Address != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address has existed")
	}

	_, err := collection.InsertOne(context.Background(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, req)
}
