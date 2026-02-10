# TransactionNoID

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Date** | **time.Time** |  | 
**Description** | Pointer to **string** |  | [optional] 
**Place** | Pointer to **string** |  | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 
**PartnerName** | Pointer to **string** |  | [optional] 
**PartnerAccount** | Pointer to **string** |  | [optional] 
**PartnerInternalId** | Pointer to **string** | Internal bank&#39;s ID to be able to match later if necessary | [optional] 
**Extra** | Pointer to **string** | Stores extra data about transaction. For example could hold \&quot;variable symbol\&quot; to distinguish payment for the same account, but with different meaning | [optional] 
**UnprocessedSources** | Pointer to **string** | Stores FULL unprocessed transactions which was source of this transaction. Could be used later for detailed analysis | [optional] 
**ExternalIds** | Pointer to **[]string** | IDs of unprocessed transaction - to match later | [optional] 
**Movements** | [**[]Movement**](Movement.md) |  | 
**MatcherId** | Pointer to **string** | ID of the matcher used for this conversion (if any) | [optional] 
**IsAuto** | Pointer to **bool** | If true, this transaction was converted automatically by the matcher | [optional] 
**SuspiciousReasons** | Pointer to **[]string** |  | [optional] 
**MergedIntoId** | Pointer to **string** | ID of the transaction this one was merged into (if any) | [optional] 
**MergedAt** | Pointer to **time.Time** | When this transaction was merged | [optional] 
**AutoMatchSkipReason** | Pointer to **string** | Reason why auto-match was skipped for this transaction | [optional] 
**DuplicateDismissed** | Pointer to **bool** | If true, user has dismissed the duplicate detected for this transaction | [optional] [default to false]
**MergedTransactionIds** | Pointer to **[]string** | List of transaction IDs that were merged into this transaction (soft-deleted duplicates pointing here via mergedIntoId) | [optional] 
**DuplicateTransactionIds** | Pointer to **[]string** | List of transaction IDs that are potential duplicates of this one (from separate junction table) | [optional] 

## Methods

### NewTransactionNoID

`func NewTransactionNoID(date time.Time, movements []Movement, ) *TransactionNoID`

NewTransactionNoID instantiates a new TransactionNoID object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransactionNoIDWithDefaults

`func NewTransactionNoIDWithDefaults() *TransactionNoID`

NewTransactionNoIDWithDefaults instantiates a new TransactionNoID object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *TransactionNoID) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *TransactionNoID) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *TransactionNoID) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetDescription

`func (o *TransactionNoID) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *TransactionNoID) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *TransactionNoID) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *TransactionNoID) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPlace

`func (o *TransactionNoID) GetPlace() string`

GetPlace returns the Place field if non-nil, zero value otherwise.

### GetPlaceOk

`func (o *TransactionNoID) GetPlaceOk() (*string, bool)`

GetPlaceOk returns a tuple with the Place field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlace

`func (o *TransactionNoID) SetPlace(v string)`

SetPlace sets Place field to given value.

### HasPlace

`func (o *TransactionNoID) HasPlace() bool`

HasPlace returns a boolean if a field has been set.

### GetTags

`func (o *TransactionNoID) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *TransactionNoID) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *TransactionNoID) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *TransactionNoID) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetPartnerName

`func (o *TransactionNoID) GetPartnerName() string`

GetPartnerName returns the PartnerName field if non-nil, zero value otherwise.

### GetPartnerNameOk

`func (o *TransactionNoID) GetPartnerNameOk() (*string, bool)`

GetPartnerNameOk returns a tuple with the PartnerName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerName

`func (o *TransactionNoID) SetPartnerName(v string)`

SetPartnerName sets PartnerName field to given value.

### HasPartnerName

`func (o *TransactionNoID) HasPartnerName() bool`

HasPartnerName returns a boolean if a field has been set.

### GetPartnerAccount

`func (o *TransactionNoID) GetPartnerAccount() string`

GetPartnerAccount returns the PartnerAccount field if non-nil, zero value otherwise.

### GetPartnerAccountOk

`func (o *TransactionNoID) GetPartnerAccountOk() (*string, bool)`

GetPartnerAccountOk returns a tuple with the PartnerAccount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerAccount

`func (o *TransactionNoID) SetPartnerAccount(v string)`

SetPartnerAccount sets PartnerAccount field to given value.

### HasPartnerAccount

`func (o *TransactionNoID) HasPartnerAccount() bool`

HasPartnerAccount returns a boolean if a field has been set.

### GetPartnerInternalId

`func (o *TransactionNoID) GetPartnerInternalId() string`

GetPartnerInternalId returns the PartnerInternalId field if non-nil, zero value otherwise.

### GetPartnerInternalIdOk

`func (o *TransactionNoID) GetPartnerInternalIdOk() (*string, bool)`

GetPartnerInternalIdOk returns a tuple with the PartnerInternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartnerInternalId

`func (o *TransactionNoID) SetPartnerInternalId(v string)`

SetPartnerInternalId sets PartnerInternalId field to given value.

### HasPartnerInternalId

`func (o *TransactionNoID) HasPartnerInternalId() bool`

HasPartnerInternalId returns a boolean if a field has been set.

### GetExtra

`func (o *TransactionNoID) GetExtra() string`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *TransactionNoID) GetExtraOk() (*string, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *TransactionNoID) SetExtra(v string)`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *TransactionNoID) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetUnprocessedSources

`func (o *TransactionNoID) GetUnprocessedSources() string`

GetUnprocessedSources returns the UnprocessedSources field if non-nil, zero value otherwise.

### GetUnprocessedSourcesOk

`func (o *TransactionNoID) GetUnprocessedSourcesOk() (*string, bool)`

GetUnprocessedSourcesOk returns a tuple with the UnprocessedSources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnprocessedSources

`func (o *TransactionNoID) SetUnprocessedSources(v string)`

SetUnprocessedSources sets UnprocessedSources field to given value.

### HasUnprocessedSources

`func (o *TransactionNoID) HasUnprocessedSources() bool`

HasUnprocessedSources returns a boolean if a field has been set.

### GetExternalIds

`func (o *TransactionNoID) GetExternalIds() []string`

GetExternalIds returns the ExternalIds field if non-nil, zero value otherwise.

### GetExternalIdsOk

`func (o *TransactionNoID) GetExternalIdsOk() (*[]string, bool)`

GetExternalIdsOk returns a tuple with the ExternalIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalIds

`func (o *TransactionNoID) SetExternalIds(v []string)`

SetExternalIds sets ExternalIds field to given value.

### HasExternalIds

`func (o *TransactionNoID) HasExternalIds() bool`

HasExternalIds returns a boolean if a field has been set.

### GetMovements

`func (o *TransactionNoID) GetMovements() []Movement`

GetMovements returns the Movements field if non-nil, zero value otherwise.

### GetMovementsOk

`func (o *TransactionNoID) GetMovementsOk() (*[]Movement, bool)`

GetMovementsOk returns a tuple with the Movements field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMovements

`func (o *TransactionNoID) SetMovements(v []Movement)`

SetMovements sets Movements field to given value.


### GetMatcherId

`func (o *TransactionNoID) GetMatcherId() string`

GetMatcherId returns the MatcherId field if non-nil, zero value otherwise.

### GetMatcherIdOk

`func (o *TransactionNoID) GetMatcherIdOk() (*string, bool)`

GetMatcherIdOk returns a tuple with the MatcherId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatcherId

`func (o *TransactionNoID) SetMatcherId(v string)`

SetMatcherId sets MatcherId field to given value.

### HasMatcherId

`func (o *TransactionNoID) HasMatcherId() bool`

HasMatcherId returns a boolean if a field has been set.

### GetIsAuto

`func (o *TransactionNoID) GetIsAuto() bool`

GetIsAuto returns the IsAuto field if non-nil, zero value otherwise.

### GetIsAutoOk

`func (o *TransactionNoID) GetIsAutoOk() (*bool, bool)`

GetIsAutoOk returns a tuple with the IsAuto field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsAuto

`func (o *TransactionNoID) SetIsAuto(v bool)`

SetIsAuto sets IsAuto field to given value.

### HasIsAuto

`func (o *TransactionNoID) HasIsAuto() bool`

HasIsAuto returns a boolean if a field has been set.

### GetSuspiciousReasons

`func (o *TransactionNoID) GetSuspiciousReasons() []string`

GetSuspiciousReasons returns the SuspiciousReasons field if non-nil, zero value otherwise.

### GetSuspiciousReasonsOk

`func (o *TransactionNoID) GetSuspiciousReasonsOk() (*[]string, bool)`

GetSuspiciousReasonsOk returns a tuple with the SuspiciousReasons field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuspiciousReasons

`func (o *TransactionNoID) SetSuspiciousReasons(v []string)`

SetSuspiciousReasons sets SuspiciousReasons field to given value.

### HasSuspiciousReasons

`func (o *TransactionNoID) HasSuspiciousReasons() bool`

HasSuspiciousReasons returns a boolean if a field has been set.

### GetMergedIntoId

`func (o *TransactionNoID) GetMergedIntoId() string`

GetMergedIntoId returns the MergedIntoId field if non-nil, zero value otherwise.

### GetMergedIntoIdOk

`func (o *TransactionNoID) GetMergedIntoIdOk() (*string, bool)`

GetMergedIntoIdOk returns a tuple with the MergedIntoId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedIntoId

`func (o *TransactionNoID) SetMergedIntoId(v string)`

SetMergedIntoId sets MergedIntoId field to given value.

### HasMergedIntoId

`func (o *TransactionNoID) HasMergedIntoId() bool`

HasMergedIntoId returns a boolean if a field has been set.

### GetMergedAt

`func (o *TransactionNoID) GetMergedAt() time.Time`

GetMergedAt returns the MergedAt field if non-nil, zero value otherwise.

### GetMergedAtOk

`func (o *TransactionNoID) GetMergedAtOk() (*time.Time, bool)`

GetMergedAtOk returns a tuple with the MergedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedAt

`func (o *TransactionNoID) SetMergedAt(v time.Time)`

SetMergedAt sets MergedAt field to given value.

### HasMergedAt

`func (o *TransactionNoID) HasMergedAt() bool`

HasMergedAt returns a boolean if a field has been set.

### GetAutoMatchSkipReason

`func (o *TransactionNoID) GetAutoMatchSkipReason() string`

GetAutoMatchSkipReason returns the AutoMatchSkipReason field if non-nil, zero value otherwise.

### GetAutoMatchSkipReasonOk

`func (o *TransactionNoID) GetAutoMatchSkipReasonOk() (*string, bool)`

GetAutoMatchSkipReasonOk returns a tuple with the AutoMatchSkipReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoMatchSkipReason

`func (o *TransactionNoID) SetAutoMatchSkipReason(v string)`

SetAutoMatchSkipReason sets AutoMatchSkipReason field to given value.

### HasAutoMatchSkipReason

`func (o *TransactionNoID) HasAutoMatchSkipReason() bool`

HasAutoMatchSkipReason returns a boolean if a field has been set.

### GetDuplicateDismissed

`func (o *TransactionNoID) GetDuplicateDismissed() bool`

GetDuplicateDismissed returns the DuplicateDismissed field if non-nil, zero value otherwise.

### GetDuplicateDismissedOk

`func (o *TransactionNoID) GetDuplicateDismissedOk() (*bool, bool)`

GetDuplicateDismissedOk returns a tuple with the DuplicateDismissed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateDismissed

`func (o *TransactionNoID) SetDuplicateDismissed(v bool)`

SetDuplicateDismissed sets DuplicateDismissed field to given value.

### HasDuplicateDismissed

`func (o *TransactionNoID) HasDuplicateDismissed() bool`

HasDuplicateDismissed returns a boolean if a field has been set.

### GetMergedTransactionIds

`func (o *TransactionNoID) GetMergedTransactionIds() []string`

GetMergedTransactionIds returns the MergedTransactionIds field if non-nil, zero value otherwise.

### GetMergedTransactionIdsOk

`func (o *TransactionNoID) GetMergedTransactionIdsOk() (*[]string, bool)`

GetMergedTransactionIdsOk returns a tuple with the MergedTransactionIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMergedTransactionIds

`func (o *TransactionNoID) SetMergedTransactionIds(v []string)`

SetMergedTransactionIds sets MergedTransactionIds field to given value.

### HasMergedTransactionIds

`func (o *TransactionNoID) HasMergedTransactionIds() bool`

HasMergedTransactionIds returns a boolean if a field has been set.

### GetDuplicateTransactionIds

`func (o *TransactionNoID) GetDuplicateTransactionIds() []string`

GetDuplicateTransactionIds returns the DuplicateTransactionIds field if non-nil, zero value otherwise.

### GetDuplicateTransactionIdsOk

`func (o *TransactionNoID) GetDuplicateTransactionIdsOk() (*[]string, bool)`

GetDuplicateTransactionIdsOk returns a tuple with the DuplicateTransactionIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateTransactionIds

`func (o *TransactionNoID) SetDuplicateTransactionIds(v []string)`

SetDuplicateTransactionIds sets DuplicateTransactionIds field to given value.

### HasDuplicateTransactionIds

`func (o *TransactionNoID) HasDuplicateTransactionIds() bool`

HasDuplicateTransactionIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


