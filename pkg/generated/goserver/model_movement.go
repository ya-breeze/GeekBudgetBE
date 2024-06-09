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

type Movement struct {
	Amount float64 `json:"amount"`

	CurrencyID string `json:"currencyID"`

	AccountID string `json:"accountID"`

	Description string `json:"description,omitempty"`
}

// AssertMovementRequired checks if the required fields are not zero-ed
func AssertMovementRequired(obj Movement) error {
	elements := map[string]interface{}{
		"amount":     obj.Amount,
		"currencyID": obj.CurrencyID,
		"accountID":  obj.AccountID,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertMovementConstraints checks if the values respects the defined constraints
func AssertMovementConstraints(obj Movement) error {
	return nil
}
