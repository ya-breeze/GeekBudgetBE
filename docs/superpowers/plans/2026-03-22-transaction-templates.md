# Transaction Templates Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add transaction templates that users can create from existing transactions and use to quickly pre-fill new transactions, filterable by account.

**Architecture:** New `TransactionTemplate` entity following the existing OpenAPI-first pattern — add schemas and endpoints to `api/openapi.yaml`, run code generation, implement `TemplateStorage` interface + handler, then build Angular feature module with a shared picker component.

**Tech Stack:** Go 1.24 (GORM, SQLite, Gorilla Mux, Ginkgo/Gomega), Angular 20 (standalone components, Angular Material, RxJS signals), OpenAPI 3.0 generator.

**Spec:** `docs/superpowers/specs/2026-03-22-transaction-templates-design.md`

---

## File Map

### New files
| File | Purpose |
|---|---|
| `backend/pkg/database/models/template.go` | `TransactionTemplate` GORM model + `TemplateToDB`/`FromDB` helpers |
| `backend/pkg/database/storage_template.go` | `TemplateStorage` implementation (CRUD + account filter) |
| `backend/pkg/database/storage_template_test.go` | Ginkgo storage-layer tests |
| `backend/pkg/server/api/api_templates.go` | HTTP handlers implementing generated interface |
| `backend/pkg/server/api/api_templates_test.go` | Ginkgo handler tests with MockStorage |
| `frontend/src/app/features/templates/templates.component.ts` | Templates list page |
| `frontend/src/app/features/templates/templates.component.html` | List page template |
| `frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.ts` | Create/edit dialog |
| `frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.html` | Dialog template |
| `frontend/src/app/features/templates/template-picker/template-picker.component.ts` | Shared picker (used in form + transaction list) |
| `frontend/src/app/features/templates/template-picker/template-picker.component.html` | Picker template |
| `frontend/src/app/features/templates/template-picker/template-picker-dialog.component.ts` | Thin dialog wrapper around the picker (created in Task 10) |
| `frontend/src/app/features/templates/services/template.service.ts` | Angular service + `templateToTransactionNoId` converter |

### Modified files
| File | Change |
|---|---|
| `api/openapi.yaml` | Add `TransactionTemplate`, `TransactionTemplateNoId` schemas + `/v1/templates` endpoints |
| `backend/pkg/database/storage.go` | Add `TemplateStorage` interface + compose into `Storage` |
| `backend/pkg/database/migration.go` | Register `&models.TransactionTemplate{}` in `AutoMigrate` |
| `backend/pkg/server/server.go` (or wherever handlers are wired) | Register `TemplatesAPIServiceImpl` |
| `frontend/src/app/app.routes.ts` | Add `/templates` route |
| `frontend/src/app/layout/sidebar/sidebar.component.ts` | Add Templates menu item |
| `frontend/src/app/features/transactions/transactions.component.ts` | Add "New from template" button |
| `frontend/src/app/features/transactions/transaction-form-dialog/transaction-form-dialog.component.ts` | Add "Use template" button |
| `frontend/src/app/features/transactions/transaction-detail/transaction-detail.component.ts` | Add "Save as template" action |

---

## Task 1: OpenAPI schemas and endpoints

**Files:**
- Modify: `api/openapi.yaml`

- [ ] **Step 1: Add `TransactionTemplateNoId` schema**

Find the `components/schemas` section (after `MatcherNoID` is a good insertion point). Add:

```yaml
  TransactionTemplateNoId:
    type: object
    required:
      - name
      - movements
    properties:
      name:
        type: string
        description: User-given label for the template
      description:
        type: string
      place:
        type: string
      tags:
        type: array
        items:
          type: string
      partnerName:
        type: string
      extra:
        type: string
      movements:
        type: array
        minItems: 1
        items:
          $ref: "#/components/schemas/Movement"

  TransactionTemplate:
    type: object
    allOf:
      - $ref: "#/components/schemas/Entity"
      - $ref: "#/components/schemas/TransactionTemplateNoId"
```

- [ ] **Step 2: Add `/v1/templates` endpoints**

Find the paths section (add after `/v1/matchers`):

```yaml
  /v1/templates:
    get:
      tags: [templates]
      summary: get all templates
      operationId: getTemplates
      parameters:
        - name: accountId
          in: query
          required: false
          schema:
            type: string
            format: uuid
      responses:
        "200":
          description: templates
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/TransactionTemplate"
    post:
      tags: [templates]
      summary: create new template
      operationId: createTemplate
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/TransactionTemplateNoId"
      responses:
        "200":
          description: created template
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/TransactionTemplate"

  /v1/templates/{id}:
    put:
      tags: [templates]
      summary: update template
      operationId: updateTemplate
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/TransactionTemplateNoId"
      responses:
        "200":
          description: updated template
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/TransactionTemplate"
    delete:
      tags: [templates]
      summary: delete template
      operationId: deleteTemplate
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        "204":
          description: deleted
```

- [ ] **Step 3: Run code generation**

```bash
cd /path/to/GeekBudgetBE
make generate
```

Expected: no errors, new files appear in `backend/pkg/generated/goserver/` — look for `api_templates_service.go` and `model_transaction_template*.go`. Also new generated Angular client files in `frontend/src/app/core/api/`.

- [ ] **Step 4: Commit**

```bash
git add api/openapi.yaml backend/pkg/generated/ frontend/src/app/core/api/
git commit -m "feat: add TransactionTemplate OpenAPI schemas and /v1/templates endpoints"
```

---

## Task 2: DB model

**Files:**
- Create: `backend/pkg/database/models/template.go`
- Modify: `backend/pkg/database/migration.go`

- [ ] **Step 1: Write failing test**

Create `backend/pkg/database/storage_template_test.go` with just enough to test the model can be created:

```go
package database_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/ya-breeze/geekbudgetbe/pkg/database"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestTemplateStorage(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "TemplateStorage Suite")
}

var _ = Describe("TemplateStorage", func() {
    var db database.Storage

    BeforeEach(func() {
        var err error
        db, err = database.NewTestStorage()
        Expect(err).NotTo(HaveOccurred())
    })

    Describe("CreateTemplate", func() {
        It("creates a template and returns it with an ID", func() {
            tpl, err := db.CreateTemplate("user1", &goserver.TransactionTemplateNoId{
                Name: "Rent",
                Movements: []goserver.Movement{
                    {Amount: 1000, CurrencyId: "some-currency-id", AccountId: "some-account-id"},
                },
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(tpl.Id).NotTo(BeEmpty())
            Expect(tpl.Name).To(Equal("Rent"))
        })
    })
})
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd backend && go test ./pkg/database/... -run TestTemplateStorage -v
```

Expected: compile error — `db.CreateTemplate` method does not exist yet.

- [ ] **Step 3: Create the model file**

```go
// backend/pkg/database/models/template.go
package models

import (
    "time"

    "github.com/google/uuid"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type TransactionTemplate struct {
    ID          uuid.UUID           `gorm:"type:uuid;primaryKey"`
    Name        string
    Description string
    Place       string
    Tags        []string            `gorm:"serializer:json"`
    PartnerName string
    Extra       string
    Movements   []goserver.Movement `gorm:"serializer:json"`
    UserID      string              `gorm:"index"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (t *TransactionTemplate) FromDB() goserver.TransactionTemplate {
    return goserver.TransactionTemplate{
        Id:          t.ID.String(),
        Name:        t.Name,
        Description: t.Description,
        Place:       t.Place,
        Tags:        t.Tags,
        PartnerName: t.PartnerName,
        Extra:       t.Extra,
        Movements:   t.Movements,
        CreatedAt:   t.CreatedAt,
        UpdatedAt:   t.UpdatedAt,
    }
}

func TemplateToDB(t *goserver.TransactionTemplateNoId, userID string) *TransactionTemplate {
    tags := t.Tags
    if tags == nil {
        tags = make([]string, 0)
    }
    movements := t.Movements
    if movements == nil {
        movements = make([]goserver.Movement, 0)
    }
    return &TransactionTemplate{
        UserID:      userID,
        Name:        t.Name,
        Description: t.Description,
        Place:       t.Place,
        Tags:        tags,
        PartnerName: t.PartnerName,
        Extra:       t.Extra,
        Movements:   movements,
    }
}
```

> **Note:** Check the exact field names generated by `make generate` for `goserver.TransactionTemplate` (e.g. `CreatedAt` might be `CreatedAt time.Time` or a `*time.Time`). Adjust `FromDB()` accordingly.

- [ ] **Step 4: Register model in AutoMigrate**

Open `backend/pkg/database/migration.go`. Find the `db.AutoMigrate(...)` call and add `&models.TransactionTemplate{}` to the list:

```go
// Add alongside the other models:
&models.TransactionTemplate{},
```

- [ ] **Step 5: Run test to verify it passes**

```bash
cd backend && go test ./pkg/database/... -run TestTemplateStorage -v
```

Expected: PASS (the model can be created)

- [ ] **Step 6: Run gofumpt on new files**

```bash
go tool mvdan.cc/gofumpt -w backend/pkg/database/models/template.go backend/pkg/database/migration.go
```

- [ ] **Step 7: Commit**

```bash
git add backend/pkg/database/models/template.go backend/pkg/database/migration.go backend/pkg/database/storage_template_test.go
git commit -m "feat: add TransactionTemplate DB model and migration registration"
```

---

## Task 3: Storage interface + implementation

**Files:**
- Modify: `backend/pkg/database/storage.go`
- Create: `backend/pkg/database/storage_template.go`

- [ ] **Step 1: Add `TemplateStorage` interface to `storage.go`**

Open `backend/pkg/database/storage.go`. Find the `Storage` interface (the large `type Storage interface { ... }` block). Add:

1. New sub-interface (add before or after `MatcherStorage`):

```go
type TemplateStorage interface {
    CreateTemplate(userID string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
    GetTemplates(userID string, accountID *string) ([]goserver.TransactionTemplate, error)
    UpdateTemplate(userID string, id string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
    DeleteTemplate(userID string, id string) error
}
```

2. Compose it into `Storage`:

```go
type Storage interface {
    // ... existing interfaces ...
    TemplateStorage
    // ...
}
```

- [ ] **Step 2: Verify compilation breaks cleanly**

```bash
cd backend && go build ./...
```

Expected: build fails with `storage does not implement Storage (missing CreateTemplate method)` — this is correct, the implementation is missing.

- [ ] **Step 3: Implement `storage_template.go`**

```go
// backend/pkg/database/storage_template.go
package database

import (
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/ya-breeze/geekbudgetbe/pkg/database/models"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
    "gorm.io/gorm"
)

func (s *storage) CreateTemplate(userID string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error) {
    tpl := models.TemplateToDB(t, userID)
    tpl.ID = uuid.New()
    if err := s.db.Create(&tpl).Error; err != nil {
        return goserver.TransactionTemplate{}, fmt.Errorf(StorageError, err)
    }

    if err := s.recordAuditLog(s.db, userID, "TransactionTemplate", tpl.ID.String(), "CREATED", nil, &tpl); err != nil {
        s.log.Error("Failed to record audit log", "error", err)
    }

    return tpl.FromDB(), nil
}

func (s *storage) GetTemplates(userID string, accountID *string) ([]goserver.TransactionTemplate, error) {
    var records []models.TransactionTemplate
    if err := s.db.Where("user_id = ?", userID).Order("name").Find(&records).Error; err != nil {
        return nil, fmt.Errorf(StorageError, err)
    }

    result := make([]goserver.TransactionTemplate, 0, len(records))
    for _, r := range records {
        if accountID != nil {
            matched := false
            for _, m := range r.Movements {
                if m.AccountId == *accountID {
                    matched = true
                    break
                }
            }
            if !matched {
                continue
            }
        }
        result = append(result, r.FromDB())
    }

    return result, nil
}

func (s *storage) UpdateTemplate(userID string, id string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error) {
    return performUpdate[models.TransactionTemplate, *goserver.TransactionTemplateNoId, goserver.TransactionTemplate](
        s, userID, "TransactionTemplate", id, t,
        func(t *goserver.TransactionTemplateNoId, userID string) *models.TransactionTemplate {
            return models.TemplateToDB(t, userID)
        },
        func(m *models.TransactionTemplate) goserver.TransactionTemplate { return m.FromDB() },
        func(m *models.TransactionTemplate, id uuid.UUID) { m.ID = id },
    )
}

func (s *storage) DeleteTemplate(userID string, id string) error {
    var tpl models.TransactionTemplate
    if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&tpl).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrNotFound
        }
        return fmt.Errorf(StorageError, err)
    }

    if err := s.recordAuditLog(s.db, userID, "TransactionTemplate", id, "DELETED", &tpl, nil); err != nil {
        s.log.Error("Failed to record audit log", "error", err)
    }

    return s.db.Delete(&tpl).Error
}
```

- [ ] **Step 4: Verify build passes**

```bash
cd backend && go build ./...
```

Expected: compiles cleanly.

- [ ] **Step 5: Regenerate mocks**

```bash
make generate_mocks
```

Expected: `backend/pkg/database/mocks/mock_storage.go` is updated with `CreateTemplate`, `GetTemplates`, `UpdateTemplate`, `DeleteTemplate` methods.

- [ ] **Step 6: Run all storage tests**

```bash
cd backend && go test ./pkg/database/... -v
```

Expected: all tests pass including the `TestTemplateStorage` suite.

- [ ] **Step 7: Run gofumpt**

```bash
go tool mvdan.cc/gofumpt -w backend/pkg/database/storage.go backend/pkg/database/storage_template.go
```

- [ ] **Step 8: Commit**

```bash
git add backend/pkg/database/storage.go backend/pkg/database/storage_template.go backend/pkg/database/mocks/mock_storage.go
git commit -m "feat: add TemplateStorage interface and implementation"
```

---

## Task 4: Storage tests (full suite)

**Files:**
- Modify: `backend/pkg/database/storage_template_test.go`

- [ ] **Step 1: Write the full test suite**

Replace the placeholder test from Task 2 with the full suite:

```go
package database_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/ya-breeze/geekbudgetbe/pkg/database"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestTemplateStorage(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "TemplateStorage Suite")
}

var _ = Describe("TemplateStorage", func() {
    var db database.Storage
    const userID = "user1"
    const otherUserID = "user2"

    movement := goserver.Movement{
        Amount:     100,
        CurrencyId: "currency-1",
        AccountId:  "account-1",
    }

    BeforeEach(func() {
        var err error
        db, err = database.NewTestStorage()
        Expect(err).NotTo(HaveOccurred())
    })

    Describe("CreateTemplate", func() {
        It("creates a template and returns it with an ID", func() {
            tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(tpl.Id).NotTo(BeEmpty())
            Expect(tpl.Name).To(Equal("Rent"))
        })

        It("stores description, place, tags, partnerName, extra", func() {
            tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:        "Groceries",
                Description: "Weekly shopping",
                Place:       "Tesco",
                Tags:        []string{"food", "weekly"},
                PartnerName: "Tesco PLC",
                Extra:       "ref:123",
                Movements:   []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(tpl.Description).To(Equal("Weekly shopping"))
            Expect(tpl.Place).To(Equal("Tesco"))
            Expect(tpl.Tags).To(ConsistOf("food", "weekly"))
            Expect(tpl.PartnerName).To(Equal("Tesco PLC"))
            Expect(tpl.Extra).To(Equal("ref:123"))
        })
    })

    Describe("GetTemplates", func() {
        BeforeEach(func() {
            _, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())

            _, err = db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:      "Salary",
                Movements: []goserver.Movement{{Amount: 2000, CurrencyId: "currency-1", AccountId: "account-2"}},
            })
            Expect(err).NotTo(HaveOccurred())

            _, err = db.CreateTemplate(otherUserID, &goserver.TransactionTemplateNoId{
                Name:      "OtherUserTemplate",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
        })

        It("returns only templates for the requesting user", func() {
            templates, err := db.GetTemplates(userID, nil)
            Expect(err).NotTo(HaveOccurred())
            Expect(templates).To(HaveLen(2))
            names := []string{templates[0].Name, templates[1].Name}
            Expect(names).To(ConsistOf("Rent", "Salary"))
        })

        It("filters by accountId when provided", func() {
            accountID := "account-1"
            templates, err := db.GetTemplates(userID, &accountID)
            Expect(err).NotTo(HaveOccurred())
            Expect(templates).To(HaveLen(1))
            Expect(templates[0].Name).To(Equal("Rent"))
        })

        It("returns empty slice when no templates match the accountId filter", func() {
            accountID := "account-999"
            templates, err := db.GetTemplates(userID, &accountID)
            Expect(err).NotTo(HaveOccurred())
            Expect(templates).To(BeEmpty())
        })
    })

    Describe("UpdateTemplate", func() {
        var templateID string

        BeforeEach(func() {
            tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
            templateID = tpl.Id
        })

        It("updates the template and returns the updated version", func() {
            updated, err := db.UpdateTemplate(userID, templateID, &goserver.TransactionTemplateNoId{
                Name:      "Rent Updated",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(updated.Name).To(Equal("Rent Updated"))
        })

        It("returns ErrNotFound when template does not exist", func() {
            _, err := db.UpdateTemplate(userID, "00000000-0000-0000-0000-000000000000", &goserver.TransactionTemplateNoId{
                Name:      "X",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).To(MatchError(database.ErrNotFound))
        })

        It("returns ErrNotFound when template belongs to another user", func() {
            _, err := db.UpdateTemplate(otherUserID, templateID, &goserver.TransactionTemplateNoId{
                Name:      "Hack",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).To(MatchError(database.ErrNotFound))
        })
    })

    Describe("DeleteTemplate", func() {
        var templateID string

        BeforeEach(func() {
            tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{movement},
            })
            Expect(err).NotTo(HaveOccurred())
            templateID = tpl.Id
        })

        It("deletes the template", func() {
            err := db.DeleteTemplate(userID, templateID)
            Expect(err).NotTo(HaveOccurred())

            templates, err := db.GetTemplates(userID, nil)
            Expect(err).NotTo(HaveOccurred())
            Expect(templates).To(BeEmpty())
        })

        It("returns ErrNotFound for non-existent template", func() {
            err := db.DeleteTemplate(userID, "00000000-0000-0000-0000-000000000000")
            Expect(err).To(MatchError(database.ErrNotFound))
        })

        It("returns ErrNotFound when template belongs to another user", func() {
            err := db.DeleteTemplate(otherUserID, templateID)
            Expect(err).To(MatchError(database.ErrNotFound))
        })
    })
})
```

> **Note:** `database.NewTestStorage()` — check how other storage tests create test storage. Look at `storage_account_test.go` or similar for the exact constructor name and pattern used in `database_test` package.

- [ ] **Step 2: Run tests**

```bash
cd backend && go test ./pkg/database/... -v
```

Expected: all tests PASS.

- [ ] **Step 3: Run gofumpt**

```bash
go tool mvdan.cc/gofumpt -w backend/pkg/database/storage_template_test.go
```

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/database/storage_template_test.go
git commit -m "test: add full TemplateStorage test suite"
```

---

## Task 5: HTTP handlers

**Files:**
- Create: `backend/pkg/server/api/api_templates.go`
- Modify: `backend/pkg/server/server.go` (or wherever `MatchersAPIServiceImpl` is wired up)

- [ ] **Step 1: Find the generated interface**

Open `backend/pkg/generated/goserver/api_templates_service.go` (generated in Task 1). Note the exact method signatures the interface requires — typically:

```go
type TemplatesAPIServicer interface {
    GetTemplates(ctx context.Context, accountId string) (ImplResponse, error)
    CreateTemplate(ctx context.Context, transactionTemplateNoId TransactionTemplateNoId) (ImplResponse, error)
    UpdateTemplate(ctx context.Context, id string, transactionTemplateNoId TransactionTemplateNoId) (ImplResponse, error)
    DeleteTemplate(ctx context.Context, id string) (ImplResponse, error)
}
```

Adjust the handler code below to match exactly.

- [ ] **Step 2: Write the failing handler test**

Create `backend/pkg/server/api/api_templates_test.go`:

```go
package api_test

import (
    "context"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/ya-breeze/geekbudgetbe/pkg/constants"
    "github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
    "github.com/ya-breeze/geekbudgetbe/pkg/server/api"
    "go.uber.org/mock/gomock"
    "log/slog"
    "net/http"
    "os"
)

var _ = Describe("TemplatesAPI", func() {
    var (
        ctrl    *gomock.Controller
        mockDB  *mocks.MockStorage
        handler *api.TemplatesAPIServiceImpl
        ctx     context.Context
    )

    BeforeEach(func() {
        ctrl = gomock.NewController(GinkgoT())
        mockDB = mocks.NewMockStorage(ctrl)
        handler = api.NewTemplatesAPIServiceImpl(
            slog.New(slog.NewTextHandler(os.Stderr, nil)),
            mockDB,
        )
        ctx = context.WithValue(context.Background(), constants.UserIDKey, "user1")
    })

    AfterEach(func() {
        ctrl.Finish()
    })

    Describe("GetTemplates", func() {
        It("returns 200 with list of templates", func() {
            mockDB.EXPECT().GetTemplates("user1", nil).Return([]goserver.TransactionTemplate{
                {Id: "tpl-1", Name: "Rent"},
            }, nil)

            resp, err := handler.GetTemplates(ctx, "")
            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Code).To(Equal(http.StatusOK))
        })

        It("returns 500 when userID is missing from context", func() {
            resp, err := handler.GetTemplates(context.Background(), "")
            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Code).To(Equal(http.StatusInternalServerError))
        })
    })

    Describe("CreateTemplate", func() {
        It("returns 400 when movements is empty", func() {
            resp, err := handler.CreateTemplate(ctx, goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{},
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Code).To(Equal(http.StatusBadRequest))
        })

        It("returns 200 with created template", func() {
            input := goserver.TransactionTemplateNoId{
                Name:      "Rent",
                Movements: []goserver.Movement{{Amount: 1000, CurrencyId: "c1", AccountId: "a1"}},
            }
            mockDB.EXPECT().CreateTemplate("user1", &input).Return(goserver.TransactionTemplate{
                Id: "new-id", Name: "Rent",
            }, nil)

            resp, err := handler.CreateTemplate(ctx, input)
            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Code).To(Equal(http.StatusOK))
        })
    })

    Describe("DeleteTemplate", func() {
        It("returns 204 on success", func() {
            mockDB.EXPECT().DeleteTemplate("user1", "tpl-1").Return(nil)
            resp, err := handler.DeleteTemplate(ctx, "tpl-1")
            Expect(err).NotTo(HaveOccurred())
            Expect(resp.Code).To(Equal(http.StatusNoContent))
        })
    })
})
```

- [ ] **Step 3: Run test to verify it fails**

```bash
cd backend && go test ./pkg/server/api/... -run "TemplatesAPI" -v
```

Expected: compile error — `api.TemplatesAPIServiceImpl` and `api.NewTemplatesAPIServiceImpl` don't exist yet.

- [ ] **Step 4: Implement handlers**

```go
// backend/pkg/server/api/api_templates.go
package api

import (
    "context"
    "log/slog"
    "net/http"

    "github.com/ya-breeze/geekbudgetbe/pkg/constants"
    "github.com/ya-breeze/geekbudgetbe/pkg/database"
    "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type TemplatesAPIServiceImpl struct {
    logger *slog.Logger
    db     database.Storage
}

func NewTemplatesAPIServiceImpl(logger *slog.Logger, db database.Storage) *TemplatesAPIServiceImpl {
    return &TemplatesAPIServiceImpl{logger: logger, db: db}
}

func (s *TemplatesAPIServiceImpl) GetTemplates(ctx context.Context, accountId string) (goserver.ImplResponse, error) {
    userID, ok := ctx.Value(constants.UserIDKey).(string)
    if !ok {
        s.logger.Error("UserID not found in context")
        return goserver.Response(http.StatusInternalServerError, nil), nil
    }

    var accountIDPtr *string
    if accountId != "" {
        accountIDPtr = &accountId
    }

    templates, err := s.db.GetTemplates(userID, accountIDPtr)
    if err != nil {
        s.logger.With("error", err).Error("Failed to get templates")
        return goserver.Response(http.StatusInternalServerError, nil), nil
    }

    return goserver.Response(http.StatusOK, templates), nil
}

func (s *TemplatesAPIServiceImpl) CreateTemplate(ctx context.Context, t goserver.TransactionTemplateNoId) (goserver.ImplResponse, error) {
    userID, ok := ctx.Value(constants.UserIDKey).(string)
    if !ok {
        s.logger.Error("UserID not found in context")
        return goserver.Response(http.StatusInternalServerError, nil), nil
    }

    if len(t.Movements) == 0 {
        return goserver.Response(http.StatusBadRequest, "movements must not be empty"), nil
    }

    result, err := s.db.CreateTemplate(userID, &t)
    if err != nil {
        s.logger.With("error", err).Error("Failed to create template")
        return goserver.Response(http.StatusInternalServerError, nil), nil
    }

    return goserver.Response(http.StatusOK, result), nil
}

func (s *TemplatesAPIServiceImpl) UpdateTemplate(ctx context.Context, id string, t goserver.TransactionTemplateNoId) (goserver.ImplResponse, error) {
    result, _, err := updateEntity[*goserver.TransactionTemplateNoId, goserver.TransactionTemplate](ctx, s.logger, "template", id, &t, s.db.UpdateTemplate)
    if err != nil {
        return mapErrorToResponse(err), nil
    }

    return goserver.Response(http.StatusOK, result), nil
}

func (s *TemplatesAPIServiceImpl) DeleteTemplate(ctx context.Context, id string) (goserver.ImplResponse, error) {
    _, err := deleteEntity(ctx, s.logger, "template", id, s.db.DeleteTemplate)
    if err != nil {
        return mapErrorToResponse(err), nil
    }

    return goserver.Response(http.StatusNoContent, nil), nil
}
```

> **Note:** Check `api_accounts.go` or `api_matchers.go` for the exact signature of `updateEntity` and `deleteEntity` helpers. If those generic helpers don't support the template types, implement directly following the `CreateAccount` / `DeleteAccount` patterns from `api_accounts.go` instead.

- [ ] **Step 5: Wire up the handler in the server**

Find where `MatchersAPIServiceImpl` is registered (search for `NewMatchersAPIServiceImpl` in `server.go` or similar). Add alongside it:

```go
templatesService := api.NewTemplatesAPIServiceImpl(logger, db)
goserver.NewRouter(
    // ... existing services ...
    goserver.NewTemplatesAPIController(templatesService),
)
```

- [ ] **Step 6: Run handler tests**

```bash
cd backend && go test ./pkg/server/api/... -run "TemplatesAPI" -v
```

Expected: all PASS.

- [ ] **Step 7: Build and run all backend tests**

```bash
cd backend && go build ./... && go tool github.com/onsi/ginkgo/v2/ginkgo -r
```

Expected: all pass.

- [ ] **Step 8: Run gofumpt**

```bash
go tool mvdan.cc/gofumpt -w backend/pkg/server/api/api_templates.go backend/pkg/server/api/api_templates_test.go
```

- [ ] **Step 9: Commit**

```bash
git add backend/pkg/server/api/api_templates.go backend/pkg/server/api/api_templates_test.go backend/pkg/server/
git commit -m "feat: add template HTTP handlers"
```

---

## Task 6: Angular service

**Files:**
- Create: `frontend/src/app/features/templates/services/template.service.ts`

After `make generate` in Task 1, the Angular client has new generated functions in `frontend/src/app/core/api/fn/templates/`.

- [ ] **Step 1: Verify generated Angular client files exist**

```bash
ls frontend/src/app/core/api/fn/templates/
```

Expected: files like `get-templates.ts`, `create-template.ts`, `update-template.ts`, `delete-template.ts`.

- [ ] **Step 2: Create the service**

```typescript
// frontend/src/app/features/templates/services/template.service.ts
import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map, tap } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';
import { TransactionTemplateNoId } from '../../../core/api/models/transaction-template-no-id';
import { TransactionNoId } from '../../../core/api/models/transaction-no-id';
import { getTemplates } from '../../../core/api/fn/templates/get-templates';
import { createTemplate } from '../../../core/api/fn/templates/create-template';
import { updateTemplate } from '../../../core/api/fn/templates/update-template';
import { deleteTemplate } from '../../../core/api/fn/templates/delete-template';

@Injectable({
    providedIn: 'root',
})
export class TemplateService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly templates = signal<TransactionTemplate[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadTemplates(accountId?: string): Observable<TransactionTemplate[]> {
        this.loading.set(true);
        this.error.set(null);

        return getTemplates(this.http, this.apiConfig.rootUrl, accountId ? { params: { accountId } } : {}).pipe(
            map((response) => response.body),
            tap({
                next: (templates) => {
                    this.templates.set(templates);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load templates');
                    this.loading.set(false);
                },
            }),
        );
    }

    create(template: TransactionTemplateNoId): Observable<TransactionTemplate> {
        this.loading.set(true);
        this.error.set(null);

        return createTemplate(this.http, this.apiConfig.rootUrl, { body: template }).pipe(
            map((response) => response.body),
            tap({
                next: (t) => {
                    this.templates.update((templates) => [...templates, t]);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to create template');
                    this.loading.set(false);
                },
            }),
        );
    }

    update(id: string, template: TransactionTemplateNoId): Observable<TransactionTemplate> {
        this.loading.set(true);
        this.error.set(null);

        return updateTemplate(this.http, this.apiConfig.rootUrl, { id, body: template }).pipe(
            map((response) => response.body),
            tap({
                next: (updated) => {
                    this.templates.update((templates) =>
                        templates.map((t) => (t.id === id ? updated : t)),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to update template');
                    this.loading.set(false);
                },
            }),
        );
    }

    delete(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return deleteTemplate(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.templates.update((templates) => templates.filter((t) => t.id !== id));
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to delete template');
                    this.loading.set(false);
                },
            }),
        );
    }

    /**
     * Converts a template to a TransactionNoId payload suitable for POSTing to /v1/transactions.
     * Sets date to today. Zeroes out all import-only fields.
     */
    templateToTransactionNoId(template: TransactionTemplate): TransactionNoId {
        return {
            date: new Date().toISOString(),
            description: template.description ?? '',
            place: template.place ?? '',
            tags: template.tags ?? [],
            partnerName: template.partnerName ?? '',
            extra: template.extra ?? '',
            movements: template.movements ?? [],
            // Explicitly zero out import-only fields
            externalIds: [],
            isAuto: false,
            matcherId: undefined,
            suspiciousReasons: [],
            duplicateDismissed: false,
            mergedIntoId: undefined,
        };
    }
}
```

> **Note:** Check the exact import paths for the generated functions — they may differ from `fn/templates/get-templates`. Look at what `make generate` actually created in `frontend/src/app/core/api/fn/templates/`. Also verify the `getTemplates` function signature for the `accountId` query parameter.

- [ ] **Step 3: Build the frontend to check for compile errors**

```bash
cd frontend && npm run build -- --configuration development 2>&1 | head -50
```

Expected: no TypeScript errors related to the template service. Other pre-existing errors are acceptable.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/app/features/templates/services/template.service.ts
git commit -m "feat: add TemplateService with CRUD and templateToTransactionNoId converter"
```

---

## Task 7: Template picker component

This shared component is used in both the transaction list ("New from template") and the create form ("Use template").

**Files:**
- Create: `frontend/src/app/features/templates/template-picker/template-picker.component.ts`
- Create: `frontend/src/app/features/templates/template-picker/template-picker.component.html`

- [ ] **Step 1: Create the picker component TS**

```typescript
// frontend/src/app/features/templates/template-picker/template-picker.component.ts
import { Component, inject, input, output, OnInit, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatListModule } from '@angular/material/list';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TemplateService } from '../services/template.service';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';

@Component({
    selector: 'app-template-picker',
    standalone: true,
    imports: [
        FormsModule,
        MatListModule,
        MatInputModule,
        MatFormFieldModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './template-picker.component.html',
})
export class TemplatePickerComponent implements OnInit {
    private readonly templateService = inject(TemplateService);

    /** Optional: pre-filter templates to those containing this accountId in their movements */
    accountId = input<string | undefined>(undefined);

    /** Emitted when the user selects a template */
    templateSelected = output<TransactionTemplate>();

    protected readonly loading = this.templateService.loading;
    protected readonly searchQuery = signal('');

    protected readonly filteredTemplates = computed(() => {
        const q = this.searchQuery().toLowerCase();
        return this.templateService.templates().filter((t) =>
            t.name.toLowerCase().includes(q),
        );
    });

    ngOnInit(): void {
        this.templateService.loadTemplates(this.accountId()).subscribe();
    }

    protected select(template: TransactionTemplate): void {
        this.templateSelected.emit(template);
    }
}
```

- [ ] **Step 2: Create the picker template**

```html
<!-- frontend/src/app/features/templates/template-picker/template-picker.component.html -->
<div class="template-picker">
  <mat-form-field appearance="outline" class="search-field">
    <mat-label>Search templates</mat-label>
    <input matInput [(ngModel)]="searchQuery" placeholder="Type to filter..." />
    <mat-icon matSuffix>search</mat-icon>
  </mat-form-field>

  @if (loading()) {
    <mat-spinner diameter="32" />
  } @else if (filteredTemplates().length === 0) {
    <p class="no-templates">No templates found.</p>
  } @else {
    <mat-selection-list [multiple]="false">
      @for (template of filteredTemplates(); track template.id) {
        <mat-list-option (click)="select(template)">
          <span matListItemTitle>{{ template.name }}</span>
          @if (template.description) {
            <span matListItemLine>{{ template.description }}</span>
          }
        </mat-list-option>
      }
    </mat-selection-list>
  }
</div>
```

- [ ] **Step 3: Build check**

```bash
cd frontend && npm run build -- --configuration development 2>&1 | head -50
```

Expected: no new TypeScript errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/app/features/templates/template-picker/
git commit -m "feat: add shared TemplatePickerComponent"
```

---

## Task 8: Template management page

**Files:**
- Create: `frontend/src/app/features/templates/templates.component.ts`
- Create: `frontend/src/app/features/templates/templates.component.html`
- Create: `frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.ts`
- Create: `frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.html`
- Modify: `frontend/src/app/app.routes.ts`
- Modify: `frontend/src/app/layout/sidebar/sidebar.component.ts`

- [ ] **Step 1: Create the edit dialog TS**

```typescript
// frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.ts
import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatSelectModule } from '@angular/material/select';
import { TemplateService } from '../services/template.service';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';
import { TransactionTemplateNoId } from '../../../core/api/models/transaction-template-no-id';

export interface TemplateInitialValues {
    name?: string;
    description?: string;
    place?: string;
    partnerName?: string;
    extra?: string;
    movements?: goserver.Movement[];
}

export interface TemplateEditDialogData {
    /** Provide when editing an existing template. Sets isEditMode = true. */
    template?: TransactionTemplate;
    /** Provide when creating a new template pre-populated from a transaction. isEditMode stays false. */
    initialValues?: TemplateInitialValues;
}

@Component({
    selector: 'app-template-edit-dialog',
    standalone: true,
    imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatIconModule,
        MatChipsModule,
        MatSelectModule,
    ],
    templateUrl: './template-edit-dialog.component.html',
})
export class TemplateEditDialogComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly dialogRef = inject(MatDialogRef<TemplateEditDialogComponent>);
    protected readonly data = inject<TemplateEditDialogData>(MAT_DIALOG_DATA);
    private readonly templateService = inject(TemplateService);
    protected readonly accountService = inject(AccountService);
    protected readonly currencyService = inject(CurrencyService);

    protected readonly isEditMode = !!this.data?.template;

    protected readonly form = this.fb.group({
        name: ['', Validators.required],
        description: [''],
        place: [''],
        partnerName: [''],
        extra: [''],
        movements: this.fb.array([]),
    });

    get movements(): FormArray {
        return this.form.get('movements') as FormArray;
    }

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();

        // Editing an existing template
        if (this.data?.template) {
            const t = this.data.template;
            this.form.patchValue({
                name: t.name,
                description: t.description ?? '',
                place: t.place ?? '',
                partnerName: t.partnerName ?? '',
                extra: t.extra ?? '',
            });
            (t.movements ?? []).forEach((m) => this.addMovement(m));
        // Creating a new template pre-populated from a transaction (initialValues has no id → isEditMode stays false)
        } else if (this.data?.initialValues) {
            const v = this.data.initialValues;
            this.form.patchValue({
                name: v.name ?? '',
                description: v.description ?? '',
                place: v.place ?? '',
                partnerName: v.partnerName ?? '',
                extra: v.extra ?? '',
            });
            (v.movements ?? []).forEach((m) => this.addMovement(m));
            if (!v.movements?.length) this.addMovement();
        } else {
            this.addMovement();
        }
    }

    protected addMovement(movement?: { amount?: number; currencyId?: string; accountId?: string }): void {
        this.movements.push(this.fb.group({
            amount: [movement?.amount ?? null, Validators.required],
            currencyId: [movement?.currencyId ?? '', Validators.required],
            accountId: [movement?.accountId ?? ''],
        }));
    }

    protected removeMovement(index: number): void {
        if (this.movements.length > 1) {
            this.movements.removeAt(index);
        }
    }

    protected save(): void {
        if (this.form.invalid) return;

        const value = this.form.getRawValue();
        const payload: TransactionTemplateNoId = {
            name: value.name!,
            description: value.description ?? undefined,
            place: value.place ?? undefined,
            partnerName: value.partnerName ?? undefined,
            extra: value.extra ?? undefined,
            movements: value.movements.map((m: any) => ({
                amount: Number(m.amount),
                currencyId: m.currencyId,
                accountId: m.accountId || undefined,
            })),
        };

        const obs = this.isEditMode
            ? this.templateService.update(this.data.template!.id, payload)
            : this.templateService.create(payload);

        obs.subscribe(() => this.dialogRef.close(true));
    }

    protected cancel(): void {
        this.dialogRef.close(false);
    }
}
```

- [ ] **Step 2: Create the edit dialog HTML**

```html
<!-- frontend/src/app/features/templates/template-edit-dialog/template-edit-dialog.component.html -->
<h2 mat-dialog-title>{{ isEditMode ? 'Edit Template' : 'New Template' }}</h2>

<mat-dialog-content [formGroup]="form">
  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Template name *</mat-label>
    <input matInput formControlName="name" placeholder="e.g. Monthly rent" />
  </mat-form-field>

  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Description</mat-label>
    <input matInput formControlName="description" />
  </mat-form-field>

  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Place</mat-label>
    <input matInput formControlName="place" />
  </mat-form-field>

  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Partner name</mat-label>
    <input matInput formControlName="partnerName" />
  </mat-form-field>

  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Extra / reference</mat-label>
    <input matInput formControlName="extra" />
  </mat-form-field>

  <h3>Movements</h3>
  <div formArrayName="movements">
    @for (movement of movements.controls; track $index; let i = $index) {
      <div [formGroupName]="i" class="movement-row">
        <mat-form-field appearance="outline">
          <mat-label>Amount</mat-label>
          <input matInput type="number" formControlName="amount" />
        </mat-form-field>

        <mat-form-field appearance="outline">
          <mat-label>Currency</mat-label>
          <mat-select formControlName="currencyId">
            @for (c of currencyService.currencies(); track c.id) {
              <mat-option [value]="c.id">{{ c.name }}</mat-option>
            }
          </mat-select>
        </mat-form-field>

        <mat-form-field appearance="outline">
          <mat-label>Account</mat-label>
          <mat-select formControlName="accountId">
            <mat-option value="">— none —</mat-option>
            @for (a of accountService.accounts(); track a.id) {
              <mat-option [value]="a.id">{{ a.name }}</mat-option>
            }
          </mat-select>
        </mat-form-field>

        <button mat-icon-button color="warn" (click)="removeMovement(i)" [disabled]="movements.length === 1">
          <mat-icon>delete</mat-icon>
        </button>
      </div>
    }
  </div>

  <button mat-stroked-button (click)="addMovement()">
    <mat-icon>add</mat-icon> Add movement
  </button>
</mat-dialog-content>

<mat-dialog-actions align="end">
  <button mat-button (click)="cancel()">Cancel</button>
  <button mat-raised-button color="primary" (click)="save()" [disabled]="form.invalid">
    {{ isEditMode ? 'Update' : 'Create' }}
  </button>
</mat-dialog-actions>
```

- [ ] **Step 3: Create the list page TS**

```typescript
// frontend/src/app/features/templates/templates.component.ts
import { Component, inject, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TemplateService } from './services/template.service';
import { AccountService } from '../accounts/services/account.service';
import { TemplateEditDialogComponent } from './template-edit-dialog/template-edit-dialog.component';
import { TransactionTemplate } from '../../core/api/models/transaction-template';

@Component({
    selector: 'app-templates',
    standalone: true,
    imports: [
        MatButtonModule,
        MatIconModule,
        MatTableModule,
        MatDialogModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './templates.component.html',
})
export class TemplatesComponent implements OnInit {
    private readonly templateService = inject(TemplateService);
    private readonly accountService = inject(AccountService);
    private readonly dialog = inject(MatDialog);

    protected readonly templates = this.templateService.templates;
    protected readonly loading = this.templateService.loading;
    protected readonly displayedColumns = ['name', 'description', 'accounts', 'actions'];

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.templateService.loadTemplates().subscribe();
    }

    protected getAccountNames(template: TransactionTemplate): string {
        const accounts = this.accountService.accounts();
        const accountMap = new Map(accounts.map((a) => [a.id, a.name]));
        const ids = [...new Set((template.movements ?? []).map((m) => m.accountId).filter(Boolean))];
        return ids.map((id) => accountMap.get(id!) ?? id!).join(', ');
    }

    protected openCreateDialog(): void {
        this.dialog.open(TemplateEditDialogComponent, {
            width: '600px',
            data: {},
            disableClose: true,
        });
    }

    protected openEditDialog(template: TransactionTemplate): void {
        this.dialog.open(TemplateEditDialogComponent, {
            width: '600px',
            data: { template },
            disableClose: true,
        });
    }

    protected delete(template: TransactionTemplate): void {
        if (confirm(`Delete template "${template.name}"?`)) {
            this.templateService.delete(template.id).subscribe();
        }
    }
}
```

- [ ] **Step 4: Create the list page HTML**

```html
<!-- frontend/src/app/features/templates/templates.component.html -->
<div class="page-container">
  <div class="page-header">
    <h1>Templates</h1>
    <button mat-raised-button color="primary" (click)="openCreateDialog()">
      <mat-icon>add</mat-icon> New template
    </button>
  </div>

  @if (loading()) {
    <mat-spinner />
  } @else if (templates().length === 0) {
    <p class="empty-state">No templates yet. Create one to get started.</p>
  } @else {
    <table mat-table [dataSource]="templates()" class="mat-elevation-z2">
      <ng-container matColumnDef="name">
        <th mat-header-cell *matHeaderCellDef>Name</th>
        <td mat-cell *matCellDef="let t">{{ t.name }}</td>
      </ng-container>

      <ng-container matColumnDef="description">
        <th mat-header-cell *matHeaderCellDef>Description</th>
        <td mat-cell *matCellDef="let t">{{ t.description }}</td>
      </ng-container>

      <ng-container matColumnDef="accounts">
        <th mat-header-cell *matHeaderCellDef>Accounts</th>
        <td mat-cell *matCellDef="let t">{{ getAccountNames(t) }}</td>
      </ng-container>

      <ng-container matColumnDef="actions">
        <th mat-header-cell *matHeaderCellDef></th>
        <td mat-cell *matCellDef="let t">
          <button mat-icon-button (click)="openEditDialog(t)">
            <mat-icon>edit</mat-icon>
          </button>
          <button mat-icon-button color="warn" (click)="delete(t)">
            <mat-icon>delete</mat-icon>
          </button>
        </td>
      </ng-container>

      <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
      <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
    </table>
  }
</div>
```

- [ ] **Step 5: Add route to `app.routes.ts`**

Open `frontend/src/app/app.routes.ts`. Inside the `children` array of the `LayoutComponent` route, add:

```typescript
{
    path: 'templates',
    loadComponent: () =>
        import('./features/templates/templates.component').then(
            (m) => m.TemplatesComponent,
        ),
},
```

- [ ] **Step 6: Add sidebar menu item**

Open `frontend/src/app/layout/sidebar/sidebar.component.ts`. In the `menuItems` signal array, add after `Matchers`:

```typescript
{ label: 'Templates', icon: 'content_copy', route: '/templates' },
```

- [ ] **Step 7: Build check**

```bash
cd frontend && npm run build -- --configuration development 2>&1 | head -80
```

Expected: no new TypeScript errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/src/app/features/templates/ frontend/src/app/app.routes.ts frontend/src/app/layout/sidebar/sidebar.component.ts
git commit -m "feat: add templates management page with create/edit/delete"
```

---

## Task 9: "Save as template" from transaction detail

**Files:**
- Modify: `frontend/src/app/features/transactions/transaction-detail/transaction-detail.component.ts`

- [ ] **Step 1: Open and read the transaction detail component**

Read `frontend/src/app/features/transactions/transaction-detail/transaction-detail.component.ts` to understand its current structure before modifying it.

- [ ] **Step 2: Add "Save as template" button logic**

Add a `saveAsTemplate()` method that opens `TemplateEditDialogComponent` pre-populated from the current transaction:

```typescript
// Add to imports
import { TemplateEditDialogComponent } from '../../../features/templates/template-edit-dialog/template-edit-dialog.component';
import { MatDialog } from '@angular/material/dialog';

// Add to class
private readonly dialog = inject(MatDialog);

protected saveAsTemplate(): void {
    const tx = this.transaction(); // adjust to however the component exposes the current transaction
    // Use `initialValues` (NOT `template`) so isEditMode stays false — dialog will call create(), not update()
    this.dialog.open(TemplateEditDialogComponent, {
        width: '600px',
        data: {
            initialValues: {
                name: tx.description || '',
                description: tx.description,
                place: tx.place,
                partnerName: tx.partnerName,
                extra: tx.extra,
                movements: tx.movements,
            },
        },
        disableClose: true,
    });
}
```

> **Note:** Adjust signal/property access to match how the component currently holds the transaction. Read the component first.

- [ ] **Step 3: Add button to template**

In `transaction-detail.component.html`, add a "Save as template" button in the actions section:

```html
<button mat-stroked-button (click)="saveAsTemplate()">
  <mat-icon>content_copy</mat-icon> Save as template
</button>
```

- [ ] **Step 4: Build check**

```bash
cd frontend && npm run build -- --configuration development 2>&1 | head -50
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/features/transactions/transaction-detail/
git commit -m "feat: add Save as template action to transaction detail"
```

---

## Task 10: "New from template" and "Use template" in transaction form

**Files:**
- Modify: `frontend/src/app/features/transactions/transactions.component.ts` (transaction list)
- Modify: `frontend/src/app/features/transactions/transaction-form-dialog/transaction-form-dialog.component.ts`

- [ ] **Step 1: Read both files before modifying**

Read both component files to understand their current structure.

- [ ] **Step 2: Add "New from template" to transaction list**

In `transactions.component.ts`:
- Import `MatMenuModule`, `TemplatePickerComponent`, `TemplateService`, `MatDialog`
- Add a method `newFromTemplate()` that:
  1. Opens a dialog/bottom sheet containing `TemplatePickerComponent`
  2. On `templateSelected`, calls `templateService.templateToTransactionNoId(template)` and opens the transaction form dialog pre-populated

The simplest approach is to open a small `MatDialog` containing just the picker:

```typescript
protected newFromTemplate(): void {
    const dialogRef = this.dialog.open(TemplatePickerDialogComponent, {
        width: '400px',
        data: { accountId: this.selectedAccountId() }, // adjust to current filter
    });
    dialogRef.afterClosed().subscribe((template) => {
        if (template) {
            const txNoId = this.templateService.templateToTransactionNoId(template);
            this.openCreateDialog(txNoId);
        }
    });
}
```

Create a thin wrapper `TemplatePickerDialogComponent` that hosts `TemplatePickerComponent` and closes the dialog with the selected template:

```typescript
// frontend/src/app/features/templates/template-picker/template-picker-dialog.component.ts
import { Component, inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TemplatePickerComponent } from './template-picker.component';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';

@Component({
    selector: 'app-template-picker-dialog',
    standalone: true,
    imports: [MatDialogModule, TemplatePickerComponent],
    template: `
        <h2 mat-dialog-title>Choose a template</h2>
        <mat-dialog-content>
            <app-template-picker
                [accountId]="data?.accountId"
                (templateSelected)="onSelect($event)"
            />
        </mat-dialog-content>
        <mat-dialog-actions align="end">
            <button mat-button mat-dialog-close>Cancel</button>
        </mat-dialog-actions>
    `,
})
export class TemplatePickerDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<TemplatePickerDialogComponent>);
    protected readonly data = inject<{ accountId?: string }>(MAT_DIALOG_DATA);

    protected onSelect(template: TransactionTemplate): void {
        this.dialogRef.close(template);
    }
}
```

Add a "New from template" button in `transactions.component.html` next to the existing "New transaction" button:

```html
<button mat-stroked-button (click)="newFromTemplate()">
  <mat-icon>content_copy</mat-icon> New from template
</button>
```

- [ ] **Step 3: Add "Use template" inside the create/edit form**

In `transaction-form-dialog.component.ts`, add:
- Import `MatDialog`, `TemplatePickerDialogComponent`, `TemplateService`
- Add a `useTemplate()` method:

```typescript
protected useTemplate(): void {
    const dialogRef = this.dialog.open(TemplatePickerDialogComponent, { width: '400px' });
    dialogRef.afterClosed().subscribe((template) => {
        if (template) {
            const tx = this.templateService.templateToTransactionNoId(template);
            // Patch the form fields
            this.form.patchValue({
                description: tx.description,
                place: tx.place,
                partnerName: tx.partnerName,
                extra: tx.extra,
            });
            // Reset and repopulate movements FormArray
            this.movements.clear();
            (tx.movements ?? []).forEach((m) => this.addMovement(m));
        }
    });
}
```

Add a "Use template" button at the top of `transaction-form-dialog.component.html` (only shown in create mode):

```html
@if (!isEditMode) {
  <button mat-stroked-button type="button" (click)="useTemplate()">
    <mat-icon>content_copy</mat-icon> Use template
  </button>
}
```

> **Note:** Adjust to match the form's actual structure — `movements` FormArray name, `addMovement` method signature, etc.

- [ ] **Step 4: Build check**

```bash
cd frontend && npm run build -- --configuration development 2>&1 | head -80
```

Expected: no new errors.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/features/transactions/ frontend/src/app/features/templates/template-picker/template-picker-dialog.component.ts
git commit -m "feat: add New from template and Use template entry points"
```

---

## Task 11: Full pipeline validation

- [ ] **Step 1: Run full backend test suite**

```bash
# Run from the repo root (where the Makefile lives)
make all
```

Expected: build, tests, OpenAPI validation, and lint all pass.

- [ ] **Step 2: Start the backend and frontend**

```bash
# Terminal 1
make run-backend

# Terminal 2
make run-frontend
```

- [ ] **Step 3: Smoke test the feature**

Using `test@test.com` / `test` credentials:

1. Navigate to `/templates` — should see empty list with "New template" button
2. Create a template with name "Rent", add a movement with amount and account
3. Template appears in the list
4. Go to `/transactions` — click "New from template", picker opens, select "Rent", create form pre-fills
5. Go to an existing transaction detail — "Save as template" button appears, click it, name dialog opens
6. In the create transaction form, click "Use template" at the top — picker opens

- [ ] **Step 4: Run lint**

```bash
make lint
```

Expected: no lint errors.

- [ ] **Step 5: Final commit if any lint fixes were needed**

```bash
git add -p
git commit -m "fix: address lint issues in template feature"
```
