package main

import (
	"context"

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
	ping, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	println(ping)
}
