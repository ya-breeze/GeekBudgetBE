# \AggregationsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetBalances**](AggregationsAPI.md#GetBalances) | **Get** /v1/balances | get balance for filtered transactions
[**GetExpenses**](AggregationsAPI.md#GetExpenses) | **Get** /v1/expenses | get expenses for filtered transactions
[**GetIncomes**](AggregationsAPI.md#GetIncomes) | **Get** /v1/incomes | get incomes for filtered transactions



## GetBalances

> Aggregation GetBalances(ctx).From(from).To(to).OutputCurrencyId(outputCurrencyId).Execute()

get balance for filtered transactions

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
	from := time.Now() // time.Time | Uses transactions from this date (optional)
	to := time.Now() // time.Time | Uses transactions to this date (optional)
	outputCurrencyId := "123e4567-e89b-12d3-a456-426614174000" // string | Converts all transactions to this currency (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AggregationsAPI.GetBalances(context.Background()).From(from).To(to).OutputCurrencyId(outputCurrencyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AggregationsAPI.GetBalances``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBalances`: Aggregation
	fmt.Fprintf(os.Stdout, "Response from `AggregationsAPI.GetBalances`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetBalancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **from** | **time.Time** | Uses transactions from this date | 
 **to** | **time.Time** | Uses transactions to this date | 
 **outputCurrencyId** | **string** | Converts all transactions to this currency | 

### Return type

[**Aggregation**](Aggregation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetExpenses

> Aggregation GetExpenses(ctx).From(from).To(to).OutputCurrencyId(outputCurrencyId).Granularity(granularity).Execute()

get expenses for filtered transactions

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
	from := time.Now() // time.Time | Uses transactions from this date (optional)
	to := time.Now() // time.Time | Uses transactions to this date (optional)
	outputCurrencyId := "outputCurrencyId_example" // string | Converts all transactions to this currency (optional)
	granularity := "granularity_example" // string | Granularity of expenses (month or year) (optional) (default to "month")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AggregationsAPI.GetExpenses(context.Background()).From(from).To(to).OutputCurrencyId(outputCurrencyId).Granularity(granularity).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AggregationsAPI.GetExpenses``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetExpenses`: Aggregation
	fmt.Fprintf(os.Stdout, "Response from `AggregationsAPI.GetExpenses`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetExpensesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **from** | **time.Time** | Uses transactions from this date | 
 **to** | **time.Time** | Uses transactions to this date | 
 **outputCurrencyId** | **string** | Converts all transactions to this currency | 
 **granularity** | **string** | Granularity of expenses (month or year) | [default to &quot;month&quot;]

### Return type

[**Aggregation**](Aggregation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetIncomes

> Aggregation GetIncomes(ctx).From(from).To(to).OutputCurrencyId(outputCurrencyId).Execute()

get incomes for filtered transactions

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
	from := time.Now() // time.Time | Uses transactions from this date (optional)
	to := time.Now() // time.Time | Uses transactions to this date (optional)
	outputCurrencyId := "outputCurrencyId_example" // string | Converts all transactions to this currency (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AggregationsAPI.GetIncomes(context.Background()).From(from).To(to).OutputCurrencyId(outputCurrencyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AggregationsAPI.GetIncomes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetIncomes`: Aggregation
	fmt.Fprintf(os.Stdout, "Response from `AggregationsAPI.GetIncomes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetIncomesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **from** | **time.Time** | Uses transactions from this date | 
 **to** | **time.Time** | Uses transactions to this date | 
 **outputCurrencyId** | **string** | Converts all transactions to this currency | 

### Return type

[**Aggregation**](Aggregation.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

