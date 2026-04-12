package webapp

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	kinauth "github.com/ya-breeze/kin-core/auth"
	"github.com/ya-breeze/kin-core/authdb"
	kincookies "github.com/ya-breeze/kin-core/cookies"
)

const (
	webAccessTokenTTL  = 15 * time.Minute
	webRefreshTokenTTL = 365 * 24 * time.Hour
)

//nolint:cyclop,funlen // refactor
func (r *WebAppRouter) loginHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := req.Form.Get("username")
	password := req.Form.Get("password")

	if username == "" || password == "" {
		r.RespondError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Timing-safe credential verification
	hash := kinauth.DummyHash
	user, err := r.db.GetUserByUsername(username)
	if err == nil {
		hash = user.PasswordHash
	}
	if !kinauth.VerifyPassword(password, hash) || err != nil {
		if errors.Is(err, database.ErrNotFound) {
			r.logger.Warn("User not found", "username", username)
		} else if err != nil {
			r.logger.Error("Failed to get user", "username", username, "error", err)
		} else {
			r.logger.Warn("Invalid password", "username", username)
		}
		r.RespondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	familyID := user.FamilyID
	accessToken, err := kinauth.GenerateAccessToken(user.ID, &familyID, []byte(r.cfg.JWTSecret), webAccessTokenTTL)
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookieCfg := kincookies.Config{Secure: r.cfg.CookieSecure}

	// Create refresh token
	if r.gormDB != nil {
		rt, rtErr := authdb.CreateRefreshToken(r.gormDB, user.ID, webRefreshTokenTTL)
		if rtErr != nil {
			r.logger.Warn("Failed to create refresh token", "error", rtErr)
		} else {
			kincookies.SetRefreshCookie(w, rt.Token, int(webRefreshTokenTTL.Seconds()), cookieCfg)
		}
	}
	kincookies.SetAccessCookie(w, accessToken, int(webAccessTokenTTL.Seconds()), cookieCfg)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *WebAppRouter) logoutHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookieCfg := kincookies.Config{Secure: r.cfg.CookieSecure}

	if tokenStr := kincookies.GetAccessToken(req); tokenStr != "" {
		if claims, parseErr := kinauth.ParseToken(tokenStr, []byte(r.cfg.JWTSecret)); parseErr == nil {
			if r.gormDB != nil {
				if err := authdb.BlacklistToken(r.gormDB, tokenStr, claims.ExpiresAt.Time); err != nil {
					r.logger.Warn("Failed to blacklist token on logout", "error", err)
				}
			}
		}
	}
	if rtStr := kincookies.GetRefreshToken(req); rtStr != "" && r.gormDB != nil {
		if err := authdb.RevokeRefreshToken(r.gormDB, rtStr); err != nil {
			r.logger.Warn("Failed to revoke refresh token on logout", "error", err)
		}
	}

	kincookies.ClearAuthCookies(w, cookieCfg)

	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "login.tpl", nil); err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetFamilyIDFromRequest extracts the family UUID from the request.
// It first checks the context (set by auth middleware), then falls back to the kin_access cookie.
func (r *WebAppRouter) GetFamilyIDFromRequest(req *http.Request) (uuid.UUID, int, error) {
	// Try context first (set by middleware for non-web routes)
	if familyID, ok := constants.GetFamilyID(req.Context()); ok {
		return familyID, http.StatusOK, nil
	}

	// Try kin-core access cookie
	tokenStr := kincookies.GetAccessToken(req)
	if tokenStr != "" {
		claims, err := kinauth.ParseToken(tokenStr, []byte(r.cfg.JWTSecret))
		if err == nil && claims.FamilyID != nil {
			return *claims.FamilyID, http.StatusOK, nil
		}
	}

	return uuid.UUID{}, http.StatusUnauthorized, errors.New("not authenticated")
}

func (r *WebAppRouter) ValidateUserID(
	tmpl *template.Template, w http.ResponseWriter, req *http.Request,
) (uuid.UUID, error) {
	familyID, _, err := r.GetFamilyIDFromRequest(req)
	if err != nil {
		if errTmpl := tmpl.ExecuteTemplate(w, "login.tpl", nil); errTmpl != nil {
			r.logger.Warn("failed to execute login template", "error", errTmpl)
			r.RespondError(w, errTmpl.Error(), http.StatusInternalServerError)
		}
		return uuid.UUID{}, fmt.Errorf("not authenticated: %w", err)
	}

	return familyID, nil
}
