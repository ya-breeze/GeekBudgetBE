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

type BankAccountInfo struct {
	AccountId string `json:"accountId,omitempty"`

	BankId string `json:"bankId,omitempty"`

	OpeningBalance float64 `json:"openingBalance,omitempty"`

	ClosingBalance float64 `json:"closingBalance,omitempty"`
}

// AssertBankAccountInfoRequired checks if the required fields are not zero-ed
func AssertBankAccountInfoRequired(obj BankAccountInfo) error {
	return nil
}

// AssertBankAccountInfoConstraints checks if the values respects the defined constraints
func AssertBankAccountInfoConstraints(obj BankAccountInfo) error {
	return nil
}