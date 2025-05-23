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

// checks if the BankImporter type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BankImporter{}

// BankImporter struct for BankImporter
type BankImporter struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	// ID of account which is used to for movements from this bank importer
	AccountId string `json:"accountId"`
	// ID of account which is used for fee movements from this bank importer
	FeeAccountId *string `json:"feeAccountId,omitempty"`
	// Stores extra data about bank importer. For example could hold \"bank account number\" to be able to distinguish between different bank accounts, or it could hold token for bank API
	Extra *string `json:"extra,omitempty"`
	// If true, importer will fetch all transactions from the bank, if false, it will fetch only recent transactions
	FetchAll *bool `json:"fetchAll,omitempty"`
	// Type of bank importer. It's used to distinguish between different banks. For example, FIO bank or KB bank.
	Type *string `json:"type,omitempty"`
	// Date of last successful import.
	LastSuccessfulImport *time.Time `json:"lastSuccessfulImport,omitempty"`
	// List of last imports. It could be shown to user to explain what was imported recently
	LastImports []ImportResult `json:"lastImports,omitempty"`
	// List of mappings which are used to enrich transactions with additional tags
	Mappings []BankImporterNoIDMappingsInner `json:"mappings,omitempty"`
}

type _BankImporter BankImporter

// NewBankImporter instantiates a new BankImporter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBankImporter(id string, name string, accountId string) *BankImporter {
	this := BankImporter{}
	this.Id = id
	this.Name = name
	this.AccountId = accountId
	return &this
}

// NewBankImporterWithDefaults instantiates a new BankImporter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBankImporterWithDefaults() *BankImporter {
	this := BankImporter{}
	return &this
}

// GetId returns the Id field value
func (o *BankImporter) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *BankImporter) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *BankImporter) SetId(v string) {
	o.Id = v
}

// GetName returns the Name field value
func (o *BankImporter) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *BankImporter) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *BankImporter) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *BankImporter) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *BankImporter) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *BankImporter) SetDescription(v string) {
	o.Description = &v
}

// GetAccountId returns the AccountId field value
func (o *BankImporter) GetAccountId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value
// and a boolean to check if the value has been set.
func (o *BankImporter) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccountId, true
}

// SetAccountId sets field value
func (o *BankImporter) SetAccountId(v string) {
	o.AccountId = v
}

// GetFeeAccountId returns the FeeAccountId field value if set, zero value otherwise.
func (o *BankImporter) GetFeeAccountId() string {
	if o == nil || IsNil(o.FeeAccountId) {
		var ret string
		return ret
	}
	return *o.FeeAccountId
}

// GetFeeAccountIdOk returns a tuple with the FeeAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetFeeAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.FeeAccountId) {
		return nil, false
	}
	return o.FeeAccountId, true
}

// HasFeeAccountId returns a boolean if a field has been set.
func (o *BankImporter) HasFeeAccountId() bool {
	if o != nil && !IsNil(o.FeeAccountId) {
		return true
	}

	return false
}

// SetFeeAccountId gets a reference to the given string and assigns it to the FeeAccountId field.
func (o *BankImporter) SetFeeAccountId(v string) {
	o.FeeAccountId = &v
}

// GetExtra returns the Extra field value if set, zero value otherwise.
func (o *BankImporter) GetExtra() string {
	if o == nil || IsNil(o.Extra) {
		var ret string
		return ret
	}
	return *o.Extra
}

// GetExtraOk returns a tuple with the Extra field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetExtraOk() (*string, bool) {
	if o == nil || IsNil(o.Extra) {
		return nil, false
	}
	return o.Extra, true
}

// HasExtra returns a boolean if a field has been set.
func (o *BankImporter) HasExtra() bool {
	if o != nil && !IsNil(o.Extra) {
		return true
	}

	return false
}

// SetExtra gets a reference to the given string and assigns it to the Extra field.
func (o *BankImporter) SetExtra(v string) {
	o.Extra = &v
}

// GetFetchAll returns the FetchAll field value if set, zero value otherwise.
func (o *BankImporter) GetFetchAll() bool {
	if o == nil || IsNil(o.FetchAll) {
		var ret bool
		return ret
	}
	return *o.FetchAll
}

// GetFetchAllOk returns a tuple with the FetchAll field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetFetchAllOk() (*bool, bool) {
	if o == nil || IsNil(o.FetchAll) {
		return nil, false
	}
	return o.FetchAll, true
}

// HasFetchAll returns a boolean if a field has been set.
func (o *BankImporter) HasFetchAll() bool {
	if o != nil && !IsNil(o.FetchAll) {
		return true
	}

	return false
}

// SetFetchAll gets a reference to the given bool and assigns it to the FetchAll field.
func (o *BankImporter) SetFetchAll(v bool) {
	o.FetchAll = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *BankImporter) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *BankImporter) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *BankImporter) SetType(v string) {
	o.Type = &v
}

// GetLastSuccessfulImport returns the LastSuccessfulImport field value if set, zero value otherwise.
func (o *BankImporter) GetLastSuccessfulImport() time.Time {
	if o == nil || IsNil(o.LastSuccessfulImport) {
		var ret time.Time
		return ret
	}
	return *o.LastSuccessfulImport
}

// GetLastSuccessfulImportOk returns a tuple with the LastSuccessfulImport field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetLastSuccessfulImportOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastSuccessfulImport) {
		return nil, false
	}
	return o.LastSuccessfulImport, true
}

// HasLastSuccessfulImport returns a boolean if a field has been set.
func (o *BankImporter) HasLastSuccessfulImport() bool {
	if o != nil && !IsNil(o.LastSuccessfulImport) {
		return true
	}

	return false
}

// SetLastSuccessfulImport gets a reference to the given time.Time and assigns it to the LastSuccessfulImport field.
func (o *BankImporter) SetLastSuccessfulImport(v time.Time) {
	o.LastSuccessfulImport = &v
}

// GetLastImports returns the LastImports field value if set, zero value otherwise.
func (o *BankImporter) GetLastImports() []ImportResult {
	if o == nil || IsNil(o.LastImports) {
		var ret []ImportResult
		return ret
	}
	return o.LastImports
}

// GetLastImportsOk returns a tuple with the LastImports field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetLastImportsOk() ([]ImportResult, bool) {
	if o == nil || IsNil(o.LastImports) {
		return nil, false
	}
	return o.LastImports, true
}

// HasLastImports returns a boolean if a field has been set.
func (o *BankImporter) HasLastImports() bool {
	if o != nil && !IsNil(o.LastImports) {
		return true
	}

	return false
}

// SetLastImports gets a reference to the given []ImportResult and assigns it to the LastImports field.
func (o *BankImporter) SetLastImports(v []ImportResult) {
	o.LastImports = v
}

// GetMappings returns the Mappings field value if set, zero value otherwise.
func (o *BankImporter) GetMappings() []BankImporterNoIDMappingsInner {
	if o == nil || IsNil(o.Mappings) {
		var ret []BankImporterNoIDMappingsInner
		return ret
	}
	return o.Mappings
}

// GetMappingsOk returns a tuple with the Mappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BankImporter) GetMappingsOk() ([]BankImporterNoIDMappingsInner, bool) {
	if o == nil || IsNil(o.Mappings) {
		return nil, false
	}
	return o.Mappings, true
}

// HasMappings returns a boolean if a field has been set.
func (o *BankImporter) HasMappings() bool {
	if o != nil && !IsNil(o.Mappings) {
		return true
	}

	return false
}

// SetMappings gets a reference to the given []BankImporterNoIDMappingsInner and assigns it to the Mappings field.
func (o *BankImporter) SetMappings(v []BankImporterNoIDMappingsInner) {
	o.Mappings = v
}

func (o BankImporter) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BankImporter) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
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
	if !IsNil(o.FetchAll) {
		toSerialize["fetchAll"] = o.FetchAll
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

func (o *BankImporter) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
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

	varBankImporter := _BankImporter{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varBankImporter)

	if err != nil {
		return err
	}

	*o = BankImporter(varBankImporter)

	return err
}

type NullableBankImporter struct {
	value *BankImporter
	isSet bool
}

func (v NullableBankImporter) Get() *BankImporter {
	return v.value
}

func (v *NullableBankImporter) Set(val *BankImporter) {
	v.value = val
	v.isSet = true
}

func (v NullableBankImporter) IsSet() bool {
	return v.isSet
}

func (v *NullableBankImporter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBankImporter(val *BankImporter) *NullableBankImporter {
	return &NullableBankImporter{value: val, isSet: true}
}

func (v NullableBankImporter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBankImporter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
