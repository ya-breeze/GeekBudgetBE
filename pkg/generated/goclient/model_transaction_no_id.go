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

// checks if the TransactionNoID type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TransactionNoID{}

// TransactionNoID struct for TransactionNoID
type TransactionNoID struct {
	Date           time.Time `json:"date"`
	Description    *string   `json:"description,omitempty"`
	Place          *string   `json:"place,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
	PartnerName    *string   `json:"partnerName,omitempty"`
	PartnerAccount *string   `json:"partnerAccount,omitempty"`
	// Internal bank's ID to be able to match later if necessary
	PartnerInternalID *string `json:"partnerInternalID,omitempty"`
	// Stores extra data about transaction. For example could hold \"variable symbol\" to distinguish payment for the same account, but with different meaning
	Extra *string `json:"extra,omitempty"`
	// Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis
	UnprocessedSources *string `json:"unprocessedSources,omitempty"`
	// IDs of unprocessed transaction - to match later
	ExternalIDs *string    `json:"externalIDs,omitempty"`
	Movements   []Movement `json:"movements"`
}

type _TransactionNoID TransactionNoID

// NewTransactionNoID instantiates a new TransactionNoID object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTransactionNoID(date time.Time, movements []Movement) *TransactionNoID {
	this := TransactionNoID{}
	this.Date = date
	this.Movements = movements
	return &this
}

// NewTransactionNoIDWithDefaults instantiates a new TransactionNoID object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTransactionNoIDWithDefaults() *TransactionNoID {
	this := TransactionNoID{}
	return &this
}

// GetDate returns the Date field value
func (o *TransactionNoID) GetDate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Date
}

// GetDateOk returns a tuple with the Date field value
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Date, true
}

// SetDate sets field value
func (o *TransactionNoID) SetDate(v time.Time) {
	o.Date = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *TransactionNoID) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *TransactionNoID) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *TransactionNoID) SetDescription(v string) {
	o.Description = &v
}

// GetPlace returns the Place field value if set, zero value otherwise.
func (o *TransactionNoID) GetPlace() string {
	if o == nil || IsNil(o.Place) {
		var ret string
		return ret
	}
	return *o.Place
}

// GetPlaceOk returns a tuple with the Place field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetPlaceOk() (*string, bool) {
	if o == nil || IsNil(o.Place) {
		return nil, false
	}
	return o.Place, true
}

// HasPlace returns a boolean if a field has been set.
func (o *TransactionNoID) HasPlace() bool {
	if o != nil && !IsNil(o.Place) {
		return true
	}

	return false
}

// SetPlace gets a reference to the given string and assigns it to the Place field.
func (o *TransactionNoID) SetPlace(v string) {
	o.Place = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *TransactionNoID) GetTags() []string {
	if o == nil || IsNil(o.Tags) {
		var ret []string
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetTagsOk() ([]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *TransactionNoID) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *TransactionNoID) SetTags(v []string) {
	o.Tags = v
}

// GetPartnerName returns the PartnerName field value if set, zero value otherwise.
func (o *TransactionNoID) GetPartnerName() string {
	if o == nil || IsNil(o.PartnerName) {
		var ret string
		return ret
	}
	return *o.PartnerName
}

// GetPartnerNameOk returns a tuple with the PartnerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetPartnerNameOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerName) {
		return nil, false
	}
	return o.PartnerName, true
}

// HasPartnerName returns a boolean if a field has been set.
func (o *TransactionNoID) HasPartnerName() bool {
	if o != nil && !IsNil(o.PartnerName) {
		return true
	}

	return false
}

// SetPartnerName gets a reference to the given string and assigns it to the PartnerName field.
func (o *TransactionNoID) SetPartnerName(v string) {
	o.PartnerName = &v
}

// GetPartnerAccount returns the PartnerAccount field value if set, zero value otherwise.
func (o *TransactionNoID) GetPartnerAccount() string {
	if o == nil || IsNil(o.PartnerAccount) {
		var ret string
		return ret
	}
	return *o.PartnerAccount
}

// GetPartnerAccountOk returns a tuple with the PartnerAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetPartnerAccountOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerAccount) {
		return nil, false
	}
	return o.PartnerAccount, true
}

// HasPartnerAccount returns a boolean if a field has been set.
func (o *TransactionNoID) HasPartnerAccount() bool {
	if o != nil && !IsNil(o.PartnerAccount) {
		return true
	}

	return false
}

// SetPartnerAccount gets a reference to the given string and assigns it to the PartnerAccount field.
func (o *TransactionNoID) SetPartnerAccount(v string) {
	o.PartnerAccount = &v
}

// GetPartnerInternalID returns the PartnerInternalID field value if set, zero value otherwise.
func (o *TransactionNoID) GetPartnerInternalID() string {
	if o == nil || IsNil(o.PartnerInternalID) {
		var ret string
		return ret
	}
	return *o.PartnerInternalID
}

// GetPartnerInternalIDOk returns a tuple with the PartnerInternalID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetPartnerInternalIDOk() (*string, bool) {
	if o == nil || IsNil(o.PartnerInternalID) {
		return nil, false
	}
	return o.PartnerInternalID, true
}

// HasPartnerInternalID returns a boolean if a field has been set.
func (o *TransactionNoID) HasPartnerInternalID() bool {
	if o != nil && !IsNil(o.PartnerInternalID) {
		return true
	}

	return false
}

// SetPartnerInternalID gets a reference to the given string and assigns it to the PartnerInternalID field.
func (o *TransactionNoID) SetPartnerInternalID(v string) {
	o.PartnerInternalID = &v
}

// GetExtra returns the Extra field value if set, zero value otherwise.
func (o *TransactionNoID) GetExtra() string {
	if o == nil || IsNil(o.Extra) {
		var ret string
		return ret
	}
	return *o.Extra
}

// GetExtraOk returns a tuple with the Extra field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetExtraOk() (*string, bool) {
	if o == nil || IsNil(o.Extra) {
		return nil, false
	}
	return o.Extra, true
}

// HasExtra returns a boolean if a field has been set.
func (o *TransactionNoID) HasExtra() bool {
	if o != nil && !IsNil(o.Extra) {
		return true
	}

	return false
}

// SetExtra gets a reference to the given string and assigns it to the Extra field.
func (o *TransactionNoID) SetExtra(v string) {
	o.Extra = &v
}

// GetUnprocessedSources returns the UnprocessedSources field value if set, zero value otherwise.
func (o *TransactionNoID) GetUnprocessedSources() string {
	if o == nil || IsNil(o.UnprocessedSources) {
		var ret string
		return ret
	}
	return *o.UnprocessedSources
}

// GetUnprocessedSourcesOk returns a tuple with the UnprocessedSources field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetUnprocessedSourcesOk() (*string, bool) {
	if o == nil || IsNil(o.UnprocessedSources) {
		return nil, false
	}
	return o.UnprocessedSources, true
}

// HasUnprocessedSources returns a boolean if a field has been set.
func (o *TransactionNoID) HasUnprocessedSources() bool {
	if o != nil && !IsNil(o.UnprocessedSources) {
		return true
	}

	return false
}

// SetUnprocessedSources gets a reference to the given string and assigns it to the UnprocessedSources field.
func (o *TransactionNoID) SetUnprocessedSources(v string) {
	o.UnprocessedSources = &v
}

// GetExternalIDs returns the ExternalIDs field value if set, zero value otherwise.
func (o *TransactionNoID) GetExternalIDs() string {
	if o == nil || IsNil(o.ExternalIDs) {
		var ret string
		return ret
	}
	return *o.ExternalIDs
}

// GetExternalIDsOk returns a tuple with the ExternalIDs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetExternalIDsOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalIDs) {
		return nil, false
	}
	return o.ExternalIDs, true
}

// HasExternalIDs returns a boolean if a field has been set.
func (o *TransactionNoID) HasExternalIDs() bool {
	if o != nil && !IsNil(o.ExternalIDs) {
		return true
	}

	return false
}

// SetExternalIDs gets a reference to the given string and assigns it to the ExternalIDs field.
func (o *TransactionNoID) SetExternalIDs(v string) {
	o.ExternalIDs = &v
}

// GetMovements returns the Movements field value
func (o *TransactionNoID) GetMovements() []Movement {
	if o == nil {
		var ret []Movement
		return ret
	}

	return o.Movements
}

// GetMovementsOk returns a tuple with the Movements field value
// and a boolean to check if the value has been set.
func (o *TransactionNoID) GetMovementsOk() ([]Movement, bool) {
	if o == nil {
		return nil, false
	}
	return o.Movements, true
}

// SetMovements sets field value
func (o *TransactionNoID) SetMovements(v []Movement) {
	o.Movements = v
}

func (o TransactionNoID) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TransactionNoID) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
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
	if !IsNil(o.PartnerInternalID) {
		toSerialize["partnerInternalID"] = o.PartnerInternalID
	}
	if !IsNil(o.Extra) {
		toSerialize["extra"] = o.Extra
	}
	if !IsNil(o.UnprocessedSources) {
		toSerialize["unprocessedSources"] = o.UnprocessedSources
	}
	if !IsNil(o.ExternalIDs) {
		toSerialize["externalIDs"] = o.ExternalIDs
	}
	toSerialize["movements"] = o.Movements
	return toSerialize, nil
}

func (o *TransactionNoID) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
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

	varTransactionNoID := _TransactionNoID{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varTransactionNoID)

	if err != nil {
		return err
	}

	*o = TransactionNoID(varTransactionNoID)

	return err
}

type NullableTransactionNoID struct {
	value *TransactionNoID
	isSet bool
}

func (v NullableTransactionNoID) Get() *TransactionNoID {
	return v.value
}

func (v *NullableTransactionNoID) Set(val *TransactionNoID) {
	v.value = val
	v.isSet = true
}

func (v NullableTransactionNoID) IsSet() bool {
	return v.isSet
}

func (v *NullableTransactionNoID) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTransactionNoID(val *TransactionNoID) *NullableTransactionNoID {
	return &NullableTransactionNoID{value: val, isSet: true}
}

func (v NullableTransactionNoID) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTransactionNoID) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
