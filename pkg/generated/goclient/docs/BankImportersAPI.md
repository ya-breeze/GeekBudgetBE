# \BankImportersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateBankImporter**](BankImportersAPI.md#CreateBankImporter) | **Post** /v1/bankImporters | create new bank importer
[**DeleteBankImporter**](BankImportersAPI.md#DeleteBankImporter) | **Delete** /v1/bankImporters/{id} | delete bank importer
[**GetBankImporters**](BankImportersAPI.md#GetBankImporters) | **Get** /v1/bankImporters | get all bank importers
[**UpdateBankImporter**](BankImportersAPI.md#UpdateBankImporter) | **Put** /v1/bankImporters/{id} | update bank importer



## CreateBankImporter

> BankImporter CreateBankImporter(ctx).BankImporterNoID(bankImporterNoID).Execute()

create new bank importer

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
	bankImporterNoID := *openapiclient.NewBankImporterNoID() // BankImporterNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BankImportersAPI.CreateBankImporter(context.Background()).BankImporterNoID(bankImporterNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BankImportersAPI.CreateBankImporter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateBankImporter`: BankImporter
	fmt.Fprintf(os.Stdout, "Response from `BankImportersAPI.CreateBankImporter`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateBankImporterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bankImporterNoID** | [**BankImporterNoID**](BankImporterNoID.md) |  | 

### Return type

[**BankImporter**](BankImporter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteBankImporter

> DeleteBankImporter(ctx, id).Execute()

delete bank importer

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the bankimporter

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.BankImportersAPI.DeleteBankImporter(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BankImportersAPI.DeleteBankImporter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the bankimporter | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteBankImporterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


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


## GetBankImporters

> []BankImporter GetBankImporters(ctx).Execute()

get all bank importers

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
	resp, r, err := apiClient.BankImportersAPI.GetBankImporters(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BankImportersAPI.GetBankImporters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBankImporters`: []BankImporter
	fmt.Fprintf(os.Stdout, "Response from `BankImportersAPI.GetBankImporters`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetBankImportersRequest struct via the builder pattern


### Return type

[**[]BankImporter**](BankImporter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBankImporter

> BankImporter UpdateBankImporter(ctx, id).BankImporterNoID(bankImporterNoID).Execute()

update bank importer

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the bank importer
	bankImporterNoID := *openapiclient.NewBankImporterNoID() // BankImporterNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BankImportersAPI.UpdateBankImporter(context.Background(), id).BankImporterNoID(bankImporterNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BankImportersAPI.UpdateBankImporter``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateBankImporter`: BankImporter
	fmt.Fprintf(os.Stdout, "Response from `BankImportersAPI.UpdateBankImporter`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the bank importer | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBankImporterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **bankImporterNoID** | [**BankImporterNoID**](BankImporterNoID.md) |  | 

### Return type

[**BankImporter**](BankImporter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

