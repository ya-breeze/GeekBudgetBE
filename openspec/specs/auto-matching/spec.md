# auto-matching Specification

## Purpose

Trusted matchers can convert unprocessed transactions automatically, without user action. Trust is
earned through a clean confirmation history ("perfect match"), and conversion is guarded against
creating duplicates.

## Requirements

### Requirement: Perfect-match threshold

A matcher SHALL qualify for auto-matching only when its `ConfirmationHistory` has at least 10 entries
and every entry is `true`.

#### Scenario: Matcher below threshold does not auto-match
- **GIVEN** a matcher with fewer than 10 confirmations
- **THEN** it does not auto-convert any transactions

#### Scenario: Matcher with any failure does not auto-match
- **GIVEN** a matcher with 10+ confirmations where at least one is `false`
- **THEN** it does not auto-convert any transactions

#### Scenario: Perfect matcher auto-converts
- **GIVEN** a matcher with 10+ confirmations all `true`
- **WHEN** auto-matching runs for it
- **THEN** unprocessed transactions it matches are converted automatically

### Requirement: Triggered after manual confirmation

After a user manually confirms a conversion against a matcher, the system SHALL run that matcher
against all remaining unprocessed transactions for the family.

#### Scenario: Confirmation cascades to similar transactions
- **GIVEN** a perfect matcher
- **WHEN** a user confirms one conversion with that matcher
- **THEN** the matcher is run against all other unprocessed transactions and converts those it matches

### Requirement: Auto-converted transactions are marked

A transaction converted by auto-matching SHALL be flagged `IsAuto = true` and record the converting
`MatcherID`.

#### Scenario: Auto conversion sets provenance
- **WHEN** a transaction is auto-converted
- **THEN** it is stored with `IsAuto = true` and the converting matcher's id

### Requirement: Duplicate guard during auto-matching

Before auto-converting, the system SHALL check for a potential duplicate (similar amount and date);
if one is found, conversion SHALL be skipped and `AutoMatchSkipReason` recorded instead.

#### Scenario: Skip auto-conversion on potential duplicate
- **GIVEN** an unprocessed transaction matched by a perfect matcher
- **AND** a similar existing transaction is detected
- **THEN** the transaction is not auto-converted and `AutoMatchSkipReason` is set

### Requirement: Respect ignore-before cutoffs

Auto-matching SHALL skip unprocessed transactions excluded by an account's
`IgnoreUnprocessedBefore` cutoff.

#### Scenario: Cutoff excludes old transactions from auto-match
- **GIVEN** an account with `IgnoreUnprocessedBefore` set
- **AND** an unprocessed transaction dated before the cutoff
- **THEN** auto-matching does not convert it
