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

type AccountNoId struct {
	Name string `json:"name"`

	Description string `json:"description,omitempty"`

	Type string `json:"type"`

	BankInfo BankAccountInfo `json:"bankInfo,omitempty"`
}

type AccountNoIdInterface interface {
	GetName() string
	GetDescription() string
	GetType() string
	GetBankInfo() BankAccountInfo
}

func (c *AccountNoId) GetName() string {
	return c.Name
}
func (c *AccountNoId) GetDescription() string {
	return c.Description
}
func (c *AccountNoId) GetType() string {
	return c.Type
}
func (c *AccountNoId) GetBankInfo() BankAccountInfo {
	return c.BankInfo
}

// AssertAccountNoIdRequired checks if the required fields are not zero-ed
func AssertAccountNoIdRequired(obj AccountNoId) error {
	elements := map[string]interface{}{
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

// AssertAccountNoIdConstraints checks if the values respects the defined constraints
func AssertAccountNoIdConstraints(obj AccountNoId) error {
	if err := AssertBankAccountInfoConstraints(obj.BankInfo); err != nil {
		return err
	}
	return nil
}
