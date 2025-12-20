package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config) mux.MiddlewareFunc {
	// Re-create cookie store here or pass it in. Ideally pass it in, but for now we recreate it
	// based on the key "SESSION_KEY" usage in webapp.go.
	// TODO: Refactor to pass cookie store consistently or use a shared constant/config.
	// In webapp.go: cookies: sessions.NewCookieStore([]byte("SESSION_KEY")),
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))

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

			// Skip authorization for the authorize endpoint - there is no way to do it with
			// go-server openapi templates now :(
			if req.URL.Path == "/v1/authorize" {
				next.ServeHTTP(writer, req)
				return
			}

			if strings.HasPrefix(req.URL.Path, "/images/") {
				checkCookie(logger, cfg.Issuer, cfg.JWTSecret, cfg.CookieName, store, next, writer, req)
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

func checkCookie(
	logger *slog.Logger, issuer, jwtSecret, cookieName string, store *sessions.CookieStore,
	next http.Handler, writer http.ResponseWriter, req *http.Request,
) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		logger.With("err", err).Warn("Failed to get session")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		logger.Warn("Token not found in session")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := auth.CheckJWT(token, issuer, jwtSecret)
	if err != nil {
		logger.With("err", err).Warn("Invalid token in cookie")
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), common.UserIDKey, userID))
	next.ServeHTTP(writer, req)
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


