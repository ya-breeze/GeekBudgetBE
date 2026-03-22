# TransactionTemplate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** | User-given label for the template | 
**Description** | Pointer to **string** |  | [optional] 
**Place** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**PartnerName** | Pointer to **string** |  | [optional] 
**Extra** | Pointer to **string** |  | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 

## Methods

### NewTransactionTemplate

`func NewTransactionTemplate(id string, name string, movements []Movement, ) *TransactionTemplate`

NewTransactionTemplate instantiates a new TransactionTemplate object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionTemplateWithDefaults

`func NewTransactionTemplateWithDefaults() *TransactionTemplate`

NewTransactionTemplateWithDefaults instantiates a new TransactionTemplate object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TransactionTemplate) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TransactionTemplate) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TransactionTemplate) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *TransactionTemplate) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TransactionTemplate) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TransactionTemplate) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *TransactionTemplate) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *TransactionTemplate) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *TransactionTemplate) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *TransactionTemplate) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPlace

`func (o *TransactionTemplate) GetPlace() string`

GetPlace returns the Place field if non-nil, zero value otherwise.

### GetPlaceOk

`func (o *TransactionTemplate) GetPlaceOk() (*string, bool)`

GetPlaceOk returns a tuple with the Place field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlace

`func (o *TransactionTemplate) SetPlace(v string)`

SetPlace sets Place field to given value.

### HasPlace

`func (o *TransactionTemplate) HasPlace() bool`

HasPlace returns a boolean if a field has been set.

### GetTags

`func (o *TransactionTemplate) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *TransactionTemplate) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *TransactionTemplate) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *TransactionTemplate) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetPartnerName

`func (o *TransactionTemplate) GetPartnerName() string`

GetPartnerName returns the PartnerName field if non-nil, zero value otherwise.

### GetPartnerNameOk

`func (o *TransactionTemplate) GetPartnerNameOk() (*string, bool)`

GetPartnerNameOk returns a tuple with the PartnerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerName

`func (o *TransactionTemplate) SetPartnerName(v string)`

SetPartnerName sets PartnerName field to given value.

### HasPartnerName

`func (o *TransactionTemplate) HasPartnerName() bool`

HasPartnerName returns a boolean if a field has been set.

### GetExtra

`func (o *TransactionTemplate) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *TransactionTemplate) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *TransactionTemplate) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *TransactionTemplate) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetMovements

`func (o *TransactionTemplate) GetMovements() []Movement`

GetMovements returns the Movements field if non-nil, zero value otherwise.

### GetMovementsOk

`func (o *TransactionTemplate) GetMovementsOk() (*[]Movement, bool)`

GetMovementsOk returns a tuple with the Movements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMovements

`func (o *TransactionTemplate) SetMovements(v []Movement)`

SetMovements sets Movements field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


