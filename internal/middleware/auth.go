package middleware

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) (*http.Request, bool) {
	idCookie, idErr := r.Cookie("userId")
	if idErr != nil {
		return nil, false
	}

	tokenCookie, sessionErr := r.Cookie("session")
	if sessionErr != nil {
		return nil, false
	}

	ctx := r.Context()
	redisSession, redisErr := redisClient.Get(ctx, idCookie.Value).Result()
	if redisErr == redis.Nil || redisErr != nil {
		return nil, false
	}

	if redisSession != tokenCookie.Value {
		return nil, false
	}
	ctx = context.WithValue(ctx, "userId", idCookie.Value)

	return r.WithContext(ctx), true
}
