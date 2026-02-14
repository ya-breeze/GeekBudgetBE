# DisbalanceAnalysis

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Delta** | [**decimal.Decimal**](decimal.Decimal.md) |  | 
**TransactionCount** | **int32** |  | 
**Candidates** | [**[]DisbalanceCandidate**](DisbalanceCandidate.md) |  | 

## Methods

### NewDisbalanceAnalysis

`func NewDisbalanceAnalysis(delta decimal.Decimal, transactionCount int32, candidates []DisbalanceCandidate, ) *DisbalanceAnalysis`

NewDisbalanceAnalysis instantiates a new DisbalanceAnalysis object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDisbalanceAnalysisWithDefaults

`func NewDisbalanceAnalysisWithDefaults() *DisbalanceAnalysis`

NewDisbalanceAnalysisWithDefaults instantiates a new DisbalanceAnalysis object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDelta

`func (o *DisbalanceAnalysis) GetDelta() decimal.Decimal`

GetDelta returns the Delta field if non-nil, zero value otherwise.

### GetDeltaOk

`func (o *DisbalanceAnalysis) GetDeltaOk() (*decimal.Decimal, bool)`

GetDeltaOk returns a tuple with the Delta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelta

`func (o *DisbalanceAnalysis) SetDelta(v decimal.Decimal)`

SetDelta sets Delta field to given value.


### GetTransactionCount

`func (o *DisbalanceAnalysis) GetTransactionCount() int32`

GetTransactionCount returns the TransactionCount field if non-nil, zero value otherwise.

### GetTransactionCountOk

`func (o *DisbalanceAnalysis) GetTransactionCountOk() (*int32, bool)`

GetTransactionCountOk returns a tuple with the TransactionCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactionCount

`func (o *DisbalanceAnalysis) SetTransactionCount(v int32)`

SetTransactionCount sets TransactionCount field to given value.


### GetCandidates

`func (o *DisbalanceAnalysis) GetCandidates() []DisbalanceCandidate`

GetCandidates returns the Candidates field if non-nil, zero value otherwise.

### GetCandidatesOk

`func (o *DisbalanceAnalysis) GetCandidatesOk() (*[]DisbalanceCandidate, bool)`

GetCandidatesOk returns a tuple with the Candidates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCandidates

`func (o *DisbalanceAnalysis) SetCandidates(v []DisbalanceCandidate)`

SetCandidates sets Candidates field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


