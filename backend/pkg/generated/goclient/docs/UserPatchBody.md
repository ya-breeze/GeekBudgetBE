# UserPatchBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FavoriteCurrencyId** | Pointer to **string** | ID of the user&#39;s favorite currency. By default this currency will be used to convert other currencies. | [optional] 

## Methods

### NewUserPatchBody

`func NewUserPatchBody() *UserPatchBody`

NewUserPatchBody instantiates a new UserPatchBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserPatchBodyWithDefaults

`func NewUserPatchBodyWithDefaults() *UserPatchBody`

NewUserPatchBodyWithDefaults instantiates a new UserPatchBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFavoriteCurrencyId

`func (o *UserPatchBody) GetFavoriteCurrencyId() string`

GetFavoriteCurrencyId returns the FavoriteCurrencyId field if non-nil, zero value otherwise.

### GetFavoriteCurrencyIdOk

`func (o *UserPatchBody) GetFavoriteCurrencyIdOk() (*string, bool)`

GetFavoriteCurrencyIdOk returns a tuple with the FavoriteCurrencyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFavoriteCurrencyId

`func (o *UserPatchBody) SetFavoriteCurrencyId(v string)`

SetFavoriteCurrencyId sets FavoriteCurrencyId field to given value.

### HasFavoriteCurrencyId

`func (o *UserPatchBody) HasFavoriteCurrencyId() bool`

HasFavoriteCurrencyId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


