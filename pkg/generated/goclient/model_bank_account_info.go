/*
Geek Budget - OpenAPI 3.0

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.0.1
Contact: ilya.korolev@outlook.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package goclient

import (
	"encoding/json"
)

// checks if the BankAccountInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BankAccountInfo{}

// BankAccountInfo struct for BankAccountInfo
type BankAccountInfo struct {
	AccountId      *string  `json:"accountId,omitempty"`
	BankId         *string  `json:"bankId,omitempty"`
	OpeningBalance *float64 `json:"openingBalance,omitempty"`
	ClosingBalance *float64 `json:"closingBalance,omitempty"`
}

// NewBankAccountInfo instantiates a new BankAccountInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBankAccountInfo() *BankAccountInfo {
	this := BankAccountInfo{}
	return &this
}

// NewBankAccountInfoWithDefaults instantiates a new BankAccountInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBankAccountInfoWithDefaults() *BankAccountInfo {
	this := BankAccountInfo{}
	return &this
}

// GetAccountId returns the AccountId field value if set, zero value otherwise.
func (o *BankAccountInfo) GetAccountId() string {
	if o == nil || IsNil(o.AccountId) {
		var ret string
		return ret
	}
	return *o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankAccountInfo) GetAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.AccountId) {
		return nil, false
	}
	return o.AccountId, true
}

// HasAccountId returns a boolean if a field has been set.
func (o *BankAccountInfo) HasAccountId() bool {
	if o != nil && !IsNil(o.AccountId) {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given string and assigns it to the AccountId field.
func (o *BankAccountInfo) SetAccountId(v string) {
	o.AccountId = &v
}

// GetBankId returns the BankId field value if set, zero value otherwise.
func (o *BankAccountInfo) GetBankId() string {
	if o == nil || IsNil(o.BankId) {
		var ret string
		return ret
	}
	return *o.BankId
}

// GetBankIdOk returns a tuple with the BankId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankAccountInfo) GetBankIdOk() (*string, bool) {
	if o == nil || IsNil(o.BankId) {
		return nil, false
	}
	return o.BankId, true
}

// HasBankId returns a boolean if a field has been set.
func (o *BankAccountInfo) HasBankId() bool {
	if o != nil && !IsNil(o.BankId) {
		return true
	}

	return false
}

// SetBankId gets a reference to the given string and assigns it to the BankId field.
func (o *BankAccountInfo) SetBankId(v string) {
	o.BankId = &v
}

// GetOpeningBalance returns the OpeningBalance field value if set, zero value otherwise.
func (o *BankAccountInfo) GetOpeningBalance() float64 {
	if o == nil || IsNil(o.OpeningBalance) {
		var ret float64
		return ret
	}
	return *o.OpeningBalance
}

// GetOpeningBalanceOk returns a tuple with the OpeningBalance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankAccountInfo) GetOpeningBalanceOk() (*float64, bool) {
	if o == nil || IsNil(o.OpeningBalance) {
		return nil, false
	}
	return o.OpeningBalance, true
}

// HasOpeningBalance returns a boolean if a field has been set.
func (o *BankAccountInfo) HasOpeningBalance() bool {
	if o != nil && !IsNil(o.OpeningBalance) {
		return true
	}

	return false
}

// SetOpeningBalance gets a reference to the given float64 and assigns it to the OpeningBalance field.
func (o *BankAccountInfo) SetOpeningBalance(v float64) {
	o.OpeningBalance = &v
}

// GetClosingBalance returns the ClosingBalance field value if set, zero value otherwise.
func (o *BankAccountInfo) GetClosingBalance() float64 {
	if o == nil || IsNil(o.ClosingBalance) {
		var ret float64
		return ret
	}
	return *o.ClosingBalance
}

// GetClosingBalanceOk returns a tuple with the ClosingBalance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankAccountInfo) GetClosingBalanceOk() (*float64, bool) {
	if o == nil || IsNil(o.ClosingBalance) {
		return nil, false
	}
	return o.ClosingBalance, true
}

// HasClosingBalance returns a boolean if a field has been set.
func (o *BankAccountInfo) HasClosingBalance() bool {
	if o != nil && !IsNil(o.ClosingBalance) {
		return true
	}

	return false
}

// SetClosingBalance gets a reference to the given float64 and assigns it to the ClosingBalance field.
func (o *BankAccountInfo) SetClosingBalance(v float64) {
	o.ClosingBalance = &v
}

func (o BankAccountInfo) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BankAccountInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AccountId) {
		toSerialize["accountId"] = o.AccountId
	}
	if !IsNil(o.BankId) {
		toSerialize["bankId"] = o.BankId
	}
	if !IsNil(o.OpeningBalance) {
		toSerialize["openingBalance"] = o.OpeningBalance
	}
	if !IsNil(o.ClosingBalance) {
		toSerialize["closingBalance"] = o.ClosingBalance
	}
	return toSerialize, nil
}

type NullableBankAccountInfo struct {
	value *BankAccountInfo
	isSet bool
}

func (v NullableBankAccountInfo) Get() *BankAccountInfo {
	return v.value
}

func (v *NullableBankAccountInfo) Set(val *BankAccountInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableBankAccountInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableBankAccountInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBankAccountInfo(val *BankAccountInfo) *NullableBankAccountInfo {
	return &NullableBankAccountInfo{value: val, isSet: true}
}

func (v NullableBankAccountInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBankAccountInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
