# UpdateMatcher200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Matcher** | [**Matcher**](Matcher.md) |  | 
**AutoProcessedIds** | Pointer to **[]string** | IDs of unprocessed transactions that were automatically matched because of this matcher update | [optional] 

## Methods

### NewUpdateMatcher200Response

`func NewUpdateMatcher200Response(matcher Matcher, ) *UpdateMatcher200Response`

NewUpdateMatcher200Response instantiates a new UpdateMatcher200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateMatcher200ResponseWithDefaults

`func NewUpdateMatcher200ResponseWithDefaults() *UpdateMatcher200Response`

NewUpdateMatcher200ResponseWithDefaults instantiates a new UpdateMatcher200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMatcher

`func (o *UpdateMatcher200Response) GetMatcher() Matcher`

GetMatcher returns the Matcher field if non-nil, zero value otherwise.

### GetMatcherOk

`func (o *UpdateMatcher200Response) GetMatcherOk() (*Matcher, bool)`

GetMatcherOk returns a tuple with the Matcher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatcher

`func (o *UpdateMatcher200Response) SetMatcher(v Matcher)`

SetMatcher sets Matcher field to given value.


### GetAutoProcessedIds

`func (o *UpdateMatcher200Response) GetAutoProcessedIds() []string`

GetAutoProcessedIds returns the AutoProcessedIds field if non-nil, zero value otherwise.

### GetAutoProcessedIdsOk

`func (o *UpdateMatcher200Response) GetAutoProcessedIdsOk() (*[]string, bool)`

GetAutoProcessedIdsOk returns a tuple with the AutoProcessedIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoProcessedIds

`func (o *UpdateMatcher200Response) SetAutoProcessedIds(v []string)`

SetAutoProcessedIds sets AutoProcessedIds field to given value.

### HasAutoProcessedIds

`func (o *UpdateMatcher200Response) HasAutoProcessedIds() bool`

HasAutoProcessedIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


