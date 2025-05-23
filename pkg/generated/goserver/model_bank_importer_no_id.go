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

type BankImporterNoId struct {
	Name string `json:"name"`

	Description string `json:"description,omitempty"`

	// ID of account which is used to for movements from this bank importer
	AccountId string `json:"accountId"`

	// ID of account which is used for fee movements from this bank importer
	FeeAccountId string `json:"feeAccountId,omitempty"`

	// Stores extra data about bank importer. For example could hold \"bank account number\" to be able to distinguish between different bank accounts, or it could hold token for bank API
	Extra string `json:"extra,omitempty"`

	// If true, importer will fetch all transactions from the bank, if false, it will fetch only recent transactions
	FetchAll bool `json:"fetchAll,omitempty"`

	// Type of bank importer. It's used to distinguish between different banks. For example, FIO bank or KB bank.
	Type string `json:"type,omitempty"`

	// Date of last successful import.
	LastSuccessfulImport time.Time `json:"lastSuccessfulImport,omitempty"`

	// List of last imports. It could be shown to user to explain what was imported recently
	LastImports []ImportResult `json:"lastImports,omitempty"`

	// List of mappings which are used to enrich transactions with additional tags
	Mappings []BankImporterNoIdMappingsInner `json:"mappings,omitempty"`
}

type BankImporterNoIdInterface interface {
	GetName() string
	GetDescription() string
	GetAccountId() string
	GetFeeAccountId() string
	GetExtra() string
	GetFetchAll() bool
	GetType() string
	GetLastSuccessfulImport() time.Time
	GetLastImports() []ImportResult
	GetMappings() []BankImporterNoIdMappingsInner
}

func (c *BankImporterNoId) GetName() string {
	return c.Name
}
func (c *BankImporterNoId) GetDescription() string {
	return c.Description
}
func (c *BankImporterNoId) GetAccountId() string {
	return c.AccountId
}
func (c *BankImporterNoId) GetFeeAccountId() string {
	return c.FeeAccountId
}
func (c *BankImporterNoId) GetExtra() string {
	return c.Extra
}
func (c *BankImporterNoId) GetFetchAll() bool {
	return c.FetchAll
}
func (c *BankImporterNoId) GetType() string {
	return c.Type
}
func (c *BankImporterNoId) GetLastSuccessfulImport() time.Time {
	return c.LastSuccessfulImport
}
func (c *BankImporterNoId) GetLastImports() []ImportResult {
	return c.LastImports
}
func (c *BankImporterNoId) GetMappings() []BankImporterNoIdMappingsInner {
	return c.Mappings
}

// AssertBankImporterNoIdRequired checks if the required fields are not zero-ed
func AssertBankImporterNoIdRequired(obj BankImporterNoId) error {
	elements := map[string]interface{}{
		"name":      obj.Name,
		"accountId": obj.AccountId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.LastImports {
		if err := AssertImportResultRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Mappings {
		if err := AssertBankImporterNoIdMappingsInnerRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertBankImporterNoIdConstraints checks if the values respects the defined constraints
func AssertBankImporterNoIdConstraints(obj BankImporterNoId) error {
	for _, el := range obj.LastImports {
		if err := AssertImportResultConstraints(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Mappings {
		if err := AssertBankImporterNoIdMappingsInnerConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
