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
