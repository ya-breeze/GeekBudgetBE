# BankAccountInfoBalancesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CurrencyId** | Pointer to **string** |  | [optional] 
**OpeningBalance** | Pointer to [**decimal.Decimal**](decimal.Decimal.md) |  | [optional] 
**ClosingBalance** | Pointer to [**decimal.Decimal**](decimal.Decimal.md) |  | [optional] 
**LastUpdatedAt** | Pointer to **NullableTime** |  | [optional] 

## Methods

### NewBankAccountInfoBalancesInner

`func NewBankAccountInfoBalancesInner() *BankAccountInfoBalancesInner`

NewBankAccountInfoBalancesInner instantiates a new BankAccountInfoBalancesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBankAccountInfoBalancesInnerWithDefaults

`func NewBankAccountInfoBalancesInnerWithDefaults() *BankAccountInfoBalancesInner`

NewBankAccountInfoBalancesInnerWithDefaults instantiates a new BankAccountInfoBalancesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCurrencyId

`func (o *BankAccountInfoBalancesInner) GetCurrencyId() string`

GetCurrencyId returns the CurrencyId field if non-nil, zero value otherwise.

### GetCurrencyIdOk

`func (o *BankAccountInfoBalancesInner) GetCurrencyIdOk() (*string, bool)`

GetCurrencyIdOk returns a tuple with the CurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyId

`func (o *BankAccountInfoBalancesInner) SetCurrencyId(v string)`

SetCurrencyId sets CurrencyId field to given value.

### HasCurrencyId

`func (o *BankAccountInfoBalancesInner) HasCurrencyId() bool`

HasCurrencyId returns a boolean if a field has been set.

### GetOpeningBalance

`func (o *BankAccountInfoBalancesInner) GetOpeningBalance() decimal.Decimal`

GetOpeningBalance returns the OpeningBalance field if non-nil, zero value otherwise.

### GetOpeningBalanceOk

`func (o *BankAccountInfoBalancesInner) GetOpeningBalanceOk() (*decimal.Decimal, bool)`

GetOpeningBalanceOk returns a tuple with the OpeningBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOpeningBalance

`func (o *BankAccountInfoBalancesInner) SetOpeningBalance(v decimal.Decimal)`

SetOpeningBalance sets OpeningBalance field to given value.

### HasOpeningBalance

`func (o *BankAccountInfoBalancesInner) HasOpeningBalance() bool`

HasOpeningBalance returns a boolean if a field has been set.

### GetClosingBalance

`func (o *BankAccountInfoBalancesInner) GetClosingBalance() decimal.Decimal`

GetClosingBalance returns the ClosingBalance field if non-nil, zero value otherwise.

### GetClosingBalanceOk

`func (o *BankAccountInfoBalancesInner) GetClosingBalanceOk() (*decimal.Decimal, bool)`

GetClosingBalanceOk returns a tuple with the ClosingBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClosingBalance

`func (o *BankAccountInfoBalancesInner) SetClosingBalance(v decimal.Decimal)`

SetClosingBalance sets ClosingBalance field to given value.

### HasClosingBalance

`func (o *BankAccountInfoBalancesInner) HasClosingBalance() bool`

HasClosingBalance returns a boolean if a field has been set.

### GetLastUpdatedAt

`func (o *BankAccountInfoBalancesInner) GetLastUpdatedAt() time.Time`

GetLastUpdatedAt returns the LastUpdatedAt field if non-nil, zero value otherwise.

### GetLastUpdatedAtOk

`func (o *BankAccountInfoBalancesInner) GetLastUpdatedAtOk() (*time.Time, bool)`

GetLastUpdatedAtOk returns a tuple with the LastUpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastUpdatedAt

`func (o *BankAccountInfoBalancesInner) SetLastUpdatedAt(v time.Time)`

SetLastUpdatedAt sets LastUpdatedAt field to given value.

### HasLastUpdatedAt

`func (o *BankAccountInfoBalancesInner) HasLastUpdatedAt() bool`

HasLastUpdatedAt returns a boolean if a field has been set.

### SetLastUpdatedAtNil

`func (o *BankAccountInfoBalancesInner) SetLastUpdatedAtNil(b bool)`

 SetLastUpdatedAtNil sets the value for LastUpdatedAt to be an explicit nil

### UnsetLastUpdatedAt
`func (o *BankAccountInfoBalancesInner) UnsetLastUpdatedAt()`

UnsetLastUpdatedAt ensures that no value is present for LastUpdatedAt, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


