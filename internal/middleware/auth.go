package middleware

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) (*http.Request, bool) {
	idCookie, err := r.Cookie("userId")
	if err != nil {
		return nil, false
	}

	tokenCookie, err := r.Cookie("session")
	if err != nil {
		return nil, false
	}

	ctx := r.Context()
	redisSession, err := redisClient.Get(ctx, idCookie.Value).Result()
	if err == redis.Nil || err != nil {
		return nil, false
	}

	if redisSession != tokenCookie.Value {
		return nil, false
	}
	ctx = context.WithValue(ctx, "userId", idCookie.Value)

	return r.WithContext(ctx), true
}
