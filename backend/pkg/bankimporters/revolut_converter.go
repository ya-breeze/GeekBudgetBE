package bankimporters

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

const (
	RevolutIndexType          = 0
	RevolutIndexProduct       = 1
	RevolutIndexStartedDate   = 2
	RevolutIndexCompletedDate = 3
	RevolutIndexDescription   = 4
	RevolutIndexAmount        = 5
	RevolutIndexFee           = 6
	RevolutIndexCurrency      = 7
	RevolutIndexState         = 8
	RevolutIndexBalance       = 9

	RevolutExchangePrefix   = "EXCHANGE: Exchanged to "
	RevolutExchangePrefixRU = "Обмен валюты: Обменено на "
)

//nolint:gochecknoglobals // const list of fields in Revolut file
var revolutCSVFieldsEN = []string{
	"Type", "Product", "Started Date",
	"Completed Date", "Description",
	"Amount", "Fee", "Currency", "State", "Balance",
}

//nolint:gochecknoglobals // const list of fields in Revolut file (Russian)
var revolutCSVFieldsRU = []string{
	"Тип", "Продукт", "Дата начала",
	"Дата выполнения", "Описание",
	"Сумма", "Комиссия", "Валюта", "State", "Остаток средств",
}

type RevolutConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	location     *time.Location
	cp           CurrencyProvider
}

func NewRevolutConverter(logger *slog.Logger, bankImporter goserver.BankImporter, cp CurrencyProvider,
) (*RevolutConverter, error) {
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &RevolutConverter{
		logger:       logger,
		bankImporter: bankImporter,
		location:     loc,
		cp:           cp,
	}, nil
}

func (fc *RevolutConverter) ParseAndImport(
	format, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	return fc.ParseTransactions(context.Background(), format, data)
}

func (fc *RevolutConverter) ParseTransactions(ctx context.Context, format, data string,
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
	type currencyState struct {
		firstBalance decimal.Decimal
		firstAmount  decimal.Decimal
		firstDate    time.Time
		lastBalance  decimal.Decimal
		lastAmount   decimal.Decimal
		lastDate     time.Time
		currencyName string
	}
	states := make(map[string]*currencyState)

	info := goserver.BankAccountInfo{}
	res := make([]goserver.TransactionNoId, 0, len(records))
	for i, record := range records {
		if fc.shouldSkipRecord(i, record) {
			continue
		}

		var tr goserver.TransactionNoId
		tr, err = fc.convertToTransaction(ctx, fc.bankImporter, record)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert Revolut transaction: %w", err)
		}
		res = append(res, tr)

		balance, err := decimal.NewFromString(record[RevolutIndexBalance])
		if err != nil {
			return nil, nil, fmt.Errorf("can't parse balance (%v): %w", record, err)
		}
		amount, err := decimal.NewFromString(record[RevolutIndexAmount])
		if err != nil {
			return nil, nil, fmt.Errorf("can't parse amount (%v): %w", record, err)
		}
		curr := record[RevolutIndexCurrency]

		state, ok := states[curr]
		if !ok {
			state = &currencyState{
				firstBalance: balance,
				firstAmount:  amount,
				firstDate:    tr.Date,
				currencyName: curr,
			}
			states[curr] = state
		}
		state.lastBalance = balance
		state.lastAmount = amount
		state.lastDate = tr.Date
	}

	for curr, s := range states {
		currencyID, err := fc.cp.GetCurrencyIdByName(ctx, curr)
		if err != nil {
			return nil, nil, fmt.Errorf("can't resolve currency %q: %w", curr, err)
		}

		var closing, opening decimal.Decimal
		var lastUpdatedDate time.Time
		if s.firstDate.After(s.lastDate) || (s.firstDate.Equal(s.lastDate)) {
			// Assume newest first if first date is after or equal to last date
			closing = s.firstBalance
			opening = s.lastBalance.Sub(s.lastAmount)
			lastUpdatedDate = s.firstDate
		} else {
			// Oldest first
			closing = s.lastBalance
			opening = s.firstBalance.Sub(s.firstAmount)
			lastUpdatedDate = s.lastDate
		}
		bal := goserver.BankAccountInfoBalancesInner{
			CurrencyId:     currencyID,
			ClosingBalance: closing,
			OpeningBalance: opening,
		}
		if !lastUpdatedDate.IsZero() {
			bal.LastUpdatedAt = &lastUpdatedDate
		}
		info.Balances = append(info.Balances, bal)
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

	if len(records[0]) != len(revolutCSVFieldsEN) {
		return fmt.Errorf(
			"revolut record has unexpected number of columns (%d), expected number is %d",
			len(records[0]), len(revolutCSVFieldsEN))
	}

	isEN := true
	isRU := true
	for i := range records[0] {
		if records[0][i] != revolutCSVFieldsEN[i] {
			isEN = false
		}
		if records[0][i] != revolutCSVFieldsRU[i] {
			isRU = false
		}
	}

	if !isEN && !isRU {
		return fmt.Errorf(
			"revolut record has unexpected columns: %v", records[0])
	}

	return nil
}

func (fc *RevolutConverter) shouldSkipRecord(i int, record []string) bool {
	// skip row with header
	if i == 0 {
		return true
	}

	state := record[RevolutIndexState]
	if state != "COMPLETED" && state != "ВЫПОЛНЕНО" {
		fc.logger.Info("Skipping transaction because of state", "state", state)
		return true
	}

	if len(record[RevolutIndexBalance]) == 0 {
		fc.logger.Info("Skipping transaction without balance", "record", record)
		return true
	}

	return false
}

func (fc *RevolutConverter) convertToTransaction(ctx context.Context, _ goserver.BankImporter, record []string,
) (goserver.TransactionNoId, error) {
	var err error
	var res goserver.TransactionNoId

	strCurrencyID, err := fc.cp.GetCurrencyIdByName(ctx, record[RevolutIndexCurrency])
	if err != nil {
		return res, fmt.Errorf("can't resolve currency %q: %w", record[RevolutIndexCurrency], err)
	}

	res.Date, err = time.ParseInLocation("2006-01-02 15:04:05", record[RevolutIndexStartedDate], fc.location)
	if err != nil {
		res.Date, err = time.ParseInLocation("1/2/06 15:04", record[RevolutIndexStartedDate], fc.location)
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", record[RevolutIndexStartedDate], err)
		}
	}

	res.Description = record[RevolutIndexType] + ": " + record[RevolutIndexDescription]

	amount, err := decimal.NewFromString(record[RevolutIndexAmount])
	if err != nil {
		return res, fmt.Errorf("can't parse amount %q: %w", record[RevolutIndexAmount], err)
	}
	feeAmount, err := decimal.NewFromString(record[RevolutIndexFee])
	if err != nil {
		return res, fmt.Errorf("can't parse fee %q: %w", record[RevolutIndexFee], err)
	}

	res.Movements = make([]goserver.Movement, 0, 3)
	if !amount.IsZero() {
		res.Movements = append(res.Movements, goserver.Movement{
			Amount:     amount.Neg(),
			CurrencyId: strCurrencyID,
		})
	}
	remainingAmount := amount.Sub(feeAmount)
	if !remainingAmount.IsZero() {
		res.Movements = append(res.Movements, goserver.Movement{
			AccountId:  fc.bankImporter.AccountId,
			Amount:     remainingAmount,
			CurrencyId: strCurrencyID,
		})
	}
	if !feeAmount.IsZero() {
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
	res.ExternalIds = append(res.ExternalIds, HashString(res.UnprocessedSources))

	return res, nil
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
		isExchange := strings.HasPrefix(transactions[i].Description, RevolutExchangePrefix)
		isExchangeRU := strings.HasPrefix(transactions[i].Description, RevolutExchangePrefixRU)
		if !isExchange && !isExchangeRU {
			res = append(res, transactions[i])
			continue
		}
		prefix := RevolutExchangePrefix
		if isExchangeRU {
			prefix = RevolutExchangePrefixRU
		}
		cur := strings.TrimPrefix(transactions[i].Description, prefix)
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
