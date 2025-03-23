package bankimporters

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func FetchFioTransactions(logger *slog.Logger, ctx context.Context, token string, fetchAll bool) ([]byte, error) {
	logger.With("fetchAll", fetchAll).Info("Fetching FIO transactions")
	if false {
		// read from file for testing
		data, err := os.ReadFile("transactions.json")
		if err != nil {
			return nil, fmt.Errorf("can't read file: %w", err)
		}

		return data, nil
	}

	// Prepare today and 90 days ago
	today := time.Now().Format("2006-01-02")
	from := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	if fetchAll {
		// TODO: fetch from the beginning of time, but in my case it requires more
		// currencies, i.e. to support that, currencies should be added automatically
		// when an unknown currency is found in the transactions
		from = time.Now().AddDate(0, 0, -365*2).Format("2006-01-02")
	}
	logger.With("to", today).With("from", from).Info("Fetching transactions")

	// fetch from URL 2024-09-01
	url := fmt.Sprintf("https://fioapi.fio.cz/v1/rest/periods/%s/%s/%s/transactions.json", token, from, today)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}
	defer resp.Body.Close()

	// Read all data from the io.ReadCloser
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d - %s", resp.StatusCode, body)
	}

	logger.Info("Fetched FIO transactions")
	return body, nil
}
