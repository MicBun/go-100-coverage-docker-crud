package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MicBun/go-100-coverage-docker-crud/core"
	"github.com/MicBun/go-100-coverage-docker-crud/service"
	"github.com/MicBun/go-100-coverage-docker-crud/util/jwtAuth"
	"github.com/MicBun/go-100-coverage-docker-crud/web"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var adminHeader = map[string]string{
	"Authorization": "Bearer " + func() string {
		token, _ := jwtAuth.GenerateToken(1)
		return token
	}(),
}

var userHeader = map[string]string{
	"Authorization": "Bearer " + func() string {
		token, _ := jwtAuth.GenerateToken(2)
		return token
	}(),
}

func TestHelloEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodGet, "/hello", nil)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Message string
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Hello", resp.Message)
	})
}

func TestRegisterEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodPost, "/user/register", nil, userHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, err := json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Message string
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "User registered", resp.Message)

		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"UNIQUE constraint failed: users.username"}`, w.Body.String())

		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", nil, adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodPut, "/user/update/1", nil, userHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		updatedUser := core.User{
			Username: "updated@bar.com",
			Password: "updatedPassword",
			Name:     "Updated Bar",
		}
		jsonBody, _ = json.Marshal(updatedUser)
		w, err = web.MakeRequest(c.Web, http.MethodPut, "/user/update/1", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodPut, "/user/update/1", nil, adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodPut, "/user/update/2", bytes.NewReader(jsonBody), adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"record not found"}`, w.Body.String())
	})
}

func TestDeleteEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodDelete, "/user/delete/1", nil, userHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodDelete, "/user/delete/1", nil, adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"message":"User deleted"}`, w.Body.String())

		w, err = web.MakeRequest(c.Web, http.MethodDelete, "/user/delete/1", nil, adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"record not found"}`, w.Body.String())
	})
}

func TestGetUserByIDEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodGet, "/user/get/1", nil, userHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/get/1", nil, adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/get/1", nil, adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Message string
			User    core.User
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, resp.User.Name)
		assert.Equal(t, "User retrieved", resp.Message)
	})
}

func TestGetUserByTokenEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodGet, "/user/get", nil, adminHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		user.Username = "bar@foo.com"
		user.Name = "Bar Foo"
		jsonBody, _ = json.Marshal(user)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/get", nil, userHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Message string
			User    core.User
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, resp.User.Name)
		assert.Equal(t, "User retrieved", resp.Message)
	})
}

func TestGetUsersEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		w, err := web.MakeRequest(c.Web, http.MethodGet, "/user/list", nil, userHeader)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/list", nil, adminHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"no users found"}`, w.Body.String())

		var users []core.User
		for i := 0; i < 10; i++ {
			user := core.User{
				Username: fmt.Sprintf("foo%d@bar.com", i),
				Password: "securePassword",
				Name:     "Foo Bar" + strconv.Itoa(i),
			}
			users = append(users, user)
		}
		for _, user := range users {
			jsonBody, _ := json.Marshal(user)
			w, err := web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), adminHeader)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/list", nil, adminHeader)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Message string
			Users   []core.User
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Users retrieved", resp.Message)
		assert.Equal(t, len(users), len(resp.Users))
	})
}

func TestLoginEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err := web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), map[string]string{
			"Authorization": "Bearer " + func() string {
				token, _ := jwtAuth.GenerateToken(1)
				return token
			}(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var req struct {
			Username string
			Password string
		}
		req.Username = user.Username
		req.Password = user.Password
		jsonBody, _ = json.Marshal(req)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/login", bytes.NewReader(jsonBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodPost, "/login", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		req.Password = "wrongPassword"
		jsonBody, _ = json.Marshal(req)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/login", bytes.NewReader(jsonBody))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRefreshTokenEndpoint(t *testing.T) {
	web.RunTest(func(c *service.Container) {
		user := core.User{
			Username: "foo@bar.com",
			Password: "securePassword",
			Name:     "Foo Bar",
		}
		jsonBody, _ := json.Marshal(user)
		w, err := web.MakeRequest(c.Web, http.MethodPost, "/user/register", bytes.NewReader(jsonBody), map[string]string{
			"Authorization": "Bearer " + func() string {
				token, _ := jwtAuth.GenerateToken(1)
				return token
			}(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var req struct {
			Username string
			Password string
		}
		req.Username = user.Username
		req.Password = user.Password
		jsonBody, _ = json.Marshal(req)
		w, err = web.MakeRequest(c.Web, http.MethodPost, "/login", bytes.NewReader(jsonBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Message string
			Token   string
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		newHeader := map[string]string{
			"Authorization": "Bearer " + resp.Token,
		}
		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/refresh", bytes.NewReader(jsonBody), newHeader)
		assert.NoError(t, err)
		fmt.Println(w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)

		w, err = web.MakeRequest(c.Web, http.MethodGet, "/user/refresh", bytes.NewReader(jsonBody), userHeader)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})
}
