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

type Matcher struct {
	Id string `json:"id"`

	Name string `json:"name"`

	OutputDescription string `json:"outputDescription,omitempty"`

	Amount float64 `json:"amount,omitempty"`

	CurrencyRegExp string `json:"currencyRegExp,omitempty"`

	PartnerNameRegExp string `json:"partnerNameRegExp,omitempty"`

	PartnerAccountNumber string `json:"partnerAccountNumber,omitempty"`

	DescriptionRegExp string `json:"descriptionRegExp,omitempty"`

	ExtraRegExp string `json:"extraRegExp,omitempty"`

	OutputMovements []Movement `json:"outputMovements,omitempty"`
}

type MatcherInterface interface {
	GetId() string
	GetName() string
	GetOutputDescription() string
	GetAmount() float64
	GetCurrencyRegExp() string
	GetPartnerNameRegExp() string
	GetPartnerAccountNumber() string
	GetDescriptionRegExp() string
	GetExtraRegExp() string
	GetOutputMovements() []Movement
}

func (c *Matcher) GetId() string {
	return c.Id
}
func (c *Matcher) GetName() string {
	return c.Name
}
func (c *Matcher) GetOutputDescription() string {
	return c.OutputDescription
}
func (c *Matcher) GetAmount() float64 {
	return c.Amount
}
func (c *Matcher) GetCurrencyRegExp() string {
	return c.CurrencyRegExp
}
func (c *Matcher) GetPartnerNameRegExp() string {
	return c.PartnerNameRegExp
}
func (c *Matcher) GetPartnerAccountNumber() string {
	return c.PartnerAccountNumber
}
func (c *Matcher) GetDescriptionRegExp() string {
	return c.DescriptionRegExp
}
func (c *Matcher) GetExtraRegExp() string {
	return c.ExtraRegExp
}
func (c *Matcher) GetOutputMovements() []Movement {
	return c.OutputMovements
}

// AssertMatcherRequired checks if the required fields are not zero-ed
func AssertMatcherRequired(obj Matcher) error {
	elements := map[string]interface{}{
		"id":   obj.Id,
		"name": obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.OutputMovements {
		if err := AssertMovementRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertMatcherConstraints checks if the values respects the defined constraints
func AssertMatcherConstraints(obj Matcher) error {
	for _, el := range obj.OutputMovements {
		if err := AssertMovementConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
