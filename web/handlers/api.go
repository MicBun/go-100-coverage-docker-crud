package handlers

import (
	"github.com/MicBun/go-100-coverage-docker-crud/service"
	"github.com/MicBun/go-100-coverage-docker-crud/util/jwtAuth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type apiHandler struct {
	container *service.Container
}

type ApiHandlerInterface interface {
	Hello(c *gin.Context)
	RegisterUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserByToken(c *gin.Context)
	ListUsers(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

func NewApiHandler(container *service.Container) ApiHandlerInterface {
	return &apiHandler{
		container: container,
	}
}

// Hello godoc
// @Summary Hello
// @Description Hello
// @Tags Hello
// @Success 200 {object} map[string]interface{}
// @Router /hello [get]
func (h *apiHandler) Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello"})
	return
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// RegisterUser godoc
// @Summary Register User
// @Description Register User
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body RegisterUserRequest true "User"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/register [post]
func (h *apiHandler) RegisterUser(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	_, err := h.container.Admin.RegisterUser(req.Username, req.Password, req.Name)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User registered"})
	return
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// UpdateUser godoc
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/update/{id} [put]
func (h *apiHandler) UpdateUser(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	updatedUser, err := h.container.Admin.UpdateUser(uint(id), req.Username, req.Password, req.Name)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User updated", "user": updatedUser})
	return
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/delete/{id} [delete]
func (h *apiHandler) DeleteUser(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.container.Admin.DeleteUser(uint(id))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
	return
}

// GetUserByID godoc
// @Summary Get User By ID
// @Description Get User By ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/get/{id} [get]
func (h *apiHandler) GetUserByID(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.container.Admin.GetUser(uint(id))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User retrieved", "user": user})
	return
}

// GetUserByToken godoc
// @Summary Get User By Token
// @Description Get User By Token
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/get [get]
func (h *apiHandler) GetUserByToken(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	id, _ := jwtAuth.ExtractTokenID(c)
	user, _ := h.container.Admin.GetUser(id)
	c.JSON(200, gin.H{"message": "User retrieved", "user": user})
	return
}

// ListUsers godoc
// @Summary List Users
// @Description List Users
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/list [get]
func (h *apiHandler) ListUsers(c *gin.Context) {
	role, _ := jwtAuth.ExtractTokenRole(c)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	users, err := h.container.Admin.ListUsers()
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	var userList []map[string]interface{}
	for _, user := range users {
		userList = append(userList, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"link":     "/user/get/" + strconv.Itoa(int(user.ID)),
		})
	}
	c.JSON(200, gin.H{"message": "Users retrieved", "users": userList})
	return
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body LoginRequest true "User"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /login [post]
func (h *apiHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	user, err := h.container.Admin.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	token, _ := jwtAuth.GenerateToken(user.ID)
	h.container.Admin.SaveToken(user.ID, token)
	c.JSON(200, gin.H{"message": "User logged in", "user": user, "token": token})
	return
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh Token
// @Tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /user/refresh [get]
func (h *apiHandler) RefreshToken(c *gin.Context) {
	userID, _ := jwtAuth.ExtractTokenID(c)
	user, _ := h.container.Admin.GetUser(userID)
	if user.Token != jwtAuth.ExtractToken(c) {
		c.JSON(400, gin.H{"message": "Previous token is not valid"})
		return
	}
	token, _ := jwtAuth.GenerateToken(user.ID)
	h.container.Admin.SaveToken(user.ID, token)
	c.JSON(200, gin.H{"message": "Token refreshed", "token": token})
}
