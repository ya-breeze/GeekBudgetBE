# BudgetItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Date** | **time.Time** |  | 
**AccountId** | **string** |  | 
**Amount** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Description** | Pointer to **string** |  | [optional] 

## Methods

### NewBudgetItem

`func NewBudgetItem(id string, date time.Time, accountId string, amount decimal.Decimal, ) *BudgetItem`

NewBudgetItem instantiates a new BudgetItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBudgetItemWithDefaults

`func NewBudgetItemWithDefaults() *BudgetItem`

NewBudgetItemWithDefaults instantiates a new BudgetItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BudgetItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BudgetItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BudgetItem) SetId(v string)`

SetId sets Id field to given value.


### GetDate

`func (o *BudgetItem) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *BudgetItem) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *BudgetItem) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetAccountId

`func (o *BudgetItem) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BudgetItem) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BudgetItem) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetAmount

`func (o *BudgetItem) GetAmount() decimal.Decimal`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *BudgetItem) GetAmountOk() (*decimal.Decimal, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *BudgetItem) SetAmount(v decimal.Decimal)`

SetAmount sets Amount field to given value.


### GetDescription

`func (o *BudgetItem) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BudgetItem) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BudgetItem) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BudgetItem) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


