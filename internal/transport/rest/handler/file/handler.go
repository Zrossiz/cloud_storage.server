package file

import (
	service "cloudStorage/internal/services"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func FileHandler(db *gorm.DB, redis *redis.Client, minioStorage *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "upload-file/") {
			service.UploadFile(w, r, redis, minioStorage)
		}
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "create-folder/") {
			service.CreateFolder(w, r, redis, minioStorage)
		}
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "find/") {
			service.FindFiles(w, r, redis, minioStorage)
		}
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "update/") {
			service.UpdateFile(w, r, redis, minioStorage)
		}
		if r.Method == http.MethodGet && strings.Contains(r.URL.String(), "get/") {
			service.GetAllByPath(w, r, redis, minioStorage)
		}
		if r.Method == http.MethodDelete && strings.Contains(r.URL.String(), "delete/") {
			service.DeleteObj(w, r, redis, minioStorage)
		}
	}
}
