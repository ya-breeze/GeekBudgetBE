# \ReconciliationAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AnalyzeDisbalance**](ReconciliationAPI.md#AnalyzeDisbalance) | **Post** /v1/accounts/{id}/analyze-disbalance | find transactions that might explain the disbalance
[**EnableAccountReconciliation**](ReconciliationAPI.md#EnableAccountReconciliation) | **Post** /v1/accounts/{id}/enable-reconciliation | enable manual reconciliation for accounts without bank importer
[**GetReconciliationHistory**](ReconciliationAPI.md#GetReconciliationHistory) | **Get** /v1/accounts/{id}/reconciliation-history | return all reconciliation records for an account+currency pair
[**GetReconciliationStatus**](ReconciliationAPI.md#GetReconciliationStatus) | **Get** /v1/reconciliation/status | get reconciliation status for all asset accounts
[**GetTransactionsSinceReconciliation**](ReconciliationAPI.md#GetTransactionsSinceReconciliation) | **Get** /v1/accounts/{id}/transactions-since-reconciliation | return transactions since last reconciliation
[**ReconcileAccount**](ReconciliationAPI.md#ReconcileAccount) | **Post** /v1/accounts/{id}/reconcile | manually mark an account as reconciled



## AnalyzeDisbalance

> DisbalanceAnalysis AnalyzeDisbalance(ctx, id).AnalyzeDisbalanceRequest(analyzeDisbalanceRequest).Execute()

find transactions that might explain the disbalance

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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	analyzeDisbalanceRequest := *openapiclient.NewAnalyzeDisbalanceRequest("CurrencyId_example", "TODO") // AnalyzeDisbalanceRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReconciliationAPI.AnalyzeDisbalance(context.Background(), id).AnalyzeDisbalanceRequest(analyzeDisbalanceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.AnalyzeDisbalance``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AnalyzeDisbalance`: DisbalanceAnalysis
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.AnalyzeDisbalance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAnalyzeDisbalanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **analyzeDisbalanceRequest** | [**AnalyzeDisbalanceRequest**](AnalyzeDisbalanceRequest.md) |  | 

### Return type

[**DisbalanceAnalysis**](DisbalanceAnalysis.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EnableAccountReconciliation

> Reconciliation EnableAccountReconciliation(ctx, id).EnableReconciliationRequest(enableReconciliationRequest).Execute()

enable manual reconciliation for accounts without bank importer

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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	enableReconciliationRequest := *openapiclient.NewEnableReconciliationRequest("CurrencyId_example", "TODO") // EnableReconciliationRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReconciliationAPI.EnableAccountReconciliation(context.Background(), id).EnableReconciliationRequest(enableReconciliationRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.EnableAccountReconciliation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EnableAccountReconciliation`: Reconciliation
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.EnableAccountReconciliation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiEnableAccountReconciliationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **enableReconciliationRequest** | [**EnableReconciliationRequest**](EnableReconciliationRequest.md) |  | 

### Return type

[**Reconciliation**](Reconciliation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetReconciliationHistory

> []Reconciliation GetReconciliationHistory(ctx, id).CurrencyId(currencyId).Execute()

return all reconciliation records for an account+currency pair

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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	currencyId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReconciliationAPI.GetReconciliationHistory(context.Background(), id).CurrencyId(currencyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.GetReconciliationHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetReconciliationHistory`: []Reconciliation
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.GetReconciliationHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetReconciliationHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **currencyId** | **string** |  | 

### Return type

[**[]Reconciliation**](Reconciliation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetReconciliationStatus

> []ReconciliationStatus GetReconciliationStatus(ctx).Execute()

get reconciliation status for all asset accounts

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
	resp, r, err := apiClient.ReconciliationAPI.GetReconciliationStatus(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.GetReconciliationStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetReconciliationStatus`: []ReconciliationStatus
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.GetReconciliationStatus`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetReconciliationStatusRequest struct via the builder pattern


### Return type

[**[]ReconciliationStatus**](ReconciliationStatus.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTransactionsSinceReconciliation

> []Transaction GetTransactionsSinceReconciliation(ctx, id).CurrencyId(currencyId).Execute()

return transactions since last reconciliation

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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	currencyId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReconciliationAPI.GetTransactionsSinceReconciliation(context.Background(), id).CurrencyId(currencyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.GetTransactionsSinceReconciliation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTransactionsSinceReconciliation`: []Transaction
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.GetTransactionsSinceReconciliation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTransactionsSinceReconciliationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **currencyId** | **string** |  | 

### Return type

[**[]Transaction**](Transaction.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReconcileAccount

> Reconciliation ReconcileAccount(ctx, id).ReconcileAccountRequest(reconcileAccountRequest).Execute()

manually mark an account as reconciled

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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	reconcileAccountRequest := *openapiclient.NewReconcileAccountRequest("CurrencyId_example") // ReconcileAccountRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReconciliationAPI.ReconcileAccount(context.Background(), id).ReconcileAccountRequest(reconcileAccountRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReconciliationAPI.ReconcileAccount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ReconcileAccount`: Reconciliation
	fmt.Fprintf(os.Stdout, "Response from `ReconciliationAPI.ReconcileAccount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiReconcileAccountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **reconcileAccountRequest** | [**ReconcileAccountRequest**](ReconcileAccountRequest.md) |  | 

### Return type

[**Reconciliation**](Reconciliation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

