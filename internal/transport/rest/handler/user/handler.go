package user

import (
	"cloudStorage/internal/models"
	"cloudStorage/internal/transport/rest/response"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB, redis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.SendError(w, http.StatusMethodNotAllowed, "Invalid request method")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Error reading request body")
			return
		}
		defer r.Body.Close()

		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			response.SendError(w, http.StatusBadRequest, "Error parsing JSON")
			return
		}

		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			response.SendError(w, http.StatusConflict, "Email already in use")
			return
		}
		if err := db.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
			response.SendError(w, http.StatusConflict, "Username already in use")
			return
		}

		if err := user.SetPassword(user.Password); err != nil {
			response.SendError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}

		if err := db.Create(&user).Error; err != nil {
			response.SendError(w, http.StatusInternalServerError, "Error saving user")
			return
		}
		token := uuid.New().String()
		redisUserId := fmt.Sprintf("%d", user.ID)
		redisSession := fmt.Sprintf("%s", token)
		ctx := context.Background()

		err = redis.Set(ctx, redisUserId, redisSession, 10*time.Second).Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.SendError(w, http.StatusInternalServerError, "Error saving session")
			return
		}

		cookie := http.Cookie{
			Name:    "session",
			Value:   token,
			Expires: time.Now().Add(10 * time.Second),
			Path:    "/",
		}
		http.SetCookie(w, &cookie)

		response.SendData(w, http.StatusCreated, user)
	}
}
