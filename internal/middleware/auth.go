package middleware

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) bool {
	idCookie, idErr := r.Cookie("userId")
	if idErr != nil {
		return false
	}

	tokenCookie, sessionErr := r.Cookie("session")
	if sessionErr != nil {
		return false
	}

	ctx := context.Background()
	redisSession, redisErr := redisClient.Get(ctx, idCookie.Value).Result()
	if redisErr == redis.Nil || redisErr != nil {
		return false
	}

	if redisSession != tokenCookie.Value {
		return false
	}

	return true
}
