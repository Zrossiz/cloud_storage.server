package main

import (
	db "cloudStorage/internal/database"
	"cloudStorage/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var user models.User

		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO users (name, email, password) 
			VALUES ($1, $2, $3) 
			RETURNING id, name, email, password, created_at, updated_at, deleted_at
		`

		err = db.QueryRow(query, user.Name, user.Email, user.Password).Scan(
			&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(userJSON)
	}
}

func main() {
	database, err := db.InitConnect()

	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	http.HandleFunc("/api/user/", createUserHandler(database))

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Successfully connected to the database and applied migrations!")

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
