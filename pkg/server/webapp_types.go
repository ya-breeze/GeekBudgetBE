package server

import "time"

type WebAggregation struct {
	From        time.Time                `json:"from"`
	To          time.Time                `json:"to"`
	Granularity string                   `json:"granularity"`
	Intervals   []time.Time              `json:"intervals"`
	Currencies  []WebCurrencyAggregation `json:"currencies"`
}

type WebCurrencyAggregation struct {
	CurrencyId   string `json:"currencyId"`
	CurrencyName string `json:"currencyName"`

	Intervals []time.Time          `json:"intervals"`
	Accounts  []AccountAggregation `json:"accounts"`
}

type AccountAggregation struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`

	Amounts []float64 `json:"amounts"`
}
