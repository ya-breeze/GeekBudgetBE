# DisbalanceCandidate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Transactions** | [**[]DisbalanceCandidateTransaction**](DisbalanceCandidateTransaction.md) |  | 
**Sum** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Difference** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**Type** | **string** |  | 

## Methods

### NewDisbalanceCandidate

`func NewDisbalanceCandidate(transactions []DisbalanceCandidateTransaction, sum decimal.Decimal, difference decimal.Decimal, type_ string, ) *DisbalanceCandidate`

NewDisbalanceCandidate instantiates a new DisbalanceCandidate object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDisbalanceCandidateWithDefaults

`func NewDisbalanceCandidateWithDefaults() *DisbalanceCandidate`

NewDisbalanceCandidateWithDefaults instantiates a new DisbalanceCandidate object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransactions

`func (o *DisbalanceCandidate) GetTransactions() []DisbalanceCandidateTransaction`

GetTransactions returns the Transactions field if non-nil, zero value otherwise.

### GetTransactionsOk

`func (o *DisbalanceCandidate) GetTransactionsOk() (*[]DisbalanceCandidateTransaction, bool)`

GetTransactionsOk returns a tuple with the Transactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactions

`func (o *DisbalanceCandidate) SetTransactions(v []DisbalanceCandidateTransaction)`

SetTransactions sets Transactions field to given value.


### GetSum

`func (o *DisbalanceCandidate) GetSum() decimal.Decimal`

GetSum returns the Sum field if non-nil, zero value otherwise.

### GetSumOk

`func (o *DisbalanceCandidate) GetSumOk() (*decimal.Decimal, bool)`

GetSumOk returns a tuple with the Sum field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSum

`func (o *DisbalanceCandidate) SetSum(v decimal.Decimal)`

SetSum sets Sum field to given value.


### GetDifference

`func (o *DisbalanceCandidate) GetDifference() decimal.Decimal`

GetDifference returns the Difference field if non-nil, zero value otherwise.

### GetDifferenceOk

`func (o *DisbalanceCandidate) GetDifferenceOk() (*decimal.Decimal, bool)`

GetDifferenceOk returns a tuple with the Difference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDifference

`func (o *DisbalanceCandidate) SetDifference(v decimal.Decimal)`

SetDifference sets Difference field to given value.


### GetType

`func (o *DisbalanceCandidate) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DisbalanceCandidate) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DisbalanceCandidate) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


