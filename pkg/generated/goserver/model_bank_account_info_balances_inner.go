// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Geek Budget - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

type BankAccountInfoBalancesInner struct {
	CurrencyId string `json:"currencyId,omitempty"`

	OpeningBalance float64 `json:"openingBalance,omitempty"`

	ClosingBalance float64 `json:"closingBalance,omitempty"`
}

type BankAccountInfoBalancesInnerInterface interface {
	GetCurrencyId() string
	GetOpeningBalance() float64
	GetClosingBalance() float64
}

func (c *BankAccountInfoBalancesInner) GetCurrencyId() string {
	return c.CurrencyId
}
func (c *BankAccountInfoBalancesInner) GetOpeningBalance() float64 {
	return c.OpeningBalance
}
func (c *BankAccountInfoBalancesInner) GetClosingBalance() float64 {
	return c.ClosingBalance
}

// AssertBankAccountInfoBalancesInnerRequired checks if the required fields are not zero-ed
func AssertBankAccountInfoBalancesInnerRequired(obj BankAccountInfoBalancesInner) error {
	return nil
}

// AssertBankAccountInfoBalancesInnerConstraints checks if the values respects the defined constraints
func AssertBankAccountInfoBalancesInnerConstraints(obj BankAccountInfoBalancesInner) error {
	return nil
}