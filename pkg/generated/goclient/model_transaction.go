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

// checks if the Transaction type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Transaction{}

// Transaction struct for Transaction
type Transaction struct {
	Id             string    `json:"id"`
	Date           time.Time `json:"date"`
	Description    *string   `json:"description,omitempty"`
	Place          *string   `json:"place,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
	PartnerName    *string   `json:"partnerName,omitempty"`
	PartnerAccount *string   `json:"partnerAccount,omitempty"`
	// Internal bank's ID to be able to match later if necessary
	PartnerInternalId *string `json:"partnerInternalId,omitempty"`
	// Stores extra data about transaction. For example could hold \"variable symbol\" to distinguish payment for the same account, but with different meaning
	Extra *string `json:"extra,omitempty"`
	// Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis
	UnprocessedSources *string `json:"unprocessedSources,omitempty"`
	// IDs of unprocessed transaction - to match later
	ExternalIds *string    `json:"externalIds,omitempty"`
	Movements   []Movement `json:"movements"`
}

type _Transaction Transaction

// NewTransaction instantiates a new Transaction object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTransaction(id string, date time.Time, movements []Movement) *Transaction {
	this := Transaction{}
	this.Id = id
	this.Date = date
	this.Movements = movements
	return &this
}

// NewTransactionWithDefaults instantiates a new Transaction object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTransactionWithDefaults() *Transaction {
	this := Transaction{}
	return &this
}

// GetId returns the Id field value
func (o *Transaction) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Transaction) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Transaction) SetId(v string) {
	o.Id = v
}

// GetDate returns the Date field value
func (o *Transaction) GetDate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Date
}

// GetDateOk returns a tuple with the Date field value
// and a boolean to check if the value has been set.
func (o *Transaction) GetDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Date, true
}

// SetDate sets field value
func (o *Transaction) SetDate(v time.Time) {
	o.Date = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Transaction) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Transaction) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Transaction) SetDescription(v string) {
	o.Description = &v
}

// GetPlace returns the Place field value if set, zero value otherwise.
func (o *Transaction) GetPlace() string {
	if o == nil || IsNil(o.Place) {
		var ret string
		return ret
	}
	return *o.Place
}

// GetPlaceOk returns a tuple with the Place field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetPlaceOk() (*string, bool) {
	if o == nil || IsNil(o.Place) {
		return nil, false
	}
	return o.Place, true
}

// HasPlace returns a boolean if a field has been set.
func (o *Transaction) HasPlace() bool {
	if o != nil && !IsNil(o.Place) {
		return true
	}

	return false
}

// SetPlace gets a reference to the given string and assigns it to the Place field.
func (o *Transaction) SetPlace(v string) {
	o.Place = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *Transaction) GetTags() []string {
	if o == nil || IsNil(o.Tags) {
		var ret []string
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetTagsOk() ([]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *Transaction) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *Transaction) SetTags(v []string) {
	o.Tags = v
}

// GetPartnerName returns the PartnerName field value if set, zero value otherwise.
func (o *Transaction) GetPartnerName() string {
	if o == nil || IsNil(o.PartnerName) {
		var ret string
		return ret
	}
	return *o.PartnerName
}

// GetPartnerNameOk returns a tuple with the PartnerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetPartnerNameOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerName) {
		return nil, false
	}
	return o.PartnerName, true
}

// HasPartnerName returns a boolean if a field has been set.
func (o *Transaction) HasPartnerName() bool {
	if o != nil && !IsNil(o.PartnerName) {
		return true
	}

	return false
}

// SetPartnerName gets a reference to the given string and assigns it to the PartnerName field.
func (o *Transaction) SetPartnerName(v string) {
	o.PartnerName = &v
}

// GetPartnerAccount returns the PartnerAccount field value if set, zero value otherwise.
func (o *Transaction) GetPartnerAccount() string {
	if o == nil || IsNil(o.PartnerAccount) {
		var ret string
		return ret
	}
	return *o.PartnerAccount
}

// GetPartnerAccountOk returns a tuple with the PartnerAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetPartnerAccountOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerAccount) {
		return nil, false
	}
	return o.PartnerAccount, true
}

// HasPartnerAccount returns a boolean if a field has been set.
func (o *Transaction) HasPartnerAccount() bool {
	if o != nil && !IsNil(o.PartnerAccount) {
		return true
	}

	return false
}

// SetPartnerAccount gets a reference to the given string and assigns it to the PartnerAccount field.
func (o *Transaction) SetPartnerAccount(v string) {
	o.PartnerAccount = &v
}

// GetPartnerInternalId returns the PartnerInternalId field value if set, zero value otherwise.
func (o *Transaction) GetPartnerInternalId() string {
	if o == nil || IsNil(o.PartnerInternalId) {
		var ret string
		return ret
	}
	return *o.PartnerInternalId
}

// GetPartnerInternalIdOk returns a tuple with the PartnerInternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetPartnerInternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerInternalId) {
		return nil, false
	}
	return o.PartnerInternalId, true
}

// HasPartnerInternalId returns a boolean if a field has been set.
func (o *Transaction) HasPartnerInternalId() bool {
	if o != nil && !IsNil(o.PartnerInternalId) {
		return true
	}

	return false
}

// SetPartnerInternalId gets a reference to the given string and assigns it to the PartnerInternalId field.
func (o *Transaction) SetPartnerInternalId(v string) {
	o.PartnerInternalId = &v
}

// GetExtra returns the Extra field value if set, zero value otherwise.
func (o *Transaction) GetExtra() string {
	if o == nil || IsNil(o.Extra) {
		var ret string
		return ret
	}
	return *o.Extra
}

// GetExtraOk returns a tuple with the Extra field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetExtraOk() (*string, bool) {
	if o == nil || IsNil(o.Extra) {
		return nil, false
	}
	return o.Extra, true
}

// HasExtra returns a boolean if a field has been set.
func (o *Transaction) HasExtra() bool {
	if o != nil && !IsNil(o.Extra) {
		return true
	}

	return false
}

// SetExtra gets a reference to the given string and assigns it to the Extra field.
func (o *Transaction) SetExtra(v string) {
	o.Extra = &v
}

// GetUnprocessedSources returns the UnprocessedSources field value if set, zero value otherwise.
func (o *Transaction) GetUnprocessedSources() string {
	if o == nil || IsNil(o.UnprocessedSources) {
		var ret string
		return ret
	}
	return *o.UnprocessedSources
}

// GetUnprocessedSourcesOk returns a tuple with the UnprocessedSources field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetUnprocessedSourcesOk() (*string, bool) {
	if o == nil || IsNil(o.UnprocessedSources) {
		return nil, false
	}
	return o.UnprocessedSources, true
}

// HasUnprocessedSources returns a boolean if a field has been set.
func (o *Transaction) HasUnprocessedSources() bool {
	if o != nil && !IsNil(o.UnprocessedSources) {
		return true
	}

	return false
}

// SetUnprocessedSources gets a reference to the given string and assigns it to the UnprocessedSources field.
func (o *Transaction) SetUnprocessedSources(v string) {
	o.UnprocessedSources = &v
}

// GetExternalIds returns the ExternalIds field value if set, zero value otherwise.
func (o *Transaction) GetExternalIds() string {
	if o == nil || IsNil(o.ExternalIds) {
		var ret string
		return ret
	}
	return *o.ExternalIds
}

// GetExternalIdsOk returns a tuple with the ExternalIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetExternalIdsOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalIds) {
		return nil, false
	}
	return o.ExternalIds, true
}

// HasExternalIds returns a boolean if a field has been set.
func (o *Transaction) HasExternalIds() bool {
	if o != nil && !IsNil(o.ExternalIds) {
		return true
	}

	return false
}

// SetExternalIds gets a reference to the given string and assigns it to the ExternalIds field.
func (o *Transaction) SetExternalIds(v string) {
	o.ExternalIds = &v
}

// GetMovements returns the Movements field value
func (o *Transaction) GetMovements() []Movement {
	if o == nil {
		var ret []Movement
		return ret
	}

	return o.Movements
}

// GetMovementsOk returns a tuple with the Movements field value
// and a boolean to check if the value has been set.
func (o *Transaction) GetMovementsOk() ([]Movement, bool) {
	if o == nil {
		return nil, false
	}
	return o.Movements, true
}

// SetMovements sets field value
func (o *Transaction) SetMovements(v []Movement) {
	o.Movements = v
}

func (o Transaction) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Transaction) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["date"] = o.Date
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Place) {
		toSerialize["place"] = o.Place
	}
	if !IsNil(o.Tags) {
		toSerialize["tags"] = o.Tags
	}
	if !IsNil(o.PartnerName) {
		toSerialize["partnerName"] = o.PartnerName
	}
	if !IsNil(o.PartnerAccount) {
		toSerialize["partnerAccount"] = o.PartnerAccount
	}
	if !IsNil(o.PartnerInternalId) {
		toSerialize["partnerInternalId"] = o.PartnerInternalId
	}
	if !IsNil(o.Extra) {
		toSerialize["extra"] = o.Extra
	}
	if !IsNil(o.UnprocessedSources) {
		toSerialize["unprocessedSources"] = o.UnprocessedSources
	}
	if !IsNil(o.ExternalIds) {
		toSerialize["externalIds"] = o.ExternalIds
	}
	toSerialize["movements"] = o.Movements
	return toSerialize, nil
}

func (o *Transaction) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"date",
		"movements",
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

	varTransaction := _Transaction{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varTransaction)

	if err != nil {
		return err
	}

	*o = Transaction(varTransaction)

	return err
}

type NullableTransaction struct {
	value *Transaction
	isSet bool
}

func (v NullableTransaction) Get() *Transaction {
	return v.value
}

func (v *NullableTransaction) Set(val *Transaction) {
	v.value = val
	v.isSet = true
}

func (v NullableTransaction) IsSet() bool {
	return v.isSet
}

func (v *NullableTransaction) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTransaction(val *Transaction) *NullableTransaction {
	return &NullableTransaction{value: val, isSet: true}
}

func (v NullableTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTransaction) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
