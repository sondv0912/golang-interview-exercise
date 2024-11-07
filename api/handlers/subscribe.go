package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetSubscribe(c echo.Context) error {
    // Dữ liệu block giả lập (có thể thay bằng dữ liệu thực từ blockchain)
    blockData := map[string]interface{}{
        "block_number": 123456,
        "block_hash":   "0xabc123...",
        "timestamp":    "2024-11-07T12:34:56Z",
    }

    // Trả về dữ liệu dưới dạng JSON
    return c.JSON(http.StatusOK, blockData)
}