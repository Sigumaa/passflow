package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func main() {
	rdb = initRedis()
	defer rdb.Close()

	e := initEcho()

	go func() {
		addr := ":1323"
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initRedis() *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	return r
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ping", ping)

	return e
}

func ping(c echo.Context) error {
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ping)
}
