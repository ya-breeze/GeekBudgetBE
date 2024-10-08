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
	AccountId *string `json:"accountId,omitempty"`
	BankId    *string `json:"bankId,omitempty"`
	// List of balances for this account. It's an array since one account could hold multiple currencies, for example, cash account could hold EUR, USD and CZK. Or one bank account could hold multiple currencies.
	Balances []BankAccountInfoBalancesInner `json:"balances,omitempty"`
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

// GetBalances returns the Balances field value if set, zero value otherwise.
func (o *BankAccountInfo) GetBalances() []BankAccountInfoBalancesInner {
	if o == nil || IsNil(o.Balances) {
		var ret []BankAccountInfoBalancesInner
		return ret
	}
	return o.Balances
}

// GetBalancesOk returns a tuple with the Balances field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankAccountInfo) GetBalancesOk() ([]BankAccountInfoBalancesInner, bool) {
	if o == nil || IsNil(o.Balances) {
		return nil, false
	}
	return o.Balances, true
}

// HasBalances returns a boolean if a field has been set.
func (o *BankAccountInfo) HasBalances() bool {
	if o != nil && !IsNil(o.Balances) {
		return true
	}

	return false
}

// SetBalances gets a reference to the given []BankAccountInfoBalancesInner and assigns it to the Balances field.
func (o *BankAccountInfo) SetBalances(v []BankAccountInfoBalancesInner) {
	o.Balances = v
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
	if !IsNil(o.Balances) {
		toSerialize["balances"] = o.Balances
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
