package app

import (
	db "cloudStorage/internal/database"
	"cloudStorage/internal/models"
	"cloudStorage/internal/transport/rest/router"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func InitApp(redisClient *redis.Client) {
	database, err := db.InitConnect()
	if err != nil {
		log.Fatal(err)
	}

	if errDB := database.AutoMigrate(&models.User{}); errDB != nil {
		log.Fatalf("Failed to migrate database: %v", errDB)
	}

	minioStorage := db.InitMinio()

	router := router.NewRouter(database, redisClient, minioStorage)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Println("Successfully connected to the database and applied migrations!")
	fmt.Println("Starting server at port 8080")
	if errApp := http.ListenAndServe(":8080", handler); errApp != nil {
		log.Fatal(errApp)
	}
}
