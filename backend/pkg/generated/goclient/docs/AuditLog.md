# AuditLog

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**UserId** | **string** |  | 
**EntityType** | **string** |  | 
**EntityId** | **string** |  | 
**Action** | **string** |  | 
**ChangeSource** | **string** |  | 
**Before** | Pointer to **NullableString** |  | [optional] 
**After** | Pointer to **NullableString** |  | [optional] 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewAuditLog

`func NewAuditLog(id string, userId string, entityType string, entityId string, action string, changeSource string, createdAt time.Time, ) *AuditLog`

NewAuditLog instantiates a new AuditLog object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuditLogWithDefaults

`func NewAuditLogWithDefaults() *AuditLog`

NewAuditLogWithDefaults instantiates a new AuditLog object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AuditLog) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AuditLog) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AuditLog) SetId(v string)`

SetId sets Id field to given value.


### GetUserId

`func (o *AuditLog) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *AuditLog) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *AuditLog) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetEntityType

`func (o *AuditLog) GetEntityType() string`

GetEntityType returns the EntityType field if non-nil, zero value otherwise.

### GetEntityTypeOk

`func (o *AuditLog) GetEntityTypeOk() (*string, bool)`

GetEntityTypeOk returns a tuple with the EntityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntityType

`func (o *AuditLog) SetEntityType(v string)`

SetEntityType sets EntityType field to given value.


### GetEntityId

`func (o *AuditLog) GetEntityId() string`

GetEntityId returns the EntityId field if non-nil, zero value otherwise.

### GetEntityIdOk

`func (o *AuditLog) GetEntityIdOk() (*string, bool)`

GetEntityIdOk returns a tuple with the EntityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntityId

`func (o *AuditLog) SetEntityId(v string)`

SetEntityId sets EntityId field to given value.


### GetAction

`func (o *AuditLog) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *AuditLog) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *AuditLog) SetAction(v string)`

SetAction sets Action field to given value.


### GetChangeSource

`func (o *AuditLog) GetChangeSource() string`

GetChangeSource returns the ChangeSource field if non-nil, zero value otherwise.

### GetChangeSourceOk

`func (o *AuditLog) GetChangeSourceOk() (*string, bool)`

GetChangeSourceOk returns a tuple with the ChangeSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChangeSource

`func (o *AuditLog) SetChangeSource(v string)`

SetChangeSource sets ChangeSource field to given value.


### GetBefore

`func (o *AuditLog) GetBefore() string`

GetBefore returns the Before field if non-nil, zero value otherwise.

### GetBeforeOk

`func (o *AuditLog) GetBeforeOk() (*string, bool)`

GetBeforeOk returns a tuple with the Before field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBefore

`func (o *AuditLog) SetBefore(v string)`

SetBefore sets Before field to given value.

### HasBefore

`func (o *AuditLog) HasBefore() bool`

HasBefore returns a boolean if a field has been set.

### SetBeforeNil

`func (o *AuditLog) SetBeforeNil(b bool)`

 SetBeforeNil sets the value for Before to be an explicit nil

### UnsetBefore
`func (o *AuditLog) UnsetBefore()`

UnsetBefore ensures that no value is present for Before, not even an explicit nil
### GetAfter

`func (o *AuditLog) GetAfter() string`

GetAfter returns the After field if non-nil, zero value otherwise.

### GetAfterOk

`func (o *AuditLog) GetAfterOk() (*string, bool)`

GetAfterOk returns a tuple with the After field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAfter

`func (o *AuditLog) SetAfter(v string)`

SetAfter sets After field to given value.

### HasAfter

`func (o *AuditLog) HasAfter() bool`

HasAfter returns a boolean if a field has been set.

### SetAfterNil

`func (o *AuditLog) SetAfterNil(b bool)`

 SetAfterNil sets the value for After to be an explicit nil

### UnsetAfter
`func (o *AuditLog) UnsetAfter()`

UnsetAfter ensures that no value is present for After, not even an explicit nil
### GetCreatedAt

`func (o *AuditLog) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AuditLog) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AuditLog) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


