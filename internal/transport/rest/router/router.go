package router

import (
	"cloudStorage/internal/transport/rest/handler/user"
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/user", user.UserHandler(db))

	return mux
}
