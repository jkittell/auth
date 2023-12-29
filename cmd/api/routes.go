package main

import (
	"context"
	"fmt"
	"github.com/andrewpillar/query"
	"github.com/gin-gonic/gin"
	"github.com/jkittell/auth/model"
	"github.com/jkittell/data/database"
	"net/http"
)

type UsersHandler struct {
	users *database.PosgresDB[*model.User]
}

func NewUsersHandler(users *database.PosgresDB[*model.User]) *UsersHandler {
	return &UsersHandler{
		users: users,
	}
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

func (h UsersHandler) Authenticate(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate the user against the database
	user, ok, err := h.users.Get(context.TODO(), query.Where("email", "=", query.Arg(requestPayload.Email)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	res := struct {
		Error   bool
		Message string
		Data    *model.User
	}{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	c.JSON(200, res)
}

func (h UsersHandler) CreateUser(c *gin.Context) {
	var user *model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.users.Create(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h UsersHandler) ListUsers(c *gin.Context) {
	res, err := h.users.All(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}
func (h UsersHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, ok, err := h.users.Get(context.TODO(), query.Where("id", "=", query.Arg(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}

	c.JSON(200, user)
}
func (h UsersHandler) UpdateUser(c *gin.Context) {
	var user *model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.users.Update(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h UsersHandler) DeleteUser(c *gin.Context) {
	var user *model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.users.Delete(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
