# DisbalanceCandidateTransaction

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Date** | **time.Time** |  | 
**Description** | **string** |  | 
**Amount** | [**decimal.Decimal**](decimal.Decimal.md) |  | 

## Methods

### NewDisbalanceCandidateTransaction

`func NewDisbalanceCandidateTransaction(id string, date time.Time, description string, amount decimal.Decimal, ) *DisbalanceCandidateTransaction`

NewDisbalanceCandidateTransaction instantiates a new DisbalanceCandidateTransaction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDisbalanceCandidateTransactionWithDefaults

`func NewDisbalanceCandidateTransactionWithDefaults() *DisbalanceCandidateTransaction`

NewDisbalanceCandidateTransactionWithDefaults instantiates a new DisbalanceCandidateTransaction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DisbalanceCandidateTransaction) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DisbalanceCandidateTransaction) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DisbalanceCandidateTransaction) SetId(v string)`

SetId sets Id field to given value.


### GetDate

`func (o *DisbalanceCandidateTransaction) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *DisbalanceCandidateTransaction) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *DisbalanceCandidateTransaction) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetDescription

`func (o *DisbalanceCandidateTransaction) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *DisbalanceCandidateTransaction) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *DisbalanceCandidateTransaction) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetAmount

`func (o *DisbalanceCandidateTransaction) GetAmount() decimal.Decimal`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *DisbalanceCandidateTransaction) GetAmountOk() (*decimal.Decimal, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *DisbalanceCandidateTransaction) SetAmount(v decimal.Decimal)`

SetAmount sets Amount field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


