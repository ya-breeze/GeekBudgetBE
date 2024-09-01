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

type Account struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Description string `json:"description,omitempty"`

	Type string `json:"type"`

	BankInfo BankAccountInfo `json:"bankInfo,omitempty"`
}

// AssertAccountRequired checks if the required fields are not zero-ed
func AssertAccountRequired(obj Account) error {
	elements := map[string]interface{}{
		"id":   obj.Id,
		"name": obj.Name,
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertBankAccountInfoRequired(obj.BankInfo); err != nil {
		return err
	}
	return nil
}

// AssertAccountConstraints checks if the values respects the defined constraints
func AssertAccountConstraints(obj Account) error {
	if err := AssertBankAccountInfoConstraints(obj.BankInfo); err != nil {
		return err
	}
	return nil
}
