# ImportResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Date** | Pointer to **time.Time** | Date of import | [optional] 
**Status** | Pointer to **string** | Status of import | [optional] 
**Description** | Pointer to **string** | Details of import | [optional] 

## Methods

### NewImportResult

`func NewImportResult() *ImportResult`

NewImportResult instantiates a new ImportResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewImportResultWithDefaults

`func NewImportResultWithDefaults() *ImportResult`

NewImportResultWithDefaults instantiates a new ImportResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *ImportResult) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *ImportResult) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *ImportResult) SetDate(v time.Time)`

SetDate sets Date field to given value.

### HasDate

`func (o *ImportResult) HasDate() bool`

HasDate returns a boolean if a field has been set.

### GetStatus

`func (o *ImportResult) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ImportResult) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ImportResult) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ImportResult) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetDescription

`func (o *ImportResult) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ImportResult) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ImportResult) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ImportResult) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


