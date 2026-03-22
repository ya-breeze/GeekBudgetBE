# Transaction Templates — Design Spec

**Date:** 2026-03-22
**Status:** Approved
**Scope:** Angular frontend + Go backend

---

## Overview

Allow users to save transaction templates from existing transactions, then use those templates to quickly create new transactions with pre-filled values. Templates are linked to accounts via their movements, enabling account-based filtering.

---

## Data Model

New `TransactionTemplate` DB model in `backend/pkg/database/models/template.go`:

| Field | Type | GORM tag | Notes |
|---|---|---|---|
| `ID` | `uuid.UUID` | `gorm:"primaryKey"` | Primary key |
| `Name` | `string` | | User-given label (required) |
| `Description` | `string` | | Pre-fills transaction description |
| `Place` | `string` | | Pre-fills place |
| `Tags` | `[]string` | `gorm:"serializer:json"` | JSON serialized |
| `PartnerName` | `string` | | Pre-fills partner name |
| `Extra` | `string` | | Pre-fills extra/reference |
| `Movements` | `[]goserver.Movement` | `gorm:"serializer:json"` | JSON serialized — same pattern as `Transaction.Movements` |
| `UserID` | `string` | `gorm:"index"` | Multi-user isolation |
| `CreatedAt` | `time.Time` | | GORM auto |
| `UpdatedAt` | `time.Time` | | GORM auto |

**Omitted intentionally:** `Date` (set at creation time), `ExternalIDs`, `IsAuto`, `MatcherID`, `SuspiciousReasons`, `DuplicateDismissed` — runtime/import fields not relevant to templates.

**Account association:** Derived from `Movements[*].AccountId`. A template appears in filters for any account present in its movements. No denormalized column needed — the dataset is small and always user-scoped.

---

## API

New resource `/v1/templates` added to `api/openapi.yaml`.

### Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/templates` | List templates; optional `?accountId=` filter |
| `POST` | `/v1/templates` | Create template |
| `PUT` | `/v1/templates/{id}` | Update template |
| `DELETE` | `/v1/templates/{id}` | Delete template |

> `GET /v1/templates/{id}` is intentionally omitted — no frontend use case exists. The list endpoint is sufficient for all picker and management workflows.

### Schemas

- **`TransactionTemplate`** — full response object: `id`, `name` (required), `description`, `place`, `tags`, `partnerName`, `extra`, `movements`, `createdAt`, `updatedAt`
- **`TransactionTemplateNoId`** — create/update request body: same fields as above minus `id` and timestamps. **`name` is `required` in the OpenAPI schema.** Movements must contain at least one entry.

### Query Parameters

- `GET /v1/templates?accountId={uuid}` — server-side filter; returns only templates whose movements contain the given accountId. Filtering done in Go (not SQL) — fetch all user templates, filter in-memory.

All endpoints require JWT authentication. `UserID` extracted from context via `constants.UserIDKey`.

---

## Backend Implementation

### Storage

New `TemplateStorage` interface added to `backend/pkg/database/storage.go` and composed into the main `Storage` interface.

**All ID parameters use `string` at the interface boundary** (consistent with every other storage interface — `uuid.Parse` is called inside the implementation):

```go
type TemplateStorage interface {
    CreateTemplate(userID string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
    GetTemplates(userID string, accountID *string) ([]goserver.TransactionTemplate, error)
    UpdateTemplate(userID string, id string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
    DeleteTemplate(userID string, id string) error
}
```

Return types use **value types** (not pointers), consistent with every other storage interface method in the project (e.g. `CreateAccount`, `GetAccounts`, `UpdateAccount`).

Implementation in new file `backend/pkg/database/storage_template.go`, following the same patterns as `storage_account.go`:
- GORM queries scoped by `UserID`
- `TemplateToDB` / `TemplateFromDB` conversion helpers

### Migration

Add `&models.TransactionTemplate{}` to the explicit `db.AutoMigrate(...)` call in `backend/pkg/database/migration.go`. GORM does **not** auto-discover models — every model must be registered manually.

### Mock Regeneration

After updating `storage.go`, run `make generate_mocks` to regenerate `backend/pkg/database/mocks/mock_storage.go`. Without this, all existing tests using `MockStorage` will fail to compile.

### Input Validation

The **handler** (`api_templates.go`) validates that the movements array is non-empty before calling storage. Returns `400 Bad Request` if empty. Consistent with the handler-level validation pattern used elsewhere.

### Handlers

New file `backend/pkg/server/api/api_templates.go` implementing the generated interface. Handlers call `s.db` directly (no separate service layer — consistent with accounts, currencies, matchers). No `WithContext` call needed — all existing HTTP handlers in `pkg/server/api/` call storage directly on `s.db` (which holds a background context from construction). As a result, template CRUD mutations will be audit-logged as `ChangeSourceSystem`, consistent with every other API handler in the project (accounts, currencies, matchers, etc.). This is a known limitation of the current codebase — the request context with `ChangeSourceUser` from the auth middleware is not threaded through to storage in any handler.

### Code Generation

1. Update `api/openapi.yaml` with new schemas and endpoints
2. Run `make generate` to produce generated interface and models in `backend/pkg/generated/`
3. Run `make generate_mocks` to update `MockStorage`
4. Implement the generated interface in `api_templates.go`

---

## Angular Frontend

### 1. Template Management Page

**Location:** `frontend/src/app/features/templates/`

New feature module with:
- Template list view (name, accounts involved, amounts)
- Create/edit template dialog
- Delete confirmation
- "Save as template" action on existing transaction detail view — opens dialog pre-populated from transaction fields, user sets a name

**Navigation:** New sidebar menu item (alongside Transactions, Accounts, etc.)

### 2. "New from template" Entry Point

On the transaction list page: a dropdown/split button next to the existing "New transaction" button. Opens the shared `TemplatePickerComponent`, selecting a template opens the create transaction form pre-filled.

### 3. "Use template" Inside Create Form

A "Use template" button at the top of the new transaction form. Opens the same `TemplatePickerComponent`. Selecting a template populates all form fields; user edits as needed before submitting.

### Shared Component: `TemplatePickerComponent`

- Accepts optional `@Input() accountId: string` to pre-filter the list
- Searchable list of templates (filter by name)
- Used in both entry points (transaction list and create form)
- Calls `GET /v1/templates?accountId=` when accountId is provided

### Angular Service

New `TemplatesService` in the templates feature module, using the generated API client functions from `core/api/fn/templates/` (same pattern as `MatcherService`, `BudgetItemService`):

- `getTemplates(accountId?: string): Observable<TransactionTemplate[]>`
- `createTemplate(t: TransactionTemplateNoId): Observable<TransactionTemplate>`
- `updateTemplate(id: string, t: TransactionTemplateNoId): Observable<TransactionTemplate>`
- `deleteTemplate(id: string): Observable<void>`
- `templateToTransactionNoId(t: TransactionTemplate): TransactionNoId` — converts template to create-transaction payload:
  - Sets `date` to today
  - Copies `description`, `place`, `tags`, `partnerName`, `extra`, `movements`
  - Explicitly zeroes out all import-only fields: `externalIds: []`, `isAuto: false`, `matcherId: undefined`, `suspiciousReasons: []`, `duplicateDismissed: false`, `mergedIntoId: undefined`

---

## Error Handling

- 404 on `PUT/DELETE /v1/templates/{id}` if template not found or belongs to different user
- 400 if movements array is empty on create/update (validated in handler before storage call)
- Frontend shows snackbar errors consistent with existing Angular patterns

---

## Testing

- **Backend unit:** Ginkgo specs in `backend/pkg/database/storage_template_test.go` covering CRUD + account filter (uses real SQLite in-memory DB, consistent with other storage tests)
- **Backend handler:** Ginkgo specs in `backend/pkg/server/api/api_templates_test.go` with `MockStorage`
- **Frontend:** Karma/Jasmine specs for `TemplatesService` and `TemplatePickerComponent`

---

## Out of Scope

- Next.js frontend (deferred until Angular implementation is complete)
- Template ordering / favorites
- Template sharing between users
- Recurring transaction scheduling from templates
- `GET /v1/templates/{id}` single-fetch endpoint (no frontend use case)
