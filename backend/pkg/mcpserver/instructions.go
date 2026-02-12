package mcpserver

const instructions = `GeekBudget is a personal finance management application.

## Data Model

- **Accounts** have a type: "asset" (bank accounts, cash), "expense" (categories like groceries, rent), or "income" (salary, etc.)
- **Currencies** are user-defined (e.g. CZK, EUR, USD). All amounts reference a currency ID.
- **Transactions** represent financial events. Each transaction has a date, description, optional place/tags/partner info, and a list of **Movements**.
- **Movements** are the core of double-entry bookkeeping: each movement transfers an amount in a specific currency to/from an account. A transaction typically has 2+ movements that balance out (e.g. -100 CZK from "Cash" account, +100 CZK to "Groceries" account).
- **Matchers** are regex-based rules that auto-categorize imported bank transactions. They match on description, partner name, partner account number, currency, place, or keywords. Matchers have a confirmation history tracking their accuracy.
- **Budget Items** define expected spending per account/category.
- **Bank Importers** connect to banks (FIO, Revolut, KB) to fetch transactions automatically.
- **Reconciliation** compares the app's computed balance against the bank's reported balance for asset accounts.

## Key Relationships

- Transactions → Movements → Accounts + Currencies (each movement references one account and one currency)
- Matchers → Accounts (output account for auto-categorization)
- Bank Importers → Accounts (target account for imported transactions)
- Reconciliation → Accounts + Currencies

## Conventions

- All IDs are UUIDs
- Money amounts are decimal numbers (never floating point)
- Dates are ISO 8601 (YYYY-MM-DD)
- Tags are string arrays on transactions for custom labeling
- "Suspicious" transactions have issues flagged by the system (e.g. unbalanced movements)
- "Unprocessed" transactions are imported but not yet categorized
- Duplicate detection flags transactions with similar dates (±2 days) and amounts from different import sources

## Common Queries

- To understand spending: list transactions filtered by date, look at movements going to expense accounts
- To check balances: use get_account_balance for a specific account+currency, or financial_summary for an overview
- To find categorization issues: list_transactions with onlySuspicious=true, or check matchers
- To verify bank sync: get_reconciliation_status shows delta between app and bank balances
`
