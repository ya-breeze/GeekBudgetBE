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

type WebMovement struct {
	Amount       float64 `json:"amount"`
	AccountID    string  `json:"accountID"`
	AccountName  string  `json:"accountName"`
	CurrencyID   string  `json:"currencyID"`
	CurrencyName string  `json:"currencyName"`
}

type WebTransaction struct {
	ID             string        `json:"id"`
	Date           time.Time     `json:"date"`
	Description    string        `json:"description,omitempty"`
	Place          string        `json:"place,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	PartnerName    string        `json:"partnerName,omitempty"`
	PartnerAccount string        `json:"partnerAccount,omitempty"`
	Movements      []WebMovement `json:"movements"`
}

type WebMatcherAndTransaction struct {
	MatcherId   string         `json:"matcherId"`
	Transaction WebTransaction `json:"transaction"`
}

type WebUnprocessedTransaction struct {
	Transaction WebTransaction             `json:"transaction"`
	Matched     []WebMatcherAndTransaction `json:"matched"`
	Duplicates  []WebTransaction           `json:"duplicates"`
}
