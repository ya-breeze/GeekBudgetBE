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

type BankImporterNoIdLastImportsInner struct {

	// Date of import
	Date time.Time `json:"date,omitempty"`

	// Status of import
	Status string `json:"status,omitempty"`
}

type BankImporterNoIdLastImportsInnerInterface interface {
	GetDate() time.Time
	GetStatus() string
}

func (c *BankImporterNoIdLastImportsInner) GetDate() time.Time {
	return c.Date
}
func (c *BankImporterNoIdLastImportsInner) GetStatus() string {
	return c.Status
}

// AssertBankImporterNoIdLastImportsInnerRequired checks if the required fields are not zero-ed
func AssertBankImporterNoIdLastImportsInnerRequired(obj BankImporterNoIdLastImportsInner) error {
	return nil
}

// AssertBankImporterNoIdLastImportsInnerConstraints checks if the values respects the defined constraints
func AssertBankImporterNoIdLastImportsInnerConstraints(obj BankImporterNoIdLastImportsInner) error {
	return nil
}
