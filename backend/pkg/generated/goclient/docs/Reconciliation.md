# Reconciliation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ReconciliationId** | Pointer to **string** |  | [optional] [readonly] 
**AccountId** | **string** |  | 
**CurrencyId** | **string** |  | 
**ReconciledBalance** | **float64** |  | 
**ReconciledAt** | **time.Time** |  | 
**ExpectedBalance** | Pointer to **float64** |  | [optional] 
**IsManual** | Pointer to **bool** |  | [optional] 

## Methods

### NewReconciliation

`func NewReconciliation(accountId string, currencyId string, reconciledBalance float64, reconciledAt time.Time, ) *Reconciliation`

NewReconciliation instantiates a new Reconciliation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReconciliationWithDefaults

`func NewReconciliationWithDefaults() *Reconciliation`

NewReconciliationWithDefaults instantiates a new Reconciliation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReconciliationId

`func (o *Reconciliation) GetReconciliationId() string`

GetReconciliationId returns the ReconciliationId field if non-nil, zero value otherwise.

### GetReconciliationIdOk

`func (o *Reconciliation) GetReconciliationIdOk() (*string, bool)`

GetReconciliationIdOk returns a tuple with the ReconciliationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReconciliationId

`func (o *Reconciliation) SetReconciliationId(v string)`

SetReconciliationId sets ReconciliationId field to given value.

### HasReconciliationId

`func (o *Reconciliation) HasReconciliationId() bool`

HasReconciliationId returns a boolean if a field has been set.

### GetAccountId

`func (o *Reconciliation) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *Reconciliation) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *Reconciliation) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetCurrencyId

`func (o *Reconciliation) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *Reconciliation) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *Reconciliation) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetReconciledBalance

`func (o *Reconciliation) GetReconciledBalance() float64`

GetReconciledBalance returns the ReconciledBalance field if non-nil, zero value otherwise.

### GetReconciledBalanceOk

`func (o *Reconciliation) GetReconciledBalanceOk() (*float64, bool)`

GetReconciledBalanceOk returns a tuple with the ReconciledBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReconciledBalance

`func (o *Reconciliation) SetReconciledBalance(v float64)`

SetReconciledBalance sets ReconciledBalance field to given value.


### GetReconciledAt

`func (o *Reconciliation) GetReconciledAt() time.Time`

GetReconciledAt returns the ReconciledAt field if non-nil, zero value otherwise.

### GetReconciledAtOk

`func (o *Reconciliation) GetReconciledAtOk() (*time.Time, bool)`

GetReconciledAtOk returns a tuple with the ReconciledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReconciledAt

`func (o *Reconciliation) SetReconciledAt(v time.Time)`

SetReconciledAt sets ReconciledAt field to given value.


### GetExpectedBalance

`func (o *Reconciliation) GetExpectedBalance() float64`

GetExpectedBalance returns the ExpectedBalance field if non-nil, zero value otherwise.

### GetExpectedBalanceOk

`func (o *Reconciliation) GetExpectedBalanceOk() (*float64, bool)`

GetExpectedBalanceOk returns a tuple with the ExpectedBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedBalance

`func (o *Reconciliation) SetExpectedBalance(v float64)`

SetExpectedBalance sets ExpectedBalance field to given value.

### HasExpectedBalance

`func (o *Reconciliation) HasExpectedBalance() bool`

HasExpectedBalance returns a boolean if a field has been set.

### GetIsManual

`func (o *Reconciliation) GetIsManual() bool`

GetIsManual returns the IsManual field if non-nil, zero value otherwise.

### GetIsManualOk

`func (o *Reconciliation) GetIsManualOk() (*bool, bool)`

GetIsManualOk returns a tuple with the IsManual field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsManual

`func (o *Reconciliation) SetIsManual(v bool)`

SetIsManual sets IsManual field to given value.

### HasIsManual

`func (o *Reconciliation) HasIsManual() bool`

HasIsManual returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


