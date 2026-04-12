package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	kincookies "github.com/ya-breeze/kin-core/cookies"
	kinmiddleware "github.com/ya-breeze/kin-core/middleware"
	"gorm.io/gorm"
)

// skipAuthPaths lists paths that don't require authentication.
var skipAuthPaths = map[string]bool{
	"/v1/authorize": true,
	"/auth/refresh": true,
}

func AuthMiddleware(logger *slog.Logger, cfg *config.Config, db *gorm.DB) mux.MiddlewareFunc {
	kinCfg := kinmiddleware.Config{
		JWTSecret: []byte(cfg.JWTSecret),
		DB:        db,
		CookieCfg: kincookies.Config{Secure: cfg.CookieSecure},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Skip auth for OPTIONS (CORS preflight), root, web assets, and public endpoints
			if req.Method == "OPTIONS" || req.URL.Path == "/" ||
				strings.HasPrefix(req.URL.Path, "/web/") || skipAuthPaths[req.URL.Path] {
				next.ServeHTTP(writer, req)
				return
			}

			// If no kin_access cookie, fall back to Authorization: Bearer header
			// (needed for API clients like the generated goclient used in tests)
			if bearer := req.Header.Get("Authorization"); bearer != "" &&
				strings.HasPrefix(bearer, "Bearer ") {
				if c, _ := req.Cookie("kin_access"); c == nil {
					token := strings.TrimPrefix(bearer, "Bearer ")
					req.AddCookie(&http.Cookie{Name: "kin_access", Value: token})
				}
			}

			claims, err := kinmiddleware.ValidateRequest(req, kinCfg)
			if err != nil {
				logger.Warn("Unauthorized request", "path", req.URL.Path, "error", err)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), constants.UserIDKey, claims.UserID)
			if claims.FamilyID != nil {
				ctx = context.WithValue(ctx, constants.FamilyIDKey, *claims.FamilyID)
			}
			ctx = context.WithValue(ctx, constants.ChangeSourceKey, constants.ChangeSourceUser)
			next.ServeHTTP(writer, req.WithContext(ctx))
		})
	}
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Add CORS headers to all responses
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			// Handle OPTIONS requests for CORS preflight
			if req.Method == "OPTIONS" {
				writer.Header().Set("Access-Control-Max-Age", "86400")
				writer.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(writer, req)
		})
	}
}

func ForcedImportMiddleware(logger *slog.Logger, forcedImports chan<- common.ForcedImport) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Store forced import in the context
			ctx := context.WithValue(req.Context(), constants.ForcedImportKey, forcedImports)
			req = req.WithContext(ctx)

			next.ServeHTTP(writer, req)
		})
	}
}
