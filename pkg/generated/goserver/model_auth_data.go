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

type AuthData struct {
	Email string `json:"email"`

	Password string `json:"password"`
}

type AuthDataInterface interface {
	GetEmail() string
	GetPassword() string
}

func (c *AuthData) GetEmail() string {
	return c.Email
}
func (c *AuthData) GetPassword() string {
	return c.Password
}

// AssertAuthDataRequired checks if the required fields are not zero-ed
func AssertAuthDataRequired(obj AuthData) error {
	elements := map[string]interface{}{
		"email":    obj.Email,
		"password": obj.Password,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertAuthDataConstraints checks if the values respects the defined constraints
func AssertAuthDataConstraints(obj AuthData) error {
	return nil
}
