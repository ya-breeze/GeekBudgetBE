package utils

import (
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

func GetCurrency(currencyID string, currencies []goserver.Currency) goserver.Currency {
	for i := range currencies {
		if currencies[i].Id == currencyID {
			return currencies[i]
		}
	}
	return goserver.Currency{}
}

func GetAccount(accountID string, accounts []goserver.Account) goserver.Account {
	for i := range accounts {
		if accounts[i].Id == accountID {
			return accounts[i]
		}
	}
	return goserver.Account{}
}
