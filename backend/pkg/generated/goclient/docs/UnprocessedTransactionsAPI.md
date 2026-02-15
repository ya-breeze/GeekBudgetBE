# \UnprocessedTransactionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConvertUnprocessedTransaction**](UnprocessedTransactionsAPI.md#ConvertUnprocessedTransaction) | **Post** /v1/unprocessedTransactions/{id}/convert | convert unprocessed transactions into normal transaction
[**GetUnprocessedTransaction**](UnprocessedTransactionsAPI.md#GetUnprocessedTransaction) | **Get** /v1/unprocessedTransactions/{id} | get unprocessed transaction
[**GetUnprocessedTransactions**](UnprocessedTransactionsAPI.md#GetUnprocessedTransactions) | **Get** /v1/unprocessedTransactions | get all unprocessed transactions



## ConvertUnprocessedTransaction

> ConvertUnprocessedTransaction200Response ConvertUnprocessedTransaction(ctx, id).TransactionNoID(transactionNoID).MatcherId(matcherId).Execute()

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
	transactionNoID := *openapiclient.NewTransactionNoID(time.Now(), []openapiclient.Movement{*openapiclient.NewMovement("TODO", "CurrencyId_example")}) // TransactionNoID | 
	matcherId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | ID of the matcher used for this conversion (if any) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UnprocessedTransactionsAPI.ConvertUnprocessedTransaction(context.Background(), id).TransactionNoID(transactionNoID).MatcherId(matcherId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UnprocessedTransactionsAPI.ConvertUnprocessedTransaction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ConvertUnprocessedTransaction`: ConvertUnprocessedTransaction200Response
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
 **matcherId** | **string** | ID of the matcher used for this conversion (if any) | 

### Return type

[**ConvertUnprocessedTransaction200Response**](ConvertUnprocessedTransaction200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUnprocessedTransaction

> UnprocessedTransaction GetUnprocessedTransaction(ctx, id).Execute()

get unprocessed transaction

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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UnprocessedTransactionsAPI.GetUnprocessedTransaction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UnprocessedTransactionsAPI.GetUnprocessedTransaction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetUnprocessedTransaction`: UnprocessedTransaction
	fmt.Fprintf(os.Stdout, "Response from `UnprocessedTransactionsAPI.GetUnprocessedTransaction`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUnprocessedTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UnprocessedTransaction**](UnprocessedTransaction.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

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

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

