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

	// Stores extra data about bank importer. For example could hold \"bank account number\" to be able to distinguish between different bank accounts, or it could hold token for bank API
	Extra string `json:"extra,omitempty"`

	// Type of bank importer. It's used to distinguish between different banks. For example, FIO bank or KB bank.
	Type string `json:"type,omitempty"`

	// Date of last successful import.
	LastSuccessfulImport time.Time `json:"lastSuccessfulImport,omitempty"`

	LastImports []BankImporterNoIdLastImportsInner `json:"lastImports,omitempty"`
}

// AssertBankImporterNoIdRequired checks if the required fields are not zero-ed
func AssertBankImporterNoIdRequired(obj BankImporterNoId) error {
	elements := map[string]interface{}{
		"name": obj.Name,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.LastImports {
		if err := AssertBankImporterNoIdLastImportsInnerRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertBankImporterNoIdConstraints checks if the values respects the defined constraints
func AssertBankImporterNoIdConstraints(obj BankImporterNoId) error {
	for _, el := range obj.LastImports {
		if err := AssertBankImporterNoIdLastImportsInnerConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
