# CurrencyAggregation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CurrencyID** | **string** |  | 
**Accounts** | [**[]AccountAggregation**](AccountAggregation.md) |  | 

## Methods

### NewCurrencyAggregation

`func NewCurrencyAggregation(currencyID string, accounts []AccountAggregation, ) *CurrencyAggregation`

NewCurrencyAggregation instantiates a new CurrencyAggregation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCurrencyAggregationWithDefaults

`func NewCurrencyAggregationWithDefaults() *CurrencyAggregation`

NewCurrencyAggregationWithDefaults instantiates a new CurrencyAggregation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrencyID

`func (o *CurrencyAggregation) GetCurrencyID() string`

GetCurrencyID returns the CurrencyID field if non-nil, zero value otherwise.

### GetCurrencyIDOk

`func (o *CurrencyAggregation) GetCurrencyIDOk() (*string, bool)`

GetCurrencyIDOk returns a tuple with the CurrencyID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyID

`func (o *CurrencyAggregation) SetCurrencyID(v string)`

SetCurrencyID sets CurrencyID field to given value.


### GetAccounts

`func (o *CurrencyAggregation) GetAccounts() []AccountAggregation`

GetAccounts returns the Accounts field if non-nil, zero value otherwise.

### GetAccountsOk

`func (o *CurrencyAggregation) GetAccountsOk() (*[]AccountAggregation, bool)`

GetAccountsOk returns a tuple with the Accounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccounts

`func (o *CurrencyAggregation) SetAccounts(v []AccountAggregation)`

SetAccounts sets Accounts field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


