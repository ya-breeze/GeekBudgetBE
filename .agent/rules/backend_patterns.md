# Backend Patterns & Best Practices

## Background Tasks

When implementing long-running or periodic background tasks (e.g. importers, cleanup jobs), follow these patterns to ensure robustness and testability.

### Avoid Busy Loops
**Do not use** `select` with a `default` branch for periodic tasks. If logic within the loop inadvertently skips a delay (e.g. via `continue`), it can cause a tight busy loop that consumes 100% CPU.

### Use `time.Timer` for Scheduling
Instead, use an explicit `time.Timer`:
1. Initialize the timer (e.g. to fire immediately or after a delay).
2. Wait on `timer.C`, context cancellation, or other triggers (e.g. forced run channels).
3. Perform the work.
4. Reset the timer with the next delay.

**Example:**
```go
timer := time.NewTimer(0)
defer timer.Stop()

for {
    select {
    case <-ctx.Done():
        return
    case <-forcedRun:
        if !timer.Stop() {
            select {
            case <-timer.C:
            default:
            }
        }
    case <-timer.C:
    }

    // Do work...
    timer.Reset(nextDelay)
}
```

## Update Operations

### Generic Update Helper (API Layer)
Use the `updateEntity` helper in `backend/pkg/server/api/helpers.go` to handle the standard update flow (UserID extraction, DB update call, and error mapping). This reduces boilerplate in API handlers.

### Generic Update Helper (Storage Layer)
Use the `performUpdate` helper in `backend/pkg/database/storage.go` for standard CRUD updates. It handles the `UUID parsing -> Record fetching -> Model conversion -> Save -> Result conversion` sequence generically.

## Testing

### Mocking Storage
To test background tasks without relying on actual database timing or complex setup:
1. Create a `MockStorage` struct that embeds `database.Storage`.
2. Override only the methods needed for the test.
3. Use channels to signal when methods are called, allowing synchronization without `time.Sleep`.

**Example:**
```go
type MockStorage struct {
    database.Storage
    OnFetchData func() (Data, error)
}

func (m *MockStorage) FetchData() (Data, error) {
    if m.OnFetchData != nil {
        return m.OnFetchData()
    }
    return m.Storage.FetchData()
}
```

### Decimal Literals
When writing tests in Go, avoid using untyped float constants (e.g., `100.50`) where a `decimal.Decimal` is expected in a struct literal. This causes compiler errors. Always use `decimal.NewFromFloat(100.50)` or similar.
