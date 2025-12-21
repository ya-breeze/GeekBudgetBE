# \BudgetItemsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateBudgetItem**](BudgetItemsAPI.md#CreateBudgetItem) | **Post** /v1/budgetItems | create new budgetItem
[**DeleteBudgetItem**](BudgetItemsAPI.md#DeleteBudgetItem) | **Delete** /v1/budgetItems/{id} | delete budgetItem
[**GetBudgetItem**](BudgetItemsAPI.md#GetBudgetItem) | **Get** /v1/budgetItems/{id} | get budgetItem
[**GetBudgetItems**](BudgetItemsAPI.md#GetBudgetItems) | **Get** /v1/budgetItems | get all budgetItems
[**GetBudgetStatus**](BudgetItemsAPI.md#GetBudgetStatus) | **Get** /v1/budgets/status | get budget status with rollover
[**UpdateBudgetItem**](BudgetItemsAPI.md#UpdateBudgetItem) | **Put** /v1/budgetItems/{id} | update budgetItem



## CreateBudgetItem

> BudgetItem CreateBudgetItem(ctx).BudgetItemNoID(budgetItemNoID).Execute()

create new budgetItem

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
	budgetItemNoID := *openapiclient.NewBudgetItemNoID(time.Now(), "AccountId_example", float64(123)) // BudgetItemNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetItemsAPI.CreateBudgetItem(context.Background()).BudgetItemNoID(budgetItemNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.CreateBudgetItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateBudgetItem`: BudgetItem
	fmt.Fprintf(os.Stdout, "Response from `BudgetItemsAPI.CreateBudgetItem`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateBudgetItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **budgetItemNoID** | [**BudgetItemNoID**](BudgetItemNoID.md) |  | 

### Return type

[**BudgetItem**](BudgetItem.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteBudgetItem

> DeleteBudgetItem(ctx, id).Execute()

delete budgetItem

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the budgetItem

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.BudgetItemsAPI.DeleteBudgetItem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.DeleteBudgetItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the budgetItem | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteBudgetItemRequest struct via the builder pattern


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


## GetBudgetItem

> BudgetItem GetBudgetItem(ctx, id).Execute()

get budgetItem

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the budgetItem

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetItemsAPI.GetBudgetItem(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.GetBudgetItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBudgetItem`: BudgetItem
	fmt.Fprintf(os.Stdout, "Response from `BudgetItemsAPI.GetBudgetItem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the budgetItem | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetBudgetItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BudgetItem**](BudgetItem.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBudgetItems

> []BudgetItem GetBudgetItems(ctx).Execute()

get all budgetItems

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
	resp, r, err := apiClient.BudgetItemsAPI.GetBudgetItems(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.GetBudgetItems``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBudgetItems`: []BudgetItem
	fmt.Fprintf(os.Stdout, "Response from `BudgetItemsAPI.GetBudgetItems`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetBudgetItemsRequest struct via the builder pattern


### Return type

[**[]BudgetItem**](BudgetItem.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBudgetStatus

> []BudgetStatus GetBudgetStatus(ctx).From(from).To(to).OutputCurrencyId(outputCurrencyId).IncludeHidden(includeHidden).Execute()

get budget status with rollover

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
	from := time.Now() // time.Time | Start date (inclusive) (optional)
	to := time.Now() // time.Time | End date (exclusive) (optional)
	outputCurrencyId := "123e4567-e89b-12d3-a456-426614174000" // string | Converts all amounts to this currency (optional)
	includeHidden := true // bool | If true, include hidden accounts (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetItemsAPI.GetBudgetStatus(context.Background()).From(from).To(to).OutputCurrencyId(outputCurrencyId).IncludeHidden(includeHidden).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.GetBudgetStatus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBudgetStatus`: []BudgetStatus
	fmt.Fprintf(os.Stdout, "Response from `BudgetItemsAPI.GetBudgetStatus`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetBudgetStatusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **from** | **time.Time** | Start date (inclusive) | 
 **to** | **time.Time** | End date (exclusive) | 
 **outputCurrencyId** | **string** | Converts all amounts to this currency | 
 **includeHidden** | **bool** | If true, include hidden accounts | [default to false]

### Return type

[**[]BudgetStatus**](BudgetStatus.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBudgetItem

> BudgetItem UpdateBudgetItem(ctx, id).BudgetItemNoID(budgetItemNoID).Execute()

update budgetItem

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the budgetItem
	budgetItemNoID := *openapiclient.NewBudgetItemNoID(time.Now(), "AccountId_example", float64(123)) // BudgetItemNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetItemsAPI.UpdateBudgetItem(context.Background(), id).BudgetItemNoID(budgetItemNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetItemsAPI.UpdateBudgetItem``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBudgetItem`: BudgetItem
	fmt.Fprintf(os.Stdout, "Response from `BudgetItemsAPI.UpdateBudgetItem`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the budgetItem | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBudgetItemRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **budgetItemNoID** | [**BudgetItemNoID**](BudgetItemNoID.md) |  | 

### Return type

[**BudgetItem**](BudgetItem.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

