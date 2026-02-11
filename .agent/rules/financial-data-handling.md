# Financial Data Handling

## Decimal Migration
The project uses `github.com/shopspring/decimal` for all financial data to ensure precision.

### 1. Backend Serialization
- **CRITICAL**: The backend MUST serialize `decimal.Decimal` as numbers in JSON, not strings. 
- This is globally configured in `backend/cmd/main.go` using `decimal.MarshalJSONWithoutQuotes = true`.
- If you use a different entry point or tool, ensure this setting is applied, otherwise the frontend will receive strings and calculations may break.

### 2. Frontend Handling
- **Defensive Parsing**: Always use `Number()` when performing arithmetic on monetary fields in the frontend.
- Example: `const amount = Number(movement.amount);`
- This prevents string concatenation bugs if the API ever sends a string representation.

### 3. Database Storage
- Decimals are stored as numeric types in SQLite where possible.
- GORM handles the mapping between `decimal.Decimal` and the database.

### 4. Testing
- Do NOT use direct equality (`==`) for comparing decimals in tests.
- Use the `.Equal()` method.
- **Backend Test Pattern**: `Expect(val1.Equal(val2)).To(BeTrue())`.
- Avoid `Expect(val1).To(Equal(val2))` as it may compare internal scale which might differ even if the value is equivalent (e.g., `1.0` vs `1.00`).
