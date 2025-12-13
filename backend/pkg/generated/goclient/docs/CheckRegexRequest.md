# CheckRegexRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Regex** | **string** |  | 
**TestString** | **string** |  | 

## Methods

### NewCheckRegexRequest

`func NewCheckRegexRequest(regex string, testString string, ) *CheckRegexRequest`

NewCheckRegexRequest instantiates a new CheckRegexRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCheckRegexRequestWithDefaults

`func NewCheckRegexRequestWithDefaults() *CheckRegexRequest`

NewCheckRegexRequestWithDefaults instantiates a new CheckRegexRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegex

`func (o *CheckRegexRequest) GetRegex() string`

GetRegex returns the Regex field if non-nil, zero value otherwise.

### GetRegexOk

`func (o *CheckRegexRequest) GetRegexOk() (*string, bool)`

GetRegexOk returns a tuple with the Regex field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegex

`func (o *CheckRegexRequest) SetRegex(v string)`

SetRegex sets Regex field to given value.


### GetTestString

`func (o *CheckRegexRequest) GetTestString() string`

GetTestString returns the TestString field if non-nil, zero value otherwise.

### GetTestStringOk

`func (o *CheckRegexRequest) GetTestStringOk() (*string, bool)`

GetTestStringOk returns a tuple with the TestString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTestString

`func (o *CheckRegexRequest) SetTestString(v string)`

SetTestString sets TestString field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


