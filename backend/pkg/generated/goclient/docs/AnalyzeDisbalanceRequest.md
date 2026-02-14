# AnalyzeDisbalanceRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CurrencyId** | **string** |  | 
**TargetDelta** | [**decimal.Decimal**](decimal.Decimal.md) |  | 

## Methods

### NewAnalyzeDisbalanceRequest

`func NewAnalyzeDisbalanceRequest(currencyId string, targetDelta decimal.Decimal, ) *AnalyzeDisbalanceRequest`

NewAnalyzeDisbalanceRequest instantiates a new AnalyzeDisbalanceRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnalyzeDisbalanceRequestWithDefaults

`func NewAnalyzeDisbalanceRequestWithDefaults() *AnalyzeDisbalanceRequest`

NewAnalyzeDisbalanceRequestWithDefaults instantiates a new AnalyzeDisbalanceRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrencyId

`func (o *AnalyzeDisbalanceRequest) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *AnalyzeDisbalanceRequest) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *AnalyzeDisbalanceRequest) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.


### GetTargetDelta

`func (o *AnalyzeDisbalanceRequest) GetTargetDelta() decimal.Decimal`

GetTargetDelta returns the TargetDelta field if non-nil, zero value otherwise.

### GetTargetDeltaOk

`func (o *AnalyzeDisbalanceRequest) GetTargetDeltaOk() (*decimal.Decimal, bool)`

GetTargetDeltaOk returns a tuple with the TargetDelta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetDelta

`func (o *AnalyzeDisbalanceRequest) SetTargetDelta(v decimal.Decimal)`

SetTargetDelta sets TargetDelta field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


