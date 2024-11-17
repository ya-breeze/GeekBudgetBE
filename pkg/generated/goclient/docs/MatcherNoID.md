# MatcherNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**OutputDescription** | **string** |  | 
**OutputAccountId** | **string** |  | 
**OutputTags** | Pointer to **[]string** |  | [optional] 
**CurrencyRegExp** | Pointer to **string** |  | [optional] 
**PartnerNameRegExp** | Pointer to **string** |  | [optional] 
**PartnerAccountNumber** | Pointer to **string** |  | [optional] 
**DescriptionRegExp** | Pointer to **string** |  | [optional] 
**ExtraRegExp** | Pointer to **string** |  | [optional] 

## Methods

### NewMatcherNoID

`func NewMatcherNoID(name string, outputDescription string, outputAccountId string, ) *MatcherNoID`

NewMatcherNoID instantiates a new MatcherNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMatcherNoIDWithDefaults

`func NewMatcherNoIDWithDefaults() *MatcherNoID`

NewMatcherNoIDWithDefaults instantiates a new MatcherNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *MatcherNoID) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *MatcherNoID) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *MatcherNoID) SetName(v string)`

SetName sets Name field to given value.


### GetOutputDescription

`func (o *MatcherNoID) GetOutputDescription() string`

GetOutputDescription returns the OutputDescription field if non-nil, zero value otherwise.

### GetOutputDescriptionOk

`func (o *MatcherNoID) GetOutputDescriptionOk() (*string, bool)`

GetOutputDescriptionOk returns a tuple with the OutputDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputDescription

`func (o *MatcherNoID) SetOutputDescription(v string)`

SetOutputDescription sets OutputDescription field to given value.


### GetOutputAccountId

`func (o *MatcherNoID) GetOutputAccountId() string`

GetOutputAccountId returns the OutputAccountId field if non-nil, zero value otherwise.

### GetOutputAccountIdOk

`func (o *MatcherNoID) GetOutputAccountIdOk() (*string, bool)`

GetOutputAccountIdOk returns a tuple with the OutputAccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputAccountId

`func (o *MatcherNoID) SetOutputAccountId(v string)`

SetOutputAccountId sets OutputAccountId field to given value.


### GetOutputTags

`func (o *MatcherNoID) GetOutputTags() []string`

GetOutputTags returns the OutputTags field if non-nil, zero value otherwise.

### GetOutputTagsOk

`func (o *MatcherNoID) GetOutputTagsOk() (*[]string, bool)`

GetOutputTagsOk returns a tuple with the OutputTags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputTags

`func (o *MatcherNoID) SetOutputTags(v []string)`

SetOutputTags sets OutputTags field to given value.

### HasOutputTags

`func (o *MatcherNoID) HasOutputTags() bool`

HasOutputTags returns a boolean if a field has been set.

### GetCurrencyRegExp

`func (o *MatcherNoID) GetCurrencyRegExp() string`

GetCurrencyRegExp returns the CurrencyRegExp field if non-nil, zero value otherwise.

### GetCurrencyRegExpOk

`func (o *MatcherNoID) GetCurrencyRegExpOk() (*string, bool)`

GetCurrencyRegExpOk returns a tuple with the CurrencyRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyRegExp

`func (o *MatcherNoID) SetCurrencyRegExp(v string)`

SetCurrencyRegExp sets CurrencyRegExp field to given value.

### HasCurrencyRegExp

`func (o *MatcherNoID) HasCurrencyRegExp() bool`

HasCurrencyRegExp returns a boolean if a field has been set.

### GetPartnerNameRegExp

`func (o *MatcherNoID) GetPartnerNameRegExp() string`

GetPartnerNameRegExp returns the PartnerNameRegExp field if non-nil, zero value otherwise.

### GetPartnerNameRegExpOk

`func (o *MatcherNoID) GetPartnerNameRegExpOk() (*string, bool)`

GetPartnerNameRegExpOk returns a tuple with the PartnerNameRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerNameRegExp

`func (o *MatcherNoID) SetPartnerNameRegExp(v string)`

SetPartnerNameRegExp sets PartnerNameRegExp field to given value.

### HasPartnerNameRegExp

`func (o *MatcherNoID) HasPartnerNameRegExp() bool`

HasPartnerNameRegExp returns a boolean if a field has been set.

### GetPartnerAccountNumber

`func (o *MatcherNoID) GetPartnerAccountNumber() string`

GetPartnerAccountNumber returns the PartnerAccountNumber field if non-nil, zero value otherwise.

### GetPartnerAccountNumberOk

`func (o *MatcherNoID) GetPartnerAccountNumberOk() (*string, bool)`

GetPartnerAccountNumberOk returns a tuple with the PartnerAccountNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerAccountNumber

`func (o *MatcherNoID) SetPartnerAccountNumber(v string)`

SetPartnerAccountNumber sets PartnerAccountNumber field to given value.

### HasPartnerAccountNumber

`func (o *MatcherNoID) HasPartnerAccountNumber() bool`

HasPartnerAccountNumber returns a boolean if a field has been set.

### GetDescriptionRegExp

`func (o *MatcherNoID) GetDescriptionRegExp() string`

GetDescriptionRegExp returns the DescriptionRegExp field if non-nil, zero value otherwise.

### GetDescriptionRegExpOk

`func (o *MatcherNoID) GetDescriptionRegExpOk() (*string, bool)`

GetDescriptionRegExpOk returns a tuple with the DescriptionRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescriptionRegExp

`func (o *MatcherNoID) SetDescriptionRegExp(v string)`

SetDescriptionRegExp sets DescriptionRegExp field to given value.

### HasDescriptionRegExp

`func (o *MatcherNoID) HasDescriptionRegExp() bool`

HasDescriptionRegExp returns a boolean if a field has been set.

### GetExtraRegExp

`func (o *MatcherNoID) GetExtraRegExp() string`

GetExtraRegExp returns the ExtraRegExp field if non-nil, zero value otherwise.

### GetExtraRegExpOk

`func (o *MatcherNoID) GetExtraRegExpOk() (*string, bool)`

GetExtraRegExpOk returns a tuple with the ExtraRegExp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtraRegExp

`func (o *MatcherNoID) SetExtraRegExp(v string)`

SetExtraRegExp sets ExtraRegExp field to given value.

### HasExtraRegExp

`func (o *MatcherNoID) HasExtraRegExp() bool`

HasExtraRegExp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


