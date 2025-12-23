# MergedTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Transaction** | [**Transaction**](Transaction.md) |  | 
**MergedInto** | [**Transaction**](Transaction.md) |  | 
**MergedAt** | **time.Time** | When this transaction was merged | 

## Methods

### NewMergedTransaction

`func NewMergedTransaction(transaction Transaction, mergedInto Transaction, mergedAt time.Time, ) *MergedTransaction`

NewMergedTransaction instantiates a new MergedTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMergedTransactionWithDefaults

`func NewMergedTransactionWithDefaults() *MergedTransaction`

NewMergedTransactionWithDefaults instantiates a new MergedTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransaction

`func (o *MergedTransaction) GetTransaction() Transaction`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *MergedTransaction) GetTransactionOk() (*Transaction, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *MergedTransaction) SetTransaction(v Transaction)`

SetTransaction sets Transaction field to given value.


### GetMergedInto

`func (o *MergedTransaction) GetMergedInto() Transaction`

GetMergedInto returns the MergedInto field if non-nil, zero value otherwise.

### GetMergedIntoOk

`func (o *MergedTransaction) GetMergedIntoOk() (*Transaction, bool)`

GetMergedIntoOk returns a tuple with the MergedInto field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedInto

`func (o *MergedTransaction) SetMergedInto(v Transaction)`

SetMergedInto sets MergedInto field to given value.


### GetMergedAt

`func (o *MergedTransaction) GetMergedAt() time.Time`

GetMergedAt returns the MergedAt field if non-nil, zero value otherwise.

### GetMergedAtOk

`func (o *MergedTransaction) GetMergedAtOk() (*time.Time, bool)`

GetMergedAtOk returns a tuple with the MergedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedAt

`func (o *MergedTransaction) SetMergedAt(v time.Time)`

SetMergedAt sets MergedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


