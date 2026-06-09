# accounts Specification

## Purpose

Accounts are the nodes of the double-entry ledger. Every Movement references exactly one account.
Accounts are scoped to a family (`FamilyID`), identified by UUID, and classified by type to drive
ledger semantics, reporting, and reconciliation.

## Requirements

### Requirement: Account types

An account SHALL have a `Type` of `asset` (bank accounts, cash), `expense` (spending categories), or
`income` (inflows), classifying its role in the ledger.

#### Scenario: Create an asset account
- **WHEN** a user creates an account with type `asset`
- **THEN** the account is created with a generated UUID and scoped to the user's family

#### Scenario: Account type drives reconciliation eligibility
- **WHEN** reconciliation status is requested
- **THEN** only `asset` accounts are considered for balance reconciliation

### Requirement: Display and reporting flags

An account SHALL carry visibility flags `ShowInDashboardSummary`, `ShowInReconciliation`, and
`HideFromReports` that control where it appears.

#### Scenario: Hide an account from reports
- **GIVEN** an account with `HideFromReports = true`
- **THEN** the account's balances are excluded from report aggregations

#### Scenario: Dashboard summary selection
- **GIVEN** an account with `ShowInDashboardSummary = true`
- **THEN** the account appears in the dashboard summary view

### Requirement: Lifecycle dates

An account SHALL support optional `OpeningDate` and `ClosingDate`.

#### Scenario: Open and close dates are optional
- **WHEN** an account is created without opening or closing dates
- **THEN** the account is valid and both dates are unset

### Requirement: Ignore-unprocessed-before cutoff

An account MAY set `IgnoreUnprocessedBefore`, and unprocessed transactions dated before it SHALL be
excluded from the unprocessed view and from auto-matching.

#### Scenario: Old unprocessed transactions are ignored
- **GIVEN** an account with `IgnoreUnprocessedBefore` set to a date
- **AND** an unprocessed transaction dated before that cutoff exists with a movement on this account
- **THEN** that transaction is not returned in the unprocessed-transactions list
- **AND** it is not considered for auto-matching

### Requirement: Bank account info and balances

An account SHALL store `BankInfo` including per-currency `Balances` with a bank-reported
`ClosingBalance`, updated by bank importers and used for balance checks.

#### Scenario: Importer updates closing balance
- **WHEN** a bank import completes successfully
- **THEN** the account's `BankInfo.Balances` are updated with the bank-reported closing balance per currency

### Requirement: Deletion with reassignment

Deleting an account referenced by movements SHALL require a replacement account to absorb those
movements.

#### Scenario: Delete account in use without replacement
- **GIVEN** an account referenced by at least one movement
- **WHEN** the account is deleted without a replacement account
- **THEN** the operation fails with an "account is in use" error

#### Scenario: Delete account in use with replacement
- **GIVEN** an account referenced by movements
- **WHEN** the account is deleted with a replacement account specified
- **THEN** affected movements are reassigned to the replacement and the account is removed
