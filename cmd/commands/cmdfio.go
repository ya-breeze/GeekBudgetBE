//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func CmdFio() *cobra.Command {
	res := &cobra.Command{
		Use:   "fio",
		Short: "Work with FIO API",
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	res.AddCommand(fetch())

	return res
}

func fetch() *cobra.Command {
	var tokenFile string
	res := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch transactions from FIO API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			var token string

			if tokenFile != "" {
				var data []byte
				data, err = os.ReadFile(tokenFile)
				if err != nil {
					return fmt.Errorf("can't read token from file %q: %w", tokenFile, err)
				}
				token = strings.TrimSpace(string(data))
			} else {
				// read token from keyboard
				fmt.Print("Enter FIO API token: ")
				_, err = fmt.Scanln(&token)
				if err != nil {
					return fmt.Errorf("can't read token from keyboard: %w", err)
				}
			}

			// Prepare today and 90 days ago
			today := time.Now().Format("2006-01-02")
			ago90 := time.Now().AddDate(0, 0, -90).Format("2006-01-02")

			// fetch from URL 2024-09-01
			url := fmt.Sprintf("https://fioapi.fio.cz/v1/rest/periods/%s/%s/%s/transactions.json", token, ago90, today)
			req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("can't create request: %w", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("can't send request: %w", err)
			}
			defer resp.Body.Close()

			// Read all data from the io.ReadCloser
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("can't read response body: %w", err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code %d - %s", resp.StatusCode, body)
			}

			fmt.Print(string(body))

			return nil
		},
	}
	res.Flags().StringVarP(&tokenFile, "token-file", "f", "", "File with FIO API token")

	return res
}
