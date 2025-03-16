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

// checks if the WholeUserData type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WholeUserData{}

// WholeUserData struct for WholeUserData
type WholeUserData struct {
	User          *User          `json:"user,omitempty"`
	Currencies    []Currency     `json:"currencies,omitempty"`
	Accounts      []Account      `json:"accounts,omitempty"`
	Transactions  []Transaction  `json:"transactions,omitempty"`
	Matchers      []Matcher      `json:"matchers,omitempty"`
	BankImporters []BankImporter `json:"bankImporters,omitempty"`
}

// NewWholeUserData instantiates a new WholeUserData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWholeUserData() *WholeUserData {
	this := WholeUserData{}
	return &this
}

// NewWholeUserDataWithDefaults instantiates a new WholeUserData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWholeUserDataWithDefaults() *WholeUserData {
	this := WholeUserData{}
	return &this
}

// GetUser returns the User field value if set, zero value otherwise.
func (o *WholeUserData) GetUser() User {
	if o == nil || IsNil(o.User) {
		var ret User
		return ret
	}
	return *o.User
}

// GetUserOk returns a tuple with the User field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetUserOk() (*User, bool) {
	if o == nil || IsNil(o.User) {
		return nil, false
	}
	return o.User, true
}

// HasUser returns a boolean if a field has been set.
func (o *WholeUserData) HasUser() bool {
	if o != nil && !IsNil(o.User) {
		return true
	}

	return false
}

// SetUser gets a reference to the given User and assigns it to the User field.
func (o *WholeUserData) SetUser(v User) {
	o.User = &v
}

// GetCurrencies returns the Currencies field value if set, zero value otherwise.
func (o *WholeUserData) GetCurrencies() []Currency {
	if o == nil || IsNil(o.Currencies) {
		var ret []Currency
		return ret
	}
	return o.Currencies
}

// GetCurrenciesOk returns a tuple with the Currencies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetCurrenciesOk() ([]Currency, bool) {
	if o == nil || IsNil(o.Currencies) {
		return nil, false
	}
	return o.Currencies, true
}

// HasCurrencies returns a boolean if a field has been set.
func (o *WholeUserData) HasCurrencies() bool {
	if o != nil && !IsNil(o.Currencies) {
		return true
	}

	return false
}

// SetCurrencies gets a reference to the given []Currency and assigns it to the Currencies field.
func (o *WholeUserData) SetCurrencies(v []Currency) {
	o.Currencies = v
}

// GetAccounts returns the Accounts field value if set, zero value otherwise.
func (o *WholeUserData) GetAccounts() []Account {
	if o == nil || IsNil(o.Accounts) {
		var ret []Account
		return ret
	}
	return o.Accounts
}

// GetAccountsOk returns a tuple with the Accounts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetAccountsOk() ([]Account, bool) {
	if o == nil || IsNil(o.Accounts) {
		return nil, false
	}
	return o.Accounts, true
}

// HasAccounts returns a boolean if a field has been set.
func (o *WholeUserData) HasAccounts() bool {
	if o != nil && !IsNil(o.Accounts) {
		return true
	}

	return false
}

// SetAccounts gets a reference to the given []Account and assigns it to the Accounts field.
func (o *WholeUserData) SetAccounts(v []Account) {
	o.Accounts = v
}

// GetTransactions returns the Transactions field value if set, zero value otherwise.
func (o *WholeUserData) GetTransactions() []Transaction {
	if o == nil || IsNil(o.Transactions) {
		var ret []Transaction
		return ret
	}
	return o.Transactions
}

// GetTransactionsOk returns a tuple with the Transactions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetTransactionsOk() ([]Transaction, bool) {
	if o == nil || IsNil(o.Transactions) {
		return nil, false
	}
	return o.Transactions, true
}

// HasTransactions returns a boolean if a field has been set.
func (o *WholeUserData) HasTransactions() bool {
	if o != nil && !IsNil(o.Transactions) {
		return true
	}

	return false
}

// SetTransactions gets a reference to the given []Transaction and assigns it to the Transactions field.
func (o *WholeUserData) SetTransactions(v []Transaction) {
	o.Transactions = v
}

// GetMatchers returns the Matchers field value if set, zero value otherwise.
func (o *WholeUserData) GetMatchers() []Matcher {
	if o == nil || IsNil(o.Matchers) {
		var ret []Matcher
		return ret
	}
	return o.Matchers
}

// GetMatchersOk returns a tuple with the Matchers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetMatchersOk() ([]Matcher, bool) {
	if o == nil || IsNil(o.Matchers) {
		return nil, false
	}
	return o.Matchers, true
}

// HasMatchers returns a boolean if a field has been set.
func (o *WholeUserData) HasMatchers() bool {
	if o != nil && !IsNil(o.Matchers) {
		return true
	}

	return false
}

// SetMatchers gets a reference to the given []Matcher and assigns it to the Matchers field.
func (o *WholeUserData) SetMatchers(v []Matcher) {
	o.Matchers = v
}

// GetBankImporters returns the BankImporters field value if set, zero value otherwise.
func (o *WholeUserData) GetBankImporters() []BankImporter {
	if o == nil || IsNil(o.BankImporters) {
		var ret []BankImporter
		return ret
	}
	return o.BankImporters
}

// GetBankImportersOk returns a tuple with the BankImporters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WholeUserData) GetBankImportersOk() ([]BankImporter, bool) {
	if o == nil || IsNil(o.BankImporters) {
		return nil, false
	}
	return o.BankImporters, true
}

// HasBankImporters returns a boolean if a field has been set.
func (o *WholeUserData) HasBankImporters() bool {
	if o != nil && !IsNil(o.BankImporters) {
		return true
	}

	return false
}

// SetBankImporters gets a reference to the given []BankImporter and assigns it to the BankImporters field.
func (o *WholeUserData) SetBankImporters(v []BankImporter) {
	o.BankImporters = v
}

func (o WholeUserData) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WholeUserData) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.User) {
		toSerialize["user"] = o.User
	}
	if !IsNil(o.Currencies) {
		toSerialize["currencies"] = o.Currencies
	}
	if !IsNil(o.Accounts) {
		toSerialize["accounts"] = o.Accounts
	}
	if !IsNil(o.Transactions) {
		toSerialize["transactions"] = o.Transactions
	}
	if !IsNil(o.Matchers) {
		toSerialize["matchers"] = o.Matchers
	}
	if !IsNil(o.BankImporters) {
		toSerialize["bankImporters"] = o.BankImporters
	}
	return toSerialize, nil
}

type NullableWholeUserData struct {
	value *WholeUserData
	isSet bool
}

func (v NullableWholeUserData) Get() *WholeUserData {
	return v.value
}

func (v *NullableWholeUserData) Set(val *WholeUserData) {
	v.value = val
	v.isSet = true
}

func (v NullableWholeUserData) IsSet() bool {
	return v.isSet
}

func (v *NullableWholeUserData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWholeUserData(val *WholeUserData) *NullableWholeUserData {
	return &NullableWholeUserData{value: val, isSet: true}
}

func (v NullableWholeUserData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWholeUserData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
