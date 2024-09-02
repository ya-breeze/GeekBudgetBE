package bankimporters

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type FioConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	r            *regexp.Regexp
	location     *time.Location
}

func NewFioConverter(logger *slog.Logger, bankImporter goserver.BankImporter) (*FioConverter, error) {
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

	return &FioConverter{
		logger:       logger,
		bankImporter: bankImporter,
		r:            r,
		location:     loc,
	}, nil
}

func (fc *FioConverter) Import(ctx context.Context) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	// Fetch transactions from FIO
	body, err := FetchFioTransactions(fc.logger, ctx, fc.bankImporter.Extra)
	if err != nil {
		return nil, nil, fmt.Errorf("can't fetch FIO transactions: %w", err)
	}

	return fc.ParseTransactions(body)
}

func (fc *FioConverter) ParseTransactions(data []byte) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	var fio FioTransactions
	if err := json.Unmarshal(data, &fio); err != nil {
		return nil, nil, fmt.Errorf("can't unmarshal FIO transactions: %w", err)
	}

	// Convert transactions
	res := make([]goserver.TransactionNoId, 0, len(fio.AccountStatement.TransactionList.Transaction))
	for _, t := range fio.AccountStatement.TransactionList.Transaction {
		tr, err := fc.ConvertFioToTransaction(fc.bankImporter, t)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert FIO transaction: %w", err)
		}
		res = append(res, tr)
	}

	return &fio.AccountStatement.Info, res, nil
}

//nolint:funlen,cyclop // to be refactored
func (fc *FioConverter) ConvertFioToTransaction(bi goserver.BankImporter, fio FioTransaction,
) (goserver.TransactionNoId, error) {
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
					AccountId:  fc.bankImporter.AccountId,
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
					AccountId:  fc.bankImporter.AccountId,
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

	// Iterate mappings
	for _, m := range bi.Mappings {
		if m.FieldToMatch == "user" && m.ValueToMatch == fio.User.Value {
			res.Tags = append(res.Tags, m.TagToSet)
		}
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
