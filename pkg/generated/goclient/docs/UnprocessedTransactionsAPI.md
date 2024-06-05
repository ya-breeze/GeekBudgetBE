# \UnprocessedTransactionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConvertUnprocessedTransaction**](UnprocessedTransactionsAPI.md#ConvertUnprocessedTransaction) | **Post** /v1/unprocessedTransactions/{id}/convert | convert unprocessed transactions into normal transaction
[**DeleteUnprocessedTransaction**](UnprocessedTransactionsAPI.md#DeleteUnprocessedTransaction) | **Delete** /v1/unprocessedTransactions/{id} | delete unprocessed transaction
[**GetUnprocessedTransactions**](UnprocessedTransactionsAPI.md#GetUnprocessedTransactions) | **Get** /v1/unprocessedTransactions | get all unprocessed transactions



## ConvertUnprocessedTransaction

> Transaction ConvertUnprocessedTransaction(ctx, id).TransactionNoID(transactionNoID).Execute()

convert unprocessed transactions into normal transaction

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "123e4567-e89b-12d3-a456-426614174000" // string | 
	transactionNoID := *openapiclient.NewTransactionNoID(time.Now(), []openapiclient.Movement{*openapiclient.NewMovement(float64(123), "CurrencyID_example", "AccountID_example")}) // TransactionNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UnprocessedTransactionsAPI.ConvertUnprocessedTransaction(context.Background(), id).TransactionNoID(transactionNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UnprocessedTransactionsAPI.ConvertUnprocessedTransaction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ConvertUnprocessedTransaction`: Transaction
	fmt.Fprintf(os.Stdout, "Response from `UnprocessedTransactionsAPI.ConvertUnprocessedTransaction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiConvertUnprocessedTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **transactionNoID** | [**TransactionNoID**](TransactionNoID.md) |  | 

### Return type

[**Transaction**](Transaction.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUnprocessedTransaction

> DeleteUnprocessedTransaction(ctx, id).DuplicateOf(duplicateOf).Execute()

delete unprocessed transaction

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "123e4567-e89b-12d3-a456-426614174000" // string | 
	duplicateOf := "123e4567-e89b-12d3-a456-426614174000" // string | ID of transaction which is duplicate of this one (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.UnprocessedTransactionsAPI.DeleteUnprocessedTransaction(context.Background(), id).DuplicateOf(duplicateOf).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UnprocessedTransactionsAPI.DeleteUnprocessedTransaction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUnprocessedTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **duplicateOf** | **string** | ID of transaction which is duplicate of this one | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUnprocessedTransactions

> []UnprocessedTransaction GetUnprocessedTransactions(ctx).Execute()

get all unprocessed transactions

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UnprocessedTransactionsAPI.GetUnprocessedTransactions(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UnprocessedTransactionsAPI.GetUnprocessedTransactions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetUnprocessedTransactions`: []UnprocessedTransaction
	fmt.Fprintf(os.Stdout, "Response from `UnprocessedTransactionsAPI.GetUnprocessedTransactions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetUnprocessedTransactionsRequest struct via the builder pattern


### Return type

[**[]UnprocessedTransaction**](UnprocessedTransaction.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

