# MatcherAndTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MatcherID** | **string** |  | 
**Transaction** | [**TransactionNoID**](TransactionNoID.md) |  | 

## Methods

### NewMatcherAndTransaction

`func NewMatcherAndTransaction(matcherID string, transaction TransactionNoID, ) *MatcherAndTransaction`

NewMatcherAndTransaction instantiates a new MatcherAndTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMatcherAndTransactionWithDefaults

`func NewMatcherAndTransactionWithDefaults() *MatcherAndTransaction`

NewMatcherAndTransactionWithDefaults instantiates a new MatcherAndTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMatcherID

`func (o *MatcherAndTransaction) GetMatcherID() string`

GetMatcherID returns the MatcherID field if non-nil, zero value otherwise.

### GetMatcherIDOk

`func (o *MatcherAndTransaction) GetMatcherIDOk() (*string, bool)`

GetMatcherIDOk returns a tuple with the MatcherID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatcherID

`func (o *MatcherAndTransaction) SetMatcherID(v string)`

SetMatcherID sets MatcherID field to given value.


### GetTransaction

`func (o *MatcherAndTransaction) GetTransaction() TransactionNoID`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *MatcherAndTransaction) GetTransactionOk() (*TransactionNoID, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *MatcherAndTransaction) SetTransaction(v TransactionNoID)`

SetTransaction sets Transaction field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


