package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config) mux.MiddlewareFunc {
	// Use the configured session secret for the cookie store
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Skip authorization for OPTIONS requests (CORS preflight)
			if req.Method == "OPTIONS" {
				next.ServeHTTP(writer, req)
				return
			}

			// Skip authorization for the root endpoint
			if req.URL.Path == "/" || strings.HasPrefix(req.URL.Path, "/web/") {
				next.ServeHTTP(writer, req)
				return
			}

			// Skip authorization for the authorize endpoint
			if req.URL.Path == "/v1/authorize" {
				next.ServeHTTP(writer, req)
				return
			}

			// 1. Try Cookie Authentication
			userID, err := getUserIDFromCookie(req, store, cfg.CookieName, cfg.Issuer, cfg.JWTSecret)
			if err == nil {
				req = req.WithContext(context.WithValue(req.Context(), common.UserIDKey, userID))
				next.ServeHTTP(writer, req)
				return
			}
			// Only log real errors, not just missing cookies
			if !errors.Is(err, http.ErrNoCookie) && !strings.Contains(err.Error(), "session") {
				logger.Debug("Cookie auth failed", "error", err)
			}

			// 2. Try Bearer Token Authentication
			userID, err = getUserIDFromToken(req, cfg.Issuer, cfg.JWTSecret)
			if err == nil {
				req = req.WithContext(context.WithValue(req.Context(), common.UserIDKey, userID))
				next.ServeHTTP(writer, req)
				return
			}
			logger.Debug("Bearer token auth failed", "error", err)

			// 3. Unauthorized
			logger.Warn("Unauthorized access attempt", "path", req.URL.Path)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

func getUserIDFromToken(req *http.Request, issuer, jwtSecret string) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	bearerToken := authHeaderParts[1]

	return auth.CheckJWT(bearerToken, issuer, jwtSecret)
}

func getUserIDFromCookie(
	req *http.Request, store *sessions.CookieStore, cookieName, issuer, jwtSecret string,
) (string, error) {
	// If cookieName is empty, use default from controller (though cfg should handle this)
	if cookieName == "" {
		cookieName = "geekbudgetcookie"
	}

	session, err := store.Get(req, cookieName)
	if err != nil {
		return "", err
	}

	token, ok := session.Values["token"].(string)
	if !ok || token == "" {
		return "", errors.New("token not found in session")
	}

	return auth.CheckJWT(token, issuer, jwtSecret)
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Add CORS headers to all responses
			// Note: For cookie-based auth cross-origin, Access-Control-Allow-Origin cannot be "*"
			// and Access-Control-Allow-Credentials must be true.
			// Ideally, we restrict Origin to the frontend domain.
			// For now, we keep "*" but strictly this breaks withCredentials=true.
			// Since we act as same-site (mostly), we'll start with this.
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
			ctx := context.WithValue(req.Context(), common.ForcedImportKey, forcedImports)
			req = req.WithContext(ctx)

			next.ServeHTTP(writer, req)
		})
	}
}
