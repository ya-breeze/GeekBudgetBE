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

type CheckMatcher200Response struct {
	Result bool `json:"result,omitempty"`
}

type CheckMatcher200ResponseInterface interface {
	GetResult() bool
}

func (c *CheckMatcher200Response) GetResult() bool {
	return c.Result
}

// AssertCheckMatcher200ResponseRequired checks if the required fields are not zero-ed
func AssertCheckMatcher200ResponseRequired(obj CheckMatcher200Response) error {
	return nil
}

// AssertCheckMatcher200ResponseConstraints checks if the values respects the defined constraints
func AssertCheckMatcher200ResponseConstraints(obj CheckMatcher200Response) error {
	return nil
}
