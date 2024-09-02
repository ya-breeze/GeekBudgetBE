//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CmdFio() *cobra.Command {
	res := &cobra.Command{
		Use:   "fio",
		Short: "Work with FIO API",
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	res.AddCommand(fetch())
	res.AddCommand(parse())

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

			fmt.Println(bankimporters.FetchFioTransactions(
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				cmd.Context(), token))

			return nil
		},
	}
	res.Flags().StringVarP(&tokenFile, "token-file", "f", "", "File with FIO API token")

	return res
}

func parse() *cobra.Command {
	var jsonFile string
	res := &cobra.Command{
		Use:   "parse",
		Short: "Parse FIO transactions from JSON",
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
				slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				goserver.BankImporter{
					AccountId: "<accountID>",
				})
			if err != nil {
				return fmt.Errorf("can't create FioConverter: %w", err)
			}

			_, transactions, err := fc.ParseTransactions(data)
			if err != nil {
				return fmt.Errorf("can't parse FIO transactions: %w", err)
			}
			// fmt.Printf("Opening balance: %v\n", info.OpeningBalance)
			// fmt.Printf("Closing balance: %v\n", info.ClosingBalance)
			for _, t := range transactions {
				printTransactionNoID(t)
				fmt.Println()
			}

			return nil
		},
	}
	res.Flags().StringVarP(&jsonFile, "json-file", "f", "", "File with FIO transactions in JSON")

	return res
}

func printTransactionNoID(t goserver.TransactionNoId) {
	t.UnprocessedSources = "<replaced>"
	fmt.Printf("TransactionNoId: %+v\n", t)
}
