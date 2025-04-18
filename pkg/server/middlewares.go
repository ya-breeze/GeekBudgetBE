package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Skip authorization for the root endpoint
			if req.URL.Path == "/" || strings.HasPrefix(req.URL.Path, "/web/") {
				next.ServeHTTP(writer, req)
				return
			}

			// Skip authorization for the authorize endpoint - there is no way to do it with
			// go-server openapi templates now :(
			if req.URL.Path == "/v1/authorize" {
				next.ServeHTTP(writer, req)
				return
			}

			checkToken(logger, cfg.Issuer, cfg.JWTSecret, next, writer, req)
		})
	}
}

func checkToken(
	logger *slog.Logger, issuer, jwtSecret string, next http.Handler,
	writer http.ResponseWriter, req *http.Request,
) {
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
	userID, err := auth.CheckJWT(bearerToken, issuer, jwtSecret)
	if err != nil {
		logger.With("err", err).Warn("Invalid token")
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), common.UserIDKey, userID))
	next.ServeHTTP(writer, req)
}

func ForcedImportMiddleware(logger *slog.Logger, forcedImports chan<- background.ForcedImport) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Store forced import in the context
			ctx := context.WithValue(req.Context(), background.ForcedImportKey, forcedImports)
			req = req.WithContext(ctx)

			next.ServeHTTP(writer, req)
		})
	}
}
