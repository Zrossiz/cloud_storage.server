package service

import (
	"cloudStorage/internal/middleware"
	"cloudStorage/internal/transport/rest/response"
	"context"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func UploadFile(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "UserId not found in context", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to get file from request")
	}
	defer file.Close()

	objectName := handler.Filename
	filePath := "user-" + userId + "files/" + objectName
	contentType := handler.Header.Get("Content-Type")

	_, err = minioStorage.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), filePath, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to upload file to storage")
		return
	}

	response.SendData(w, http.StatusOK, "File uploaded successfully")
}
