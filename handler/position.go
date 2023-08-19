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

type ResPos struct {
	Cnt   int        `json:"cnt"`
	Users []Position `json:"users"`
}

type Message struct {
	Message string `json:"message"`
}

// reqでid,lat,lonが来る
// resですれ違い人数(cnt),すれ違った人のid,lat,lonを配列で返す
// "cnt": 2,
// "users": [
//
//	{
//	  "id": "user1",
//	  "lat": 35.123456,
//	  "lon": 135.123456
//	},
//	...
//
// ]
// すれ違い人数が0の場合は
// "cnt": 0,
// "users": []
// を返す
// すれ違ったユーザーの名前をFriendに追加する
// 一度すれ違ったユーザーはFriendに追加しないかつ、すれ違い人数にも含まない
func SetUserPos(rad float64, rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

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

		c.Logger().Printf("Added: ID: %s, Lat: %f, Lon: %f", u.ID, u.Lat, u.Lon)

		res, err := getNearbyUsers(u.ID, rad, rdb)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if len(res) == 0 {
			return c.JSON(http.StatusOK, ResPos{
				Cnt:   0,
				Users: []Position{},
			})
		}

		name := u.ID
		friend := GetFriends(name)

		resIDs := []string{}
		for _, v := range res {
			resIDs = append(resIDs, v.ID)
		}

		diffIDs := diff(friend, resIDs)

		if len(diffIDs) == 0 {
			return c.JSON(http.StatusOK, ResPos{
				Cnt:   0,
				Users: []Position{},
			})
		}

		for _, v := range diffIDs {
			AddFriend(name, v)
			IncrementCollections(name, v)
			IncrementRecord(name)
		}

		// diffIDsのみを返す。
		// diffIDsに存在するIDのlat,lonをresから取得し、nresに追加する
		nres := []Position{}
		for _, v := range diffIDs {
			for _, vv := range res {
				if v == vv.ID {
					nres = append(nres, vv)
				}
			}
		}

		return c.JSON(http.StatusOK, ResPos{
			Cnt:   len(nres),
			Users: nres,
		})
	}
}

func getNearbyUsers(u string, rad float64, rdb *redis.Client) ([]Position, error) {
	res, err := rdb.GeoRadiusByMember(context.Background(), "users", u, &redis.GeoRadiusQuery{
		Radius:      rad,
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

func diff(a, b []string) []string {
	m := make(map[string]bool)

	for _, v := range a {
		m[v] = true
	}

	diff := []string{}
	for _, v := range b {
		if !m[v] {
			diff = append(diff, v)
		}
	}

	return diff
}
