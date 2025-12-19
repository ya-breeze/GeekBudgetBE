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
**PartnerInternalId** | Pointer to **string** | Internal bank&#39;s ID to be able to match later if necessary | [optional] 
**Extra** | Pointer to **string** | Stores extra data about transaction. For example could hold \&quot;variable symbol\&quot; to distinguish payment for the same account, but with different meaning | [optional] 
**UnprocessedSources** | Pointer to **string** | Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis | [optional] 
**ExternalIds** | Pointer to **[]string** | IDs of unprocessed transaction - to match later | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 
**MatcherId** | Pointer to **string** | ID of the matcher used for this conversion (if any) | [optional] 
**IsAuto** | Pointer to **bool** | If true, this transaction was converted automatically by the matcher | [optional] 

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

### GetPartnerInternalId

`func (o *Transaction) GetPartnerInternalId() string`

GetPartnerInternalId returns the PartnerInternalId field if non-nil, zero value otherwise.

### GetPartnerInternalIdOk

`func (o *Transaction) GetPartnerInternalIdOk() (*string, bool)`

GetPartnerInternalIdOk returns a tuple with the PartnerInternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerInternalId

`func (o *Transaction) SetPartnerInternalId(v string)`

SetPartnerInternalId sets PartnerInternalId field to given value.

### HasPartnerInternalId

`func (o *Transaction) HasPartnerInternalId() bool`

HasPartnerInternalId returns a boolean if a field has been set.

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

### GetExternalIds

`func (o *Transaction) GetExternalIds() []string`

GetExternalIds returns the ExternalIds field if non-nil, zero value otherwise.

### GetExternalIdsOk

`func (o *Transaction) GetExternalIdsOk() (*[]string, bool)`

GetExternalIdsOk returns a tuple with the ExternalIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalIds

`func (o *Transaction) SetExternalIds(v []string)`

SetExternalIds sets ExternalIds field to given value.

### HasExternalIds

`func (o *Transaction) HasExternalIds() bool`

HasExternalIds returns a boolean if a field has been set.

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


### GetMatcherId

`func (o *Transaction) GetMatcherId() string`

GetMatcherId returns the MatcherId field if non-nil, zero value otherwise.

### GetMatcherIdOk

`func (o *Transaction) GetMatcherIdOk() (*string, bool)`

GetMatcherIdOk returns a tuple with the MatcherId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatcherId

`func (o *Transaction) SetMatcherId(v string)`

SetMatcherId sets MatcherId field to given value.

### HasMatcherId

`func (o *Transaction) HasMatcherId() bool`

HasMatcherId returns a boolean if a field has been set.

### GetIsAuto

`func (o *Transaction) GetIsAuto() bool`

GetIsAuto returns the IsAuto field if non-nil, zero value otherwise.

### GetIsAutoOk

`func (o *Transaction) GetIsAutoOk() (*bool, bool)`

GetIsAutoOk returns a tuple with the IsAuto field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsAuto

`func (o *Transaction) SetIsAuto(v bool)`

SetIsAuto sets IsAuto field to given value.

### HasIsAuto

`func (o *Transaction) HasIsAuto() bool`

HasIsAuto returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


