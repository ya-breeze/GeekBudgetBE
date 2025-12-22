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
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type FioConverter struct {
	logger       *slog.Logger
	bankImporter goserver.BankImporter
	r            *regexp.Regexp
	location     *time.Location
	cp           CurrencyProvider
}

func NewFioConverter(logger *slog.Logger, bankImporter goserver.BankImporter, cp CurrencyProvider,
) (*FioConverter, error) {
	// Example:
	// 0: 0 - "Nákup: IKEA ZLICIN RESTAURA,  Skandinavska 15a, Praha 13, 155 00, CZE, dne 31.8.2024, částka  383.00 CZK"
	// 0: 1 - "Nákup"
	// 0: 2 - "IKEA ZLICIN RESTAURA"
	// 0: 3 - "Skandinavska 15a, Praha 13, 155 00, CZE"
	// 0: 4 - "31.8.2024"
	// 0: 5 - "383.00"
	// 0: 6 - "CZK"
	r := regexp.MustCompile(`^([\s\p{L}]+): ([^,]+,  )?(.+), dne ([\.\d]+), částka  ([\.\d]+) (\p{L}+)$`)

	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		return nil, fmt.Errorf("can't load location: %w", err)
	}

	return &FioConverter{
		logger:       logger,
		bankImporter: bankImporter,
		r:            r,
		location:     loc,
		cp:           cp,
	}, nil
}

func (fc *FioConverter) Import(ctx context.Context) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	// Fetch transactions from FIO
	body, err := FetchFioTransactions(fc.logger, ctx, fc.bankImporter.Extra, fc.bankImporter.FetchAll)
	if err != nil {
		return nil, nil, fmt.Errorf("can't fetch FIO transactions: %w", err)
	}

	return fc.ParseTransactions(ctx, body)
}

func (fc *FioConverter) ParseTransactions(ctx context.Context, data []byte) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	fc.logger.Info("Parsing FIO transactions")

	var fio FioTransactions
	if err := json.Unmarshal(data, &fio); err != nil {
		return nil, nil, fmt.Errorf("can't unmarshal FIO transactions: %w", err)
	}

	// Convert transactions
	res := make([]goserver.TransactionNoId, 0, len(fio.AccountStatement.TransactionList.Transaction))
	for _, t := range fio.AccountStatement.TransactionList.Transaction {
		tr, err := fc.ConvertFioToTransaction(ctx, fc.bankImporter, t)
		if err != nil {
			return nil, nil, fmt.Errorf("can't convert FIO transaction: %w", err)
		}
		res = append(res, tr)
	}

	fc.logger.Info("Successfully parser FIO transactions", "count", len(res))

	return &fio.AccountStatement.Info, res, nil
}

//nolint:funlen,cyclop // to be refactored
func (fc *FioConverter) ConvertFioToTransaction(ctx context.Context, bi goserver.BankImporter, fio FioTransaction,
) (goserver.TransactionNoId, error) {
	var res goserver.TransactionNoId

	strCurrencyID, err := fc.cp.GetCurrencyIdByName(ctx, fio.Currency.Value)
	if err != nil {
		return res, fmt.Errorf("can't resolve currency %q: %w", fio.Currency.Value, err)
	}

	tokens := fc.r.FindAllStringSubmatch(fio.Comment.Value, -1)
	if tokens == nil {
		t, err := time.ParseInLocation("2006-01-02-0700", fio.Date.Value, fc.location)
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", fio.Date.Value, err)
		}

		d := fio.Comment.Value
		if fio.InfoForReceiver.Value != "" && d != fio.InfoForReceiver.Value {
			d = d + "; " + fio.InfoForReceiver.Value
		}
		if d == "" {
			d = fio.Type.Value
		} else {
			d = fio.Type.Value + ": " + d
		}

		res = goserver.TransactionNoId{
			Date:        t,
			Description: d,
			Movements: []goserver.Movement{
				{
					Amount:     -fio.Amount.Value.InexactFloat64(),
					CurrencyId: strCurrencyID,
				},
				{
					AccountId:  fc.bankImporter.AccountId,
					Amount:     fio.Amount.Value.InexactFloat64(),
					CurrencyId: strCurrencyID,
				},
			},
		}
	} else {
		t, err := time.Parse("2.1.2006", tokens[0][4])
		if err != nil {
			return res, fmt.Errorf("can't parse date %q: %w", tokens[0][4], err)
		}
		t = t.In(fc.location)

		m, err := decimal.NewFromString(tokens[0][5])
		if err != nil {
			return res, fmt.Errorf("can't parse amount %q: %w", tokens[0][5], err)
		}

		var strPaidCurrencyID string
		paidCurrency := tokens[0][6]
		amountFio := fio.Amount.Value
		amountUnknown := fio.Amount.Value.Neg()
		if paidCurrency != fio.Currency.Value {
			var err error
			strPaidCurrencyID, err = fc.cp.GetCurrencyIdByName(ctx, paidCurrency)
			if err != nil {
				return res, fmt.Errorf("can't resolve paid currency %q: %w", paidCurrency, err)
			}
			amountUnknown = m
		} else {
			strPaidCurrencyID = strCurrencyID
		}

		res = goserver.TransactionNoId{
			Date:  t,
			Place: tokens[0][3],
			//nolint:staticcheck // comma and space are from regexp
			Description: fmt.Sprintf("%s: %s", tokens[0][1], strings.Trim(tokens[0][2], ",  ")),
			Movements: []goserver.Movement{
				{
					Amount:     amountUnknown.InexactFloat64(),
					CurrencyId: strPaidCurrencyID,
				},
				{
					AccountId:  fc.bankImporter.AccountId,
					Amount:     amountFio.InexactFloat64(),
					CurrencyId: strCurrencyID,
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

	// if tokens == nil {
	// 	orig, err := json.MarshalIndent(fio, "", "  ")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	result, err := json.MarshalIndent(res, "", "  ")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	color.Red("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	// 	color.Red("Can't parse %q\n", fio.Comment.Value)
	// 	utils.PrintInTwoColumns(string(orig), string(result))
	// 	color.Red("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	// }

	b, err := json.Marshal(fio)
	if err != nil {
		return res, fmt.Errorf("can't marshal FIO transaction: %w", err)
	}
	res.UnprocessedSources = string(b)

	return res, nil
}

//nolint:unused // to be used
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
	Value decimal.Decimal `json:"value"`
	Name  string          `json:"name"`
	ID    int             `json:"id"`
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
