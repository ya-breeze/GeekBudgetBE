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

// checks if the AccountAggregation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AccountAggregation{}

// AccountAggregation struct for AccountAggregation
type AccountAggregation struct {
	AccountId string    `json:"accountId"`
	Amounts   []float64 `json:"amounts"`
}

type _AccountAggregation AccountAggregation

// NewAccountAggregation instantiates a new AccountAggregation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccountAggregation(accountId string, amounts []float64) *AccountAggregation {
	this := AccountAggregation{}
	this.AccountId = accountId
	this.Amounts = amounts
	return &this
}

// NewAccountAggregationWithDefaults instantiates a new AccountAggregation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccountAggregationWithDefaults() *AccountAggregation {
	this := AccountAggregation{}
	return &this
}

// GetAccountId returns the AccountId field value
func (o *AccountAggregation) GetAccountId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value
// and a boolean to check if the value has been set.
func (o *AccountAggregation) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccountId, true
}

// SetAccountId sets field value
func (o *AccountAggregation) SetAccountId(v string) {
	o.AccountId = v
}

// GetAmounts returns the Amounts field value
func (o *AccountAggregation) GetAmounts() []float64 {
	if o == nil {
		var ret []float64
		return ret
	}

	return o.Amounts
}

// GetAmountsOk returns a tuple with the Amounts field value
// and a boolean to check if the value has been set.
func (o *AccountAggregation) GetAmountsOk() ([]float64, bool) {
	if o == nil {
		return nil, false
	}
	return o.Amounts, true
}

// SetAmounts sets field value
func (o *AccountAggregation) SetAmounts(v []float64) {
	o.Amounts = v
}

func (o AccountAggregation) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AccountAggregation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["accountId"] = o.AccountId
	toSerialize["amounts"] = o.Amounts
	return toSerialize, nil
}

func (o *AccountAggregation) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"accountId",
		"amounts",
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

	varAccountAggregation := _AccountAggregation{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAccountAggregation)

	if err != nil {
		return err
	}

	*o = AccountAggregation(varAccountAggregation)

	return err
}

type NullableAccountAggregation struct {
	value *AccountAggregation
	isSet bool
}

func (v NullableAccountAggregation) Get() *AccountAggregation {
	return v.value
}

func (v *NullableAccountAggregation) Set(val *AccountAggregation) {
	v.value = val
	v.isSet = true
}

func (v NullableAccountAggregation) IsSet() bool {
	return v.isSet
}

func (v *NullableAccountAggregation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccountAggregation(val *AccountAggregation) *NullableAccountAggregation {
	return &NullableAccountAggregation{value: val, isSet: true}
}

func (v NullableAccountAggregation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccountAggregation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
