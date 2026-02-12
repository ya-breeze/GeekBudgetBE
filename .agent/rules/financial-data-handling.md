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

## Transaction Integrity & Validation

### 1. Strict Validation
- `CreateTransaction` and `UpdateTransaction` methods in `storage.go` enforce strict validation.
- All `Movement` entries MUST reference a valid `CurrencyId` and `AccountId` that exist in the database.
- Currency and Account records must be created before they can be referenced in a transaction.

### 2. Testing Requirements
- **Dependency Creation**: When writing tests involving transactions, you **MUST** create the required dependencies first.
- **Example**:
  ```go
  // Create Currency
  cur, _ := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "CZK"})
  
  // Create Account
  acc, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "MyAccount", Type: "asset"})
  
  // Create Transaction using actual IDs
  st.CreateTransaction(userID, &goserver.TransactionNoId{
      Movements: []goserver.Movement{
          {Amount: decimal.NewFromInt(100), CurrencyId: cur.Id, AccountId: acc.Id},
      },
  })
  ```
- **Do Not Use Hardcoded IDs**: Avoid using strings like `"CZK"` or `"acc-1"` unless you have explicitly created records with those specific IDs (which is generally discouraged in favor of letting the system generate UUIDs, except for currencies where names might act as keys in some legacy tests, but standardized UUIDs are preferred).
