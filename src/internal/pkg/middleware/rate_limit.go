package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	requests int
	duration time.Duration
}

func NewRateLimiter(redisClient *redis.Client, requests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		requests: requests,
		duration: duration,
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier := rl.getClientIdentifier(r)

		allowed, remaining, resetTime, err := rl.isAllowed(r.Context(), identifier)
		if err != nil {
			http.Error(w, "Rate limit check failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
		w.Header().Set("X-RateLimit-Remainig", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			w.Header().Set("Retry-After", strconv.FormatInt(resetTime.Unix(), 10))
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
			return 
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getClientIdentifier(r *http.Request) string {
	if claims, ok := r.Context().Value("user").(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(string); ok {
			return fmt.Sprintf("user:%s", userID)
		}
	}

	return fmt.Sprintf("ip:%s", r.RemoteAddr)
}

func (rl *RateLimiter) isAllowed(ctx context.Context, identifier string) (bool, int, time.Time, error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	now := time.Now()

	pipe := rl.redisClient.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.duration)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, 0, now, err
	}

	count, err := incrCmd.Result()
	if err != nil {
		return false, 0, now, err
	}

	remaining := rl.requests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	ttl, err := rl.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return false, remaining, now, err
	}

	resetTime := now.Add(ttl)

	return count <= int64(rl.requests), remaining, resetTime, nil
}