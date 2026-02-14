# \AuditLogsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAuditLogs**](AuditLogsAPI.md#GetAuditLogs) | **Get** /v1/auditLogs | get audit logs



## GetAuditLogs

> []AuditLog GetAuditLogs(ctx).EntityType(entityType).EntityId(entityId).UserId(userId).DateFrom(dateFrom).DateTo(dateTo).Limit(limit).Offset(offset).Execute()

get audit logs

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
	entityType := "entityType_example" // string | Filter by entity type (optional)
	entityId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by entity ID (optional)
	userId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Filter by user ID (optional)
	dateFrom := time.Now() // time.Time | Filter by date from (optional)
	dateTo := time.Now() // time.Time | Filter by date to (optional)
	limit := int32(56) // int32 | Limit number of results (optional) (default to 100)
	offset := int32(56) // int32 | Offset results (optional) (default to 0)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuditLogsAPI.GetAuditLogs(context.Background()).EntityType(entityType).EntityId(entityId).UserId(userId).DateFrom(dateFrom).DateTo(dateTo).Limit(limit).Offset(offset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuditLogsAPI.GetAuditLogs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAuditLogs`: []AuditLog
	fmt.Fprintf(os.Stdout, "Response from `AuditLogsAPI.GetAuditLogs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAuditLogsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **entityType** | **string** | Filter by entity type | 
 **entityId** | **string** | Filter by entity ID | 
 **userId** | **string** | Filter by user ID | 
 **dateFrom** | **time.Time** | Filter by date from | 
 **dateTo** | **time.Time** | Filter by date to | 
 **limit** | **int32** | Limit number of results | [default to 100]
 **offset** | **int32** | Offset results | [default to 0]

### Return type

[**[]AuditLog**](AuditLog.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

