package middleware

import (
	"cloudStorage/internal/transport/rest/response"
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) bool {
	idCookie, idErr := r.Cookie("userId")
	if idErr != nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	tokenCookie, sessionErr := r.Cookie("session")
	if sessionErr != nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	ctx := context.Background()
	redisSession, redisErr := redisClient.Get(ctx, idCookie.Value).Result()
	if redisErr == redis.Nil || redisErr != nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	if redisSession != tokenCookie.Value {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	return true
}
