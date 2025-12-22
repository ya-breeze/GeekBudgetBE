# AccountAggregation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | **string** |  | 
**Amounts** | **[]float64** |  | 
**Total** | Pointer to **float64** |  | [optional] 
**ChangePercent** | Pointer to **float64** |  | [optional] 

## Methods

### NewAccountAggregation

`func NewAccountAggregation(accountId string, amounts []float64, ) *AccountAggregation`

NewAccountAggregation instantiates a new AccountAggregation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountAggregationWithDefaults

`func NewAccountAggregationWithDefaults() *AccountAggregation`

NewAccountAggregationWithDefaults instantiates a new AccountAggregation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *AccountAggregation) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *AccountAggregation) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *AccountAggregation) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


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


### GetTotal

`func (o *AccountAggregation) GetTotal() float64`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *AccountAggregation) GetTotalOk() (*float64, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *AccountAggregation) SetTotal(v float64)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *AccountAggregation) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetChangePercent

`func (o *AccountAggregation) GetChangePercent() float64`

GetChangePercent returns the ChangePercent field if non-nil, zero value otherwise.

### GetChangePercentOk

`func (o *AccountAggregation) GetChangePercentOk() (*float64, bool)`

GetChangePercentOk returns a tuple with the ChangePercent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChangePercent

`func (o *AccountAggregation) SetChangePercent(v float64)`

SetChangePercent sets ChangePercent field to given value.

### HasChangePercent

`func (o *AccountAggregation) HasChangePercent() bool`

HasChangePercent returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


