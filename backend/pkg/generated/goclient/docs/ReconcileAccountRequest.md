# ReconcileAccountRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CurrencyId** | **string** |  | 
**Balance** | Pointer to **float64** | Manual balance to set. If 0, current Account balance is used. | [optional] 

## Methods

### NewReconcileAccountRequest

`func NewReconcileAccountRequest(currencyId string, ) *ReconcileAccountRequest`

NewReconcileAccountRequest instantiates a new ReconcileAccountRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReconcileAccountRequestWithDefaults

`func NewReconcileAccountRequestWithDefaults() *ReconcileAccountRequest`

NewReconcileAccountRequestWithDefaults instantiates a new ReconcileAccountRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrencyId

`func (o *ReconcileAccountRequest) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *ReconcileAccountRequest) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *ReconcileAccountRequest) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetBalance

`func (o *ReconcileAccountRequest) GetBalance() float64`

GetBalance returns the Balance field if non-nil, zero value otherwise.

### GetBalanceOk

`func (o *ReconcileAccountRequest) GetBalanceOk() (*float64, bool)`

GetBalanceOk returns a tuple with the Balance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalance

`func (o *ReconcileAccountRequest) SetBalance(v float64)`

SetBalance sets Balance field to given value.

### HasBalance

`func (o *ReconcileAccountRequest) HasBalance() bool`

HasBalance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


