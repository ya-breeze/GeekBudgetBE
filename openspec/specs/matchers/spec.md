# matchers Specification

## Purpose

Matchers are rules that categorize unprocessed transactions by filling the empty movement's account
and applying an output description/tags. They support a regex mode and a simplified keyword mode, and
track their own accuracy via a confirmation history.

## Requirements

### Requirement: Regex matching mode

A standard matcher SHALL match a transaction against any configured combination of regular
expressions (`DescriptionRegExp`, `PartnerNameRegExp`, `PartnerAccountNumberRegExp`,
`CurrencyRegExp`, `PlaceRegExp`, `ExtraRegExp`), requiring all configured patterns to match.

#### Scenario: Match on description
- **GIVEN** a matcher with `DescriptionRegExp` set
- **WHEN** an unprocessed transaction's description matches the pattern
- **AND** all other configured patterns also match
- **THEN** the matcher applies to the transaction

#### Scenario: Non-matching field rejects the matcher
- **GIVEN** a matcher with both description and partner-name patterns
- **WHEN** the description matches but the partner name does not
- **THEN** the matcher does not apply

### Requirement: Simplified keyword mode

A matcher with `Simplified = true` SHALL match a list of `Keywords` against the description, place,
or partner name, using the matched keyword to drive the output description and adding it as a tag.

#### Scenario: Keyword match applies output
- **GIVEN** a simplified matcher with keyword "Lidl"
- **WHEN** an unprocessed transaction's description, place, or partner name contains "Lidl"
- **THEN** the matcher applies, using the keyword's output as description and adding the keyword as a tag

### Requirement: Matcher output

When a matcher applies, empty `AccountId` movements SHALL be filled with `OutputAccountId`, and
`OutputDescription` and `OutputTags` SHALL be applied to the resulting transaction.

#### Scenario: Output fills the empty account
- **GIVEN** an unprocessed transaction with one empty-account movement
- **WHEN** a matcher applies
- **THEN** the empty movement's account is set to the matcher's `OutputAccountId`

### Requirement: Matcher suggestions

The unprocessed-transactions view SHALL run all matchers against each unprocessed transaction and
return the matching matchers as suggestions with detailed match information.

#### Scenario: Multiple matchers suggested
- **GIVEN** several matchers that match a transaction
- **WHEN** the unprocessed transaction is viewed
- **THEN** all matching matchers are returned as suggestions

### Requirement: Confirmation history

Each matcher SHALL keep a rolling `ConfirmationHistory` of booleans capped at a configurable maximum
(default 10), exposing `ConfirmationsCount` and `ConfirmationsTotal`.

#### Scenario: Confirmation recorded on accept
- **WHEN** a user accepts a matcher's suggestion to convert a transaction
- **THEN** a `true` confirmation is appended to that matcher's history

#### Scenario: History is capped
- **GIVEN** a matcher whose history is at the maximum length
- **WHEN** a new confirmation is added
- **THEN** the oldest entry is dropped so the length stays at the maximum
