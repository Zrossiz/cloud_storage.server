package main

import (
	"cloudStorage/internal/app"
	db "cloudStorage/internal/database"
	"log"
)

func main() {
	redisClient, err := db.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	app.InitApp(redisClient)
}
