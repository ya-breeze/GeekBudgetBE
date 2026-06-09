# budget Specification

## Purpose

Budget items define expected amounts for an account (typically an expense or income category) for a
period, enabling comparison of planned vs. actual spending.

## Requirements

### Requirement: Budget item structure

A budget item SHALL have a `Date`, an `AccountID`, an `Amount` (decimal), and an optional
`Description`, scoped to a family with a UUID id.

#### Scenario: Create a budget item
- **WHEN** a user creates a budget item for an account with an amount and date
- **THEN** the budget item is created with a generated UUID scoped to the family

### Requirement: Query budget items

Budget items SHALL be listable, optionally filtered by account.

#### Scenario: List budget items for an account
- **WHEN** budget items are requested filtered by an account id
- **THEN** only budget items for that account are returned

### Requirement: Decimal amounts

Budget item amounts SHALL be stored and returned as decimals without floating-point error.

#### Scenario: Amount precision preserved
- **WHEN** a budget item is created with a fractional amount
- **THEN** the amount is stored and returned without floating-point rounding error
