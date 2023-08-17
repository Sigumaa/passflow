package main

import (
	"github.com/Sigumaa/passflow/handler"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func setRoutes(rdb *redis.Client) *echo.Echo {
	e := echo.New()
	e.GET("/ping", handler.Ping(rdb))

	pos := e.Group("/pos")
	{
		pos.POST("", handler.SetUser(5, rdb))
	}

	info := e.Group("/info")
	{
		info.GET("", dummyHandler())
		info.POST("", dummyHandler())
	}

	return e
}

func dummyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "not implemented!")
	}
}
