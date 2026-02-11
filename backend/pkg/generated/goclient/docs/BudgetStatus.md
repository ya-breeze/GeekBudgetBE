# BudgetStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Date** | **time.Time** |  | 
**AccountId** | **string** |  | 
**Budgeted** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Spent** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Rollover** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Available** | [**decimal.Decimal**](decimal.Decimal.md) |  | 

## Methods

### NewBudgetStatus

`func NewBudgetStatus(date time.Time, accountId string, budgeted decimal.Decimal, spent decimal.Decimal, rollover decimal.Decimal, available decimal.Decimal, ) *BudgetStatus`

NewBudgetStatus instantiates a new BudgetStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBudgetStatusWithDefaults

`func NewBudgetStatusWithDefaults() *BudgetStatus`

NewBudgetStatusWithDefaults instantiates a new BudgetStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *BudgetStatus) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *BudgetStatus) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *BudgetStatus) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetAccountId

`func (o *BudgetStatus) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BudgetStatus) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BudgetStatus) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetBudgeted

`func (o *BudgetStatus) GetBudgeted() decimal.Decimal`

GetBudgeted returns the Budgeted field if non-nil, zero value otherwise.

### GetBudgetedOk

`func (o *BudgetStatus) GetBudgetedOk() (*decimal.Decimal, bool)`

GetBudgetedOk returns a tuple with the Budgeted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBudgeted

`func (o *BudgetStatus) SetBudgeted(v decimal.Decimal)`

SetBudgeted sets Budgeted field to given value.


### GetSpent

`func (o *BudgetStatus) GetSpent() decimal.Decimal`

GetSpent returns the Spent field if non-nil, zero value otherwise.

### GetSpentOk

`func (o *BudgetStatus) GetSpentOk() (*decimal.Decimal, bool)`

GetSpentOk returns a tuple with the Spent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpent

`func (o *BudgetStatus) SetSpent(v decimal.Decimal)`

SetSpent sets Spent field to given value.


### GetRollover

`func (o *BudgetStatus) GetRollover() decimal.Decimal`

GetRollover returns the Rollover field if non-nil, zero value otherwise.

### GetRolloverOk

`func (o *BudgetStatus) GetRolloverOk() (*decimal.Decimal, bool)`

GetRolloverOk returns a tuple with the Rollover field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRollover

`func (o *BudgetStatus) SetRollover(v decimal.Decimal)`

SetRollover sets Rollover field to given value.


### GetAvailable

`func (o *BudgetStatus) GetAvailable() decimal.Decimal`

GetAvailable returns the Available field if non-nil, zero value otherwise.

### GetAvailableOk

`func (o *BudgetStatus) GetAvailableOk() (*decimal.Decimal, bool)`

GetAvailableOk returns a tuple with the Available field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailable

`func (o *BudgetStatus) SetAvailable(v decimal.Decimal)`

SetAvailable sets Available field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


