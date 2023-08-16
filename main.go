package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
}

func main() {
	rdb := initRedis()
	defer rdb.Close()

	e := setRoutes(rdb)

	go startServer(e)

	waitForInterrupt(e)
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
