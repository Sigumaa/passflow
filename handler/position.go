package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var (
	mutex sync.Mutex
)

type Position struct {
	ID  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Message struct {
	Message string `json:"message"`
}

func GetUser(rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		mutex.Lock()
		defer mutex.Unlock()

		u := c.Param("id")

		res, err := getNearbyUsers(u, rdb)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if len(res) == 0 {
			return c.JSON(http.StatusOK, Message{Message: "No Users around you"})
		}

		c.Logger().Printf("ID: %s, Lat: %f, Lon: %f\n", u, res[0].Lat, res[0].Lon)

		return c.JSON(http.StatusOK, res)
	}
}

func PostUser(rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		mutex.Lock()
		defer mutex.Unlock()

		u := new(Position)
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

		c.Logger().Printf("ID: %s, Lat: %f, Lon: %f\n", u.ID, u.Lat, u.Lon)

		return c.JSON(http.StatusOK, u)
	}
}

func getNearbyUsers(u string, rdb *redis.Client) ([]Position, error) {
	res, err := rdb.GeoRadiusByMember(context.Background(), "users", u, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithGeoHash: false,
		WithCoord:   true,
		WithDist:    false,
		Count:       100,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}

	var userPositions []Position
	for _, v := range res {
		if v.Name == u {
			continue
		}

		up, err := getUserPosition(v.Name, rdb)
		if err != nil {
			return nil, err
		}

		userPositions = append(userPositions, up)
	}

	return userPositions, nil
}

func getUserPosition(u string, rdb *redis.Client) (Position, error) {
	ll, err := rdb.GeoPos(context.Background(), "users", u).Result()
	if err != nil {
		return Position{}, err
	}

	return Position{
		ID:  u,
		Lat: ll[0].Latitude,
		Lon: ll[0].Longitude,
	}, nil
}
