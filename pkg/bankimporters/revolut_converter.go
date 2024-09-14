package bankimporters

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

const (
	IndexType          = 0
	IndexProduct       = 1
	IndexStartedDate   = 2
	IndexCompletedDate = 3
	IndexDescription   = 4
	IndexAmount        = 5
	IndexFee           = 6
	IndexCurrency      = 7
	IndexState         = 8
	IndexBalance       = 9

	ExchangePrefix = "EXCHANGE: Exchanged to "
)

//nolint:gochecknoglobals // const list of fields in Revolut file
var revolutCSVFields = []string{
	"Type", "Product", "Started Date",
	"Completed Date", "Description",
	"Amount", "Fee", "Currency", "State", "Balance",
}

type RevolutConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	location     *time.Location
	currencies   []goserver.Currency
}

func NewRevolutConverter(logger *slog.Logger, bankImporter goserver.BankImporter, currencies []goserver.Currency,
) (*RevolutConverter, error) {
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &RevolutConverter{
		logger:       logger,
		bankImporter: bankImporter,
		location:     loc,
		currencies:   currencies,
	}, nil
}

func (fc *RevolutConverter) ParseAndImport(ctx context.Context, format, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	return fc.ParseTransactions(format, data)
}

func (fc *RevolutConverter) ParseTransactions(format, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	fc.logger.Info("Parsing Revolut transactions", "format", format)

	var err error
	var records [][]string
	if format == "csv" {
		records, err = fc.parseCSV(data)
	} else if format == "xlsx" {
		records, err = fc.parseXLSX(data)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("can't parse file: %w", err)
	}

	err = fc.checkFormat(records)
	if err != nil {
		return nil, nil, fmt.Errorf("wrong format: %w", err)
	}

	// Convert transactions
	filledOpeningBalances := make(map[string]bool)
	info := goserver.BankAccountInfo{}
	res := make([]goserver.TransactionNoId, 0, len(records))
	for i, record := range records {
		if fc.shouldSkipRecord(i, record) {
			continue
		}

		err = fc.updateBalances(&info, filledOpeningBalances, record)
		if err != nil {
			return nil, nil, fmt.Errorf("can't update balances: %w", err)
		}

		var tr goserver.TransactionNoId
		tr, err = fc.convertToTransaction(fc.bankImporter, record)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert Revolut transaction: %w", err)
		}
		res = append(res, tr)
	}

	res, err = fc.joinExchanges(res)
	if err != nil {
		return nil, nil, fmt.Errorf("can't join exchange transactions: %w", err)
	}

	fc.logger.Info("Successfully parsed Revolut transactions", "count", len(res))
	return &info, res, nil
}

func (fc *RevolutConverter) checkFormat(records [][]string) error {
	if len(records) < 2 {
		return errors.New("no records in CSV file")
	}

	if len(records[0]) != len(revolutCSVFields) {
		return fmt.Errorf(
			"revolut record has unexpected number of columns (%d), expected number is %d",
			len(records[0]), len(revolutCSVFields))
	}
	for i, field := range revolutCSVFields {
		if records[0][i] != field {
			return fmt.Errorf(
				"revolut record has unexpected column %q at position %d, expected column is %q",
				records[0][i], i, field)
		}
	}

	return nil
}

func (fc *RevolutConverter) shouldSkipRecord(i int, record []string) bool {
	// skip row with header
	if i == 0 {
		return true
	}

	if record[IndexState] != "COMPLETED" {
		fc.logger.Info("Skipping transaction because of state", "state", record[IndexState])
		return true
	}

	if len(record[IndexBalance]) == 0 {
		fc.logger.Info("Skipping transaction without balance", "record", record)
		return true
	}

	return false
}

func (fc *RevolutConverter) updateBalances(
	info *goserver.BankAccountInfo, filledOpeningBalances map[string]bool, record []string,
) error {
	balance, err := decimal.NewFromString(record[IndexBalance])
	if err != nil {
		return fmt.Errorf("can't parse balance (%v): %w", record, err)
	}

	// Fill opening/closing balances
	currencyIdx := slices.IndexFunc(fc.currencies, func(c goserver.Currency) bool {
		return c.Name == record[IndexCurrency]
	})
	if currencyIdx == -1 {
		return fmt.Errorf("can't find currency %q", record[IndexCurrency])
	}

	balanceIdx := slices.IndexFunc(info.Balances, func(b goserver.BankAccountInfoBalancesInner) bool {
		return b.CurrencyId == fc.currencies[currencyIdx].Id
	})
	if balanceIdx == -1 {
		info.Balances = append(info.Balances, goserver.BankAccountInfoBalancesInner{
			CurrencyId: fc.currencies[currencyIdx].Id,
		})
		balanceIdx = len(info.Balances) - 1
	}

	info.Balances[balanceIdx].ClosingBalance = balance.InexactFloat64()
	if !filledOpeningBalances[record[IndexCurrency]] {
		var amount decimal.Decimal
		amount, err = decimal.NewFromString(record[IndexAmount])
		if err != nil {
			return fmt.Errorf("can't parse amount (%v): %w", record, err)
		}

		info.Balances[balanceIdx].OpeningBalance = balance.Sub(amount).InexactFloat64()
		filledOpeningBalances[record[IndexCurrency]] = true
	}

	return nil
}

func (fc *RevolutConverter) convertToTransaction(_ goserver.BankImporter, record []string,
) (goserver.TransactionNoId, error) {
	var err error
	var res goserver.TransactionNoId

	currencyIdx := slices.IndexFunc(fc.currencies, func(c goserver.Currency) bool {
		return c.Name == record[IndexCurrency]
	})
	if currencyIdx == -1 {
		return res, fmt.Errorf("can't find currency %q", record[IndexCurrency])
	}
	strCurrencyID := fc.currencies[currencyIdx].Id

	res.Date, err = time.ParseInLocation("2006-01-02 15:04:05", record[IndexStartedDate], fc.location)
	if err != nil {
		res.Date, err = time.ParseInLocation("1/2/06 15:04", record[IndexStartedDate], fc.location)
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", record[IndexStartedDate], err)
		}
	}

	res.Description = record[IndexType] + ": " + record[IndexDescription]

	amount, err := strconv.ParseFloat(record[IndexAmount], 64)
	if err != nil {
		return res, fmt.Errorf("can't parse amount %q: %w", record[IndexAmount], err)
	}
	feeAmount, err := strconv.ParseFloat(record[IndexFee], 64)
	if err != nil {
		return res, fmt.Errorf("can't parse fee %q: %w", record[IndexFee], err)
	}

	res.Movements = []goserver.Movement{
		{
			Amount:     -amount,
			CurrencyId: strCurrencyID,
		},
		{
			AccountId:  fc.bankImporter.AccountId,
			Amount:     amount - feeAmount,
			CurrencyId: strCurrencyID,
		},
	}
	if feeAmount != 0 {
		res.Movements = append(res.Movements, goserver.Movement{
			AccountId:  fc.bankImporter.FeeAccountId,
			Amount:     feeAmount,
			CurrencyId: strCurrencyID,
		})
	}

	res.Tags = append(res.Tags, "revolut")

	b, err := json.Marshal(record)
	if err != nil {
		return res, fmt.Errorf("can't marshal Revolut transaction: %w", err)
	}
	res.UnprocessedSources = string(b)
	res.ExternalIds = append(res.ExternalIds, hashString(res.UnprocessedSources))

	return res, nil
}

func hashString(input string) string {
	// Create a new SHA-256 hash object
	hasher := sha256.New()

	// Write the input string to the hash object
	hasher.Write([]byte(input))

	// Get the resulting hash as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func (fc *RevolutConverter) parseCSV(data string) ([][]string, error) {
	fc.logger.Info("Parsing CSV data")
	reader := csv.NewReader(strings.NewReader(data))
	res, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read CSV file: %w", err)
	}

	return res, nil
}

func (fc *RevolutConverter) parseXLSX(data string) ([][]string, error) {
	fc.logger.Info("Parsing XLSX data")
	f, err := excelize.OpenReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("can't open XLSX file: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	fc.logger.Info("Sheets in XLSX file", "sheets", sheets)

	// Get all the rows in the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("can't get rows from XLSX file: %w", err)
	}

	return rows, nil
}

// Revolut shows exchange transactions as two separate transactions. These transactions should be joined.
func (fc *RevolutConverter) joinExchanges(transactions []goserver.TransactionNoId,
) ([]goserver.TransactionNoId, error) {
	fc.logger.Info("Joining exchange transactions")

	toSkip := make(map[int]bool)
	res := make([]goserver.TransactionNoId, 0, len(transactions))
outerLoop:
	for i := 0; i < len(transactions); i++ {
		if toSkip[i] {
			continue
		}
		if !strings.HasPrefix(transactions[i].Description, ExchangePrefix) {
			res = append(res, transactions[i])
			continue
		}
		cur := strings.TrimPrefix(transactions[i].Description, ExchangePrefix)
		fc.logger.Info("Found exchange transaction", "transaction", transactions[i], "currency", cur)

		// Find matching exchange transaction
		for j := i + 1; j < len(transactions); j++ {
			if transactions[i].Description == transactions[j].Description &&
				transactions[i].Date.Equal(transactions[j].Date) {
				fc.logger.Info("Found matching exchange transaction", "transaction", transactions[j])

				transactions[i].Movements[0] = transactions[j].Movements[1]
				transactions[i].UnprocessedSources += "; " + transactions[j].UnprocessedSources

				// Move fee from second transaction to first
				if len(transactions[j].Movements) == 3 {
					transactions[i].Movements = append(transactions[i].Movements, transactions[j].Movements[2])
				}

				fc.logger.Info("Joined transaction", "transaction", transactions[i])

				// mark second transaction to be skipped
				toSkip[j] = true
				res = append(res, transactions[i])

				continue outerLoop
			}
		}

		// Didn't find matching transaction
		return nil, fmt.Errorf("can't find matching exchange transaction for %v", transactions[i])
	}

	fc.logger.Info("Successfully joined exchange transactions")
	return res, nil
}
