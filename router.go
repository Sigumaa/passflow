package main

import (
	"github.com/Sigumaa/passflow/handler"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func setRoutes(rdb *redis.Client) *echo.Echo {
	e := echo.New()
	e.GET("/ping", handler.Ping(rdb))
	e.GET("/user/:id", handler.GetUser(rdb))
	e.POST("/user", handler.PostUser(rdb))

	return e
}
