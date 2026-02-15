---
trigger: always_on
---

- always run 'make all' after changes to make sure that it's not broken.
- Ensure decimal values in Go test struct literals use `decimal.NewFromFloat()` instead of raw float constants to avoid type mismatch errors.
- **Decimal Comparisons**: When comparing `decimal.Decimal` in tests using Gomega, use the `.Equal()` method (e.g., `Expect(val.Equal(expected)).To(BeTrue())`) instead of `Expect(val).To(Equal(expected))`. The latter is scale-sensitive and may fail if one value is `10` and the other is `10.00`.
- **OpenAPI Regeneration**: Always run `make generate` after modifying `api/openapi.yaml` or any mustache templates in `backend/pkg/generated/templates/goserver/`. Follow this with `make all` to verify that the generated code is correct and consistent.
- **GoServer Generated Models**: Be aware of field naming conventions in `pkg/generated/goserver`:
    - `Currency` uses `Id`, not `CurrencyId`.
    - `CurrencyNoId` uses `Name` and `Description`, not `Symbol`.
    - `AccountNoId` fields like `BankInfo` are complex types (e.g., `BankAccountInfo`), which may require further nested types like `BankAccountInfoBalancesInner`. Always check the generated file (e.g., `model_account_no_id.go`) for the exact structure.

## Interface Changes & Mocks
- **Generate Mocks**: When adding methods to `Storage` or other interfaces used in mocks, ALWAYS run `make generate_mocks` to regenerate the `mocks/` package. The build WILL FAIL if the mock implementation doesn't match the interface.

## Database & Models
- **UUID Parsing**: Use `uuid.Parse()` carefully. The `storage` layer often converts string IDs to UUIDs. Ensure error handling or proper parsing is in place.
- **GORM Updates**: When updating fields, `db.Save(&t)` updates ALL fields. If you only want to update specific columns, use `db.Model(&t).Updates(...)` or equivalent.
- **Transaction Preservation**: `UpdateTransaction` preserves fields by copying them from the *old* transaction retrieved from the DB. Ensure safe copies are made if necessary. `UpdateTransactionInternal` bypasses this preservation.

## Testing Strategy
- **Backend Tests**: Use `make test` to run all tests. `go test ./...` might miss integration tests if not configured correctly.
- **Integration Tests**: Located in `backend/test/`. These are critical for verifying complex interactions like auto-matching.
- **`IsAuto` Updates**: If a test expects `IsAuto` to be updated, it must use `UpdateTransactionInternal`. If it expects it preserved, use `UpdateTransaction`.
