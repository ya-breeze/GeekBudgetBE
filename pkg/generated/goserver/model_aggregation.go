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

import (
	"time"
)

type Aggregation struct {
	From time.Time `json:"from"`

	To time.Time `json:"to"`

	Granularity string `json:"granularity"`

	Intervals []time.Time `json:"intervals"`

	Currencies []CurrencyAggregation `json:"currencies"`
}

type AggregationInterface interface {
	GetFrom() time.Time
	GetTo() time.Time
	GetGranularity() string
	GetIntervals() []time.Time
	GetCurrencies() []CurrencyAggregation
}

func (c *Aggregation) GetFrom() time.Time {
	return c.From
}
func (c *Aggregation) GetTo() time.Time {
	return c.To
}
func (c *Aggregation) GetGranularity() string {
	return c.Granularity
}
func (c *Aggregation) GetIntervals() []time.Time {
	return c.Intervals
}
func (c *Aggregation) GetCurrencies() []CurrencyAggregation {
	return c.Currencies
}

// AssertAggregationRequired checks if the required fields are not zero-ed
func AssertAggregationRequired(obj Aggregation) error {
	elements := map[string]interface{}{
		"from":        obj.From,
		"to":          obj.To,
		"granularity": obj.Granularity,
		"intervals":   obj.Intervals,
		"currencies":  obj.Currencies,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Currencies {
		if err := AssertCurrencyAggregationRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertAggregationConstraints checks if the values respects the defined constraints
func AssertAggregationConstraints(obj Aggregation) error {
	for _, el := range obj.Currencies {
		if err := AssertCurrencyAggregationConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
