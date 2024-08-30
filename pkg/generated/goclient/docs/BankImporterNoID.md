# BankImporterNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Extra** | Pointer to **string** | Stores extra data about bank importer. For example could hold \&quot;bank account number\&quot; to be able to distinguish between different bank accounts, or it could hold token for bank API | [optional] 

## Methods

### NewBankImporterNoID

`func NewBankImporterNoID(name string, ) *BankImporterNoID`

NewBankImporterNoID instantiates a new BankImporterNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBankImporterNoIDWithDefaults

`func NewBankImporterNoIDWithDefaults() *BankImporterNoID`

NewBankImporterNoIDWithDefaults instantiates a new BankImporterNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BankImporterNoID) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BankImporterNoID) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BankImporterNoID) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *BankImporterNoID) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BankImporterNoID) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BankImporterNoID) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BankImporterNoID) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetExtra

`func (o *BankImporterNoID) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *BankImporterNoID) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *BankImporterNoID) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *BankImporterNoID) HasExtra() bool`

HasExtra returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


