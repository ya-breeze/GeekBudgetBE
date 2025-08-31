package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
)

const (
	testUser = "test"
	testPass = "test"
)

func TestBudgetWebIntegration(t *testing.T) {
	// Setup test server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := CreateTestLogger()

	// Hash password for test user
	hashed, err := auth.HashPassword([]byte(testPass))
	require.NoError(t, err)

	cfg := &config.Config{
		Port:       0, // Use random port
		Users:      testUser + ":" + base64.StdEncoding.EncodeToString(hashed),
		CookieName: "geekbudgetcookie",
	}

	// Create storage
	storage := database.NewStorage(logger, cfg)
	err = storage.Open()
	require.NoError(t, err)
	defer storage.Close()

	// Start server
	forcedImportChan := make(chan background.ForcedImport)
	addr, finishChan, err := server.Serve(ctx, logger, storage, cfg, forcedImportChan)
	require.NoError(t, err)

	// Cleanup server
	defer func() {
		cancel()
		<-finishChan
	}()

	baseURL := fmt.Sprintf("http://%s", addr.String())
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Test: Login and get session cookie
	sessionCookie := loginAndGetCookie(t, client, baseURL)

	// Test: GET /web/budget/plan should return 200
	planReq, err := http.NewRequest("GET", baseURL+"/web/budget/plan", nil)
	require.NoError(t, err)
	planReq.AddCookie(sessionCookie)

	planResp, err := client.Do(planReq)
	require.NoError(t, err)
	defer planResp.Body.Close()

	assert.Equal(t, http.StatusOK, planResp.StatusCode, "Budget planning page should return 200")

	// Test: GET /web/budget/compare should return 200
	compareReq, err := http.NewRequest("GET", baseURL+"/web/budget/compare", nil)
	require.NoError(t, err)
	compareReq.AddCookie(sessionCookie)

	compareResp, err := client.Do(compareReq)
	require.NoError(t, err)
	defer compareResp.Body.Close()

	assert.Equal(t, http.StatusOK, compareResp.StatusCode, "Budget comparison page should return 200")
}

func loginAndGetCookie(t *testing.T, client *http.Client, baseURL string) *http.Cookie {
	// Prepare login form data
	formData := url.Values{}
	formData.Set("username", testUser)
	formData.Set("password", testPass)

	// Create login request
	loginReq, err := http.NewRequest("POST", baseURL+"/web/login", strings.NewReader(formData.Encode()))
	require.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform login
	loginResp, err := client.Do(loginReq)
	require.NoError(t, err)
	defer loginResp.Body.Close()

	// Should redirect after successful login
	assert.True(t, loginResp.StatusCode == http.StatusSeeOther || loginResp.StatusCode == http.StatusFound,
		"Login should redirect, got status: %d", loginResp.StatusCode)

	// Extract session cookie
	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "geekbudgetcookie" { // Default cookie name from config
			sessionCookie = cookie
			break
		}
	}

	require.NotNil(t, sessionCookie, "Should receive session cookie after login")
	return sessionCookie
}

func TestBudgetWebIntegration_WithoutAuth(t *testing.T) {
	// Setup test server (same as above but without login)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := CreateTestLogger()

	hashed, err := auth.HashPassword([]byte(testPass))
	require.NoError(t, err)

	cfg := &config.Config{
		Port:  0,
		Users: testUser + ":" + base64.StdEncoding.EncodeToString(hashed),
	}

	storage := database.NewStorage(logger, cfg)
	err = storage.Open()
	require.NoError(t, err)
	defer storage.Close()

	forcedImportChan := make(chan background.ForcedImport)
	addr, finishChan, err := server.Serve(ctx, logger, storage, cfg, forcedImportChan)
	require.NoError(t, err)

	defer func() {
		cancel()
		<-finishChan
	}()

	baseURL := fmt.Sprintf("http://%s", addr.String())
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, we want to check the response
			return http.ErrUseLastResponse
		},
	}

	// Test: GET /web/budget/plan without auth should redirect to login
	planReq, err := http.NewRequest("GET", baseURL+"/web/budget/plan", nil)
	require.NoError(t, err)

	planResp, err := client.Do(planReq)
	require.NoError(t, err)
	defer planResp.Body.Close()

	// Should redirect to login (or return error page)
	assert.True(t, planResp.StatusCode == http.StatusSeeOther ||
		planResp.StatusCode == http.StatusFound ||
		planResp.StatusCode == http.StatusUnauthorized ||
		planResp.StatusCode == http.StatusInternalServerError,
		"Budget planning without auth should redirect or error, got: %d", planResp.StatusCode)

	// Test: GET /web/budget/compare without auth should redirect to login
	compareReq, err := http.NewRequest("GET", baseURL+"/web/budget/compare", nil)
	require.NoError(t, err)

	compareResp, err := client.Do(compareReq)
	require.NoError(t, err)
	defer compareResp.Body.Close()

	// Should redirect to login (or return error page)
	assert.True(t, compareResp.StatusCode == http.StatusSeeOther ||
		compareResp.StatusCode == http.StatusFound ||
		compareResp.StatusCode == http.StatusUnauthorized ||
		compareResp.StatusCode == http.StatusInternalServerError,
		"Budget comparison without auth should redirect or error, got: %d", compareResp.StatusCode)
}
