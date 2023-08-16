package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func Ping(rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ping, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.Logger().Printf("ping: %s", ping)

		return c.JSON(http.StatusOK, Message{Message: ping})
	}
}
