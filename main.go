package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
}

func main() {
	rdb = initRedis()
	defer rdb.Close()

	e := initEcho()

	go func() {
		echoAddr := fmt.Sprintf(":%s", os.Getenv("ECHO_ADDR"))
		if err := e.Start(echoAddr); err != nil && err != http.ErrServerClosed {
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
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_ADDR"))
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	pool, _ := strconv.Atoi(os.Getenv("REDIS_POOL"))

	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PWD"),
		DB:       db,
		PoolSize: pool,
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
