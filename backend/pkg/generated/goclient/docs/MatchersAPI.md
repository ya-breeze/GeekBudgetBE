# \MatchersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CheckMatcher**](MatchersAPI.md#CheckMatcher) | **Post** /v1/matchers/check | check if passed matcher matches given transaction
[**CheckRegex**](MatchersAPI.md#CheckRegex) | **Post** /v1/matchers/check-regex | check if regex is valid and matches string (using backend&#39;s regex engine)
[**CreateMatcher**](MatchersAPI.md#CreateMatcher) | **Post** /v1/matchers | create new matcher
[**DeleteMatcher**](MatchersAPI.md#DeleteMatcher) | **Delete** /v1/matchers/{id} | delete matcher
[**DeleteMatcherImage**](MatchersAPI.md#DeleteMatcherImage) | **Delete** /v1/matchers/{id}/image | delete matcher image
[**GetMatcher**](MatchersAPI.md#GetMatcher) | **Get** /v1/matchers/{id} | get matcher
[**GetMatchers**](MatchersAPI.md#GetMatchers) | **Get** /v1/matchers | get all matchers
[**UpdateMatcher**](MatchersAPI.md#UpdateMatcher) | **Put** /v1/matchers/{id} | update matcher
[**UploadMatcherImage**](MatchersAPI.md#UploadMatcherImage) | **Post** /v1/matchers/{id}/image | Upload matcher image



## CheckMatcher

> CheckMatcher200Response CheckMatcher(ctx).CheckMatcherRequest(checkMatcherRequest).Execute()

check if passed matcher matches given transaction

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
	checkMatcherRequest := *openapiclient.NewCheckMatcherRequest(*openapiclient.NewMatcherNoID("OutputAccountId_example"), *openapiclient.NewTransactionNoID(time.Now(), []openapiclient.Movement{*openapiclient.NewMovement("TODO", "CurrencyId_example")})) // CheckMatcherRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.CheckMatcher(context.Background()).CheckMatcherRequest(checkMatcherRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.CheckMatcher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CheckMatcher`: CheckMatcher200Response
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.CheckMatcher`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCheckMatcherRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **checkMatcherRequest** | [**CheckMatcherRequest**](CheckMatcherRequest.md) |  | 

### Return type

[**CheckMatcher200Response**](CheckMatcher200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CheckRegex

> CheckRegex200Response CheckRegex(ctx).CheckRegexRequest(checkRegexRequest).Execute()

check if regex is valid and matches string (using backend's regex engine)

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
	checkRegexRequest := *openapiclient.NewCheckRegexRequest("Regex_example", "TestString_example") // CheckRegexRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.CheckRegex(context.Background()).CheckRegexRequest(checkRegexRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.CheckRegex``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CheckRegex`: CheckRegex200Response
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.CheckRegex`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCheckRegexRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **checkRegexRequest** | [**CheckRegexRequest**](CheckRegexRequest.md) |  | 

### Return type

[**CheckRegex200Response**](CheckRegex200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateMatcher

> Matcher CreateMatcher(ctx).MatcherNoID(matcherNoID).Execute()

create new matcher

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
	matcherNoID := *openapiclient.NewMatcherNoID("OutputAccountId_example") // MatcherNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.CreateMatcher(context.Background()).MatcherNoID(matcherNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.CreateMatcher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateMatcher`: Matcher
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.CreateMatcher`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateMatcherRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **matcherNoID** | [**MatcherNoID**](MatcherNoID.md) |  | 

### Return type

[**Matcher**](Matcher.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteMatcher

> DeleteMatcher(ctx, id).Execute()

delete matcher

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the matcher

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MatchersAPI.DeleteMatcher(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.DeleteMatcher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the matcher | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteMatcherRequest struct via the builder pattern


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


## DeleteMatcherImage

> Matcher DeleteMatcherImage(ctx, id).Execute()

delete matcher image

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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.DeleteMatcherImage(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.DeleteMatcherImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteMatcherImage`: Matcher
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.DeleteMatcherImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteMatcherImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Matcher**](Matcher.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMatcher

> Matcher GetMatcher(ctx, id).Execute()

get matcher

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the matcher

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.GetMatcher(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.GetMatcher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMatcher`: Matcher
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.GetMatcher`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the matcher | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetMatcherRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Matcher**](Matcher.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMatchers

> []Matcher GetMatchers(ctx).Execute()

get all matchers

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
	resp, r, err := apiClient.MatchersAPI.GetMatchers(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.GetMatchers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMatchers`: []Matcher
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.GetMatchers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetMatchersRequest struct via the builder pattern


### Return type

[**[]Matcher**](Matcher.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateMatcher

> UpdateMatcher200Response UpdateMatcher(ctx, id).MatcherNoID(matcherNoID).Execute()

update matcher

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
	id := "123e4567-e89b-12d3-a456-426614174000" // string | ID of the matcher
	matcherNoID := *openapiclient.NewMatcherNoID("OutputAccountId_example") // MatcherNoID | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.UpdateMatcher(context.Background(), id).MatcherNoID(matcherNoID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.UpdateMatcher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateMatcher`: UpdateMatcher200Response
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.UpdateMatcher`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the matcher | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateMatcherRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **matcherNoID** | [**MatcherNoID**](MatcherNoID.md) |  | 

### Return type

[**UpdateMatcher200Response**](UpdateMatcher200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UploadMatcherImage

> Matcher UploadMatcherImage(ctx, id).File(file).Execute()

Upload matcher image

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
	file := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MatchersAPI.UploadMatcherImage(context.Background(), id).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MatchersAPI.UploadMatcherImage``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UploadMatcherImage`: Matcher
	fmt.Fprintf(os.Stdout, "Response from `MatchersAPI.UploadMatcherImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUploadMatcherImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **file** | ***os.File** |  | 

### Return type

[**Matcher**](Matcher.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

