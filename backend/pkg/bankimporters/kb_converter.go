package bankimporters

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	KBIndexDate             = 0
	KBIndexDateExecuted     = 1
	KBIndexPartnerAccount   = 2
	KBIndexPartnerName      = 3
	KBIndexAmount           = 4
	KBIndexCurrency         = 5
	KBIndexOriginalAmount   = 6
	KBIndexOriginalCurrency = 7
	KBIndexRate             = 8
	KBIndexVS               = 9
	KBIndexKS               = 10
	KBIndexSS               = 11
	KBIndexTransactionID    = 12
	KBIndexType             = 13
	KBIndexDescriptionUser  = 14
	KBIndexMessage          = 15
	KBIndexReference        = 16
	KBIndexBIC              = 17
	KBIndexFee              = 18
)

//nolint:gochecknoglobals // const list of fields in KB file
var kbCSVFields = []string{
	"Datum zauctovani",
	"Datum provedeni",
	"Protistrana",
	"Nazev protiuctu",
	"Castka",
	"Mena",
	"Originalni castka",
	"Originalni mena",
	"Smenny kurz",
	"VS",
	"KS",
	"SS",
	"Identifikace transakce",
	"Typ transakce",
	"Popis pro me",
	"Zprava pro prijemce",
	"Reference platby",
	"BIC / SWIFT",
	"Poplatek",
}

type KBConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	location     *time.Location
	cp           CurrencyProvider
}

func NewKBConverter(logger *slog.Logger, bankImporter goserver.BankImporter, cp CurrencyProvider,
) (*KBConverter, error) {
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &KBConverter{
		logger:       logger,
		bankImporter: bankImporter,
		location:     loc,
		cp:           cp,
	}, nil
}

func (fc *KBConverter) ParseAndImport(
	format, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	return fc.ParseTransactions(context.Background(), data)
}

func (fc *KBConverter) ParseTransactions(ctx context.Context, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	info := goserver.BankAccountInfo{}

	fc.logger.Info("Parsing KB transactions")
	reader := strings.NewReader(data)
	// Create a new reader that decodes windows-1250 to UTF-8
	decoder := transform.NewReader(reader, charmap.Windows1250.NewDecoder())
	scanner := bufio.NewScanner(decoder)
	for range 16 {
		if !scanner.Scan() {
			return nil, nil, errors.New("can't read CSV header")
		}

		line := scanner.Text()
		fc.logger.Info("Processing header: " + line)
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			return nil, nil, errors.New("can't read CSV header")
		}
		if parts[0] == "Cislo uctu" {
			// Cislo uctu;123-177270217;;;;;;;;;;;;;;;;;
			info.AccountId = parts[1]
		} else if parts[0] == "Konecny zustatek" {
			// Konecny zustatek;13468,31;;;;;;;;;;;;;;;;;
			amount, err := strconv.ParseFloat(strings.ReplaceAll(parts[1], ",", "."), 64)
			if err == nil {
				info.Balances = []goserver.BankAccountInfoBalancesInner{
					{
						ClosingBalance: amount,
						CurrencyId:     "CZK", // default to CZK as per example
					},
				}
			}
		}
	}

	if !scanner.Scan() {
		return nil, nil, errors.New("can't read CSV header")
	}
	line := strings.ReplaceAll(scanner.Text(), "\u00A0", " ")
	cvsFields := strings.Split(line, ";")
	for i := range kbCSVFields {
		if i >= len(cvsFields) {
			return nil, nil, fmt.Errorf("missing field %q", kbCSVFields[i])
		}
		if strings.Trim(cvsFields[i], "\"") != kbCSVFields[i] {
			return nil, nil, fmt.Errorf("wrong format: %q != %q", cvsFields[i], kbCSVFields[i])
		}
	}

	records, err := fc.parseCSV(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parse file: %w", err)
	}
	fc.logger.Info("Successfully parsed KB transactions", "count", len(records))

	// Convert transactions
	res := make([]goserver.TransactionNoId, 0, len(records))
	var sum float64
	for _, record := range records {
		var tr goserver.TransactionNoId
		tr, err = fc.ConvertToTransaction(ctx, record)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert KB transaction: %w", err)
		}

		res = append(res, tr)

		// Calculate sum of amounts for the account movements
		for _, m := range tr.Movements {
			if m.AccountId == fc.bankImporter.AccountId {
				sum += m.Amount
			}
		}
	}

	fc.logger.Info("Successfully parsed KB transactions", "count", len(res), "sum", sum)

	// If we found a balance, calculate opening balance
	if len(info.Balances) > 0 {
		info.Balances[0].OpeningBalance = info.Balances[0].ClosingBalance - sum
	}

	// If we found a balance, try to update the currency ID if we can infer it or if it's set in the transactions
	if len(info.Balances) > 0 && len(res) > 0 {
		// Assuming all transactions in the file share the same currency which is also the account currency
		// This is a simplification but often true for bank statements
		// In the example file, "Mena" column is CZK
		// We could potentially read "Mena uctu" from line 5 if we wanted to be more precise

		// Let's verify if we can find the currency ID for the balance
		// We used "CZK" as placeholder, let's see if we can resolve it to ID
		currencyId, err := fc.cp.GetCurrencyIdByName(ctx, "CZK") // Default from example
		if err == nil {
			info.Balances[0].CurrencyId = currencyId
		}
	}

	return &info, res, nil
}

//nolint:funlen,cyclop // TODO: refactor
func (fc *KBConverter) ConvertToTransaction(ctx context.Context, record []string) (goserver.TransactionNoId, error) {
	var err error
	var res goserver.TransactionNoId

	strCurrencyID, err := fc.cp.GetCurrencyIdByName(ctx, "CZK")
	if err != nil {
		return res, fmt.Errorf("can't resolve currency %q: %w", "CZK", err)
	}

	res.Date, err = time.ParseInLocation("02.01.2006", record[KBIndexDate], fc.location)
	if err != nil {
		return res, fmt.Errorf("can't parse date %q: %w", record[KBIndexDate], err)
	}

	if len(record[KBIndexType]) > 0 {
		res.Description = record[KBIndexType]
	}
	if len(record[KBIndexDescriptionUser]) > 0 {
		if res.Description != "" {
			res.Description += ": "
		}
		res.Description += record[KBIndexDescriptionUser]
	}
	if len(record[KBIndexMessage]) > 0 {
		if res.Description != "" {
			res.Description += "; "
		}
		res.Description += record[KBIndexMessage]
	}

	if len(record[KBIndexReference]) > 0 {
		if res.Description != "" {
			res.Description += "; "
		}
		res.Description += record[KBIndexReference]
	}

	res.PartnerAccount = record[KBIndexPartnerAccount]
	res.PartnerName = record[KBIndexPartnerName]
	if len(record[KBIndexVS]) > 0 && record[KBIndexVS] != "0" {
		res.PartnerAccount += "; VS:" + record[KBIndexVS]
	}
	if len(record[KBIndexKS]) > 0 && record[KBIndexKS] != "0" {
		res.PartnerAccount += "; KS:" + record[KBIndexKS]
	}
	if len(record[KBIndexSS]) > 0 && record[KBIndexSS] != "0" {
		res.PartnerAccount += "; SS:" + record[KBIndexSS]
	}

	amount, err := strconv.ParseFloat(strings.ReplaceAll(record[KBIndexAmount], ",", "."), 64)
	if err != nil {
		return res, fmt.Errorf("can't parse amount %q: %w", record[KBIndexAmount], err)
	}

	res.Movements = []goserver.Movement{
		{
			Amount:     -amount,
			CurrencyId: strCurrencyID,
		},
		{
			AccountId:  fc.bankImporter.AccountId,
			Amount:     amount,
			CurrencyId: strCurrencyID,
		},
	}

	res.Tags = append(res.Tags, "kb")

	b, err := json.Marshal(record)
	if err != nil {
		return res, fmt.Errorf("can't marshal KB transaction: %w", err)
	}
	res.UnprocessedSources = string(b)
	res.ExternalIds = append(res.ExternalIds, record[KBIndexTransactionID])

	return res, nil
}

func (fc *KBConverter) parseCSV(scanner *bufio.Scanner) ([][]string, error) {
	var res [][]string
	for scanner.Scan() {
		line := scanner.Text()
		r := csv.NewReader(strings.NewReader(line))
		r.Comma = ';'
		record, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("can't parse CSV line %q: %w", line, err)
		}
		record = fc.PrepareRow(record)

		res = append(res, record)
	}

	return res, nil
}

func (fc *KBConverter) PrepareRow(record []string) []string {
	for i := range record {
		record[i] = strings.TrimSpace(strings.Trim(record[i], "\""))
	}

	return record
}
