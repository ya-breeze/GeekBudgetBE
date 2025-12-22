# BankImporterNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
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
**IsStopped** | Pointer to **bool** | If true, automatic fetching is stopped for this importer. This is usually set automatically when fetch fails, and reset when user manually triggers fetch or updates the importer. | [optional] 

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


### GetFeeAccountId

`func (o *BankImporterNoID) GetFeeAccountId() string`

GetFeeAccountId returns the FeeAccountId field if non-nil, zero value otherwise.

### GetFeeAccountIdOk

`func (o *BankImporterNoID) GetFeeAccountIdOk() (*string, bool)`

GetFeeAccountIdOk returns a tuple with the FeeAccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeeAccountId

`func (o *BankImporterNoID) SetFeeAccountId(v string)`

SetFeeAccountId sets FeeAccountId field to given value.

### HasFeeAccountId

`func (o *BankImporterNoID) HasFeeAccountId() bool`

HasFeeAccountId returns a boolean if a field has been set.

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

### GetFetchAll

`func (o *BankImporterNoID) GetFetchAll() bool`

GetFetchAll returns the FetchAll field if non-nil, zero value otherwise.

### GetFetchAllOk

`func (o *BankImporterNoID) GetFetchAllOk() (*bool, bool)`

GetFetchAllOk returns a tuple with the FetchAll field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFetchAll

`func (o *BankImporterNoID) SetFetchAll(v bool)`

SetFetchAll sets FetchAll field to given value.

### HasFetchAll

`func (o *BankImporterNoID) HasFetchAll() bool`

HasFetchAll returns a boolean if a field has been set.

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

`func (o *BankImporterNoID) GetLastImports() []ImportResult`

GetLastImports returns the LastImports field if non-nil, zero value otherwise.

### GetLastImportsOk

`func (o *BankImporterNoID) GetLastImportsOk() (*[]ImportResult, bool)`

GetLastImportsOk returns a tuple with the LastImports field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastImports

`func (o *BankImporterNoID) SetLastImports(v []ImportResult)`

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

### GetIsStopped

`func (o *BankImporterNoID) GetIsStopped() bool`

GetIsStopped returns the IsStopped field if non-nil, zero value otherwise.

### GetIsStoppedOk

`func (o *BankImporterNoID) GetIsStoppedOk() (*bool, bool)`

GetIsStoppedOk returns a tuple with the IsStopped field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsStopped

`func (o *BankImporterNoID) SetIsStopped(v bool)`

SetIsStopped sets IsStopped field to given value.

### HasIsStopped

`func (o *BankImporterNoID) HasIsStopped() bool`

HasIsStopped returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


