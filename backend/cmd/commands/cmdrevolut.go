package commands

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CmdRevolut(log *slog.Logger) *cobra.Command {
	res := &cobra.Command{
		Use:   "revolut",
		Short: "Work with Revolut transactions",
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	res.AddCommand(parseRevolut(log))

	return res
}

func parseRevolut(log *slog.Logger) *cobra.Command {
	var file string
	var hideTransactions *bool
	res := &cobra.Command{
		Use:          "parse",
		Short:        "Parse Revolut transactions from file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			var data []byte

			ext := strings.Trim(strings.ToLower(filepath.Ext(file)), ".")
			if ext != "csv" && ext != "xlsx" {
				return fmt.Errorf("unsupported file extension %q", ext)
			}

			if file != "" {
				// read from file
				data, err = os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("can't read file %q: %w", file, err)
				}
			} else {
				// read from stdin
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("can't read from stdin: %w", err)
				}
			}

			cp := bankimporters.NewSimpleCurrencyProvider([]goserver.Currency{
				{Id: "__CZK_ID__", Name: "CZK"},
				{Id: "__EUR_ID__", Name: "EUR"},
				{Id: "__USD_ID__", Name: "USD"},
			})
			rc, err := bankimporters.NewRevolutConverter(
				log,
				goserver.BankImporter{
					AccountId: "__accountID__",
				}, cp)
			if err != nil {
				return fmt.Errorf("can't create Revolut converter: %w", err)
			}

			info, transactions, err := rc.ParseTransactions(cmd.Context(), ext, string(data))
			if err != nil {
				return fmt.Errorf("can't parse Revolut transactions: %w", err)
			}

			printResults(info, transactions, hideTransactions)

			return nil
		},
	}
	res.Flags().StringVarP(&file, "file", "f", "", "CSV file with Revolut transactions")
	hideTransactions = res.Flags().BoolP("hide-transactions", "q", false, "Don't print transactions")

	return res
}
