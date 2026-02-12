# Bank Importers

Guidelines and preferences for implementing and maintaining bank transaction importers.

## Date Selection

When importing transactions, always prioritize the date the transaction actually **happened** (execution date) over the date it was confirmed/booked by the bank.

### Importer Specifics

- **KB (Komerční banka)**:
  - Header: `Datum zauctovani;Datum provedeni;...`
  - Use **`Datum provedeni`** (Column 2, Index 1).
  - Do NOT use `Datum zauctovani` (Column 1, Index 0).
  - **Balance Date**: Searches for header fields like `Datum vypisu`, `Vytvoreno`, or `Datum exportu`. If missing, it uses the date of the **newest transaction** in the file.

  - **Revolut**:
  - Columns: `Type`, `Product`, `Started Date`, `Completed Date`, ...
  - Supports English and **Russian** localized headers.
  - Use **`Started Date`** (Column 3, Index 2).
  - User explicitly prefers `Started Date` over `Completed Date` for Revolut.
  - **Balance Date**: Derives the date from the **newest transaction** in the imported block for each currency.
  - **Russian Localization**:
    - Headers matched: `Тип`, `Продукт`, `Дата начала`, `Дата выполнения`, `Описание`, `Сумма`, `Комиссия`, `Валюта`, `State`, `Остаток средств`.
    - Status: `ВЫПОЛНЕНО` is treated as `COMPLETED`.
    - Exchange Prefix: `Обмен валюты: Обменено на ` is used to identify currency exchange transactions for joining.

- **Fio**:
  - Uses the single available date column for transactions.
  - **Balance Date**: Extracts `DateEnd` from the statement metadata.
  - **Precision**: If `DateEnd` is today, the importer uses `time.Now()` to include the current time in the "Balance Date" display.

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

## Frontend Configuration

Importers configuration in the frontend (e.g., in `bank-importer-upload-dialog.component.ts`) defines allowed extensions and help text.
- **Revolut**: Supports both `.xlsx` and `.csv`.
- **Format Toggle**: If an importer supports multiple formats, set `fixedFormat: false` in `IMPORTER_CONFIGS` to show the format selector.
