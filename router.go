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
		pos.GET("/:id", dummyHandler())
	}

	info := e.Group("/info")
	{
		info.GET("/:id", dummyHandler())
		info.POST("/:id", dummyHandler())
	}

	e.GET("/user/:id", handler.GetUser(rdb))
	e.POST("/user", handler.PostUser(rdb))

	return e
}

func dummyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		return c.String(200, "not implemented!"+id)
	}
}
