package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PostSubscribe(c echo.Context) error {
	address := c.FormValue("address")

	if address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address notfound")
	}

	dataResponse := map[string]interface{}{
		"address": address,
	}

	return c.JSON(http.StatusOK, dataResponse)
}
