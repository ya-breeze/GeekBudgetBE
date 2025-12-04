# CheckMatcherRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Matcher** | [**MatcherNoID**](MatcherNoID.md) |  | 
**Transaction** | [**TransactionNoID**](TransactionNoID.md) |  | 

## Methods

### NewCheckMatcherRequest

`func NewCheckMatcherRequest(matcher MatcherNoID, transaction TransactionNoID, ) *CheckMatcherRequest`

NewCheckMatcherRequest instantiates a new CheckMatcherRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCheckMatcherRequestWithDefaults

`func NewCheckMatcherRequestWithDefaults() *CheckMatcherRequest`

NewCheckMatcherRequestWithDefaults instantiates a new CheckMatcherRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMatcher

`func (o *CheckMatcherRequest) GetMatcher() MatcherNoID`

GetMatcher returns the Matcher field if non-nil, zero value otherwise.

### GetMatcherOk

`func (o *CheckMatcherRequest) GetMatcherOk() (*MatcherNoID, bool)`

GetMatcherOk returns a tuple with the Matcher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatcher

`func (o *CheckMatcherRequest) SetMatcher(v MatcherNoID)`

SetMatcher sets Matcher field to given value.


### GetTransaction

`func (o *CheckMatcherRequest) GetTransaction() TransactionNoID`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *CheckMatcherRequest) GetTransactionOk() (*TransactionNoID, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *CheckMatcherRequest) SetTransaction(v TransactionNoID)`

SetTransaction sets Transaction field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


