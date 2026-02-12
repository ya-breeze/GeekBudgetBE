# BankImporterFile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**BankImporterId** | **string** |  | 
**Filename** | **string** |  | 
**UploadDate** | **time.Time** |  | 

## Methods

### NewBankImporterFile

`func NewBankImporterFile(id string, bankImporterId string, filename string, uploadDate time.Time, ) *BankImporterFile`

NewBankImporterFile instantiates a new BankImporterFile object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBankImporterFileWithDefaults

`func NewBankImporterFileWithDefaults() *BankImporterFile`

NewBankImporterFileWithDefaults instantiates a new BankImporterFile object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BankImporterFile) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BankImporterFile) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BankImporterFile) SetId(v string)`

SetId sets Id field to given value.


### GetBankImporterId

`func (o *BankImporterFile) GetBankImporterId() string`

GetBankImporterId returns the BankImporterId field if non-nil, zero value otherwise.

### GetBankImporterIdOk

`func (o *BankImporterFile) GetBankImporterIdOk() (*string, bool)`

GetBankImporterIdOk returns a tuple with the BankImporterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankImporterId

`func (o *BankImporterFile) SetBankImporterId(v string)`

SetBankImporterId sets BankImporterId field to given value.


### GetFilename

`func (o *BankImporterFile) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *BankImporterFile) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *BankImporterFile) SetFilename(v string)`

SetFilename sets Filename field to given value.


### GetUploadDate

`func (o *BankImporterFile) GetUploadDate() time.Time`

GetUploadDate returns the UploadDate field if non-nil, zero value otherwise.

### GetUploadDateOk

`func (o *BankImporterFile) GetUploadDateOk() (*time.Time, bool)`

GetUploadDateOk returns a tuple with the UploadDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadDate

`func (o *BankImporterFile) SetUploadDate(v time.Time)`

SetUploadDate sets UploadDate field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


