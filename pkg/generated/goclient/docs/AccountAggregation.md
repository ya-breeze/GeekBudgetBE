# AccountAggregation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountID** | **string** |  | 
**Amounts** | **[]float64** |  | 

## Methods

### NewAccountAggregation

`func NewAccountAggregation(accountID string, amounts []float64, ) *AccountAggregation`

NewAccountAggregation instantiates a new AccountAggregation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountAggregationWithDefaults

`func NewAccountAggregationWithDefaults() *AccountAggregation`

NewAccountAggregationWithDefaults instantiates a new AccountAggregation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountID

`func (o *AccountAggregation) GetAccountID() string`

GetAccountID returns the AccountID field if non-nil, zero value otherwise.

### GetAccountIDOk

`func (o *AccountAggregation) GetAccountIDOk() (*string, bool)`

GetAccountIDOk returns a tuple with the AccountID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountID

`func (o *AccountAggregation) SetAccountID(v string)`

SetAccountID sets AccountID field to given value.


### GetAmounts

`func (o *AccountAggregation) GetAmounts() []float64`

GetAmounts returns the Amounts field if non-nil, zero value otherwise.

### GetAmountsOk

`func (o *AccountAggregation) GetAmountsOk() (*[]float64, bool)`

GetAmountsOk returns a tuple with the Amounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmounts

`func (o *AccountAggregation) SetAmounts(v []float64)`

SetAmounts sets Amounts field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


