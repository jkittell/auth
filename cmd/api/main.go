package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jkittell/auth/model"
	"github.com/jkittell/data/database"
	"log"
	"time"
)

var counts int64

func main() {
	log.Println("starting authentication service")

	// connect to DB
	users := connectToDB()
	if users == nil {
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

func openDB() (*database.PosgresDB[*model.User], error) {
	// databaseURL = os.Getenv("DATABASE_URL")
	databaseURL := "postgres://postgres:changeme@localhost:5432/postgres"

	// this returns connection pool
	pool, err := pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return database.NewPostgresDB[*model.User](pool, "users", func() *model.User {
		return &model.User{}
	}), nil
}

func connectToDB() *database.PosgresDB[*model.User] {
	for {
		connection, err := openDB()
		if err != nil {
			log.Println("postgres is not ready...")
			counts++
		} else {
			log.Println("connected to postgres")
			return connection
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
