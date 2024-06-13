package main

import (
	db "cloudStorage/internal/database"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func main() {
	http.HandleFunc("/", helloHandler)
	database, err := db.InitConnect()

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Successfully connected to the database and applied migrations!")

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
