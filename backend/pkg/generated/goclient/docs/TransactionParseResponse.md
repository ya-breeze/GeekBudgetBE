# TransactionParseResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Transaction** | [**TransactionNoID**](TransactionNoID.md) |  | 
**Warnings** | **[]string** | Human-readable warnings for fields that could not be parsed or were ambiguously matched | 

## Methods

### NewTransactionParseResponse

`func NewTransactionParseResponse(transaction TransactionNoID, warnings []string, ) *TransactionParseResponse`

NewTransactionParseResponse instantiates a new TransactionParseResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionParseResponseWithDefaults

`func NewTransactionParseResponseWithDefaults() *TransactionParseResponse`

NewTransactionParseResponseWithDefaults instantiates a new TransactionParseResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransaction

`func (o *TransactionParseResponse) GetTransaction() TransactionNoID`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *TransactionParseResponse) GetTransactionOk() (*TransactionNoID, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *TransactionParseResponse) SetTransaction(v TransactionNoID)`

SetTransaction sets Transaction field to given value.


### GetWarnings

`func (o *TransactionParseResponse) GetWarnings() []string`

GetWarnings returns the Warnings field if non-nil, zero value otherwise.

### GetWarningsOk

`func (o *TransactionParseResponse) GetWarningsOk() (*[]string, bool)`

GetWarningsOk returns a tuple with the Warnings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarnings

`func (o *TransactionParseResponse) SetWarnings(v []string)`

SetWarnings sets Warnings field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


