package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}


func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "[FAIL] authorization header is required", http.StatusUnauthorized)
            return
        }


        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            http.Error(w, "[FAIL] invalid authorization header", http.StatusUnauthorized)
            return
        }

        tokenStr := bearerToken[1]

        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return []byte(m.jwtSecret), nil
        })

        if err != nil {
            http.Error(w, fmt.Sprintf("[FAIL] %v", err), http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            http.Error(w, "[FAIL] Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}