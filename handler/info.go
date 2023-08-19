package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// TODO: 永続化する　当たり前
var (
	store = make(map[string]UserInfo)
)

type ReqUserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Like    string `json:"like"`
	Dislike string `json:"dislike"`
	From    string `json:"from"`
}

type ckey string

const (
	Likek    ckey = "like"
	Dislikek ckey = "dislike"
	Fromk    ckey = "from"
)

type UserInfo struct {
	ReqUserInfo
	Friends    []string                `json:"friends"`
	Collection map[ckey]map[string]int `json:"collection"`
	Record     map[string]int          `json:"record"`
}

func SetUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		mutex.Lock()
		defer mutex.Unlock()
		u := new(ReqUserInfo)
		if err := c.Bind(u); err != nil {
			return err
		}

		if u.ID == "" || u.Name == "" {
			return c.JSON(http.StatusBadRequest, Message{Message: "Bad Request"})
		}

		store[u.ID] = UserInfo{
			ReqUserInfo: *u,
		}

		friends := GetFriends(u.ID)
		if len(friends) == 0 {
			store[u.ID] = UserInfo{
				ReqUserInfo: *u,
				Friends:     []string{},
				Collection:  make(map[ckey]map[string]int),
				Record:      make(map[string]int),
			}
		}

		c.Logger().Printf("store: %v", store)

		return c.JSON(http.StatusOK, u)
	}
}

func GetUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if _, ok := store[id]; !ok {
			return c.JSON(http.StatusBadRequest, Message{Message: "Not Found"})
		}

		friends := GetFriends(id)

		c.Logger().Printf("friends: %v", friends)

		if len(friends) == 0 {
			return c.JSON(http.StatusOK, store[id])
		}

		return c.JSON(http.StatusOK, store[id])
	}
}

func GetFriends(id string) []string {
	friends := store[id].Friends
	return friends
}

func AddFriend(id string, friend string) {
	mutex.Lock()
	defer mutex.Unlock()
	before := GetFriends(id)
	store[id] = UserInfo{
		ReqUserInfo: store[id].ReqUserInfo,
		Friends:     append(before, friend),
		Collection:  store[id].Collection,
		Record:      store[id].Record,
	}
}

func IncrementCollections(id string, passedID string) {
	for _, v := range []ckey{Likek, Dislikek, Fromk} {
		IncrementCollection(id, v, passedID)
	}
}

func IncrementCollection(id string, key ckey, passedID string) {
	mutex.Lock()
	defer mutex.Unlock()

	collection := GetCollection(id)
	if _, ok := collection[key]; !ok {
		collection[key] = make(map[string]int)
	}

	var t string
	switch key {
	case Likek:
		t = store[passedID].Like
	case Dislikek:
		t = store[passedID].Dislike
	case Fromk:
		t = store[passedID].From
	}
	collection[key][t]++

	store[id] = UserInfo{
		ReqUserInfo: store[id].ReqUserInfo,
		Friends:     store[id].Friends,
		Collection:  collection,
		Record:      store[id].Record,
	}

}

func GetCollection(id string) map[ckey]map[string]int {
	collection := store[id].Collection
	return collection
}

func GetRecord(id string) map[string]int {
	record := store[id].Record
	return record
}

func IncrementRecord(id string) {
	mutex.Lock()
	defer mutex.Unlock()

	today := time.Now().Format("2006-01-02")
	record := GetRecord(id)
	record[today]++
	store[id] = UserInfo{
		ReqUserInfo: store[id].ReqUserInfo,
		Friends:     store[id].Friends,
		Collection:  store[id].Collection,
		Record:      record,
	}
}
