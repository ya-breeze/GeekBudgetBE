# Account

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Type** | **string** |  | 
**BankInfo** | Pointer to [**BankAccountInfo**](BankAccountInfo.md) |  | [optional] 
**ShowInDashboardSummary** | **bool** | If true, show this account in dashboard summary. | [default to true]
**HideFromReports** | Pointer to **bool** | If true, this account should be hidden from reports and budget. | [optional] [default to false]
**Image** | Pointer to **string** | ID of the account image | [optional] 
**IgnoreUnprocessedBefore** | Pointer to **time.Time** | If set, unprocessed transactions for this account older than this date will be ignored. | [optional] 
**OpeningDate** | Pointer to **time.Time** | If set, the account is ignored before this date. | [optional] 
**ClosingDate** | Pointer to **time.Time** | If set, the account is ignored after this date. | [optional] 
**ShowInReconciliation** | Pointer to **bool** | If true, this account is shown on the reconciliation page even if it has no bank importer. | [optional] [default to false]

## Methods

### NewAccount

`func NewAccount(id string, name string, type_ string, showInDashboardSummary bool, ) *Account`

NewAccount instantiates a new Account object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountWithDefaults

`func NewAccountWithDefaults() *Account`

NewAccountWithDefaults instantiates a new Account object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Account) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Account) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Account) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Account) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Account) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Account) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *Account) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *Account) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *Account) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *Account) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetType

`func (o *Account) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Account) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Account) SetType(v string)`

SetType sets Type field to given value.


### GetBankInfo

`func (o *Account) GetBankInfo() BankAccountInfo`

GetBankInfo returns the BankInfo field if non-nil, zero value otherwise.

### GetBankInfoOk

`func (o *Account) GetBankInfoOk() (*BankAccountInfo, bool)`

GetBankInfoOk returns a tuple with the BankInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankInfo

`func (o *Account) SetBankInfo(v BankAccountInfo)`

SetBankInfo sets BankInfo field to given value.

### HasBankInfo

`func (o *Account) HasBankInfo() bool`

HasBankInfo returns a boolean if a field has been set.

### GetShowInDashboardSummary

`func (o *Account) GetShowInDashboardSummary() bool`

GetShowInDashboardSummary returns the ShowInDashboardSummary field if non-nil, zero value otherwise.

### GetShowInDashboardSummaryOk

`func (o *Account) GetShowInDashboardSummaryOk() (*bool, bool)`

GetShowInDashboardSummaryOk returns a tuple with the ShowInDashboardSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShowInDashboardSummary

`func (o *Account) SetShowInDashboardSummary(v bool)`

SetShowInDashboardSummary sets ShowInDashboardSummary field to given value.


### GetHideFromReports

`func (o *Account) GetHideFromReports() bool`

GetHideFromReports returns the HideFromReports field if non-nil, zero value otherwise.

### GetHideFromReportsOk

`func (o *Account) GetHideFromReportsOk() (*bool, bool)`

GetHideFromReportsOk returns a tuple with the HideFromReports field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHideFromReports

`func (o *Account) SetHideFromReports(v bool)`

SetHideFromReports sets HideFromReports field to given value.

### HasHideFromReports

`func (o *Account) HasHideFromReports() bool`

HasHideFromReports returns a boolean if a field has been set.

### GetImage

`func (o *Account) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *Account) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *Account) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *Account) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetIgnoreUnprocessedBefore

`func (o *Account) GetIgnoreUnprocessedBefore() time.Time`

GetIgnoreUnprocessedBefore returns the IgnoreUnprocessedBefore field if non-nil, zero value otherwise.

### GetIgnoreUnprocessedBeforeOk

`func (o *Account) GetIgnoreUnprocessedBeforeOk() (*time.Time, bool)`

GetIgnoreUnprocessedBeforeOk returns a tuple with the IgnoreUnprocessedBefore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIgnoreUnprocessedBefore

`func (o *Account) SetIgnoreUnprocessedBefore(v time.Time)`

SetIgnoreUnprocessedBefore sets IgnoreUnprocessedBefore field to given value.

### HasIgnoreUnprocessedBefore

`func (o *Account) HasIgnoreUnprocessedBefore() bool`

HasIgnoreUnprocessedBefore returns a boolean if a field has been set.

### GetOpeningDate

`func (o *Account) GetOpeningDate() time.Time`

GetOpeningDate returns the OpeningDate field if non-nil, zero value otherwise.

### GetOpeningDateOk

`func (o *Account) GetOpeningDateOk() (*time.Time, bool)`

GetOpeningDateOk returns a tuple with the OpeningDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOpeningDate

`func (o *Account) SetOpeningDate(v time.Time)`

SetOpeningDate sets OpeningDate field to given value.

### HasOpeningDate

`func (o *Account) HasOpeningDate() bool`

HasOpeningDate returns a boolean if a field has been set.

### GetClosingDate

`func (o *Account) GetClosingDate() time.Time`

GetClosingDate returns the ClosingDate field if non-nil, zero value otherwise.

### GetClosingDateOk

`func (o *Account) GetClosingDateOk() (*time.Time, bool)`

GetClosingDateOk returns a tuple with the ClosingDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClosingDate

`func (o *Account) SetClosingDate(v time.Time)`

SetClosingDate sets ClosingDate field to given value.

### HasClosingDate

`func (o *Account) HasClosingDate() bool`

HasClosingDate returns a boolean if a field has been set.

### GetShowInReconciliation

`func (o *Account) GetShowInReconciliation() bool`

GetShowInReconciliation returns the ShowInReconciliation field if non-nil, zero value otherwise.

### GetShowInReconciliationOk

`func (o *Account) GetShowInReconciliationOk() (*bool, bool)`

GetShowInReconciliationOk returns a tuple with the ShowInReconciliation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShowInReconciliation

`func (o *Account) SetShowInReconciliation(v bool)`

SetShowInReconciliation sets ShowInReconciliation field to given value.

### HasShowInReconciliation

`func (o *Account) HasShowInReconciliation() bool`

HasShowInReconciliation returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


