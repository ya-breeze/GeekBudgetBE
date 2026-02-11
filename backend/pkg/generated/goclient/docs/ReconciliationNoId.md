# ReconciliationNoId

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | **string** |  | 
**CurrencyId** | **string** |  | 
**ReconciledBalance** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**ExpectedBalance** | Pointer to [**decimal.Decimal**](decimal.Decimal.md) |  | [optional] 
**IsManual** | Pointer to **bool** |  | [optional] 

## Methods

### NewReconciliationNoId

`func NewReconciliationNoId(accountId string, currencyId string, reconciledBalance decimal.Decimal, ) *ReconciliationNoId`

NewReconciliationNoId instantiates a new ReconciliationNoId object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReconciliationNoIdWithDefaults

`func NewReconciliationNoIdWithDefaults() *ReconciliationNoId`

NewReconciliationNoIdWithDefaults instantiates a new ReconciliationNoId object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *ReconciliationNoId) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *ReconciliationNoId) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *ReconciliationNoId) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetCurrencyId

`func (o *ReconciliationNoId) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *ReconciliationNoId) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *ReconciliationNoId) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetReconciledBalance

`func (o *ReconciliationNoId) GetReconciledBalance() decimal.Decimal`

GetReconciledBalance returns the ReconciledBalance field if non-nil, zero value otherwise.

### GetReconciledBalanceOk

`func (o *ReconciliationNoId) GetReconciledBalanceOk() (*decimal.Decimal, bool)`

GetReconciledBalanceOk returns a tuple with the ReconciledBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReconciledBalance

`func (o *ReconciliationNoId) SetReconciledBalance(v decimal.Decimal)`

SetReconciledBalance sets ReconciledBalance field to given value.


### GetExpectedBalance

`func (o *ReconciliationNoId) GetExpectedBalance() decimal.Decimal`

GetExpectedBalance returns the ExpectedBalance field if non-nil, zero value otherwise.

### GetExpectedBalanceOk

`func (o *ReconciliationNoId) GetExpectedBalanceOk() (*decimal.Decimal, bool)`

GetExpectedBalanceOk returns a tuple with the ExpectedBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedBalance

`func (o *ReconciliationNoId) SetExpectedBalance(v decimal.Decimal)`

SetExpectedBalance sets ExpectedBalance field to given value.

### HasExpectedBalance

`func (o *ReconciliationNoId) HasExpectedBalance() bool`

HasExpectedBalance returns a boolean if a field has been set.

### GetIsManual

`func (o *ReconciliationNoId) GetIsManual() bool`

GetIsManual returns the IsManual field if non-nil, zero value otherwise.

### GetIsManualOk

`func (o *ReconciliationNoId) GetIsManualOk() (*bool, bool)`

GetIsManualOk returns a tuple with the IsManual field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsManual

`func (o *ReconciliationNoId) SetIsManual(v bool)`

SetIsManual sets IsManual field to given value.

### HasIsManual

`func (o *ReconciliationNoId) HasIsManual() bool`

HasIsManual returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


