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

type CurrencyNoId struct {
	Name string `json:"name"`

	Description string `json:"description,omitempty"`
}

type CurrencyNoIdInterface interface {
	GetName() string
	GetDescription() string
}

func (c *CurrencyNoId) GetName() string {
	return c.Name
}
func (c *CurrencyNoId) GetDescription() string {
	return c.Description
}

// AssertCurrencyNoIdRequired checks if the required fields are not zero-ed
func AssertCurrencyNoIdRequired(obj CurrencyNoId) error {
	elements := map[string]interface{}{
		"name": obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertCurrencyNoIdConstraints checks if the values respects the defined constraints
func AssertCurrencyNoIdConstraints(obj CurrencyNoId) error {
	return nil
}
