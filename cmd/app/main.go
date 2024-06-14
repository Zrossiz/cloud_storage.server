package main

import (
	db "cloudStorage/internal/database"
	"cloudStorage/internal/models"
	"cloudStorage/internal/transport/rest/router"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	database, err := db.InitConnect()

	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	router := router.NewRouter(database)

	fmt.Println("Successfully connected to the database and applied migrations!")

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
