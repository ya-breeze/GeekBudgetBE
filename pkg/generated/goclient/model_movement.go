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

// checks if the Movement type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Movement{}

// Movement struct for Movement
type Movement struct {
	Amount      float64 `json:"amount"`
	CurrencyID  string  `json:"currencyID"`
	AccountID   string  `json:"accountID"`
	Description *string `json:"description,omitempty"`
}

type _Movement Movement

// NewMovement instantiates a new Movement object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMovement(amount float64, currencyID string, accountID string) *Movement {
	this := Movement{}
	this.Amount = amount
	this.CurrencyID = currencyID
	this.AccountID = accountID
	return &this
}

// NewMovementWithDefaults instantiates a new Movement object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMovementWithDefaults() *Movement {
	this := Movement{}
	return &this
}

// GetAmount returns the Amount field value
func (o *Movement) GetAmount() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Amount
}

// GetAmountOk returns a tuple with the Amount field value
// and a boolean to check if the value has been set.
func (o *Movement) GetAmountOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Amount, true
}

// SetAmount sets field value
func (o *Movement) SetAmount(v float64) {
	o.Amount = v
}

// GetCurrencyID returns the CurrencyID field value
func (o *Movement) GetCurrencyID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CurrencyID
}

// GetCurrencyIDOk returns a tuple with the CurrencyID field value
// and a boolean to check if the value has been set.
func (o *Movement) GetCurrencyIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CurrencyID, true
}

// SetCurrencyID sets field value
func (o *Movement) SetCurrencyID(v string) {
	o.CurrencyID = v
}

// GetAccountID returns the AccountID field value
func (o *Movement) GetAccountID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccountID
}

// GetAccountIDOk returns a tuple with the AccountID field value
// and a boolean to check if the value has been set.
func (o *Movement) GetAccountIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccountID, true
}

// SetAccountID sets field value
func (o *Movement) SetAccountID(v string) {
	o.AccountID = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Movement) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Movement) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Movement) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Movement) SetDescription(v string) {
	o.Description = &v
}

func (o Movement) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Movement) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["amount"] = o.Amount
	toSerialize["currencyID"] = o.CurrencyID
	toSerialize["accountID"] = o.AccountID
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	return toSerialize, nil
}

func (o *Movement) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"amount",
		"currencyID",
		"accountID",
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

	varMovement := _Movement{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varMovement)

	if err != nil {
		return err
	}

	*o = Movement(varMovement)

	return err
}

type NullableMovement struct {
	value *Movement
	isSet bool
}

func (v NullableMovement) Get() *Movement {
	return v.value
}

func (v *NullableMovement) Set(val *Movement) {
	v.value = val
	v.isSet = true
}

func (v NullableMovement) IsSet() bool {
	return v.isSet
}

func (v *NullableMovement) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMovement(val *Movement) *NullableMovement {
	return &NullableMovement{value: val, isSet: true}
}

func (v NullableMovement) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMovement) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
