package server

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

type CurrenciesRatesFetcher struct {
	logger  *slog.Logger
	storage database.Storage
	// date -> currency -> rate
	rateCache map[string]map[string]float64

	BaseURL string
}

func NewCurrenciesRatesFetcher(logger *slog.Logger, storage database.Storage) *CurrenciesRatesFetcher {
	return &CurrenciesRatesFetcher{
		logger:    logger,
		storage:   storage,
		rateCache: make(map[string]map[string]float64),
		BaseURL: "https://www.cnb.cz/" +
			"cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt",
	}
}

func (f *CurrenciesRatesFetcher) Convert(
	ctx context.Context, day time.Time, from, to string, amount float64,
) (float64, error) {
	// Format date for the URL in DD.MM.YYYY format
	dateStr := day.Format("02.01.2006")

	// Check if we have cached rates for this date
	dateKey := day.Format("2006-01-02")
	rates, ok := f.rateCache[dateKey]
	if !ok {
		// Try to get rates from DB first
		dbRates, err := f.storage.GetCNBRates(day)
		if err != nil {
			f.logger.Warn("failed to get rates from DB", "error", err, "date", dateKey)
		}

		if len(dbRates) > 0 {
			// Use rates from DB
			rates = dbRates
			f.logger.Debug("using rates from DB", "date", dateKey)
		} else {
			// Fetch rates if not in DB
			rates, err = f.fetchRates(ctx, dateStr)
			if err != nil {
				return 0, fmt.Errorf("failed to fetch currency rates: %w", err)
			}

			// Store rates in DB
			if err := f.storage.SaveCNBRates(rates, day); err != nil {
				f.logger.Warn("failed to store rates to DB", "error", err, "date", dateKey)
			}
		}

		// Cache the rates
		f.rateCache[dateKey] = rates
	}

	// Perform the currency conversion
	return f.performConversion(from, to, amount, rates)
}

// performConversion handles the actual currency conversion calculation
func (f *CurrenciesRatesFetcher) performConversion(
	from, to string, amount float64, rates map[string]float64,
) (float64, error) {
	// Convert CZK to another currency or vice versa
	switch {
	case from == "CZK":
		rate, ok := rates[to]
		if !ok {
			return 0, fmt.Errorf("currency not found: %s", to)
		}
		return amount / rate, nil
	case to == "CZK":
		rate, ok := rates[from]
		if !ok {
			return 0, fmt.Errorf("currency not found: %s", from)
		}
		return amount * rate, nil
	default:
		// For other currency pairs, convert via CZK
		fromRate, ok := rates[from]
		if !ok {
			return 0, fmt.Errorf("currency not found: %s", from)
		}
		toRate, ok := rates[to]
		if !ok {
			return 0, fmt.Errorf("currency not found: %s", to)
		}

		// First convert to CZK, then to target currency
		czk := amount * fromRate
		return czk / toRate, nil
	}
}

func (f *CurrenciesRatesFetcher) fetchRates(ctx context.Context, date string) (map[string]float64, error) {
	url := fmt.Sprintf("%s?date=%s", f.BaseURL, date)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status: %s", resp.Status)
	}

	// Parse the CNB format
	rates := make(map[string]float64)
	scanner := bufio.NewScanner(resp.Body)

	// Skip the first two lines (header)
	scanner.Scan() // Skip first line
	scanner.Scan() // Skip second line

	// Parse the rest of the lines
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")

		if len(parts) < 5 {
			continue
		}

		// Format is: Country|Currency|Amount|Code|Rate
		currencyCode := parts[3]
		amountStr := parts[2]
		rateStr := strings.Replace(parts[4], ",", ".", 1) // Replace comma with dot for decimal point

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue
		}

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			continue
		}

		// Store rate per unit of currency
		rates[currencyCode] = rate / amount
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rates, nil
}
