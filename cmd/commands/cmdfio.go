//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
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

			// parse JSON
			var transactions FioTransactions
			if err = json.Unmarshal(data, &transactions); err != nil {
				return fmt.Errorf("can't parse JSON: %w", err)
			}

			fc, err := NewFioConverter(goserver.Account{
				Id:   "123",
				Name: "Fio",
				Type: "asset`",
			})
			if err != nil {
				return fmt.Errorf("can't create FioConverter: %w", err)
			}

			fmt.Printf("Account: %s:%s\n",
				transactions.AccountStatement.Info.AccountId, transactions.AccountStatement.Info.BankId)
			fmt.Printf("Read %d transaction(s)\n:", len(transactions.AccountStatement.TransactionList.Transaction))
			for _, t := range transactions.AccountStatement.TransactionList.Transaction {
				fmt.Printf("%d - %q\n", t.ID.Value, t.Comment.Value)
				tr, err := fc.ConvertFioToTransaction(t)
				if err != nil {
					return fmt.Errorf("can't convert FIO transaction: %w", err)
				}
				printTransactionNoID(tr)
				fmt.Println()
			}

			return nil
		},
	}
	res.Flags().StringVarP(&jsonFile, "json-file", "f", "", "File with FIO transactions in JSON")

	return res
}

type FioStringColumn struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

type FioIntColumn struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

type FioFloatColumn struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
	ID    int     `json:"id"`
}
type FioTransaction struct {
	Date             FioStringColumn `json:"column0"`  // Datum
	ID               FioIntColumn    `json:"column22"` // ID pohybu
	Amount           FioFloatColumn  `json:"column1"`  // ID Objem
	Currency         FioStringColumn `json:"column14"` // Měna
	PartnerAccountID FioStringColumn `json:"column2"`  // Protiúčet
	PartnerBankID    FioStringColumn `json:"column3"`  // Kód banky protiúčtu
	VS               FioStringColumn `json:"column5"`  // VS protiúčtu (Variabilní symbol)
	PartnerName      FioStringColumn `json:"column10"` // Název protiúčtu
	PartnerBankName  FioStringColumn `json:"column12"` // Název banky protiúčtu
	Column7          FioStringColumn `json:"column7"`  // Uživatelská identifikace
	Type             FioStringColumn `json:"column8"`  // Typ
	User             FioStringColumn `json:"column9"`  // Provedl
	Comment          FioStringColumn `json:"column25"` // Komentář
	InfoForReceiver  FioStringColumn `json:"column16"` // Zpráva pro příjemce
}

type FioTransactions struct {
	AccountStatement struct {
		Info            goserver.BankAccountInfo `json:"info"`
		TransactionList struct {
			Transaction []FioTransaction `json:"transaction"`
		} `json:"transactionList"`
	} `json:"accountStatement"`
}

type FioConverter struct {
	account  goserver.Account
	r        *regexp.Regexp
	location *time.Location
}

func NewFioConverter(account goserver.Account) (*FioConverter, error) {
	// Example:
	// 0: 0 - "Nákup: IKEA ZLICIN RESTAURA,  Skandinavska 15a, Praha 13, 155 00, CZE, dne 31.8.2024, částka  383.00 CZK"
	// 0: 1 - "Nákup"
	// 0: 2 - "IKEA ZLICIN RESTAURA"
	// 0: 3 - "Skandinavska 15a, Praha 13, 155 00, CZE"
	// 0: 4 - "31.8.2024"
	// 0: 5 - "383.00"
	// 0: 6 - "CZK"
	r := regexp.MustCompile(`^(\p{L}+): ([^,]+),  (.+), dne ([\.\d]+), částka  ([\.\d]+) (\p{L}+)$`)

	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &FioConverter{account: account, r: r, location: loc}, nil
}

//nolint:funlen,cyclop // to be refactored
func (fc *FioConverter) ConvertFioToTransaction(fio FioTransaction) (goserver.TransactionNoId, error) {
	var res goserver.TransactionNoId
	tokens := fc.r.FindAllStringSubmatch(fio.Comment.Value, -1)
	if tokens == nil {
		t, err := time.ParseInLocation("2006-01-02-0700", fio.Date.Value, fc.location)
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", fio.Date.Value, err)
		}

		d := fio.Comment.Value
		if d != fio.InfoForReceiver.Value {
			d = d + "; " + fio.InfoForReceiver.Value
		}

		res = goserver.TransactionNoId{
			Date:        t,
			Description: d,
			Movements: []goserver.Movement{
				{
					Amount:     -fio.Amount.Value,
					CurrencyId: fio.Currency.Value,
				},
				{
					AccountId:  fc.account.Id,
					Amount:     fio.Amount.Value,
					CurrencyId: fio.Currency.Value,
				},
			},
		}
	} else {
		t, err := time.Parse("2.1.2006", tokens[0][4])
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", tokens[0][4], err)
		}
		t = t.In(fc.location)

		m, err := strconv.ParseFloat(tokens[0][5], 64)
		if err != nil {
			return res, fmt.Errorf("can't parse amount %q: %w", tokens[0][5], err)
		}

		res = goserver.TransactionNoId{
			Date:        t,
			Place:       tokens[0][3],
			Description: fmt.Sprintf("%s: %s", tokens[0][1], tokens[0][2]),
			Movements: []goserver.Movement{
				{
					Amount:     m,
					CurrencyId: tokens[0][6],
				},
				{
					AccountId:  fc.account.Id,
					Amount:     fio.Amount.Value,
					CurrencyId: fio.Currency.Value,
				},
			},
		}
	}

	res.Tags = append(res.Tags, "fio")
	res.ExternalIds = append(res.ExternalIds, strconv.Itoa(fio.ID.Value))
	res.PartnerAccount = fio.PartnerAccountID.Value + "/" + fio.PartnerBankID.Value
	if fio.VS.Value != "" {
		res.PartnerAccount += " vs=" + fio.VS.Value
	}
	res.PartnerName = fio.PartnerName.Value

	if fio.User.Value == "Korolev, Ilya" {
		res.Tags = append(res.Tags, "ilya")
	} else if fio.User.Value == "Koroleva, Anzhela" {
		res.Tags = append(res.Tags, "angela")
	}

	b, err := json.Marshal(fio)
	if err != nil {
		return res, fmt.Errorf("can't marshal FIO transaction: %w", err)
	}
	res.UnprocessedSources, err = compress(string(b))
	if err != nil {
		return res, fmt.Errorf("can't compress FIO transaction: %w", err)
	}

	return res, nil
}

func compress(s string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(s)); err != nil {
		return "", fmt.Errorf("can't write to gzip: %w", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("can't close gzip: %w", err)
	}

	return b.String(), nil
}

func printTransactionNoID(t goserver.TransactionNoId) {
	fmt.Printf("TransactionNoId: %+v\n", t)
}
