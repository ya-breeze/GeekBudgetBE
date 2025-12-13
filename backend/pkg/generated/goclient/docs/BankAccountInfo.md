# BankAccountInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | Pointer to **string** |  | [optional] 
**BankId** | Pointer to **string** |  | [optional] 
**Balances** | Pointer to [**[]BankAccountInfoBalancesInner**](BankAccountInfoBalancesInner.md) | List of balances for this account. It&#39;s an array since one account could hold multiple currencies, for example, cash account could hold EUR, USD and CZK. Or one bank account could hold multiple currencies. | [optional] 

## Methods

### NewBankAccountInfo

`func NewBankAccountInfo() *BankAccountInfo`

NewBankAccountInfo instantiates a new BankAccountInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBankAccountInfoWithDefaults

`func NewBankAccountInfoWithDefaults() *BankAccountInfo`

NewBankAccountInfoWithDefaults instantiates a new BankAccountInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *BankAccountInfo) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *BankAccountInfo) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *BankAccountInfo) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *BankAccountInfo) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### GetBankId

`func (o *BankAccountInfo) GetBankId() string`

GetBankId returns the BankId field if non-nil, zero value otherwise.

### GetBankIdOk

`func (o *BankAccountInfo) GetBankIdOk() (*string, bool)`

GetBankIdOk returns a tuple with the BankId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankId

`func (o *BankAccountInfo) SetBankId(v string)`

SetBankId sets BankId field to given value.

### HasBankId

`func (o *BankAccountInfo) HasBankId() bool`

HasBankId returns a boolean if a field has been set.

### GetBalances

`func (o *BankAccountInfo) GetBalances() []BankAccountInfoBalancesInner`

GetBalances returns the Balances field if non-nil, zero value otherwise.

### GetBalancesOk

`func (o *BankAccountInfo) GetBalancesOk() (*[]BankAccountInfoBalancesInner, bool)`

GetBalancesOk returns a tuple with the Balances field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalances

`func (o *BankAccountInfo) SetBalances(v []BankAccountInfoBalancesInner)`

SetBalances sets Balances field to given value.

### HasBalances

`func (o *BankAccountInfo) HasBalances() bool`

HasBalances returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


