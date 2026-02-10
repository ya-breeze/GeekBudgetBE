# Bank Importers

Guidelines and preferences for implementing and maintaining bank transaction importers.

## Date Selection

When importing transactions, always prioritize the date the transaction actually **happened** (execution date) over the date it was confirmed/booked by the bank.

### Importer Specifics

- **KB (Komerční banka)**:
  - Header: `Datum zauctovani;Datum provedeni;...`
  - Use **`Datum provedeni`** (Column 2, Index 1).
  - Do NOT use `Datum zauctovani` (Column 1, Index 0).

- **Revolut**:
  - Columns: `Type`, `Product`, `Started Date`, `Completed Date`, ...
  - Use **`Started Date`** (Column 3, Index 2).
  - User explicitly prefers `Started Date` over `Completed Date` for Revolut.

- **Fio**:
  - Uses the single available date column.

## Testing

Backend bank importer tests use Ginkgo. To run them:

```bash
cd backend
go tool github.com/onsi/ginkgo/v2/ginkgo -r ./pkg/bankimporters
```

## Implementation Details

- All converters implement the `ParseAndImport` method.
- Use `time.ParseInLocation` with "Europe/Prague" timezone for Czech banks.
- Raw CSV/XLSX records should be stored in the `UnprocessedSources` field as JSON for debugging.
