# BudgetItemNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Date** | **time.Time** |  | 
**AccountId** | **string** |  | 
**Amount** | **float64** |  | 
**Description** | Pointer to **string** |  | [optional] 

## Methods

### NewBudgetItemNoID

`func NewBudgetItemNoID(date time.Time, accountId string, amount float64, ) *BudgetItemNoID`

NewBudgetItemNoID instantiates a new BudgetItemNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBudgetItemNoIDWithDefaults

`func NewBudgetItemNoIDWithDefaults() *BudgetItemNoID`

NewBudgetItemNoIDWithDefaults instantiates a new BudgetItemNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *BudgetItemNoID) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *BudgetItemNoID) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *BudgetItemNoID) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetAccountId

`func (o *BudgetItemNoID) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BudgetItemNoID) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BudgetItemNoID) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetAmount

`func (o *BudgetItemNoID) GetAmount() float64`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *BudgetItemNoID) GetAmountOk() (*float64, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *BudgetItemNoID) SetAmount(v float64)`

SetAmount sets Amount field to given value.


### GetDescription

`func (o *BudgetItemNoID) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BudgetItemNoID) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BudgetItemNoID) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BudgetItemNoID) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


