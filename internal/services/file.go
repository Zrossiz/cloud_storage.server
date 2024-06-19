package service

import (
	"cloudStorage/internal/middleware"
	"cloudStorage/internal/transport/rest/response"
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type FileLink struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func UploadFile(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
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
		return
	}
	defer file.Close()

	objectName := handler.Filename
	filePath := "user-" + userId + "-files/" + objectName
	contentType := handler.Header.Get("Content-Type")

	_, err = minioStorage.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), filePath, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to upload file to storage")
		return
	}

	response.SendData(w, http.StatusCreated, "File uploaded successfully")
}

func FindFiles(w http.ResponseWriter, r *http.Request, db *gorm.DB, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "UserId not found in context", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	prefix := "user-" + userId + "-files/"
	searchQuery := r.URL.Query().Get("search")

	objectCh := minioStorage.ListObjects(ctx, os.Getenv("BUCKET_NAME"), minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var fileLinks []FileLink

	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}

		if searchQuery != "" && !strings.Contains(object.Key, searchQuery) {
			continue
		}

		reqParams := make(url.Values)
		presignedURL, err := minioStorage.PresignedGetObject(ctx, os.Getenv("BUCKET_NAME"), object.Key, time.Hour*24, reqParams)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to get presigned URL")
			return
		}

		fileLink := strings.Replace(presignedURL.String(), "minio:9000", os.Getenv("MINIO_PUBLIC_HOST"), 1)

		fileLinks = append(fileLinks, FileLink{
			Name: strings.TrimPrefix(object.Key, prefix),
			Link: fileLink,
		})
	}

	response.SendData(w, http.StatusOK, fileLinks)
}
