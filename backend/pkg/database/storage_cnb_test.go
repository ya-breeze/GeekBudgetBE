package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func TestCNBRatesStorage(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	date1 := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	rates1 := map[string]decimal.Decimal{
		"USD": decimal.NewFromFloat(22.5),
		"EUR": decimal.NewFromFloat(24.5),
	}

	date2 := time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)
	rates2 := map[string]decimal.Decimal{
		"USD": decimal.NewFromFloat(22.6),
		"EUR": decimal.NewFromFloat(24.6),
	}

	t.Run("Save and Get Specific Date", func(t *testing.T) {
		err := st.SaveCNBRates(rates1, date1)
		if err != nil {
			t.Fatalf("failed to save rates: %v", err)
		}

		retrieved, err := st.GetCNBRates(date1)
		if err != nil {
			t.Fatalf("failed to get rates: %v", err)
		}

		if len(retrieved) != len(rates1) {
			t.Errorf("expected %d rates, got %d", len(rates1), len(retrieved))
		}

		for k, v := range rates1 {
			if !retrieved[k].Equal(v) {
				t.Errorf("expected %s rate %v, got %v", k, v, retrieved[k])
			}
		}
	})

	t.Run("Get Latest Rates", func(t *testing.T) {
		err := st.SaveCNBRates(rates2, date2)
		if err != nil {
			t.Fatalf("failed to save second rates: %v", err)
		}

		// Passing zero time should return latest
		retrieved, err := st.GetCNBRates(time.Time{})
		if err != nil {
			t.Fatalf("failed to get latest rates: %v", err)
		}

		if len(retrieved) != len(rates2) {
			t.Errorf("expected %d rates, got %d", len(rates2), len(retrieved))
		}

		for k, v := range rates2 {
			if !retrieved[k].Equal(v) {
				t.Errorf("expected %s rate %v, got %v", k, v, retrieved[k])
			}
		}
	})

	t.Run("Overwrite Rates for Same Date", func(t *testing.T) {
		newRates1 := map[string]decimal.Decimal{
			"USD": decimal.NewFromFloat(22.9),
		}

		err := st.SaveCNBRates(newRates1, date1)
		if err != nil {
			t.Fatalf("failed to overwrite rates: %v", err)
		}

		retrieved, err := st.GetCNBRates(date1)
		if err != nil {
			t.Fatalf("failed to get overwritten rates: %v", err)
		}

		if len(retrieved) != 1 {
			t.Errorf("expected 1 rate after overwrite, got %d", len(retrieved))
		}

		if !retrieved["USD"].Equal(newRates1["USD"]) {
			t.Errorf("expected USD rate %v, got %v", newRates1["USD"], retrieved["USD"])
		}
	})
}
