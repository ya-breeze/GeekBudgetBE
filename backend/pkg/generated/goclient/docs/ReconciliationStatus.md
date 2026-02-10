# ReconciliationStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | **string** |  | 
**AccountName** | **string** |  | 
**CurrencyId** | **string** |  | 
**CurrencySymbol** | Pointer to **string** |  | [optional] 
**BankBalance** | Pointer to **float64** |  | [optional] 
**AppBalance** | **float64** |  | 
**Delta** | **float64** |  | 
**LastReconciledAt** | Pointer to **NullableTime** |  | [optional] 
**LastReconciledBalance** | Pointer to **float64** |  | [optional] 
**HasUnprocessedTransactions** | Pointer to **bool** |  | [optional] 
**HasBankImporter** | Pointer to **bool** |  | [optional] 
**IsManualReconciliationEnabled** | Pointer to **bool** |  | [optional] 
**BankBalanceAt** | Pointer to **NullableTime** |  | [optional] 
**HasTransactionsAfterBankBalance** | Pointer to **bool** |  | [optional] 

## Methods

### NewReconciliationStatus

`func NewReconciliationStatus(accountId string, accountName string, currencyId string, appBalance float64, delta float64, ) *ReconciliationStatus`

NewReconciliationStatus instantiates a new ReconciliationStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReconciliationStatusWithDefaults

`func NewReconciliationStatusWithDefaults() *ReconciliationStatus`

NewReconciliationStatusWithDefaults instantiates a new ReconciliationStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *ReconciliationStatus) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *ReconciliationStatus) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *ReconciliationStatus) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetAccountName

`func (o *ReconciliationStatus) GetAccountName() string`

GetAccountName returns the AccountName field if non-nil, zero value otherwise.

### GetAccountNameOk

`func (o *ReconciliationStatus) GetAccountNameOk() (*string, bool)`

GetAccountNameOk returns a tuple with the AccountName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountName

`func (o *ReconciliationStatus) SetAccountName(v string)`

SetAccountName sets AccountName field to given value.


### GetCurrencyId

`func (o *ReconciliationStatus) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *ReconciliationStatus) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *ReconciliationStatus) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetCurrencySymbol

`func (o *ReconciliationStatus) GetCurrencySymbol() string`

GetCurrencySymbol returns the CurrencySymbol field if non-nil, zero value otherwise.

### GetCurrencySymbolOk

`func (o *ReconciliationStatus) GetCurrencySymbolOk() (*string, bool)`

GetCurrencySymbolOk returns a tuple with the CurrencySymbol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencySymbol

`func (o *ReconciliationStatus) SetCurrencySymbol(v string)`

SetCurrencySymbol sets CurrencySymbol field to given value.

### HasCurrencySymbol

`func (o *ReconciliationStatus) HasCurrencySymbol() bool`

HasCurrencySymbol returns a boolean if a field has been set.

### GetBankBalance

`func (o *ReconciliationStatus) GetBankBalance() float64`

GetBankBalance returns the BankBalance field if non-nil, zero value otherwise.

### GetBankBalanceOk

`func (o *ReconciliationStatus) GetBankBalanceOk() (*float64, bool)`

GetBankBalanceOk returns a tuple with the BankBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankBalance

`func (o *ReconciliationStatus) SetBankBalance(v float64)`

SetBankBalance sets BankBalance field to given value.

### HasBankBalance

`func (o *ReconciliationStatus) HasBankBalance() bool`

HasBankBalance returns a boolean if a field has been set.

### GetAppBalance

`func (o *ReconciliationStatus) GetAppBalance() float64`

GetAppBalance returns the AppBalance field if non-nil, zero value otherwise.

### GetAppBalanceOk

`func (o *ReconciliationStatus) GetAppBalanceOk() (*float64, bool)`

GetAppBalanceOk returns a tuple with the AppBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppBalance

`func (o *ReconciliationStatus) SetAppBalance(v float64)`

SetAppBalance sets AppBalance field to given value.


### GetDelta

`func (o *ReconciliationStatus) GetDelta() float64`

GetDelta returns the Delta field if non-nil, zero value otherwise.

### GetDeltaOk

`func (o *ReconciliationStatus) GetDeltaOk() (*float64, bool)`

GetDeltaOk returns a tuple with the Delta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelta

`func (o *ReconciliationStatus) SetDelta(v float64)`

SetDelta sets Delta field to given value.


### GetLastReconciledAt

`func (o *ReconciliationStatus) GetLastReconciledAt() time.Time`

GetLastReconciledAt returns the LastReconciledAt field if non-nil, zero value otherwise.

### GetLastReconciledAtOk

`func (o *ReconciliationStatus) GetLastReconciledAtOk() (*time.Time, bool)`

GetLastReconciledAtOk returns a tuple with the LastReconciledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastReconciledAt

`func (o *ReconciliationStatus) SetLastReconciledAt(v time.Time)`

SetLastReconciledAt sets LastReconciledAt field to given value.

### HasLastReconciledAt

`func (o *ReconciliationStatus) HasLastReconciledAt() bool`

HasLastReconciledAt returns a boolean if a field has been set.

### SetLastReconciledAtNil

`func (o *ReconciliationStatus) SetLastReconciledAtNil(b bool)`

 SetLastReconciledAtNil sets the value for LastReconciledAt to be an explicit nil

### UnsetLastReconciledAt
`func (o *ReconciliationStatus) UnsetLastReconciledAt()`

UnsetLastReconciledAt ensures that no value is present for LastReconciledAt, not even an explicit nil
### GetLastReconciledBalance

`func (o *ReconciliationStatus) GetLastReconciledBalance() float64`

GetLastReconciledBalance returns the LastReconciledBalance field if non-nil, zero value otherwise.

### GetLastReconciledBalanceOk

`func (o *ReconciliationStatus) GetLastReconciledBalanceOk() (*float64, bool)`

GetLastReconciledBalanceOk returns a tuple with the LastReconciledBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastReconciledBalance

`func (o *ReconciliationStatus) SetLastReconciledBalance(v float64)`

SetLastReconciledBalance sets LastReconciledBalance field to given value.

### HasLastReconciledBalance

`func (o *ReconciliationStatus) HasLastReconciledBalance() bool`

HasLastReconciledBalance returns a boolean if a field has been set.

### GetHasUnprocessedTransactions

`func (o *ReconciliationStatus) GetHasUnprocessedTransactions() bool`

GetHasUnprocessedTransactions returns the HasUnprocessedTransactions field if non-nil, zero value otherwise.

### GetHasUnprocessedTransactionsOk

`func (o *ReconciliationStatus) GetHasUnprocessedTransactionsOk() (*bool, bool)`

GetHasUnprocessedTransactionsOk returns a tuple with the HasUnprocessedTransactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasUnprocessedTransactions

`func (o *ReconciliationStatus) SetHasUnprocessedTransactions(v bool)`

SetHasUnprocessedTransactions sets HasUnprocessedTransactions field to given value.

### HasHasUnprocessedTransactions

`func (o *ReconciliationStatus) HasHasUnprocessedTransactions() bool`

HasHasUnprocessedTransactions returns a boolean if a field has been set.

### GetHasBankImporter

`func (o *ReconciliationStatus) GetHasBankImporter() bool`

GetHasBankImporter returns the HasBankImporter field if non-nil, zero value otherwise.

### GetHasBankImporterOk

`func (o *ReconciliationStatus) GetHasBankImporterOk() (*bool, bool)`

GetHasBankImporterOk returns a tuple with the HasBankImporter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasBankImporter

`func (o *ReconciliationStatus) SetHasBankImporter(v bool)`

SetHasBankImporter sets HasBankImporter field to given value.

### HasHasBankImporter

`func (o *ReconciliationStatus) HasHasBankImporter() bool`

HasHasBankImporter returns a boolean if a field has been set.

### GetIsManualReconciliationEnabled

`func (o *ReconciliationStatus) GetIsManualReconciliationEnabled() bool`

GetIsManualReconciliationEnabled returns the IsManualReconciliationEnabled field if non-nil, zero value otherwise.

### GetIsManualReconciliationEnabledOk

`func (o *ReconciliationStatus) GetIsManualReconciliationEnabledOk() (*bool, bool)`

GetIsManualReconciliationEnabledOk returns a tuple with the IsManualReconciliationEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsManualReconciliationEnabled

`func (o *ReconciliationStatus) SetIsManualReconciliationEnabled(v bool)`

SetIsManualReconciliationEnabled sets IsManualReconciliationEnabled field to given value.

### HasIsManualReconciliationEnabled

`func (o *ReconciliationStatus) HasIsManualReconciliationEnabled() bool`

HasIsManualReconciliationEnabled returns a boolean if a field has been set.

### GetBankBalanceAt

`func (o *ReconciliationStatus) GetBankBalanceAt() time.Time`

GetBankBalanceAt returns the BankBalanceAt field if non-nil, zero value otherwise.

### GetBankBalanceAtOk

`func (o *ReconciliationStatus) GetBankBalanceAtOk() (*time.Time, bool)`

GetBankBalanceAtOk returns a tuple with the BankBalanceAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankBalanceAt

`func (o *ReconciliationStatus) SetBankBalanceAt(v time.Time)`

SetBankBalanceAt sets BankBalanceAt field to given value.

### HasBankBalanceAt

`func (o *ReconciliationStatus) HasBankBalanceAt() bool`

HasBankBalanceAt returns a boolean if a field has been set.

### SetBankBalanceAtNil

`func (o *ReconciliationStatus) SetBankBalanceAtNil(b bool)`

 SetBankBalanceAtNil sets the value for BankBalanceAt to be an explicit nil

### UnsetBankBalanceAt
`func (o *ReconciliationStatus) UnsetBankBalanceAt()`

UnsetBankBalanceAt ensures that no value is present for BankBalanceAt, not even an explicit nil
### GetHasTransactionsAfterBankBalance

`func (o *ReconciliationStatus) GetHasTransactionsAfterBankBalance() bool`

GetHasTransactionsAfterBankBalance returns the HasTransactionsAfterBankBalance field if non-nil, zero value otherwise.

### GetHasTransactionsAfterBankBalanceOk

`func (o *ReconciliationStatus) GetHasTransactionsAfterBankBalanceOk() (*bool, bool)`

GetHasTransactionsAfterBankBalanceOk returns a tuple with the HasTransactionsAfterBankBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasTransactionsAfterBankBalance

`func (o *ReconciliationStatus) SetHasTransactionsAfterBankBalance(v bool)`

SetHasTransactionsAfterBankBalance sets HasTransactionsAfterBankBalance field to given value.

### HasHasTransactionsAfterBankBalance

`func (o *ReconciliationStatus) HasHasTransactionsAfterBankBalance() bool`

HasHasTransactionsAfterBankBalance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


