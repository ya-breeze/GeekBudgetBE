# Movement

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | **float64** |  | 
**CurrencyId** | **string** |  | 
**AccountId** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 

## Methods

### NewMovement

`func NewMovement(amount float64, currencyId string, accountId string, ) *Movement`

NewMovement instantiates a new Movement object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMovementWithDefaults

`func NewMovementWithDefaults() *Movement`

NewMovementWithDefaults instantiates a new Movement object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *Movement) GetAmount() float64`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *Movement) GetAmountOk() (*float64, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *Movement) SetAmount(v float64)`

SetAmount sets Amount field to given value.


### GetCurrencyId

`func (o *Movement) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *Movement) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *Movement) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetAccountId

`func (o *Movement) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *Movement) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *Movement) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetDescription

`func (o *Movement) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Movement) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Movement) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Movement) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


