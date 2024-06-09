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

// checks if the CheckMatcherRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CheckMatcherRequest{}

// CheckMatcherRequest struct for CheckMatcherRequest
type CheckMatcherRequest struct {
	Matcher     MatcherNoID     `json:"matcher"`
	Transaction TransactionNoID `json:"transaction"`
}

type _CheckMatcherRequest CheckMatcherRequest

// NewCheckMatcherRequest instantiates a new CheckMatcherRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCheckMatcherRequest(matcher MatcherNoID, transaction TransactionNoID) *CheckMatcherRequest {
	this := CheckMatcherRequest{}
	this.Matcher = matcher
	this.Transaction = transaction
	return &this
}

// NewCheckMatcherRequestWithDefaults instantiates a new CheckMatcherRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCheckMatcherRequestWithDefaults() *CheckMatcherRequest {
	this := CheckMatcherRequest{}
	return &this
}

// GetMatcher returns the Matcher field value
func (o *CheckMatcherRequest) GetMatcher() MatcherNoID {
	if o == nil {
		var ret MatcherNoID
		return ret
	}

	return o.Matcher
}

// GetMatcherOk returns a tuple with the Matcher field value
// and a boolean to check if the value has been set.
func (o *CheckMatcherRequest) GetMatcherOk() (*MatcherNoID, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Matcher, true
}

// SetMatcher sets field value
func (o *CheckMatcherRequest) SetMatcher(v MatcherNoID) {
	o.Matcher = v
}

// GetTransaction returns the Transaction field value
func (o *CheckMatcherRequest) GetTransaction() TransactionNoID {
	if o == nil {
		var ret TransactionNoID
		return ret
	}

	return o.Transaction
}

// GetTransactionOk returns a tuple with the Transaction field value
// and a boolean to check if the value has been set.
func (o *CheckMatcherRequest) GetTransactionOk() (*TransactionNoID, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Transaction, true
}

// SetTransaction sets field value
func (o *CheckMatcherRequest) SetTransaction(v TransactionNoID) {
	o.Transaction = v
}

func (o CheckMatcherRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CheckMatcherRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["matcher"] = o.Matcher
	toSerialize["transaction"] = o.Transaction
	return toSerialize, nil
}

func (o *CheckMatcherRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"matcher",
		"transaction",
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

	varCheckMatcherRequest := _CheckMatcherRequest{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCheckMatcherRequest)

	if err != nil {
		return err
	}

	*o = CheckMatcherRequest(varCheckMatcherRequest)

	return err
}

type NullableCheckMatcherRequest struct {
	value *CheckMatcherRequest
	isSet bool
}

func (v NullableCheckMatcherRequest) Get() *CheckMatcherRequest {
	return v.value
}

func (v *NullableCheckMatcherRequest) Set(val *CheckMatcherRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCheckMatcherRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCheckMatcherRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCheckMatcherRequest(val *CheckMatcherRequest) *NullableCheckMatcherRequest {
	return &NullableCheckMatcherRequest{value: val, isSet: true}
}

func (v NullableCheckMatcherRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCheckMatcherRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
