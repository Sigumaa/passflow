package handler

import (
	"maps"
	"net/http"

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

type UserInfo struct {
	ReqUserInfo
	Friends        []string       `json:"friends"`
	LikeCollection map[string]int `json:"likeCollection"`
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
				ReqUserInfo:    *u,
				Friends:        []string{},
				LikeCollection: make(map[string]int),
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
	like := GetLikeCollection(id)
	store[id] = UserInfo{
		ReqUserInfo:    store[id].ReqUserInfo,
		Friends:        append(before, friend),
		LikeCollection: maps.Clone(like),
	}
}

func GetLikeCollection(id string) map[string]int {
	likeCollection := store[id].LikeCollection
	return likeCollection
}

func IncrementLikeCollection(id string, passedID string) {
	mutex.Lock()
	defer mutex.Unlock()
	likeCollection := GetLikeCollection(id)
	likeLanguage := store[passedID].Like
	likeCollection[likeLanguage]++
	store[id] = UserInfo{
		ReqUserInfo:    store[id].ReqUserInfo,
		Friends:        store[id].Friends,
		LikeCollection: likeCollection,
	}

}
