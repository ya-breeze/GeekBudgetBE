//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CmdFio(log *slog.Logger) *cobra.Command {
	res := &cobra.Command{
		Use:   "fio",
		Short: "Work with FIO API",
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	res.AddCommand(fetchFIO(log))
	res.AddCommand(parseFIO(log))

	return res
}

func fetchFIO(log *slog.Logger) *cobra.Command {
	var tokenFile, outputFile string
	res := &cobra.Command{
		Use:          "fetch",
		Short:        "Fetch transactions from FIO API",
		SilenceUsage: true,
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

			res, err := bankimporters.FetchFioTransactions(log, cmd.Context(), token)
			if err != nil {
				return fmt.Errorf("can't fetch FIO transactions: %w", err)
			}

			if outputFile != "" {
				err = os.WriteFile(outputFile, res, 0o600)
				if err != nil {
					return fmt.Errorf("can't write transactions to file %q: %w", outputFile, err)
				}
			} else {
				fmt.Println(string(res))
			}

			return nil
		},
	}
	res.Flags().StringVarP(&tokenFile, "token-file", "f", "", "File with FIO API token")
	res.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write transactions to file")

	return res
}

func parseFIO(log *slog.Logger) *cobra.Command {
	var jsonFile string
	var hideTransactions *bool
	res := &cobra.Command{
		Use:          "parse",
		Short:        "Parse FIO transactions from JSON",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			var data []byte

			if jsonFile != "" {
				// read JSON from file
				data, err = os.ReadFile(jsonFile)
				if err != nil {
					return fmt.Errorf("can't read file %q: %w", jsonFile, err)
				}
			} else {
				// read JSON from stdin
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("can't read JSON from stdin: %w", err)
				}
			}

			fc, err := bankimporters.NewFioConverter(
				log,
				goserver.BankImporter{
					AccountId: "__accountID__",
				}, []goserver.Currency{
					{Id: "__CZK_ID__", Name: "CZK"},
					{Id: "__EUR_ID__", Name: "EUR"},
					{Id: "__USD_ID__", Name: "USD"},
				})
			if err != nil {
				return fmt.Errorf("can't create FioConverter: %w", err)
			}

			info, transactions, err := fc.ParseTransactions(data)
			if err != nil {
				return fmt.Errorf("can't parse FIO transactions: %w", err)
			}
			for _, b := range info.Balances {
				fmt.Printf("Balance for %s\n", b.CurrencyId)
				fmt.Printf("- Opening balance: %v\n", b.OpeningBalance)
				fmt.Printf("- Closing balance: %v\n", b.ClosingBalance)
			}
			fmt.Printf("Parsed transactions: %d\n", len(transactions))
			if hideTransactions != nil && !*hideTransactions {
				for _, t := range transactions {
					printTransactionNoID(t)
					fmt.Println()
				}
			}

			return nil
		},
	}
	res.Flags().StringVarP(&jsonFile, "json-file", "f", "", "File with FIO transactions in JSON")
	hideTransactions = res.Flags().BoolP("hide-transactions", "q", false, "Don't print transactions")

	return res
}

func printTransactionNoID(t goserver.TransactionNoId) {
	t.UnprocessedSources = "__replaced__"
	json, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("TransactionNoId: %s\n", json)
}
