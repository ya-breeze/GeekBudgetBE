/*
Geek Budget - OpenAPI 3.0

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.0.1
Contact: ilya.korolev@outlook.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package goclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the CurrencyAggregation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CurrencyAggregation{}

// CurrencyAggregation struct for CurrencyAggregation
type CurrencyAggregation struct {
	CurrencyId string               `json:"currencyId"`
	Accounts   []AccountAggregation `json:"accounts"`
}

type _CurrencyAggregation CurrencyAggregation

// NewCurrencyAggregation instantiates a new CurrencyAggregation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCurrencyAggregation(currencyId string, accounts []AccountAggregation) *CurrencyAggregation {
	this := CurrencyAggregation{}
	this.CurrencyId = currencyId
	this.Accounts = accounts
	return &this
}

// NewCurrencyAggregationWithDefaults instantiates a new CurrencyAggregation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCurrencyAggregationWithDefaults() *CurrencyAggregation {
	this := CurrencyAggregation{}
	return &this
}

// GetCurrencyId returns the CurrencyId field value
func (o *CurrencyAggregation) GetCurrencyId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CurrencyId
}

// GetCurrencyIdOk returns a tuple with the CurrencyId field value
// and a boolean to check if the value has been set.
func (o *CurrencyAggregation) GetCurrencyIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CurrencyId, true
}

// SetCurrencyId sets field value
func (o *CurrencyAggregation) SetCurrencyId(v string) {
	o.CurrencyId = v
}

// GetAccounts returns the Accounts field value
func (o *CurrencyAggregation) GetAccounts() []AccountAggregation {
	if o == nil {
		var ret []AccountAggregation
		return ret
	}

	return o.Accounts
}

// GetAccountsOk returns a tuple with the Accounts field value
// and a boolean to check if the value has been set.
func (o *CurrencyAggregation) GetAccountsOk() ([]AccountAggregation, bool) {
	if o == nil {
		return nil, false
	}
	return o.Accounts, true
}

// SetAccounts sets field value
func (o *CurrencyAggregation) SetAccounts(v []AccountAggregation) {
	o.Accounts = v
}

func (o CurrencyAggregation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CurrencyAggregation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["currencyId"] = o.CurrencyId
	toSerialize["accounts"] = o.Accounts
	return toSerialize, nil
}

func (o *CurrencyAggregation) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"currencyId",
		"accounts",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varCurrencyAggregation := _CurrencyAggregation{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCurrencyAggregation)

	if err != nil {
		return err
	}

	*o = CurrencyAggregation(varCurrencyAggregation)

	return err
}

type NullableCurrencyAggregation struct {
	value *CurrencyAggregation
	isSet bool
}

func (v NullableCurrencyAggregation) Get() *CurrencyAggregation {
	return v.value
}

func (v *NullableCurrencyAggregation) Set(val *CurrencyAggregation) {
	v.value = val
	v.isSet = true
}

func (v NullableCurrencyAggregation) IsSet() bool {
	return v.isSet
}

func (v *NullableCurrencyAggregation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCurrencyAggregation(val *CurrencyAggregation) *NullableCurrencyAggregation {
	return &NullableCurrencyAggregation{value: val, isSet: true}
}

func (v NullableCurrencyAggregation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCurrencyAggregation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
