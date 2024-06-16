package user

import (
	service "cloudStorage/internal/services"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func UserHandler(db *gorm.DB, redis *redis.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "register/") {
			service.CreateUser(w, r, db, redis)
		}
		if r.Method == http.MethodPost && strings.Contains(r.URL.String(), "login/") {
			service.LoginUser(w, r, db, redis)
		}
	}
}
