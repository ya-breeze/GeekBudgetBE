package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

type contextKey string

const (
	UsernameKey contextKey = "username"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config, db database.Storage) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Skip authorization for the authorize endpoint - there is no way to do it with
			// go-server openapi templates now :(
			if req.URL.Path == "/v1/authorize" {
				next.ServeHTTP(writer, req)
				return
			}

			// Authorization logic
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				http.Error(writer, "Invalid authorization header", http.StatusUnauthorized)
				return
			}
			bearerToken := authHeaderParts[1]

			// Parse the token
			token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				// Make sure that the token method conforms to "SigningMethodHMAC"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				signingKey := []byte(cfg.JWTSecret)
				return signingKey, nil
			})
			if err != nil {
				logger.With("err", err).Warn("Invalid token")
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				username, err := claims.GetSubject()
				if err != nil {
					logger.With("err", err).Warn("Invalid subject")
					http.Error(writer, "Invalid token", http.StatusUnauthorized)
					return
				}
				logger.With("username", username).Info("Authorized user")

				req = req.WithContext(context.WithValue(req.Context(), UsernameKey, username))
				next.ServeHTTP(writer, req)
				return
			}

			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		})
	}
}
