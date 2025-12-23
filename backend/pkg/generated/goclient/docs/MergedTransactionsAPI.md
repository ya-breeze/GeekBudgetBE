# \MergedTransactionsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMergedTransactions**](MergedTransactionsAPI.md#GetMergedTransactions) | **Get** /v1/mergedTransactions | get all merged (deduplicated) transactions
[**UnmergeMergedTransaction**](MergedTransactionsAPI.md#UnmergeMergedTransaction) | **Post** /v1/mergedTransactions/{id}/unmerge | restore a merged transaction



## GetMergedTransactions

> []MergedTransaction GetMergedTransactions(ctx).Execute()

get all merged (deduplicated) transactions

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
	resp, r, err := apiClient.MergedTransactionsAPI.GetMergedTransactions(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MergedTransactionsAPI.GetMergedTransactions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMergedTransactions`: []MergedTransaction
	fmt.Fprintf(os.Stdout, "Response from `MergedTransactionsAPI.GetMergedTransactions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetMergedTransactionsRequest struct via the builder pattern


### Return type

[**[]MergedTransaction**](MergedTransaction.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UnmergeMergedTransaction

> UnmergeMergedTransaction(ctx, id).Execute()

restore a merged transaction

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the merged transaction to restore

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MergedTransactionsAPI.UnmergeMergedTransaction(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MergedTransactionsAPI.UnmergeMergedTransaction``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the merged transaction to restore | 

### Other Parameters

Other parameters are passed through a pointer to a apiUnmergeMergedTransactionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

