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
	"time"
)

// checks if the BankImporterNoID type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BankImporterNoID{}

// BankImporterNoID struct for BankImporterNoID
type BankImporterNoID struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	// ID of account which is used to for movements from this bank importer
	AccountId string `json:"accountId"`
	// ID of account which is used for fee movements from this bank importer
	FeeAccountId *string `json:"feeAccountId,omitempty"`
	// Stores extra data about bank importer. For example could hold \"bank account number\" to be able to distinguish between different bank accounts, or it could hold token for bank API
	Extra *string `json:"extra,omitempty"`
	// Type of bank importer. It's used to distinguish between different banks. For example, FIO bank or KB bank.
	Type *string `json:"type,omitempty"`
	// Date of last successful import.
	LastSuccessfulImport *time.Time `json:"lastSuccessfulImport,omitempty"`
	// List of last imports. It could be shown to user to explain what was imported recently
	LastImports []ImportResult `json:"lastImports,omitempty"`
	// List of mappings which are used to enrich transactions with additional tags
	Mappings []BankImporterNoIDMappingsInner `json:"mappings,omitempty"`
}

type _BankImporterNoID BankImporterNoID

// NewBankImporterNoID instantiates a new BankImporterNoID object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBankImporterNoID(name string, accountId string) *BankImporterNoID {
	this := BankImporterNoID{}
	this.Name = name
	this.AccountId = accountId
	return &this
}

// NewBankImporterNoIDWithDefaults instantiates a new BankImporterNoID object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBankImporterNoIDWithDefaults() *BankImporterNoID {
	this := BankImporterNoID{}
	return &this
}

// GetName returns the Name field value
func (o *BankImporterNoID) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *BankImporterNoID) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *BankImporterNoID) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *BankImporterNoID) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *BankImporterNoID) SetDescription(v string) {
	o.Description = &v
}

// GetAccountId returns the AccountId field value
func (o *BankImporterNoID) GetAccountId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccountId, true
}

// SetAccountId sets field value
func (o *BankImporterNoID) SetAccountId(v string) {
	o.AccountId = v
}

// GetFeeAccountId returns the FeeAccountId field value if set, zero value otherwise.
func (o *BankImporterNoID) GetFeeAccountId() string {
	if o == nil || IsNil(o.FeeAccountId) {
		var ret string
		return ret
	}
	return *o.FeeAccountId
}

// GetFeeAccountIdOk returns a tuple with the FeeAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetFeeAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.FeeAccountId) {
		return nil, false
	}
	return o.FeeAccountId, true
}

// HasFeeAccountId returns a boolean if a field has been set.
func (o *BankImporterNoID) HasFeeAccountId() bool {
	if o != nil && !IsNil(o.FeeAccountId) {
		return true
	}

	return false
}

// SetFeeAccountId gets a reference to the given string and assigns it to the FeeAccountId field.
func (o *BankImporterNoID) SetFeeAccountId(v string) {
	o.FeeAccountId = &v
}

// GetExtra returns the Extra field value if set, zero value otherwise.
func (o *BankImporterNoID) GetExtra() string {
	if o == nil || IsNil(o.Extra) {
		var ret string
		return ret
	}
	return *o.Extra
}

// GetExtraOk returns a tuple with the Extra field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetExtraOk() (*string, bool) {
	if o == nil || IsNil(o.Extra) {
		return nil, false
	}
	return o.Extra, true
}

// HasExtra returns a boolean if a field has been set.
func (o *BankImporterNoID) HasExtra() bool {
	if o != nil && !IsNil(o.Extra) {
		return true
	}

	return false
}

// SetExtra gets a reference to the given string and assigns it to the Extra field.
func (o *BankImporterNoID) SetExtra(v string) {
	o.Extra = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *BankImporterNoID) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *BankImporterNoID) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *BankImporterNoID) SetType(v string) {
	o.Type = &v
}

// GetLastSuccessfulImport returns the LastSuccessfulImport field value if set, zero value otherwise.
func (o *BankImporterNoID) GetLastSuccessfulImport() time.Time {
	if o == nil || IsNil(o.LastSuccessfulImport) {
		var ret time.Time
		return ret
	}
	return *o.LastSuccessfulImport
}

// GetLastSuccessfulImportOk returns a tuple with the LastSuccessfulImport field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetLastSuccessfulImportOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastSuccessfulImport) {
		return nil, false
	}
	return o.LastSuccessfulImport, true
}

// HasLastSuccessfulImport returns a boolean if a field has been set.
func (o *BankImporterNoID) HasLastSuccessfulImport() bool {
	if o != nil && !IsNil(o.LastSuccessfulImport) {
		return true
	}

	return false
}

// SetLastSuccessfulImport gets a reference to the given time.Time and assigns it to the LastSuccessfulImport field.
func (o *BankImporterNoID) SetLastSuccessfulImport(v time.Time) {
	o.LastSuccessfulImport = &v
}

// GetLastImports returns the LastImports field value if set, zero value otherwise.
func (o *BankImporterNoID) GetLastImports() []ImportResult {
	if o == nil || IsNil(o.LastImports) {
		var ret []ImportResult
		return ret
	}
	return o.LastImports
}

// GetLastImportsOk returns a tuple with the LastImports field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetLastImportsOk() ([]ImportResult, bool) {
	if o == nil || IsNil(o.LastImports) {
		return nil, false
	}
	return o.LastImports, true
}

// HasLastImports returns a boolean if a field has been set.
func (o *BankImporterNoID) HasLastImports() bool {
	if o != nil && !IsNil(o.LastImports) {
		return true
	}

	return false
}

// SetLastImports gets a reference to the given []ImportResult and assigns it to the LastImports field.
func (o *BankImporterNoID) SetLastImports(v []ImportResult) {
	o.LastImports = v
}

// GetMappings returns the Mappings field value if set, zero value otherwise.
func (o *BankImporterNoID) GetMappings() []BankImporterNoIDMappingsInner {
	if o == nil || IsNil(o.Mappings) {
		var ret []BankImporterNoIDMappingsInner
		return ret
	}
	return o.Mappings
}

// GetMappingsOk returns a tuple with the Mappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporterNoID) GetMappingsOk() ([]BankImporterNoIDMappingsInner, bool) {
	if o == nil || IsNil(o.Mappings) {
		return nil, false
	}
	return o.Mappings, true
}

// HasMappings returns a boolean if a field has been set.
func (o *BankImporterNoID) HasMappings() bool {
	if o != nil && !IsNil(o.Mappings) {
		return true
	}

	return false
}

// SetMappings gets a reference to the given []BankImporterNoIDMappingsInner and assigns it to the Mappings field.
func (o *BankImporterNoID) SetMappings(v []BankImporterNoIDMappingsInner) {
	o.Mappings = v
}

func (o BankImporterNoID) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BankImporterNoID) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	toSerialize["accountId"] = o.AccountId
	if !IsNil(o.FeeAccountId) {
		toSerialize["feeAccountId"] = o.FeeAccountId
	}
	if !IsNil(o.Extra) {
		toSerialize["extra"] = o.Extra
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.LastSuccessfulImport) {
		toSerialize["lastSuccessfulImport"] = o.LastSuccessfulImport
	}
	if !IsNil(o.LastImports) {
		toSerialize["lastImports"] = o.LastImports
	}
	if !IsNil(o.Mappings) {
		toSerialize["mappings"] = o.Mappings
	}
	return toSerialize, nil
}

func (o *BankImporterNoID) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"accountId",
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

	varBankImporterNoID := _BankImporterNoID{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varBankImporterNoID)

	if err != nil {
		return err
	}

	*o = BankImporterNoID(varBankImporterNoID)

	return err
}

type NullableBankImporterNoID struct {
	value *BankImporterNoID
	isSet bool
}

func (v NullableBankImporterNoID) Get() *BankImporterNoID {
	return v.value
}

func (v *NullableBankImporterNoID) Set(val *BankImporterNoID) {
	v.value = val
	v.isSet = true
}

func (v NullableBankImporterNoID) IsSet() bool {
	return v.isSet
}

func (v *NullableBankImporterNoID) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBankImporterNoID(val *BankImporterNoID) *NullableBankImporterNoID {
	return &NullableBankImporterNoID{value: val, isSet: true}
}

func (v NullableBankImporterNoID) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBankImporterNoID) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
