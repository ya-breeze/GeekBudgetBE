# currencies Specification

## Purpose

Currencies are user-defined (e.g. CZK, EUR, USD) and referenced by every Movement. Exchange rates
against CZK are fetched from the Czech National Bank (CNB) in the background to support conversion.

## Requirements

### Requirement: User-defined currencies

A currency SHALL have a `Name` and optional `Description`, be scoped to a family, and be identified
by UUID.

#### Scenario: Create a currency
- **WHEN** a user creates a currency named "EUR"
- **THEN** the currency is created with a generated UUID scoped to the family

#### Scenario: Currency referenced by movements cannot be deleted
- **GIVEN** a currency referenced by at least one movement
- **WHEN** the currency is deleted without a replacement
- **THEN** the operation fails with a "currency is in use" error

### Requirement: CNB exchange-rate fetching

A background task SHALL fetch daily exchange rates from the Czech National Bank, stored as
`CNBCurrencyRate` records (`CurrencyCode`, `RateToCZK`, `RateDate`).

#### Scenario: Daily rate fetch
- **GIVEN** the rate fetcher is enabled
- **THEN** it fetches the current day's rates and then sleeps for 24 hours

#### Scenario: Rate fetch failure retries hourly
- **WHEN** a rate fetch fails
- **THEN** the fetcher retries after 1 hour instead of waiting a full day

#### Scenario: Rate fetching can be disabled
- **GIVEN** `DisableCurrenciesRatesFetch` is set in config
- **THEN** the rate fetcher does not start

### Requirement: Currency conversion

The system SHALL convert an amount from one currency to another for a given date using stored CNB
rates, anchoring both legs through CZK.

#### Scenario: Convert between two currencies
- **WHEN** an amount in USD is converted to EUR for a date
- **THEN** the conversion uses the CZK-anchored rates for that date
