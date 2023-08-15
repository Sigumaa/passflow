package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client

	mutex sync.Mutex
)

type UserPosition struct {
	ID  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Message struct {
	Message string `json:"message"`
}

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
	e.GET("/ping", ping)
	e.GET("/user/:id", getUser)
	e.POST("/user", postUser)

	return e
}

func ping(c echo.Context) error {
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Logger().Printf("ping: %s", ping)

	return c.JSON(http.StatusOK, Message{Message: ping})
}

func getUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	u := c.Param("id")

	res, err := rdb.GeoRadiusByMember(context.Background(), "users", u, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithGeoHash: false,
		WithCoord:   false,
		WithDist:    false,
		Count:       100,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		return err
	}

	// このresには自分自身も含まれているので、自分自身を除外する
	for i, v := range res {
		if v.Name == u {
			res = append(res[:i], res[i+1:]...)
			break
		}
	}

	if len(res) == 0 {
		return c.JSON(http.StatusOK, Message{Message: "No Users around you"})
	}

	// GeoRadiusByMemberから返ってくるデータはName以外0が返ってくる(俺の実装ミスかもしれない)
	// なので、Name以外を0から正しい値に変換する
	r := make([]UserPosition, len(res))
	for i, v := range res {
		name := v.Name

		ll, err := rdb.GeoPos(context.Background(), "users", name).Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		lat := ll[0].Latitude
		lon := ll[0].Longitude

		r[i] = UserPosition{
			ID:  name,
			Lat: lat,
			Lon: lon,
		}
	}

	c.Logger().Printf("ID: %s, Lat: %f, Lon: %f\n", u, r[0].Lat, r[0].Lon)

	return c.JSON(http.StatusOK, r)

}

func postUser(c echo.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	u := new(UserPosition)
	if err := c.Bind(u); err != nil {
		return err
	}

	if u.ID == "" || u.Lat == 0 || u.Lon == 0 {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	if err := rdb.GeoAdd(context.Background(), "users", &redis.GeoLocation{
		Name:      u.ID,
		Latitude:  u.Lat,
		Longitude: u.Lon,
	}).Err(); err != nil {
		return err
	}

	// logを出力する
	c.Logger().Printf("ID: %s, Lat: %f, Lon: %f\n", u.ID, u.Lat, u.Lon)

	// 200 u を返す
	return c.JSON(http.StatusOK, u)
}
