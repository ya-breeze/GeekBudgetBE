# ConvertUnprocessedTransaction200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Transaction** | [**Transaction**](Transaction.md) |  | 
**AutoProcessedIds** | Pointer to **[]string** | IDs of other unprocessed transactions that were automatically matched because of this conversion (if matcher became perfect) | [optional] 

## Methods

### NewConvertUnprocessedTransaction200Response

`func NewConvertUnprocessedTransaction200Response(transaction Transaction, ) *ConvertUnprocessedTransaction200Response`

NewConvertUnprocessedTransaction200Response instantiates a new ConvertUnprocessedTransaction200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConvertUnprocessedTransaction200ResponseWithDefaults

`func NewConvertUnprocessedTransaction200ResponseWithDefaults() *ConvertUnprocessedTransaction200Response`

NewConvertUnprocessedTransaction200ResponseWithDefaults instantiates a new ConvertUnprocessedTransaction200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransaction

`func (o *ConvertUnprocessedTransaction200Response) GetTransaction() Transaction`

GetTransaction returns the Transaction field if non-nil, zero value otherwise.

### GetTransactionOk

`func (o *ConvertUnprocessedTransaction200Response) GetTransactionOk() (*Transaction, bool)`

GetTransactionOk returns a tuple with the Transaction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransaction

`func (o *ConvertUnprocessedTransaction200Response) SetTransaction(v Transaction)`

SetTransaction sets Transaction field to given value.


### GetAutoProcessedIds

`func (o *ConvertUnprocessedTransaction200Response) GetAutoProcessedIds() []string`

GetAutoProcessedIds returns the AutoProcessedIds field if non-nil, zero value otherwise.

### GetAutoProcessedIdsOk

`func (o *ConvertUnprocessedTransaction200Response) GetAutoProcessedIdsOk() (*[]string, bool)`

GetAutoProcessedIdsOk returns a tuple with the AutoProcessedIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoProcessedIds

`func (o *ConvertUnprocessedTransaction200Response) SetAutoProcessedIds(v []string)`

SetAutoProcessedIds sets AutoProcessedIds field to given value.

### HasAutoProcessedIds

`func (o *ConvertUnprocessedTransaction200Response) HasAutoProcessedIds() bool`

HasAutoProcessedIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


