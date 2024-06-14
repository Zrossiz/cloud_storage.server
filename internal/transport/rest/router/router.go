package router

import (
	"cloudStorage/internal/transport/rest/handler/user"
	"net/http"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/user/", user.UserHandler(db))
	return mux
}
