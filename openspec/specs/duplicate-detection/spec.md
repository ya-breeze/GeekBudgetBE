# duplicate-detection Specification

## Purpose

A background task finds transactions that are likely the same financial event imported from two
different sources (e.g. a transfer appearing on both bank accounts) and flags them for the user to
resolve, while avoiding false positives on genuinely distinct transactions.

## Requirements

### Requirement: Scheduled background scan

Duplicate detection SHALL run every 24 hours (after a startup delay) and scan transactions from the
last 30 days per family.

#### Scenario: Periodic scan
- **GIVEN** the server has been running past the initial delay
- **THEN** duplicate detection runs and then repeats every 24 hours

#### Scenario: Scan window
- **WHEN** duplicate detection runs
- **THEN** it only considers transactions from the last 30 days

### Requirement: Duplicate criteria

Two transactions SHALL be flagged as duplicates only when all hold: dates within ±2 days, matching
movement amounts, different external ids, and opposite net directions.

#### Scenario: Cross-account transfer detected
- **GIVEN** two transactions from different importers within 2 days, equal amounts, opposite directions
- **THEN** they are linked as potential duplicates

#### Scenario: Same-source transactions are not flagged
- **GIVEN** two transactions that share an external id
- **THEN** they are not flagged, since same-source dedup is handled at import time

#### Scenario: Two same-direction purchases are not flagged
- **GIVEN** two separate outgoing purchases of the same amount on consecutive days
- **THEN** they are not flagged, because their net directions are not opposite

#### Scenario: Dismissed transactions are skipped
- **GIVEN** a transaction with `DuplicateDismissed = true`
- **THEN** it is not considered during detection

### Requirement: Linking and marking

Detected duplicates SHALL be linked bidirectionally in the `TransactionDuplicate` table and each
marked with the duplicate reason in `SuspiciousReasons`.

#### Scenario: Bidirectional link created
- **WHEN** two transactions are detected as duplicates
- **THEN** the relationship is stored in both directions and both transactions are marked suspicious

### Requirement: Detection notification

When duplicates are detected for a family, the system SHALL create a single `duplicateDetected`
notification summarizing the count.

#### Scenario: Notification on detection
- **GIVEN** duplicate detection flags one or more pairs for a family
- **THEN** a notification of type `duplicateDetected` is created with the count

### Requirement: Manual trigger

Duplicate detection SHALL also be invokable on demand via a CLI command, independent of the schedule.

#### Scenario: Manual run
- **WHEN** the duplicate-detection command is invoked
- **THEN** detection runs immediately using the same criteria
