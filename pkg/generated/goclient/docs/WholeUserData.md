# WholeUserData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | Pointer to [**User**](User.md) |  | [optional] 
**Currencies** | Pointer to [**[]Currency**](Currency.md) |  | [optional] 
**Accounts** | Pointer to [**[]Account**](Account.md) |  | [optional] 
**Transactions** | Pointer to [**[]Transaction**](Transaction.md) |  | [optional] 
**Matchers** | Pointer to [**[]Matcher**](Matcher.md) |  | [optional] 
**BankImporters** | Pointer to [**[]BankImporter**](BankImporter.md) |  | [optional] 

## Methods

### NewWholeUserData

`func NewWholeUserData() *WholeUserData`

NewWholeUserData instantiates a new WholeUserData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWholeUserDataWithDefaults

`func NewWholeUserDataWithDefaults() *WholeUserData`

NewWholeUserDataWithDefaults instantiates a new WholeUserData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *WholeUserData) GetUser() User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *WholeUserData) GetUserOk() (*User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *WholeUserData) SetUser(v User)`

SetUser sets User field to given value.

### HasUser

`func (o *WholeUserData) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetCurrencies

`func (o *WholeUserData) GetCurrencies() []Currency`

GetCurrencies returns the Currencies field if non-nil, zero value otherwise.

### GetCurrenciesOk

`func (o *WholeUserData) GetCurrenciesOk() (*[]Currency, bool)`

GetCurrenciesOk returns a tuple with the Currencies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencies

`func (o *WholeUserData) SetCurrencies(v []Currency)`

SetCurrencies sets Currencies field to given value.

### HasCurrencies

`func (o *WholeUserData) HasCurrencies() bool`

HasCurrencies returns a boolean if a field has been set.

### GetAccounts

`func (o *WholeUserData) GetAccounts() []Account`

GetAccounts returns the Accounts field if non-nil, zero value otherwise.

### GetAccountsOk

`func (o *WholeUserData) GetAccountsOk() (*[]Account, bool)`

GetAccountsOk returns a tuple with the Accounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccounts

`func (o *WholeUserData) SetAccounts(v []Account)`

SetAccounts sets Accounts field to given value.

### HasAccounts

`func (o *WholeUserData) HasAccounts() bool`

HasAccounts returns a boolean if a field has been set.

### GetTransactions

`func (o *WholeUserData) GetTransactions() []Transaction`

GetTransactions returns the Transactions field if non-nil, zero value otherwise.

### GetTransactionsOk

`func (o *WholeUserData) GetTransactionsOk() (*[]Transaction, bool)`

GetTransactionsOk returns a tuple with the Transactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactions

`func (o *WholeUserData) SetTransactions(v []Transaction)`

SetTransactions sets Transactions field to given value.

### HasTransactions

`func (o *WholeUserData) HasTransactions() bool`

HasTransactions returns a boolean if a field has been set.

### GetMatchers

`func (o *WholeUserData) GetMatchers() []Matcher`

GetMatchers returns the Matchers field if non-nil, zero value otherwise.

### GetMatchersOk

`func (o *WholeUserData) GetMatchersOk() (*[]Matcher, bool)`

GetMatchersOk returns a tuple with the Matchers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchers

`func (o *WholeUserData) SetMatchers(v []Matcher)`

SetMatchers sets Matchers field to given value.

### HasMatchers

`func (o *WholeUserData) HasMatchers() bool`

HasMatchers returns a boolean if a field has been set.

### GetBankImporters

`func (o *WholeUserData) GetBankImporters() []BankImporter`

GetBankImporters returns the BankImporters field if non-nil, zero value otherwise.

### GetBankImportersOk

`func (o *WholeUserData) GetBankImportersOk() (*[]BankImporter, bool)`

GetBankImportersOk returns a tuple with the BankImporters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBankImporters

`func (o *WholeUserData) SetBankImporters(v []BankImporter)`

SetBankImporters sets BankImporters field to given value.

### HasBankImporters

`func (o *WholeUserData) HasBankImporters() bool`

HasBankImporters returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


