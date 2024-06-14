package user

import (
	"cloudStorage/internal/models"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB) http.HandlerFunc {
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
		}

		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error generating JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(userJSON)
	}
}
