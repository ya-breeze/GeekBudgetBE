package commands

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CmdKB(log *slog.Logger) *cobra.Command {
	res := &cobra.Command{
		Use:   "kb",
		Short: "Work with KB transactions",
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	res.AddCommand(parseKB(log))

	return res
}

func parseKB(log *slog.Logger) *cobra.Command {
	var filename string
	var hideTransactions *bool
	res := &cobra.Command{
		Use:          "parse",
		Short:        "Parse KB transactions from file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			var data []byte

			if filename != "" {
				var file *os.File
				file, err = os.Open(filename)
				if err != nil {
					return fmt.Errorf("can't open file %q: %w", filename, err)
				}
				defer file.Close()

				data, err = io.ReadAll(file)
				if err != nil {
					return fmt.Errorf("can't read file %q: %w", filename, err)
				}
			} else {
				// read from stdin
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("can't read from stdin: %w", err)
				}
			}

			rc, err := bankimporters.NewKBConverter(
				log,
				goserver.BankImporter{
					AccountId: "__accountID__",
				}, []goserver.Currency{
					{Id: "__CZK_ID__", Name: "CZK"},
					{Id: "__EUR_ID__", Name: "EUR"},
					{Id: "__USD_ID__", Name: "USD"},
				})
			if err != nil {
				return fmt.Errorf("can't create KB converter: %w", err)
			}

			info, transactions, err := rc.ParseTransactions(string(data))
			if err != nil {
				return fmt.Errorf("can't parse KB transactions: %w", err)
			}

			printResults(info, transactions, hideTransactions)

			return nil
		},
	}
	res.Flags().StringVarP(&filename, "file", "f", "", "CSV file with KB transactions")
	hideTransactions = res.Flags().BoolP("hide-transactions", "q", false, "Don't print transactions")

	return res
}
