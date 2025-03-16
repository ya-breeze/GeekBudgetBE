package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/server"
)

func createMockServer() (*httptest.Server, *int) {
	callCount := 0

	// Mock CNB response with a few currencies
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		mockResponse := `14.03.2025 #53
Země|Měna|Množství|Kód|Kurz
Austrálie|dolar|1|AUD|15,123
EMU|euro|1|EUR|25,490
Japonsko|jen|100|JPY|19,651
USA|dolar|1|USD|22,758
`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))

	return s, &callCount
}

//nolint:funlen // Test function with many cases
func TestCurrenciesRatesFetcher_Convert(t *testing.T) {
	// Setup a mock HTTP cnbMockServer
	cnbMockServer, _ := createMockServer()
	defer cnbMockServer.Close()

	sut := server.NewCurrenciesRatesFetcher(nil, nil) // Using nil for logger and config in tests
	sut.BaseURL = cnbMockServer.URL
	testDate := time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	tests := []struct {
		name          string
		day           time.Time
		from          string
		to            string
		amount        float64
		expected      float64
		expectedError bool
	}{
		{
			name:     "Convert CZK to USD",
			day:      testDate,
			from:     "CZK",
			to:       "USD",
			amount:   100,
			expected: 100 / 22.758,
		},
		{
			name:     "Convert USD to CZK",
			day:      testDate,
			from:     "USD",
			to:       "CZK",
			amount:   50,
			expected: 50 * 22.758,
		},
		{
			name:     "Convert EUR to USD",
			day:      testDate,
			from:     "EUR",
			to:       "USD",
			amount:   200,
			expected: (200 * 25.490) / 22.758,
		},
		{
			name:     "Convert JPY to EUR",
			day:      testDate,
			from:     "JPY",
			to:       "EUR",
			amount:   10000,
			expected: (10000 * (19.651 / 100)) / 25.490,
		},
		{
			name:          "Invalid source currency",
			day:           testDate,
			from:          "XYZ",
			to:            "USD",
			amount:        100,
			expectedError: true,
		},
		{
			name:          "Invalid target currency",
			day:           testDate,
			from:          "USD",
			to:            "XYZ",
			amount:        100,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sut.Convert(ctx, tt.day, tt.from, tt.to, tt.amount)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Compare with a small delta to account for floating point imprecision
			if !almostEqual(result, tt.expected, 0.0001) {
				t.Errorf("expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

func TestCurrenciesRatesFetcher_FetchRates(t *testing.T) {
	// Test caching behavior
	cnbMockServer, callCount := createMockServer()
	defer cnbMockServer.Close()

	sut := server.NewCurrenciesRatesFetcher(nil, nil)
	sut.BaseURL = cnbMockServer.URL
	testDate := time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	// First call should fetch from server
	_, err := sut.Convert(ctx, testDate, "USD", "CZK", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *callCount != 1 {
		t.Errorf("expected 1 HTTP call, got %d", callCount)
	}

	// Second call with same date should use cache
	_, err = sut.Convert(ctx, testDate, "USD", "CZK", 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *callCount != 1 {
		t.Errorf("expected HTTP call count to remain 1, got %d", callCount)
	}

	// Different date should cause another HTTP call
	differentDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	_, err = sut.Convert(ctx, differentDate, "USD", "CZK", 300)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *callCount != 2 {
		t.Errorf("expected HTTP call count to increase to 2, got %d", callCount)
	}
}

// Helper function to compare floating point numbers
func almostEqual(a, b, delta float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= delta
}

func TestCurrenciesRatesFetcher_ErrorHandling(t *testing.T) {
	// Test error responses
	cnbMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	}))
	defer cnbMockServer.Close()

	sut := server.NewCurrenciesRatesFetcher(nil, nil)
	sut.BaseURL = cnbMockServer.URL
	testDate := time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	_, err := sut.Convert(ctx, testDate, "USD", "CZK", 100)
	if err == nil {
		t.Error("expected error for bad HTTP status but got nil")
	}
}

func TestCurrenciesRatesFetcher_ContextCancellation(t *testing.T) {
	// Test context cancellation
	cnbMockServer, _ := createMockServer()
	defer cnbMockServer.Close()

	sut := server.NewCurrenciesRatesFetcher(nil, nil)
	sut.BaseURL = cnbMockServer.URL
	testDate := time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC)

	// Create a context with timeout shorter than the server's response time
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := sut.Convert(ctx, testDate, "USD", "CZK", 100)
	if err == nil {
		t.Error("expected error due to context timeout but got nil")
	}
}
