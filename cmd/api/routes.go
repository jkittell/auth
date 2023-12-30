package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andrewpillar/query"
	"github.com/gin-gonic/gin"
	"github.com/jkittell/auth/model"
	"github.com/jkittell/data/database"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
		return
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if CheckPasswordHash(requestPayload.Password, user.Password) {
		res := struct {
			Error   bool
			Message string
			Data    *model.User
		}{
			Error:   false,
			Message: fmt.Sprintf("authorized user: %s", user.Email),
			Data:    user,
		}

		c.JSON(200, res)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}
}

func (h UsersHandler) CreateUser(c *gin.Context) {
	var user *model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password
	if hashedPassword, err := HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		user.Password = hashedPassword
	}

	id, err := h.users.Create(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success payload
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h UsersHandler) ListUsers(c *gin.Context) {
	res, err := h.users.All(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < res.Length(); i++ {
		r := res.Lookup(i)
		data, _ := json.Marshal(r)
		log.Println(string(data))
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
		return
	}

	c.JSON(200, user)
}

func (h UsersHandler) UpdateUser(c *gin.Context) {
	var user *model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password
	if hashedPassword, err := HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		user.Password = hashedPassword
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
