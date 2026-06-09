# reconciliation Specification

## Purpose

Reconciliation confirms that the app's computed balance for an asset account matches the bank's
reported balance. A `Reconciliation` record captures a checkpoint: `ReconciledBalance` (computed),
`ExpectedBalance` (bank or manual), and whether it was manual.

## Requirements

### Requirement: Automatic balance check after import

After an account's unprocessed transactions are cleared, a balance check SHALL compare the computed
balance to the bank's `ClosingBalance` within a tolerance of 0.01.

#### Scenario: Balance check deferred while unprocessed exist
- **GIVEN** an account still has unprocessed transactions
- **WHEN** a balance check is requested
- **THEN** the check is skipped

#### Scenario: Matching balance records reconciliation
- **GIVEN** an account with no unprocessed transactions
- **WHEN** the app balance equals the bank closing balance within tolerance (0.01)
- **THEN** a non-manual reconciliation record is created

#### Scenario: Mismatching balance notifies
- **GIVEN** an account with no unprocessed transactions
- **WHEN** the app balance differs from the bank balance by more than tolerance
- **THEN** a `balanceDoesntMatch` notification is created and no reconciliation record is written

### Requirement: Reconciliation eligibility

Only `asset` accounts SHALL participate in reconciliation status; accounts with a bank importer are
auto-reconciled and others can be manually reconciled.

#### Scenario: Status lists asset accounts
- **WHEN** reconciliation status is requested
- **THEN** it covers asset accounts, indicating which have bank importers

### Requirement: Manual reconciliation

A user SHALL be able to enable reconciliation for an account by providing an initial balance,
creating a manual reconciliation record (`IsManual = true`).

#### Scenario: Enable manual reconciliation
- **WHEN** a user enables reconciliation for an account with an initial balance
- **THEN** a manual reconciliation record is created with that balance as both reconciled and expected

### Requirement: History and transactions since checkpoint

Reconciliation history per account+currency SHALL be retrievable, as SHALL the transactions
affecting that account+currency since the latest reconciliation.

#### Scenario: Transactions since last reconciliation
- **WHEN** transactions-since-reconciliation is requested for an account+currency
- **THEN** transactions dated on/after the latest reconciliation that touch that account+currency are returned

### Requirement: Invalidation on back-dated changes

Editing a transaction dated before an existing reconciliation SHALL invalidate reconciliations from
that date forward.

#### Scenario: Back-dated edit invalidates reconciliations
- **GIVEN** a reconciliation exists for an account+currency
- **WHEN** a transaction dated before that reconciliation is changed
- **THEN** reconciliations from that date forward are invalidated

### Requirement: Disbalance analysis

When a balance does not match, the user SHALL be able to request an analysis that searches
transactions since the last reconciliation for a subset explaining a target delta.

#### Scenario: Analyze a disbalance
- **GIVEN** a known balance delta for an account+currency
- **WHEN** disbalance analysis is requested with that target delta
- **THEN** candidate transactions that could explain the delta are returned
