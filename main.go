package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jkittell/data/database"
	"log"
)

var counts int64

func main() {
	users, err := database.NewPostgresDB[*User]("users", func() *User {
		return &User{}
	})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	usersHandler := NewUsersHandler(&users)

	router.POST("/authenticate", usersHandler.Authenticate)
	router.GET("/users", usersHandler.ListUsers)
	router.POST("/users", usersHandler.CreateUser)
	router.GET("/users/:id", usersHandler.GetUser)
	router.PUT("/users/:id", usersHandler.UpdateUser)
	router.DELETE("/users", usersHandler.DeleteUser)

	err = router.Run(":80")
	if err != nil {
		log.Fatal(err)
	}
}
