# notifications Specification

## Purpose

Notifications surface system events to the user (balance results, import failures, detected
duplicates). Each is scoped to a family with a `Date`, `Type`, `Title`, `Description`, and optional
`URL`.

## Requirements

### Requirement: Notification types

A notification SHALL have one of these types: `other`, `balanceMatch`, `balanceDoesntMatch`,
`error`, `info`, `duplicateDetected`.

#### Scenario: Type is persisted and returned
- **WHEN** a notification is created with a given type
- **THEN** the same type is returned when the notification is read

### Requirement: Balance check notifications

A balance mismatch during reconciliation SHALL create a `balanceDoesntMatch` notification describing
the account, app balance, bank balance, and currency.

#### Scenario: Mismatch notification content
- **WHEN** an account balance does not match the bank balance
- **THEN** a `balanceDoesntMatch` notification is created naming the account and both balances

### Requirement: Import failure notifications

Bank import failures SHALL create `error` notifications: one when an importer is stopped, and one
when a `FetchAll` attempt fails and is reset.

#### Scenario: Stopped importer notification
- **WHEN** a scheduled importer fails and is stopped
- **THEN** an `error` notification about the stopped importer is created

#### Scenario: FetchAll failure notification
- **WHEN** a `FetchAll` import fails
- **THEN** an `error` notification is created and `FetchAll` is reset

### Requirement: Duplicate detection notification

When duplicate detection flags transactions for a family, the system SHALL create a
`duplicateDetected` notification with the count.

#### Scenario: Duplicate notification
- **WHEN** one or more duplicate pairs are detected for a family
- **THEN** a `duplicateDetected` notification summarizing the count is created

### Requirement: Listing notifications

Notifications SHALL be listable per family, ordered by date.

#### Scenario: List notifications
- **WHEN** notifications are requested
- **THEN** the family's notifications are returned ordered by date
