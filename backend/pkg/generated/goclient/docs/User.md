# User

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Email** | **string** |  | 
**StartDate** | **time.Time** |  | 
**FavoriteCurrencyId** | Pointer to **string** | ID of the user&#39;s favorite currency. By default this currency will be used to convert other currencies. | [optional] 

## Methods

### NewUser

`func NewUser(id string, email string, startDate time.Time, ) *User`

NewUser instantiates a new User object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserWithDefaults

`func NewUserWithDefaults() *User`

NewUserWithDefaults instantiates a new User object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *User) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *User) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *User) SetId(v string)`

SetId sets Id field to given value.


### GetEmail

`func (o *User) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *User) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *User) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetStartDate

`func (o *User) GetStartDate() time.Time`

GetStartDate returns the StartDate field if non-nil, zero value otherwise.

### GetStartDateOk

`func (o *User) GetStartDateOk() (*time.Time, bool)`

GetStartDateOk returns a tuple with the StartDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartDate

`func (o *User) SetStartDate(v time.Time)`

SetStartDate sets StartDate field to given value.


### GetFavoriteCurrencyId

`func (o *User) GetFavoriteCurrencyId() string`

GetFavoriteCurrencyId returns the FavoriteCurrencyId field if non-nil, zero value otherwise.

### GetFavoriteCurrencyIdOk

`func (o *User) GetFavoriteCurrencyIdOk() (*string, bool)`

GetFavoriteCurrencyIdOk returns a tuple with the FavoriteCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFavoriteCurrencyId

`func (o *User) SetFavoriteCurrencyId(v string)`

SetFavoriteCurrencyId sets FavoriteCurrencyId field to given value.

### HasFavoriteCurrencyId

`func (o *User) HasFavoriteCurrencyId() bool`

HasFavoriteCurrencyId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


