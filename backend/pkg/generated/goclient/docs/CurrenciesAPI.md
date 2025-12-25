# \CurrenciesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateCurrency**](CurrenciesAPI.md#CreateCurrency) | **Post** /v1/currencies | create new currency
[**DeleteCurrency**](CurrenciesAPI.md#DeleteCurrency) | **Delete** /v1/currencies/{id} | delete currency
[**GetCurrencies**](CurrenciesAPI.md#GetCurrencies) | **Get** /v1/currencies | get all currencies
[**UpdateCurrency**](CurrenciesAPI.md#UpdateCurrency) | **Put** /v1/currencies/{id} | update currency



## CreateCurrency

> Currency CreateCurrency(ctx).CurrencyNoID(currencyNoID).Execute()

create new currency

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
	currencyNoID := *openapiclient.NewCurrencyNoID("Name_example") // CurrencyNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CurrenciesAPI.CreateCurrency(context.Background()).CurrencyNoID(currencyNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CurrenciesAPI.CreateCurrency``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateCurrency`: Currency
	fmt.Fprintf(os.Stdout, "Response from `CurrenciesAPI.CreateCurrency`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateCurrencyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **currencyNoID** | [**CurrencyNoID**](CurrencyNoID.md) |  | 

### Return type

[**Currency**](Currency.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCurrency

> DeleteCurrency(ctx, id).ReplaceWithCurrencyId(replaceWithCurrencyId).Execute()

delete currency

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the currency
	replaceWithCurrencyId := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the currency which should be used instead of the deleted one (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CurrenciesAPI.DeleteCurrency(context.Background(), id).ReplaceWithCurrencyId(replaceWithCurrencyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CurrenciesAPI.DeleteCurrency``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the currency | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCurrencyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **replaceWithCurrencyId** | **string** | ID of the currency which should be used instead of the deleted one | 

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


## GetCurrencies

> []Currency GetCurrencies(ctx).Execute()

get all currencies

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
	resp, r, err := apiClient.CurrenciesAPI.GetCurrencies(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CurrenciesAPI.GetCurrencies``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCurrencies`: []Currency
	fmt.Fprintf(os.Stdout, "Response from `CurrenciesAPI.GetCurrencies`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetCurrenciesRequest struct via the builder pattern


### Return type

[**[]Currency**](Currency.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCurrency

> Currency UpdateCurrency(ctx, id).CurrencyNoID(currencyNoID).Execute()

update currency

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the currency
	currencyNoID := *openapiclient.NewCurrencyNoID("Name_example") // CurrencyNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CurrenciesAPI.UpdateCurrency(context.Background(), id).CurrencyNoID(currencyNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CurrenciesAPI.UpdateCurrency``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCurrency`: Currency
	fmt.Fprintf(os.Stdout, "Response from `CurrenciesAPI.UpdateCurrency`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the currency | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCurrencyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **currencyNoID** | [**CurrencyNoID**](CurrencyNoID.md) |  | 

### Return type

[**Currency**](Currency.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

