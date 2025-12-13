# \ImportAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CallImport**](ImportAPI.md#CallImport) | **Post** /v1/import | Upload and import full user&#39;s data



## CallImport

> CallImport(ctx).WholeUserData(wholeUserData).Execute()

Upload and import full user's data

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
	wholeUserData := *openapiclient.NewWholeUserData() // WholeUserData |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImportAPI.CallImport(context.Background()).WholeUserData(wholeUserData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImportAPI.CallImport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCallImportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **wholeUserData** | [**WholeUserData**](WholeUserData.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json:
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

