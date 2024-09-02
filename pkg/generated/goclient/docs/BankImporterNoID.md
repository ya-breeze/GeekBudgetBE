# BankImporterNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**AccountId** | **string** | ID of account which is used to store transactions from this bank importer | 
**Extra** | Pointer to **string** | Stores extra data about bank importer. For example could hold \&quot;bank account number\&quot; to be able to distinguish between different bank accounts, or it could hold token for bank API | [optional] 
**Type** | Pointer to **string** | Type of bank importer. It&#39;s used to distinguish between different banks. For example, FIO bank or KB bank. | [optional] 
**LastSuccessfulImport** | Pointer to **time.Time** | Date of last successful import. | [optional] 
**LastImports** | Pointer to [**[]BankImporterNoIDLastImportsInner**](BankImporterNoIDLastImportsInner.md) | List of last imports. It could be shown to user to explain what was imported recently | [optional] 
**Mappings** | Pointer to [**[]BankImporterNoIDMappingsInner**](BankImporterNoIDMappingsInner.md) | List of mappings which are used to enrich transactions with additional tags | [optional] 

## Methods

### NewBankImporterNoID

`func NewBankImporterNoID(name string, accountId string, ) *BankImporterNoID`

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

### GetAccountId

`func (o *BankImporterNoID) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BankImporterNoID) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BankImporterNoID) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


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

### GetType

`func (o *BankImporterNoID) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *BankImporterNoID) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *BankImporterNoID) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *BankImporterNoID) HasType() bool`

HasType returns a boolean if a field has been set.

### GetLastSuccessfulImport

`func (o *BankImporterNoID) GetLastSuccessfulImport() time.Time`

GetLastSuccessfulImport returns the LastSuccessfulImport field if non-nil, zero value otherwise.

### GetLastSuccessfulImportOk

`func (o *BankImporterNoID) GetLastSuccessfulImportOk() (*time.Time, bool)`

GetLastSuccessfulImportOk returns a tuple with the LastSuccessfulImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSuccessfulImport

`func (o *BankImporterNoID) SetLastSuccessfulImport(v time.Time)`

SetLastSuccessfulImport sets LastSuccessfulImport field to given value.

### HasLastSuccessfulImport

`func (o *BankImporterNoID) HasLastSuccessfulImport() bool`

HasLastSuccessfulImport returns a boolean if a field has been set.

### GetLastImports

`func (o *BankImporterNoID) GetLastImports() []BankImporterNoIDLastImportsInner`

GetLastImports returns the LastImports field if non-nil, zero value otherwise.

### GetLastImportsOk

`func (o *BankImporterNoID) GetLastImportsOk() (*[]BankImporterNoIDLastImportsInner, bool)`

GetLastImportsOk returns a tuple with the LastImports field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastImports

`func (o *BankImporterNoID) SetLastImports(v []BankImporterNoIDLastImportsInner)`

SetLastImports sets LastImports field to given value.

### HasLastImports

`func (o *BankImporterNoID) HasLastImports() bool`

HasLastImports returns a boolean if a field has been set.

### GetMappings

`func (o *BankImporterNoID) GetMappings() []BankImporterNoIDMappingsInner`

GetMappings returns the Mappings field if non-nil, zero value otherwise.

### GetMappingsOk

`func (o *BankImporterNoID) GetMappingsOk() (*[]BankImporterNoIDMappingsInner, bool)`

GetMappingsOk returns a tuple with the Mappings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMappings

`func (o *BankImporterNoID) SetMappings(v []BankImporterNoIDMappingsInner)`

SetMappings sets Mappings field to given value.

### HasMappings

`func (o *BankImporterNoID) HasMappings() bool`

HasMappings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


