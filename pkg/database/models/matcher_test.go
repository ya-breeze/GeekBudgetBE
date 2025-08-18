package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func TestMatcher_GetConfirmationPercentage(t *testing.T) {
	tests := []struct {
		name                string
		confirmationHistory []bool
		expected            float64
	}{
		{
			name:                "empty history",
			confirmationHistory: []bool{},
			expected:            0.0,
		},
		{
			name:                "all confirmed",
			confirmationHistory: []bool{true, true, true},
			expected:            100.0,
		},
		{
			name:                "all rejected",
			confirmationHistory: []bool{false, false, false},
			expected:            0.0,
		},
		{
			name:                "mixed - 50%",
			confirmationHistory: []bool{true, false, true, false},
			expected:            50.0,
		},
		{
			name:                "mixed - 75%",
			confirmationHistory: []bool{true, true, true, false},
			expected:            75.0,
		},
		{
			name:                "single confirmed",
			confirmationHistory: []bool{true},
			expected:            100.0,
		},
		{
			name:                "single rejected",
			confirmationHistory: []bool{false},
			expected:            0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &models.Matcher{
				ID:                  uuid.New(),
				ConfirmationHistory: tt.confirmationHistory,
			}

			result := matcher.GetConfirmationPercentage()
			if result != tt.expected {
				t.Errorf("GetConfirmationPercentage() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

//nolint:funlen
func TestMatcher_AddConfirmation(t *testing.T) {
	tests := []struct {
		name            string
		initialHistory  []bool
		confirmed       bool
		maxLength       int
		expectedHistory []bool
		expectedLength  int
	}{
		{
			name:            "add to empty history",
			initialHistory:  []bool{},
			confirmed:       true,
			maxLength:       5,
			expectedHistory: []bool{true},
			expectedLength:  1,
		},
		{
			name:            "add without exceeding max",
			initialHistory:  []bool{true, false},
			confirmed:       true,
			maxLength:       5,
			expectedHistory: []bool{true, false, true},
			expectedLength:  3,
		},
		{
			name:            "add exceeding max length - removes oldest",
			initialHistory:  []bool{true, false, true, false, true},
			confirmed:       false,
			maxLength:       5,
			expectedHistory: []bool{false, true, false, true, false},
			expectedLength:  5,
		},
		{
			name:            "add multiple exceeding max length",
			initialHistory:  []bool{true, false, true, false, true, true},
			confirmed:       false,
			maxLength:       3,
			expectedHistory: []bool{true, true, false},
			expectedLength:  3,
		},
		{
			name:            "max length of 1",
			initialHistory:  []bool{true, false},
			confirmed:       true,
			maxLength:       1,
			expectedHistory: []bool{true},
			expectedLength:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &models.Matcher{
				ID:                  uuid.New(),
				ConfirmationHistory: make([]bool, len(tt.initialHistory)),
			}
			copy(matcher.ConfirmationHistory, tt.initialHistory)

			matcher.AddConfirmation(tt.confirmed, tt.maxLength)

			if len(matcher.ConfirmationHistory) != tt.expectedLength {
				t.Errorf("AddConfirmation() length = %v, expected %v", len(matcher.ConfirmationHistory), tt.expectedLength)
			}

			for i, expected := range tt.expectedHistory {
				if i >= len(matcher.ConfirmationHistory) || matcher.ConfirmationHistory[i] != expected {
					t.Errorf("AddConfirmation() history[%d] = %v, expected %v", i,
						matcher.ConfirmationHistory[i], expected)
				}
			}
		})
	}
}

func TestMatcher_GetConfirmationHistoryLength(t *testing.T) {
	tests := []struct {
		name                string
		confirmationHistory []bool
		expectedLength      int
	}{
		{
			name:                "empty history",
			confirmationHistory: []bool{},
			expectedLength:      0,
		},
		{
			name:                "single item",
			confirmationHistory: []bool{true},
			expectedLength:      1,
		},
		{
			name:                "multiple items",
			confirmationHistory: []bool{true, false, true, false, true},
			expectedLength:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &models.Matcher{
				ID:                  uuid.New(),
				ConfirmationHistory: tt.confirmationHistory,
			}

			result := matcher.GetConfirmationHistoryLength()
			if result != tt.expectedLength {
				t.Errorf("GetConfirmationHistoryLength() = %v, expected %v", result, tt.expectedLength)
			}
		})
	}
}

func TestMatcher_IntegrationConfirmationWorkflow(t *testing.T) {
	matcher := &models.Matcher{
		ID:                  uuid.New(),
		Name:                "Test Matcher",
		ConfirmationHistory: make([]bool, 0),
	}

	maxLength := 3

	// Add confirmations and check percentage at each step
	matcher.AddConfirmation(true, maxLength)
	if percentage := matcher.GetConfirmationPercentage(); percentage != 100.0 {
		t.Errorf("After 1 confirmation (true), percentage = %v, expected 100.0", percentage)
	}

	matcher.AddConfirmation(false, maxLength)
	if percentage := matcher.GetConfirmationPercentage(); percentage != 50.0 {
		t.Errorf("After 2 confirmations (true, false), percentage = %v, expected 50.0", percentage)
	}

	matcher.AddConfirmation(true, maxLength)
	if percentage := matcher.GetConfirmationPercentage(); percentage != 66.66666666666666 {
		t.Errorf("After 3 confirmations (true, false, true), percentage = %v, expected ~66.67", percentage)
	}

	// Add another confirmation - should remove the oldest (first true)
	matcher.AddConfirmation(false, maxLength)
	if percentage := matcher.GetConfirmationPercentage(); percentage != 33.33333333333333 {
		t.Errorf("After 4 confirmations with max 3 (false, true, false), percentage = %v, expected ~33.33", percentage)
	}

	if length := matcher.GetConfirmationHistoryLength(); length != maxLength {
		t.Errorf("History length = %v, expected %v", length, maxLength)
	}
}

func TestMatcher_FromDBAndToDB_ConfirmationHistoryRoundTrip(t *testing.T) {
	// Prepare DB model with confirmation history
	dbMatcher := &models.Matcher{
		ID:                  uuid.New(),
		Name:                "roundtrip",
		ConfirmationHistory: []bool{true, false, true},
	}

	// Convert to API model using FromDB and back using MatcherWithoutID/MatcherToDB
	apiModel := dbMatcher.FromDB()

	// Ensure API model contains confirmation history
	if apiModel.GetConfirmationHistory() == nil {
		t.Fatalf("FromDB did not populate ConfirmationHistory")
	}

	// Create NoId model and convert back
	noID := models.MatcherWithoutID(&apiModel)
	db2 := models.MatcherToDB(noID, "user-1")

	// Compare histories
	if len(db2.ConfirmationHistory) != len(dbMatcher.ConfirmationHistory) {
		t.Fatalf("Roundtrip history length mismatch: got %d want %d",
			len(db2.ConfirmationHistory), len(dbMatcher.ConfirmationHistory))
	}
	for i := range db2.ConfirmationHistory {
		if db2.ConfirmationHistory[i] != dbMatcher.ConfirmationHistory[i] {
			t.Fatalf("Roundtrip history mismatch at %d: got %v want %v",
				i, db2.ConfirmationHistory[i], dbMatcher.ConfirmationHistory[i])
		}
	}
}
