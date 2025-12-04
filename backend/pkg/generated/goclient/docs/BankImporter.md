# BankImporter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**Description** | Pointer to **string** |  | [optional] 
**AccountId** | **string** | ID of account which is used to for movements from this bank importer | 
**FeeAccountId** | Pointer to **string** | ID of account which is used for fee movements from this bank importer | [optional] 
**Extra** | Pointer to **string** | Stores extra data about bank importer. For example could hold \&quot;bank account number\&quot; to be able to distinguish between different bank accounts, or it could hold token for bank API | [optional] 
**FetchAll** | Pointer to **bool** | If true, importer will fetch all transactions from the bank, if false, it will fetch only recent transactions | [optional] 
**Type** | Pointer to **string** | Type of bank importer. It&#39;s used to distinguish between different banks. For example, FIO bank or KB bank. | [optional] 
**LastSuccessfulImport** | Pointer to **time.Time** | Date of last successful import. | [optional] 
**LastImports** | Pointer to [**[]ImportResult**](ImportResult.md) | List of last imports. It could be shown to user to explain what was imported recently | [optional] 
**Mappings** | Pointer to [**[]BankImporterNoIDMappingsInner**](BankImporterNoIDMappingsInner.md) | List of mappings which are used to enrich transactions with additional tags | [optional] 

## Methods

### NewBankImporter

`func NewBankImporter(id string, name string, accountId string, ) *BankImporter`

NewBankImporter instantiates a new BankImporter object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBankImporterWithDefaults

`func NewBankImporterWithDefaults() *BankImporter`

NewBankImporterWithDefaults instantiates a new BankImporter object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BankImporter) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BankImporter) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BankImporter) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *BankImporter) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BankImporter) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BankImporter) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *BankImporter) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BankImporter) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BankImporter) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BankImporter) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetAccountId

`func (o *BankImporter) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BankImporter) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BankImporter) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.


### GetFeeAccountId

`func (o *BankImporter) GetFeeAccountId() string`

GetFeeAccountId returns the FeeAccountId field if non-nil, zero value otherwise.

### GetFeeAccountIdOk

`func (o *BankImporter) GetFeeAccountIdOk() (*string, bool)`

GetFeeAccountIdOk returns a tuple with the FeeAccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeeAccountId

`func (o *BankImporter) SetFeeAccountId(v string)`

SetFeeAccountId sets FeeAccountId field to given value.

### HasFeeAccountId

`func (o *BankImporter) HasFeeAccountId() bool`

HasFeeAccountId returns a boolean if a field has been set.

### GetExtra

`func (o *BankImporter) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *BankImporter) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *BankImporter) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *BankImporter) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetFetchAll

`func (o *BankImporter) GetFetchAll() bool`

GetFetchAll returns the FetchAll field if non-nil, zero value otherwise.

### GetFetchAllOk

`func (o *BankImporter) GetFetchAllOk() (*bool, bool)`

GetFetchAllOk returns a tuple with the FetchAll field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFetchAll

`func (o *BankImporter) SetFetchAll(v bool)`

SetFetchAll sets FetchAll field to given value.

### HasFetchAll

`func (o *BankImporter) HasFetchAll() bool`

HasFetchAll returns a boolean if a field has been set.

### GetType

`func (o *BankImporter) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *BankImporter) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *BankImporter) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *BankImporter) HasType() bool`

HasType returns a boolean if a field has been set.

### GetLastSuccessfulImport

`func (o *BankImporter) GetLastSuccessfulImport() time.Time`

GetLastSuccessfulImport returns the LastSuccessfulImport field if non-nil, zero value otherwise.

### GetLastSuccessfulImportOk

`func (o *BankImporter) GetLastSuccessfulImportOk() (*time.Time, bool)`

GetLastSuccessfulImportOk returns a tuple with the LastSuccessfulImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSuccessfulImport

`func (o *BankImporter) SetLastSuccessfulImport(v time.Time)`

SetLastSuccessfulImport sets LastSuccessfulImport field to given value.

### HasLastSuccessfulImport

`func (o *BankImporter) HasLastSuccessfulImport() bool`

HasLastSuccessfulImport returns a boolean if a field has been set.

### GetLastImports

`func (o *BankImporter) GetLastImports() []ImportResult`

GetLastImports returns the LastImports field if non-nil, zero value otherwise.

### GetLastImportsOk

`func (o *BankImporter) GetLastImportsOk() (*[]ImportResult, bool)`

GetLastImportsOk returns a tuple with the LastImports field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastImports

`func (o *BankImporter) SetLastImports(v []ImportResult)`

SetLastImports sets LastImports field to given value.

### HasLastImports

`func (o *BankImporter) HasLastImports() bool`

HasLastImports returns a boolean if a field has been set.

### GetMappings

`func (o *BankImporter) GetMappings() []BankImporterNoIDMappingsInner`

GetMappings returns the Mappings field if non-nil, zero value otherwise.

### GetMappingsOk

`func (o *BankImporter) GetMappingsOk() (*[]BankImporterNoIDMappingsInner, bool)`

GetMappingsOk returns a tuple with the Mappings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMappings

`func (o *BankImporter) SetMappings(v []BankImporterNoIDMappingsInner)`

SetMappings sets Mappings field to given value.

### HasMappings

`func (o *BankImporter) HasMappings() bool`

HasMappings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


