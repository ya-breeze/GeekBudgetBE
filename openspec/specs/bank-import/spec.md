# bank-import Specification

## Purpose

Bank importers bring transactions in from external sources. Each importer has a `Type` (`fio`, `kb`,
`revolut`), a target `AccountID`, an optional `FeeAccountID`, and produces **unprocessed**
transactions (the source account is known; the categorization account is empty).

## Requirements

### Requirement: Importer types and sources

FIO SHALL be fetched from the bank's API; KB and Revolut SHALL be imported by uploading a statement
file (XLSX for KB, CSV for Revolut), which is parsed by the matching converter. Uploading a file for
an unsupported importer type SHALL fail and record an error import result.

#### Scenario: FIO fetches from API
- **GIVEN** a bank importer of type `fio`
- **WHEN** an import runs
- **THEN** transactions are fetched from the FIO API

#### Scenario: KB and Revolut parse uploaded files
- **GIVEN** a bank importer of type `kb` or `revolut`
- **WHEN** a statement file is uploaded
- **THEN** the matching converter parses the file and extracts transactions

#### Scenario: Unsupported type on upload
- **GIVEN** a bank importer whose type is not `kb` or `revolut`
- **WHEN** a statement file is uploaded for it
- **THEN** the operation fails and an error import result is recorded

### Requirement: Execution-date preference

Importers SHALL record the date the transaction actually happened (execution date), not the bank's
booking date.

#### Scenario: KB uses execution date
- **WHEN** a KB statement is parsed
- **THEN** the transaction date is taken from `Datum provedeni` (execution), not `Datum zauctovani` (booking)

#### Scenario: Revolut uses started date
- **WHEN** a Revolut statement is parsed
- **THEN** the transaction date is taken from `Started Date`, not `Completed Date`

### Requirement: Deduplication by external id

Imported transactions SHALL carry `ExternalIDs`, and a transaction already present (matching external
id) SHALL NOT be imported again.

#### Scenario: Re-import skips known transactions
- **GIVEN** a transaction with a given external id already exists
- **WHEN** an import that contains that external id runs
- **THEN** the existing transaction is not duplicated

### Requirement: Scheduled FIO imports

A background task SHALL run on a 24-hour timer, fetching only importers of type `fio` (other types
are skipped). On failure to load importers it SHALL retry after 1 hour. After each scheduled batch it
SHALL process unprocessed transactions for auto-conversion.

#### Scenario: Scheduled cycle fetches FIO importers
- **WHEN** the scheduled import cycle runs
- **THEN** each `fio` importer is fetched and non-`fio` importers are skipped
- **AND** the next cycle is scheduled 24 hours later

#### Scenario: Retry on importer load failure
- **WHEN** the scheduled cycle cannot load the list of importers
- **THEN** it retries after 1 hour instead of waiting a full day

### Requirement: Triggered imports on change

Creating or updating a bank importer SHALL queue an immediate forced import for that importer.

#### Scenario: Importer create triggers immediate import
- **WHEN** a new bank importer is created
- **THEN** a forced import is queued for that importer

#### Scenario: Importer update triggers immediate import
- **WHEN** a bank importer is updated
- **THEN** a forced import is queued for that importer

### Requirement: Failure handling and stopping

On failure, a non-interactive importer with `FetchAll = false` SHALL be stopped (`IsStopped = true`)
with a notification; a failed `FetchAll` SHALL reset `FetchAll = false` and notify.

#### Scenario: Failure stops a scheduled importer
- **GIVEN** a non-interactive importer with `FetchAll = false`
- **WHEN** its import fails
- **THEN** `IsStopped` is set to true and an error notification is created

#### Scenario: A stopped importer is skipped
- **GIVEN** an importer with `IsStopped = true`
- **WHEN** a scheduled (non-interactive) import cycle runs
- **THEN** that importer is skipped

#### Scenario: Successful import clears stopped flag
- **GIVEN** an importer with `IsStopped = true`
- **WHEN** an import succeeds (e.g. interactively)
- **THEN** `IsStopped` is reset to false

### Requirement: Import results and balances

Each import SHALL append an `ImportResult` to `LastImports`, update `LastSuccessfulImport`, update
the target account's `BankInfo.Balances`, and trigger a balance check.

#### Scenario: Successful import records a result and updates balance
- **WHEN** an import succeeds
- **THEN** an import result is recorded and the account's closing balance is updated
- **AND** a balance check is performed for the account

### Requirement: Importers can be globally disabled

When `DisableImporters` is set in config, the bank import background task SHALL NOT run.

#### Scenario: Importers disabled
- **GIVEN** `DisableImporters` is true
- **THEN** the bank import background task does not start
