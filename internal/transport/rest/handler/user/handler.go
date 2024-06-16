package user

import (
	"cloudStorage/internal/dto"
	"cloudStorage/internal/models"
	"cloudStorage/internal/transport/rest/response"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB, redis *redis.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "register/") {
			create(w, r, db, redis)
		}
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "login/") {
			login(w, r, db, redis)
		}
	}
}

func create(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client) {
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

	cookie, err := createSession(int(user.ID), redis)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Error saving session")
		return
	}

	userCookie := http.Cookie{
		Name:    "userId",
		Value:   fmt.Sprintf("%d", user.ID),
		Expires: time.Now().Add(10 * time.Second),
		Path:    "/",
	}

	http.SetCookie(w, &userCookie)
	http.SetCookie(w, &cookie)

	response.SendData(w, http.StatusCreated, user)
}

func login(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Error reading request body")
		return
	}
	defer r.Body.Close()

	var userDTO dto.LoginUserDTO
	if err := json.Unmarshal(body, &userDTO); err != nil {
		response.SendError(w, http.StatusBadRequest, "Error parsing JSON")
		return
	}

	var existingUser models.User
	if err := db.Where("name = ?", userDTO.Name).First(&existingUser).Error; err != nil {
		response.SendError(w, http.StatusNotFound, "Invalid credentials")
		return
	}

	isValid := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userDTO.Password))
	if isValid != nil {
		response.SendError(w, http.StatusNotFound, "Invalid credentials")
		return
	}

	cookie, err := createSession(int(existingUser.ID), redis)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Error saving session")
		return
	}

	userCookie := http.Cookie{
		Name:    "userId",
		Value:   fmt.Sprintf("%d", existingUser.ID),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &userCookie)
	http.SetCookie(w, &cookie)

	response.SendData(w, http.StatusOK, existingUser)
}

func createSession(userId int, redis *redis.Client) (http.Cookie, error) {
	token := uuid.New().String()
	redisUserId := fmt.Sprintf("%d", userId)
	ctx := context.Background()

	err := redis.Set(ctx, redisUserId, token, 24*time.Hour).Err()
	if err != nil {
		return http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}

	return cookie, nil
}
