# AccountNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Type** | **string** |  | 
**BankInfo** | Pointer to [**BankAccountInfo**](BankAccountInfo.md) |  | [optional] 
**ShowInDashboardSummary** | Pointer to **bool** | If true, show this account in dashboard summary. | [optional] [default to true]
**Image** | Pointer to **string** | ID of the account image | [optional] 

## Methods

### NewAccountNoID

`func NewAccountNoID(name string, type_ string, ) *AccountNoID`

NewAccountNoID instantiates a new AccountNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountNoIDWithDefaults

`func NewAccountNoIDWithDefaults() *AccountNoID`

NewAccountNoIDWithDefaults instantiates a new AccountNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AccountNoID) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AccountNoID) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AccountNoID) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *AccountNoID) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *AccountNoID) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *AccountNoID) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *AccountNoID) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetType

`func (o *AccountNoID) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *AccountNoID) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *AccountNoID) SetType(v string)`

SetType sets Type field to given value.


### GetBankInfo

`func (o *AccountNoID) GetBankInfo() BankAccountInfo`

GetBankInfo returns the BankInfo field if non-nil, zero value otherwise.

### GetBankInfoOk

`func (o *AccountNoID) GetBankInfoOk() (*BankAccountInfo, bool)`

GetBankInfoOk returns a tuple with the BankInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankInfo

`func (o *AccountNoID) SetBankInfo(v BankAccountInfo)`

SetBankInfo sets BankInfo field to given value.

### HasBankInfo

`func (o *AccountNoID) HasBankInfo() bool`

HasBankInfo returns a boolean if a field has been set.

### GetShowInDashboardSummary

`func (o *AccountNoID) GetShowInDashboardSummary() bool`

GetShowInDashboardSummary returns the ShowInDashboardSummary field if non-nil, zero value otherwise.

### GetShowInDashboardSummaryOk

`func (o *AccountNoID) GetShowInDashboardSummaryOk() (*bool, bool)`

GetShowInDashboardSummaryOk returns a tuple with the ShowInDashboardSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShowInDashboardSummary

`func (o *AccountNoID) SetShowInDashboardSummary(v bool)`

SetShowInDashboardSummary sets ShowInDashboardSummary field to given value.

### HasShowInDashboardSummary

`func (o *AccountNoID) HasShowInDashboardSummary() bool`

HasShowInDashboardSummary returns a boolean if a field has been set.

### GetImage

`func (o *AccountNoID) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *AccountNoID) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *AccountNoID) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *AccountNoID) HasImage() bool`

HasImage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


