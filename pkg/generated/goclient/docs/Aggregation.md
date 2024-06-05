# Aggregation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**From** | **time.Time** |  | 
**To** | **time.Time** |  | 
**Granularity** | **string** |  | 
**Intervals** | [**[]time.Time**](time.Time.md) |  | 
**Currencies** | [**[]CurrencyAggregation**](CurrencyAggregation.md) |  | 

## Methods

### NewAggregation

`func NewAggregation(from time.Time, to time.Time, granularity string, intervals []time.Time, currencies []CurrencyAggregation, ) *Aggregation`

NewAggregation instantiates a new Aggregation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAggregationWithDefaults

`func NewAggregationWithDefaults() *Aggregation`

NewAggregationWithDefaults instantiates a new Aggregation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFrom

`func (o *Aggregation) GetFrom() time.Time`

GetFrom returns the From field if non-nil, zero value otherwise.

### GetFromOk

`func (o *Aggregation) GetFromOk() (*time.Time, bool)`

GetFromOk returns a tuple with the From field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrom

`func (o *Aggregation) SetFrom(v time.Time)`

SetFrom sets From field to given value.


### GetTo

`func (o *Aggregation) GetTo() time.Time`

GetTo returns the To field if non-nil, zero value otherwise.

### GetToOk

`func (o *Aggregation) GetToOk() (*time.Time, bool)`

GetToOk returns a tuple with the To field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTo

`func (o *Aggregation) SetTo(v time.Time)`

SetTo sets To field to given value.


### GetGranularity

`func (o *Aggregation) GetGranularity() string`

GetGranularity returns the Granularity field if non-nil, zero value otherwise.

### GetGranularityOk

`func (o *Aggregation) GetGranularityOk() (*string, bool)`

GetGranularityOk returns a tuple with the Granularity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGranularity

`func (o *Aggregation) SetGranularity(v string)`

SetGranularity sets Granularity field to given value.


### GetIntervals

`func (o *Aggregation) GetIntervals() []time.Time`

GetIntervals returns the Intervals field if non-nil, zero value otherwise.

### GetIntervalsOk

`func (o *Aggregation) GetIntervalsOk() (*[]time.Time, bool)`

GetIntervalsOk returns a tuple with the Intervals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntervals

`func (o *Aggregation) SetIntervals(v []time.Time)`

SetIntervals sets Intervals field to given value.


### GetCurrencies

`func (o *Aggregation) GetCurrencies() []CurrencyAggregation`

GetCurrencies returns the Currencies field if non-nil, zero value otherwise.

### GetCurrenciesOk

`func (o *Aggregation) GetCurrenciesOk() (*[]CurrencyAggregation, bool)`

GetCurrenciesOk returns a tuple with the Currencies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencies

`func (o *Aggregation) SetCurrencies(v []CurrencyAggregation)`

SetCurrencies sets Currencies field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


