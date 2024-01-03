package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jkittell/data/database"
	"log"
	"time"
)

var counts int64

func main() {
	log.Println("starting authentication service")

	// connect to DB
	users := connectToDB()
	if &users == nil {
		log.Fatal("cannot connect to postgres")
	}

	router := gin.New()
	usersHandler := NewUsersHandler(users)

	// Register Routes
	router.GET("/", homePage)
	router.POST("/authenticate", usersHandler.Authenticate)
	router.GET("/users", usersHandler.ListUsers)
	router.POST("/users", usersHandler.CreateUser)
	router.GET("/users/:id", usersHandler.GetUser)
	router.PUT("/users/:id", usersHandler.UpdateUser)
	router.DELETE("/users", usersHandler.DeleteUser)

	err := router.Run(":80")
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDB() *database.PosgresDB[*User] {
	for {
		connection, err := database.NewPostgresDB[*User](".env", "users", func() *User {
			return &User{}
		})
		if err != nil {
			log.Println("postgres is not ready...")
			counts++
		} else {
			log.Println("connected to postgres")
			return &connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
