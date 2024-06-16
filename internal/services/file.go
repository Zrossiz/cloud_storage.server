package service

import (
	"cloudStorage/internal/middleware"
	"cloudStorage/internal/transport/rest/response"
	"net/http"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func UploadFile(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client) {
	var isAuth bool = middleware.AuthMiddleware(w, r, redis)
	if !isAuth {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	response.SendData(w, http.StatusOK, "Upload route")
}
