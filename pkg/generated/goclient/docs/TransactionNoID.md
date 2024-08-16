# TransactionNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Date** | **time.Time** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Place** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**PartnerName** | Pointer to **string** |  | [optional] 
**PartnerAccount** | Pointer to **string** |  | [optional] 
**PartnerInternalId** | Pointer to **string** | Internal bank&#39;s ID to be able to match later if necessary | [optional] 
**Extra** | Pointer to **string** | Stores extra data about transaction. For example could hold \&quot;variable symbol\&quot; to distinguish payment for the same account, but with different meaning | [optional] 
**UnprocessedSources** | Pointer to **string** | Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis | [optional] 
**ExternalIds** | Pointer to **string** | IDs of unprocessed transaction - to match later | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 

## Methods

### NewTransactionNoID

`func NewTransactionNoID(date time.Time, movements []Movement, ) *TransactionNoID`

NewTransactionNoID instantiates a new TransactionNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionNoIDWithDefaults

`func NewTransactionNoIDWithDefaults() *TransactionNoID`

NewTransactionNoIDWithDefaults instantiates a new TransactionNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *TransactionNoID) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *TransactionNoID) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *TransactionNoID) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetDescription

`func (o *TransactionNoID) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *TransactionNoID) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *TransactionNoID) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *TransactionNoID) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPlace

`func (o *TransactionNoID) GetPlace() string`

GetPlace returns the Place field if non-nil, zero value otherwise.

### GetPlaceOk

`func (o *TransactionNoID) GetPlaceOk() (*string, bool)`

GetPlaceOk returns a tuple with the Place field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlace

`func (o *TransactionNoID) SetPlace(v string)`

SetPlace sets Place field to given value.

### HasPlace

`func (o *TransactionNoID) HasPlace() bool`

HasPlace returns a boolean if a field has been set.

### GetTags

`func (o *TransactionNoID) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *TransactionNoID) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *TransactionNoID) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *TransactionNoID) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetPartnerName

`func (o *TransactionNoID) GetPartnerName() string`

GetPartnerName returns the PartnerName field if non-nil, zero value otherwise.

### GetPartnerNameOk

`func (o *TransactionNoID) GetPartnerNameOk() (*string, bool)`

GetPartnerNameOk returns a tuple with the PartnerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerName

`func (o *TransactionNoID) SetPartnerName(v string)`

SetPartnerName sets PartnerName field to given value.

### HasPartnerName

`func (o *TransactionNoID) HasPartnerName() bool`

HasPartnerName returns a boolean if a field has been set.

### GetPartnerAccount

`func (o *TransactionNoID) GetPartnerAccount() string`

GetPartnerAccount returns the PartnerAccount field if non-nil, zero value otherwise.

### GetPartnerAccountOk

`func (o *TransactionNoID) GetPartnerAccountOk() (*string, bool)`

GetPartnerAccountOk returns a tuple with the PartnerAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerAccount

`func (o *TransactionNoID) SetPartnerAccount(v string)`

SetPartnerAccount sets PartnerAccount field to given value.

### HasPartnerAccount

`func (o *TransactionNoID) HasPartnerAccount() bool`

HasPartnerAccount returns a boolean if a field has been set.

### GetPartnerInternalId

`func (o *TransactionNoID) GetPartnerInternalId() string`

GetPartnerInternalId returns the PartnerInternalId field if non-nil, zero value otherwise.

### GetPartnerInternalIdOk

`func (o *TransactionNoID) GetPartnerInternalIdOk() (*string, bool)`

GetPartnerInternalIdOk returns a tuple with the PartnerInternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerInternalId

`func (o *TransactionNoID) SetPartnerInternalId(v string)`

SetPartnerInternalId sets PartnerInternalId field to given value.

### HasPartnerInternalId

`func (o *TransactionNoID) HasPartnerInternalId() bool`

HasPartnerInternalId returns a boolean if a field has been set.

### GetExtra

`func (o *TransactionNoID) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *TransactionNoID) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *TransactionNoID) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *TransactionNoID) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetUnprocessedSources

`func (o *TransactionNoID) GetUnprocessedSources() string`

GetUnprocessedSources returns the UnprocessedSources field if non-nil, zero value otherwise.

### GetUnprocessedSourcesOk

`func (o *TransactionNoID) GetUnprocessedSourcesOk() (*string, bool)`

GetUnprocessedSourcesOk returns a tuple with the UnprocessedSources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnprocessedSources

`func (o *TransactionNoID) SetUnprocessedSources(v string)`

SetUnprocessedSources sets UnprocessedSources field to given value.

### HasUnprocessedSources

`func (o *TransactionNoID) HasUnprocessedSources() bool`

HasUnprocessedSources returns a boolean if a field has been set.

### GetExternalIds

`func (o *TransactionNoID) GetExternalIds() string`

GetExternalIds returns the ExternalIds field if non-nil, zero value otherwise.

### GetExternalIdsOk

`func (o *TransactionNoID) GetExternalIdsOk() (*string, bool)`

GetExternalIdsOk returns a tuple with the ExternalIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalIds

`func (o *TransactionNoID) SetExternalIds(v string)`

SetExternalIds sets ExternalIds field to given value.

### HasExternalIds

`func (o *TransactionNoID) HasExternalIds() bool`

HasExternalIds returns a boolean if a field has been set.

### GetMovements

`func (o *TransactionNoID) GetMovements() []Movement`

GetMovements returns the Movements field if non-nil, zero value otherwise.

### GetMovementsOk

`func (o *TransactionNoID) GetMovementsOk() (*[]Movement, bool)`

GetMovementsOk returns a tuple with the Movements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMovements

`func (o *TransactionNoID) SetMovements(v []Movement)`

SetMovements sets Movements field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


