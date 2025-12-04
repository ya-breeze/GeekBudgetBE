package bankimporters

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	KBIndexStartedDate      = 0
	KBIndexDateOtherBank    = 1
	KBIndexPartnerAccount   = 2
	KBIndexPartnerName      = 3
	KBIndexAmount           = 4
	KBIndexOriginalAmount   = 5
	KBIndexOriginalCurrency = 6
	KBIndexKurz             = 7
	KBIndexVS               = 8
	KBIndexKS               = 9
	KBIndexSS               = 10
	KBIndexTransactionID    = 11
	KBIndexSystemNote       = 12
	KBIndexXXXNote          = 13
	KBIndexPartnerNote      = 14
	KBIndexAV1              = 15
	KBIndexAV2              = 16
	KBIndexAV3              = 17
	KBIndexAV4              = 18
)

//nolint:gochecknoglobals // const list of fields in KB file
var kbCSVFields = []string{
	"Datum splatnosti",
	"Datum odepsání z jiné banky",
	"Protiúčet a kód banky",
	"Název protiúčtu",
	"Částka",
	"Originální částka",
	"Originální měna",
	"Kurz",
	"VS",
	"KS",
	"SS",
	"Identifikace transakce",
	"Systémový popis",
	"Popis příkazce",
	"Popis pro příjemce",
	"AV pole 1",
	"AV pole 2",
	"AV pole 3",
	"AV pole 4",
}

type KBConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	location     *time.Location
	currencies   []goserver.Currency
}

func NewKBConverter(logger *slog.Logger, bankImporter goserver.BankImporter, currencies []goserver.Currency,
) (*KBConverter, error) {
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &KBConverter{
		logger:       logger,
		bankImporter: bankImporter,
		location:     loc,
		currencies:   currencies,
	}, nil
}

func (fc *KBConverter) ParseAndImport(
	format, data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	return fc.ParseTransactions(data)
}

func (fc *KBConverter) ParseTransactions(data string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	info := goserver.BankAccountInfo{}

	fc.logger.Info("Parsing KB transactions")
	reader := strings.NewReader(data)
	// Create a new reader that decodes windows-1250 to UTF-8
	decoder := transform.NewReader(reader, charmap.Windows1250.NewDecoder())
	scanner := bufio.NewScanner(decoder)
	for range 17 {
		if !scanner.Scan() {
			return nil, nil, errors.New("can't read CSV header")
		}

		line := scanner.Text()
		fc.logger.Info("Processing header: " + line)
	}

	if !scanner.Scan() {
		return nil, nil, errors.New("can't read CSV header")
	}
	line := strings.ReplaceAll(scanner.Text(), "\u00A0", " ")
	cvsFields := strings.Split(line, ";")
	for i := range kbCSVFields {
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
	for _, record := range records {
		var tr goserver.TransactionNoId
		tr, err = fc.ConvertToTransaction(record)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert KB transaction: %w", err)
		}

		res = append(res, tr)
	}

	fc.logger.Info("Successfully parsed KB transactions", "count", len(res))
	return &info, res, nil
}

//nolint:funlen,cyclop // TODO: refactor
func (fc *KBConverter) ConvertToTransaction(record []string) (goserver.TransactionNoId, error) {
	var err error
	var res goserver.TransactionNoId

	currencyIdx := slices.IndexFunc(fc.currencies, func(c goserver.Currency) bool {
		return c.Name == "CZK"
	})
	if currencyIdx == -1 {
		return res, fmt.Errorf("can't find currency %q", "CZK")
	}
	strCurrencyID := fc.currencies[currencyIdx].Id

	res.Date, err = time.ParseInLocation("02.01.2006", record[KBIndexStartedDate], fc.location)
	if err != nil {
		return res, fmt.Errorf("can't parse date %q: %w", record[KBIndexStartedDate], err)
	}

	if len(record[KBIndexSystemNote]) > 0 {
		res.Description = record[KBIndexSystemNote]
	}
	if len(record[KBIndexXXXNote]) > 0 {
		res.Description += ": " + record[KBIndexXXXNote]
	}
	if len(record[KBIndexPartnerNote]) > 0 {
		res.Description += "; " + record[KBIndexPartnerNote]
	}

	if len(record[KBIndexAV1]) > 0 && record[KBIndexAV1] != record[KBIndexXXXNote] {
		res.Description += "; " + record[KBIndexAV1]
	}
	if len(record[KBIndexAV2]) > 0 {
		res.Description += record[KBIndexAV2]
	}
	if len(record[KBIndexAV3]) > 0 {
		res.Description += record[KBIndexAV3]
	}
	if len(record[KBIndexAV4]) > 0 {
		res.Description += record[KBIndexAV4]
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
