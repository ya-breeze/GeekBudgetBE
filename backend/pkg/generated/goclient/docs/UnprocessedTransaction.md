# UnprocessedTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Transaction** | [**Transaction**](Transaction.md) |  | 
**Matched** | [**[]MatcherAndTransaction**](MatcherAndTransaction.md) |  | 
**Duplicates** | [**[]Transaction**](Transaction.md) |  | 

## Methods

### NewUnprocessedTransaction

`func NewUnprocessedTransaction(transaction Transaction, matched []MatcherAndTransaction, duplicates []Transaction, ) *UnprocessedTransaction`

NewUnprocessedTransaction instantiates a new UnprocessedTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUnprocessedTransactionWithDefaults

`func NewUnprocessedTransactionWithDefaults() *UnprocessedTransaction`

NewUnprocessedTransactionWithDefaults instantiates a new UnprocessedTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransaction

`func (o *UnprocessedTransaction) GetTransaction() Transaction`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *UnprocessedTransaction) GetTransactionOk() (*Transaction, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *UnprocessedTransaction) SetTransaction(v Transaction)`

SetTransaction sets Transaction field to given value.


### GetMatched

`func (o *UnprocessedTransaction) GetMatched() []MatcherAndTransaction`

GetMatched returns the Matched field if non-nil, zero value otherwise.

### GetMatchedOk

`func (o *UnprocessedTransaction) GetMatchedOk() (*[]MatcherAndTransaction, bool)`

GetMatchedOk returns a tuple with the Matched field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatched

`func (o *UnprocessedTransaction) SetMatched(v []MatcherAndTransaction)`

SetMatched sets Matched field to given value.


### GetDuplicates

`func (o *UnprocessedTransaction) GetDuplicates() []Transaction`

GetDuplicates returns the Duplicates field if non-nil, zero value otherwise.

### GetDuplicatesOk

`func (o *UnprocessedTransaction) GetDuplicatesOk() (*[]Transaction, bool)`

GetDuplicatesOk returns a tuple with the Duplicates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicates

`func (o *UnprocessedTransaction) SetDuplicates(v []Transaction)`

SetDuplicates sets Duplicates field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


