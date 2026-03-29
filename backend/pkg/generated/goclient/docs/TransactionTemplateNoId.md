# TransactionTemplateNoId

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | User-given label for the template | 
**Description** | Pointer to **string** |  | [optional] 
**Place** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**PartnerName** | Pointer to **string** |  | [optional] 
**Extra** | Pointer to **string** |  | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 

## Methods

### NewTransactionTemplateNoId

`func NewTransactionTemplateNoId(name string, movements []Movement, ) *TransactionTemplateNoId`

NewTransactionTemplateNoId instantiates a new TransactionTemplateNoId object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionTemplateNoIdWithDefaults

`func NewTransactionTemplateNoIdWithDefaults() *TransactionTemplateNoId`

NewTransactionTemplateNoIdWithDefaults instantiates a new TransactionTemplateNoId object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *TransactionTemplateNoId) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TransactionTemplateNoId) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TransactionTemplateNoId) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *TransactionTemplateNoId) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *TransactionTemplateNoId) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *TransactionTemplateNoId) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *TransactionTemplateNoId) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPlace

`func (o *TransactionTemplateNoId) GetPlace() string`

GetPlace returns the Place field if non-nil, zero value otherwise.

### GetPlaceOk

`func (o *TransactionTemplateNoId) GetPlaceOk() (*string, bool)`

GetPlaceOk returns a tuple with the Place field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlace

`func (o *TransactionTemplateNoId) SetPlace(v string)`

SetPlace sets Place field to given value.

### HasPlace

`func (o *TransactionTemplateNoId) HasPlace() bool`

HasPlace returns a boolean if a field has been set.

### GetTags

`func (o *TransactionTemplateNoId) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *TransactionTemplateNoId) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *TransactionTemplateNoId) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *TransactionTemplateNoId) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetPartnerName

`func (o *TransactionTemplateNoId) GetPartnerName() string`

GetPartnerName returns the PartnerName field if non-nil, zero value otherwise.

### GetPartnerNameOk

`func (o *TransactionTemplateNoId) GetPartnerNameOk() (*string, bool)`

GetPartnerNameOk returns a tuple with the PartnerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerName

`func (o *TransactionTemplateNoId) SetPartnerName(v string)`

SetPartnerName sets PartnerName field to given value.

### HasPartnerName

`func (o *TransactionTemplateNoId) HasPartnerName() bool`

HasPartnerName returns a boolean if a field has been set.

### GetExtra

`func (o *TransactionTemplateNoId) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *TransactionTemplateNoId) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *TransactionTemplateNoId) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *TransactionTemplateNoId) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetMovements

`func (o *TransactionTemplateNoId) GetMovements() []Movement`

GetMovements returns the Movements field if non-nil, zero value otherwise.

### GetMovementsOk

`func (o *TransactionTemplateNoId) GetMovementsOk() (*[]Movement, bool)`

GetMovementsOk returns a tuple with the Movements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMovements

`func (o *TransactionTemplateNoId) SetMovements(v []Movement)`

SetMovements sets Movements field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


