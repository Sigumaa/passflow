package handler

import (
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
	Friends []string `json:"friends"`
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
			return c.JSON(http.StatusBadRequest, "Bad Request")
		}

		store[u.ID] = UserInfo{
			ReqUserInfo: *u,
		}

		c.Logger().Printf("store: %v", store)

		return c.JSON(http.StatusOK, u)
	}
}

func GetUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if _, ok := store[id]; !ok {
			return c.JSON(404, "Not Found")
		}

		friends := GetFriends(id)
		if friends == nil {
			store[id] = UserInfo{
				ReqUserInfo: store[id].ReqUserInfo,
				Friends:     []string{},
			}
		}

		return c.JSON(http.StatusOK, store[id])
	}
}

func GetFriends(id string) []string {
	return store[id].Friends
}

func AddFriend(id string, friend string) {
	before := store[id].Friends
	store[id] = UserInfo{
		ReqUserInfo: store[id].ReqUserInfo,
		Friends:     append(before, friend),
	}
}
