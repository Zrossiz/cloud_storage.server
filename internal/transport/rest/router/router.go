package router

import (
	"cloudStorage/internal/transport/rest/handler/file"
	"cloudStorage/internal/transport/rest/handler/user"
	"net/http"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, redis *redis.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/user/", user.UserHandler(db, redis))
	mux.HandleFunc("/api/file/", file.FileHandler(db, redis))
	return mux
}
