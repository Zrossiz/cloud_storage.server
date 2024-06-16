package file

import (
	"cloudStorage/internal/services/middleware"
	"cloudStorage/internal/transport/rest/response"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func FileHandler(db *gorm.DB, redis *redis.Client, minioStorage *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "upload-file/") {
			uploadFile(w, r, db, redis)
		}
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client) {
	var isAuth bool = middleware.AuthMiddleware(w, r, redis)
	if !isAuth {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	response.SendData(w, http.StatusOK, "Upload route")
}
