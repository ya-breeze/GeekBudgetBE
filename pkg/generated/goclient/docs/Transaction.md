# Transaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Date** | **time.Time** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Place** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**PartnerName** | Pointer to **string** |  | [optional] 
**PartnerAccount** | Pointer to **string** |  | [optional] 
**PartnerInternalID** | Pointer to **string** | Internal bank&#39;s ID to be able to match later if necessary | [optional] 
**Extra** | Pointer to **string** | Stores extra data about transaction. For example could hold \&quot;variable symbol\&quot; to distinguish payment for the same account, but with different meaning | [optional] 
**UnprocessedSources** | Pointer to **string** | Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis | [optional] 
**ExternalIDs** | Pointer to **string** | IDs of unprocessed transaction - to match later | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 

## Methods

### NewTransaction

`func NewTransaction(id string, date time.Time, movements []Movement, ) *Transaction`

NewTransaction instantiates a new Transaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionWithDefaults

`func NewTransactionWithDefaults() *Transaction`

NewTransactionWithDefaults instantiates a new Transaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Transaction) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Transaction) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Transaction) SetId(v string)`

SetId sets Id field to given value.


### GetDate

`func (o *Transaction) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *Transaction) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *Transaction) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetDescription

`func (o *Transaction) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Transaction) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Transaction) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Transaction) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPlace

`func (o *Transaction) GetPlace() string`

GetPlace returns the Place field if non-nil, zero value otherwise.

### GetPlaceOk

`func (o *Transaction) GetPlaceOk() (*string, bool)`

GetPlaceOk returns a tuple with the Place field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlace

`func (o *Transaction) SetPlace(v string)`

SetPlace sets Place field to given value.

### HasPlace

`func (o *Transaction) HasPlace() bool`

HasPlace returns a boolean if a field has been set.

### GetTags

`func (o *Transaction) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *Transaction) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *Transaction) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *Transaction) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetPartnerName

`func (o *Transaction) GetPartnerName() string`

GetPartnerName returns the PartnerName field if non-nil, zero value otherwise.

### GetPartnerNameOk

`func (o *Transaction) GetPartnerNameOk() (*string, bool)`

GetPartnerNameOk returns a tuple with the PartnerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerName

`func (o *Transaction) SetPartnerName(v string)`

SetPartnerName sets PartnerName field to given value.

### HasPartnerName

`func (o *Transaction) HasPartnerName() bool`

HasPartnerName returns a boolean if a field has been set.

### GetPartnerAccount

`func (o *Transaction) GetPartnerAccount() string`

GetPartnerAccount returns the PartnerAccount field if non-nil, zero value otherwise.

### GetPartnerAccountOk

`func (o *Transaction) GetPartnerAccountOk() (*string, bool)`

GetPartnerAccountOk returns a tuple with the PartnerAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerAccount

`func (o *Transaction) SetPartnerAccount(v string)`

SetPartnerAccount sets PartnerAccount field to given value.

### HasPartnerAccount

`func (o *Transaction) HasPartnerAccount() bool`

HasPartnerAccount returns a boolean if a field has been set.

### GetPartnerInternalID

`func (o *Transaction) GetPartnerInternalID() string`

GetPartnerInternalID returns the PartnerInternalID field if non-nil, zero value otherwise.

### GetPartnerInternalIDOk

`func (o *Transaction) GetPartnerInternalIDOk() (*string, bool)`

GetPartnerInternalIDOk returns a tuple with the PartnerInternalID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerInternalID

`func (o *Transaction) SetPartnerInternalID(v string)`

SetPartnerInternalID sets PartnerInternalID field to given value.

### HasPartnerInternalID

`func (o *Transaction) HasPartnerInternalID() bool`

HasPartnerInternalID returns a boolean if a field has been set.

### GetExtra

`func (o *Transaction) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *Transaction) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *Transaction) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *Transaction) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetUnprocessedSources

`func (o *Transaction) GetUnprocessedSources() string`

GetUnprocessedSources returns the UnprocessedSources field if non-nil, zero value otherwise.

### GetUnprocessedSourcesOk

`func (o *Transaction) GetUnprocessedSourcesOk() (*string, bool)`

GetUnprocessedSourcesOk returns a tuple with the UnprocessedSources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnprocessedSources

`func (o *Transaction) SetUnprocessedSources(v string)`

SetUnprocessedSources sets UnprocessedSources field to given value.

### HasUnprocessedSources

`func (o *Transaction) HasUnprocessedSources() bool`

HasUnprocessedSources returns a boolean if a field has been set.

### GetExternalIDs

`func (o *Transaction) GetExternalIDs() string`

GetExternalIDs returns the ExternalIDs field if non-nil, zero value otherwise.

### GetExternalIDsOk

`func (o *Transaction) GetExternalIDsOk() (*string, bool)`

GetExternalIDsOk returns a tuple with the ExternalIDs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalIDs

`func (o *Transaction) SetExternalIDs(v string)`

SetExternalIDs sets ExternalIDs field to given value.

### HasExternalIDs

`func (o *Transaction) HasExternalIDs() bool`

HasExternalIDs returns a boolean if a field has been set.

### GetMovements

`func (o *Transaction) GetMovements() []Movement`

GetMovements returns the Movements field if non-nil, zero value otherwise.

### GetMovementsOk

`func (o *Transaction) GetMovementsOk() (*[]Movement, bool)`

GetMovementsOk returns a tuple with the Movements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMovements

`func (o *Transaction) SetMovements(v []Movement)`

SetMovements sets Movements field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


