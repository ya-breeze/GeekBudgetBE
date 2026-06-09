# transactions Specification

## Purpose

A transaction represents a financial event and owns a list of **Movements**. Each movement transfers
a signed `Amount` in one `Currency` to/from one `Account`. This double-entry model is the core of the
ledger: a normal transaction has 2+ movements that balance to zero across currencies.

## Requirements

### Requirement: Transaction structure

A transaction SHALL have a `Date`, optional descriptive fields (`Description`, `Place`, `Tags`,
`PartnerName`, `PartnerAccount`, `Extra`), a list of `Movements`, and be scoped to a family with a
UUID id.

#### Scenario: Create a balanced transaction
- **GIVEN** an asset account, an expense account, and a currency exist
- **WHEN** a transaction is created with a -100 movement on the asset account and a +100 movement on the expense account in that currency
- **THEN** the transaction is created with a generated UUID and both movements persisted

### Requirement: Movement referential integrity

Every movement SHALL reference a `CurrencyId` and `AccountId` that exist in the family, except that
an empty `AccountId` is permitted and marks the transaction as unprocessed.

#### Scenario: Movement with non-existent account is rejected
- **WHEN** a transaction is created with a movement referencing an account id that does not exist
- **THEN** the create fails with a validation error

#### Scenario: Empty account id is permitted (unprocessed)
- **WHEN** a transaction is created with one movement having a valid account and another movement with empty `AccountId`
- **THEN** the transaction is created and is treated as unprocessed

### Requirement: Strict update payloads

Update endpoints SHALL decode into a no-id transaction type and reject bodies containing an `id` or
other entity-level fields.

#### Scenario: Update with id field in body fails
- **WHEN** a transaction update request includes an `id` field in the JSON body
- **THEN** the request fails with `json: unknown field "id"`

### Requirement: Decimal money amounts

All movement amounts SHALL be `decimal.Decimal` and serialize to JSON as numbers, not strings.

#### Scenario: Amounts serialize as numbers
- **WHEN** a transaction is returned from the API
- **THEN** movement amounts are JSON numbers (not quoted strings)

### Requirement: Batch creation is atomic

Multiple transactions MAY be created in a single batch; if any fails, the whole batch SHALL be
rolled back.

#### Scenario: Batch rollback on failure
- **GIVEN** a batch of transactions where one references an invalid account
- **WHEN** the batch is created
- **THEN** no transactions from the batch are persisted

### Requirement: Transaction history snapshots

Create, update, and delete operations SHALL record a `TransactionHistory` entry holding a JSON
snapshot of the transaction state.

#### Scenario: Update records a snapshot
- **WHEN** a transaction is updated
- **THEN** a history entry is recorded capturing the prior state

### Requirement: Imported transactions cannot be deleted

A transaction that originated from a bank import (has external ids) SHALL NOT be deletable through
the standard delete path.

#### Scenario: Delete imported transaction is rejected
- **GIVEN** a transaction with external ids from an importer
- **WHEN** a delete is attempted
- **THEN** the operation fails with an "imported transaction cannot be deleted" error

### Requirement: Suspicious transactions

A transaction MAY carry `SuspiciousReasons` and SHALL be queryable filtering for only suspicious
ones.

#### Scenario: Filter suspicious transactions
- **WHEN** transactions are queried with the suspicious-only filter
- **THEN** only transactions with non-empty `SuspiciousReasons` are returned
