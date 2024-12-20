# Matcher

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**OutputDescription** | **string** |  | 
**OutputAccountId** | **string** |  | 
**OutputTags** | Pointer to **[]string** |  | [optional] 
**CurrencyRegExp** | Pointer to **string** |  | [optional] 
**PartnerNameRegExp** | Pointer to **string** |  | [optional] 
**PartnerAccountNumberRegExp** | Pointer to **string** |  | [optional] 
**DescriptionRegExp** | Pointer to **string** |  | [optional] 
**ExtraRegExp** | Pointer to **string** |  | [optional] 

## Methods

### NewMatcher

`func NewMatcher(id string, name string, outputDescription string, outputAccountId string, ) *Matcher`

NewMatcher instantiates a new Matcher object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMatcherWithDefaults

`func NewMatcherWithDefaults() *Matcher`

NewMatcherWithDefaults instantiates a new Matcher object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Matcher) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Matcher) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Matcher) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Matcher) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Matcher) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Matcher) SetName(v string)`

SetName sets Name field to given value.


### GetOutputDescription

`func (o *Matcher) GetOutputDescription() string`

GetOutputDescription returns the OutputDescription field if non-nil, zero value otherwise.

### GetOutputDescriptionOk

`func (o *Matcher) GetOutputDescriptionOk() (*string, bool)`

GetOutputDescriptionOk returns a tuple with the OutputDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputDescription

`func (o *Matcher) SetOutputDescription(v string)`

SetOutputDescription sets OutputDescription field to given value.


### GetOutputAccountId

`func (o *Matcher) GetOutputAccountId() string`

GetOutputAccountId returns the OutputAccountId field if non-nil, zero value otherwise.

### GetOutputAccountIdOk

`func (o *Matcher) GetOutputAccountIdOk() (*string, bool)`

GetOutputAccountIdOk returns a tuple with the OutputAccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputAccountId

`func (o *Matcher) SetOutputAccountId(v string)`

SetOutputAccountId sets OutputAccountId field to given value.


### GetOutputTags

`func (o *Matcher) GetOutputTags() []string`

GetOutputTags returns the OutputTags field if non-nil, zero value otherwise.

### GetOutputTagsOk

`func (o *Matcher) GetOutputTagsOk() (*[]string, bool)`

GetOutputTagsOk returns a tuple with the OutputTags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputTags

`func (o *Matcher) SetOutputTags(v []string)`

SetOutputTags sets OutputTags field to given value.

### HasOutputTags

`func (o *Matcher) HasOutputTags() bool`

HasOutputTags returns a boolean if a field has been set.

### GetCurrencyRegExp

`func (o *Matcher) GetCurrencyRegExp() string`

GetCurrencyRegExp returns the CurrencyRegExp field if non-nil, zero value otherwise.

### GetCurrencyRegExpOk

`func (o *Matcher) GetCurrencyRegExpOk() (*string, bool)`

GetCurrencyRegExpOk returns a tuple with the CurrencyRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyRegExp

`func (o *Matcher) SetCurrencyRegExp(v string)`

SetCurrencyRegExp sets CurrencyRegExp field to given value.

### HasCurrencyRegExp

`func (o *Matcher) HasCurrencyRegExp() bool`

HasCurrencyRegExp returns a boolean if a field has been set.

### GetPartnerNameRegExp

`func (o *Matcher) GetPartnerNameRegExp() string`

GetPartnerNameRegExp returns the PartnerNameRegExp field if non-nil, zero value otherwise.

### GetPartnerNameRegExpOk

`func (o *Matcher) GetPartnerNameRegExpOk() (*string, bool)`

GetPartnerNameRegExpOk returns a tuple with the PartnerNameRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerNameRegExp

`func (o *Matcher) SetPartnerNameRegExp(v string)`

SetPartnerNameRegExp sets PartnerNameRegExp field to given value.

### HasPartnerNameRegExp

`func (o *Matcher) HasPartnerNameRegExp() bool`

HasPartnerNameRegExp returns a boolean if a field has been set.

### GetPartnerAccountNumberRegExp

`func (o *Matcher) GetPartnerAccountNumberRegExp() string`

GetPartnerAccountNumberRegExp returns the PartnerAccountNumberRegExp field if non-nil, zero value otherwise.

### GetPartnerAccountNumberRegExpOk

`func (o *Matcher) GetPartnerAccountNumberRegExpOk() (*string, bool)`

GetPartnerAccountNumberRegExpOk returns a tuple with the PartnerAccountNumberRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerAccountNumberRegExp

`func (o *Matcher) SetPartnerAccountNumberRegExp(v string)`

SetPartnerAccountNumberRegExp sets PartnerAccountNumberRegExp field to given value.

### HasPartnerAccountNumberRegExp

`func (o *Matcher) HasPartnerAccountNumberRegExp() bool`

HasPartnerAccountNumberRegExp returns a boolean if a field has been set.

### GetDescriptionRegExp

`func (o *Matcher) GetDescriptionRegExp() string`

GetDescriptionRegExp returns the DescriptionRegExp field if non-nil, zero value otherwise.

### GetDescriptionRegExpOk

`func (o *Matcher) GetDescriptionRegExpOk() (*string, bool)`

GetDescriptionRegExpOk returns a tuple with the DescriptionRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptionRegExp

`func (o *Matcher) SetDescriptionRegExp(v string)`

SetDescriptionRegExp sets DescriptionRegExp field to given value.

### HasDescriptionRegExp

`func (o *Matcher) HasDescriptionRegExp() bool`

HasDescriptionRegExp returns a boolean if a field has been set.

### GetExtraRegExp

`func (o *Matcher) GetExtraRegExp() string`

GetExtraRegExp returns the ExtraRegExp field if non-nil, zero value otherwise.

### GetExtraRegExpOk

`func (o *Matcher) GetExtraRegExpOk() (*string, bool)`

GetExtraRegExpOk returns a tuple with the ExtraRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtraRegExp

`func (o *Matcher) SetExtraRegExp(v string)`

SetExtraRegExp sets ExtraRegExp field to given value.

### HasExtraRegExp

`func (o *Matcher) HasExtraRegExp() bool`

HasExtraRegExp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


