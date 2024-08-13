# Movement

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Amount** | **float64** |  | 
**CurrencyID** | **string** |  | 
**AccountID** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 

## Methods

### NewMovement

`func NewMovement(amount float64, currencyID string, accountID string, ) *Movement`

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


### GetCurrencyID

`func (o *Movement) GetCurrencyID() string`

GetCurrencyID returns the CurrencyID field if non-nil, zero value otherwise.

### GetCurrencyIDOk

`func (o *Movement) GetCurrencyIDOk() (*string, bool)`

GetCurrencyIDOk returns a tuple with the CurrencyID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyID

`func (o *Movement) SetCurrencyID(v string)`

SetCurrencyID sets CurrencyID field to given value.


### GetAccountID

`func (o *Movement) GetAccountID() string`

GetAccountID returns the AccountID field if non-nil, zero value otherwise.

### GetAccountIDOk

`func (o *Movement) GetAccountIDOk() (*string, bool)`

GetAccountIDOk returns a tuple with the AccountID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountID

`func (o *Movement) SetAccountID(v string)`

SetAccountID sets AccountID field to given value.


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


