# duplicate-resolution Specification

## Purpose

Once transactions are flagged as potential duplicates, the user resolves them by either dismissing
the flag (false positive) or merging the pair. Merging preserves the original record via an archive.

## Requirements

### Requirement: Dismiss a duplicate

Dismissing SHALL set `DuplicateDismissed = true`, clear all duplicate relationships for the
transaction, and prevent it from being re-flagged.

#### Scenario: Dismiss clears links
- **GIVEN** a transaction flagged as a potential duplicate
- **WHEN** the user dismisses it
- **THEN** `DuplicateDismissed` is set to true and its duplicate relationships are cleared

#### Scenario: Dismissed transaction is not re-flagged
- **GIVEN** a dismissed transaction
- **WHEN** duplicate detection runs again
- **THEN** it is not flagged again

### Requirement: Merge two transactions

Merging SHALL keep one transaction and remove the other: the kept transaction absorbs the other's
external ids, and the merged transaction is soft-deleted and archived.

#### Scenario: Merge transfers external ids
- **WHEN** transactions are merged with a keep id and a merge id
- **THEN** the kept transaction gains the merged transaction's external ids

#### Scenario: Merged transaction is soft-deleted and archived
- **WHEN** a transaction is merged into another
- **THEN** it is soft-deleted from the active set and a snapshot is written to the merged-transactions archive

### Requirement: Archived merged transactions are retrievable

A merged (archived) transaction SHALL NOT be returned by the standard transaction-get endpoint but
SHALL be fetchable from the merged-transactions endpoint.

#### Scenario: Standard get excludes archived
- **GIVEN** a transaction that was merged away
- **WHEN** it is requested via the standard transaction endpoint
- **THEN** it is not returned

#### Scenario: Merged endpoint returns archived
- **GIVEN** a merged transaction
- **WHEN** it is requested via the merged-transactions endpoint by id
- **THEN** the archived snapshot is returned

### Requirement: Synchronized link cleanup

When duplicate relationships are cleared, the duplicate reason SHALL be removed from any linked
transaction that no longer has remaining duplicate links.

#### Scenario: Partner cleared when no links remain
- **GIVEN** two transactions linked as duplicates with no other links
- **WHEN** the link between them is cleared
- **THEN** the duplicate reason is removed from the partner transaction's suspicious reasons
