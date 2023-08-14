package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	rdb = r
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ping", func(c echo.Context) error {
		ping, err := rdb.Ping(ctx).Result()
		if err != nil {
			c.JSON(500, err.Error())
		}
		return c.JSON(200, ping)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
