package database_test

import (
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var testFamilyID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func TestAddMatcherConfirmation(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false, MatcherConfirmationHistoryMax: 3}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	// create a matcher using Storage API
	noid := &goserver.MatcherNoId{
		ConfirmationHistory: []bool{},
	}
	created, err := st.CreateMatcher(testFamilyID, noid)
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	// helper to add confirmation and fail on error
	add := func(id string, confirmed bool) {
		if e := st.AddMatcherConfirmation(testFamilyID, id, confirmed); e != nil {
			t.Fatalf("AddMatcherConfirmation failed: %v", e)
		}
	}

	// add confirmations
	add(created.Id, true)

	loadedG, err := st.GetMatcher(testFamilyID, created.Id)
	if err != nil {
		t.Fatalf("failed to load matcher: %v", err)
	}
	if len(loadedG.ConfirmationHistory) != 1 || loadedG.ConfirmationHistory[0] != true {
		t.Fatalf("unexpected confirmation history: %v", loadedG.ConfirmationHistory)
	}

	// add more to exceed max
	add(created.Id, false)
	add(created.Id, true)
	add(created.Id, false)

	loadedG2, err := st.GetMatcher(testFamilyID, created.Id)
	if err != nil {
		t.Fatalf("failed to load matcher: %v", err)
	}

	if len(loadedG2.ConfirmationHistory) != cfg.MatcherConfirmationHistoryMax {
		t.Fatalf("history length = %d, expected %d", len(loadedG2.ConfirmationHistory), cfg.MatcherConfirmationHistoryMax)
	}

	// Check the most recent confirmation equals false
	if loadedG2.ConfirmationHistory[len(loadedG2.ConfirmationHistory)-1] != false {
		t.Fatalf("most recent confirmation expected false, got %v", loadedG2.ConfirmationHistory)
	}
}
