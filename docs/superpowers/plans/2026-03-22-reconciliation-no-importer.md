# Reconciliation: No-Importer Account Support

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Allow users to reconcile accounts that have no bank importer, removing the permanent "delta too large" block and fixing a latent backend bug that caused HTTP 400 for any non-zero balance.

**Architecture:** Backend detects no-importer accounts via `GetBankImporters` lookup (same mechanism as status endpoint), skips tolerance check, and sets `expectedBalance = balance`. Frontend removes the delta gate for no-importer accounts, adds a confirmation dialog for large deltas, and updates row coloring.

**Tech Stack:** Go (Ginkgo/Gomega/gomock), Angular 20 (standalone components, Angular Material, RxJS)

**Spec:** `docs/superpowers/specs/2026-03-22-reconciliation-no-importer-design.md`

---

## File Map

| File | Change |
|------|--------|
| `backend/pkg/server/api/api_reconciliation.go` | Restructure `ReconcileAccount` handler — importer detection + branch |
| `backend/pkg/server/api/api_reconciliation_test.go` | Add `Describe("ReconcileAccount")` with 5 new test cases |
| `frontend/src/app/features/reconciliation/reconciliation.component.ts` | Inject `MatDialog`, update `reconcile()`, `getReconcileTooltip()`, `getStatusClass()` |
| `frontend/src/app/features/reconciliation/reconciliation.component.html` | Update `[disabled]` binding and delta cell coloring |
| `frontend/src/app/features/reconciliation/reconciliation.component.spec.ts` | Add tests (create if not exists) |

---

## Task 1: Backend — Write Failing Tests for `ReconcileAccount`

**Files:**
- Modify: `backend/pkg/server/api/api_reconciliation_test.go`

The existing test file has a `Describe("GetReconciliationStatus")` block. Add a new `Describe("ReconcileAccount")` block inside the outer `Describe("Reconciliation API")`.

Key mock setup patterns from the existing file:
- `mockStorage.EXPECT().GetBankImporters("user1").Return(...)` — returns `[]goserver.BankImporter`
- `mockStorage.EXPECT().GetAccountBalance("user1", id, currencyId).Return(decimal.Decimal, nil)`
- `mockStorage.EXPECT().GetAccount("user1", id).Return(goserver.Account, nil)`
- `mockStorage.EXPECT().CreateReconciliation("user1", matcher).Return(&goserver.Reconciliation{...}, nil)`

- [ ] **Step 1: Add the `ReconcileAccount` describe block with 5 failing tests**

Add after the closing `})` of `Describe("GetReconciliationStatus")`, still inside the outer `Describe("Reconciliation API", func() {`:

```go
Describe("ReconcileAccount", func() {
    var noImporterAccount goserver.Account
    var importerAccount goserver.Account
    var importerEntry goserver.BankImporter

    BeforeEach(func() {
        noImporterAccount = goserver.Account{
            Id:   "acc_noimporter",
            Name: "Cash",
            Type: "asset",
            BankInfo: goserver.BankAccountInfo{
                Balances: []goserver.BankAccountInfoBalancesInner{
                    {CurrencyId: "USD", ClosingBalance: decimal.NewFromInt(0)},
                },
            },
        }
        importerAccount = goserver.Account{
            Id:   "acc_importer",
            Name: "Bank",
            Type: "asset",
            BankInfo: goserver.BankAccountInfo{
                Balances: []goserver.BankAccountInfoBalancesInner{
                    {CurrencyId: "USD", ClosingBalance: decimal.NewFromInt(1000)},
                },
            },
        }
        importerEntry = goserver.BankImporter{AccountId: "acc_importer"}
    })

    It("reconciles no-importer account with large delta — returns 200 with IsManual=true and delta=0 in history", func() {
        mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{}, nil)
        mockStorage.EXPECT().GetAccountBalance("user1", "acc_noimporter", "USD").
            Return(decimal.NewFromInt(500), nil)

        mockStorage.EXPECT().CreateReconciliation("user1", gomock.Any()).
            DoAndReturn(func(_ string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
                Expect(rec.ReconciledBalance).To(Equal(decimal.NewFromInt(500)))
                Expect(rec.ExpectedBalance).To(Equal(decimal.NewFromInt(500))) // delta=0 in history
                Expect(rec.IsManual).To(BeTrue())
                return goserver.Reconciliation{
                    AccountId:         rec.AccountId,
                    CurrencyId:        rec.CurrencyId,
                    ReconciledBalance: rec.ReconciledBalance,
                    ExpectedBalance:   rec.ExpectedBalance,
                    IsManual:          rec.IsManual,
                }, nil
            })

        resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
            CurrencyId: "USD",
            Balance:    decimal.NewFromInt(0), // frontend always sends 0 — "use app balance"
        })
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.Code).To(Equal(http.StatusOK))
    })

    It("reconciles no-importer account with balance=0 in request — IsManual=true even if fetched balance is zero", func() {
        mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{}, nil)
        mockStorage.EXPECT().GetAccountBalance("user1", "acc_noimporter", "USD").
            Return(decimal.NewFromInt(0), nil) // zero balance

        mockStorage.EXPECT().CreateReconciliation("user1", gomock.Any()).
            DoAndReturn(func(_ string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
                Expect(rec.IsManual).To(BeTrue()) // must be true even when balance is zero
                Expect(rec.ReconciledBalance.IsZero()).To(BeTrue())
                Expect(rec.ExpectedBalance.IsZero()).To(BeTrue())
                return goserver.Reconciliation{}, nil
            })

        resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
            CurrencyId: "USD",
            Balance:    decimal.NewFromInt(0),
        })
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.Code).To(Equal(http.StatusOK))
    })

    It("reconciles no-importer account within tolerance — returns 200 (happy path unchanged)", func() {
        mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{}, nil)

        mockStorage.EXPECT().CreateReconciliation("user1", gomock.Any()).
            DoAndReturn(func(_ string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
                Expect(rec.IsManual).To(BeTrue())
                Expect(rec.ReconciledBalance).To(Equal(decimal.NewFromFloat(500.005)))
                Expect(rec.ExpectedBalance).To(Equal(decimal.NewFromFloat(500.005))) // no-importer: expectedBalance == balance
                return goserver.Reconciliation{}, nil
            })

        resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
            CurrencyId: "USD",
            Balance:    decimal.NewFromFloat(500.005), // within tolerance (0.01)
        })
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.Code).To(Equal(http.StatusOK))
    })

    It("blocks importer account with large delta — returns 400 (existing behavior preserved)", func() {
        mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{importerEntry}, nil)
        mockStorage.EXPECT().GetAccount("user1", "acc_importer").Return(importerAccount, nil)
        // Bank balance is 1000, request balance is 500 → delta = 500 → 400
        resp, err := sut.ReconcileAccount(ctx, "acc_importer", goserver.ReconcileAccountRequest{
            CurrencyId: "USD",
            Balance:    decimal.NewFromInt(500),
        })
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.Code).To(Equal(http.StatusBadRequest))
    })

    It("blocks account with importer record but no BankInfo.Balances (never run) — treated as has-importer, returns 400", func() {
        accountNoBankInfo := goserver.Account{
            Id:   "acc_importer",
            Name: "Bank",
            Type: "asset",
            BankInfo: goserver.BankAccountInfo{
                Balances: []goserver.BankAccountInfoBalancesInner{}, // no balance data from importer
            },
        }
        mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{importerEntry}, nil)
        mockStorage.EXPECT().GetAccount("user1", "acc_importer").Return(accountNoBankInfo, nil)
        // expectedBalance = 0 (no balance entries), balance = 500 → delta = 500 → 400
        resp, err := sut.ReconcileAccount(ctx, "acc_importer", goserver.ReconcileAccountRequest{
            CurrencyId: "USD",
            Balance:    decimal.NewFromInt(500),
        })
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.Code).To(Equal(http.StatusBadRequest))
    })
})
```

- [ ] **Step 2: Run the tests to verify they all fail**

```bash
cd /Users/ek/work/GeekBudgetBE/backend && go tool github.com/onsi/ginkgo/v2/ginkgo -r ./pkg/server/api/
```

Expected: 5 new test failures (existing tests should still pass). You should see failures because `ReconcileAccount` does not yet call `GetBankImporters`.

- [ ] **Step 3: Commit the failing tests**

```bash
cd /Users/ek/work/GeekBudgetBE
git add backend/pkg/server/api/api_reconciliation_test.go
git commit -m "test: add ReconcileAccount tests for no-importer accounts"
```

---

## Task 2: Backend — Implement `ReconcileAccount` Fix

**Files:**
- Modify: `backend/pkg/server/api/api_reconciliation.go` (lines 127–180)

- [ ] **Step 1: Restructure `ReconcileAccount`**

Replace the current `ReconcileAccount` function body (lines 128–180) with:

```go
// ReconcileAccount creates a new reconciliation record
func (s *ReconciliationAPIServiceImpl) ReconcileAccount(
	ctx context.Context, id string, body goserver.ReconcileAccountRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Detect whether this account has a bank importer configured.
	// Use GetBankImporters (same as GetReconciliationStatus) rather than inspecting
	// BankInfo.Balances — an importer that has never run would have no Balances entries
	// but is still a configured importer and should enforce the tolerance check.
	importers, err := s.db.GetBankImporters(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}
	hasImporter := false
	for _, imp := range importers {
		if imp.AccountId == id {
			hasImporter = true
			break
		}
	}

	// Resolve balance: if frontend sends 0, use the current computed account balance.
	balance := body.Balance
	if balance.IsZero() {
		balance, err = s.db.GetAccountBalance(userID, id, body.CurrencyId)
		if err != nil {
			s.logger.With("error", err).Error("Failed to get account balance")
			return goserver.Response(500, nil), nil
		}
	}

	var expectedBalance decimal.Decimal
	isManual := true

	if hasImporter {
		// Importer path: derive expected balance from last import data and enforce tolerance.
		acc, accErr := s.db.GetAccount(userID, id)
		if accErr != nil {
			return goserver.Response(404, nil), nil
		}
		for _, b := range acc.BankInfo.Balances {
			if b.CurrencyId == body.CurrencyId {
				expectedBalance = b.ClosingBalance
				break
			}
		}
		if balance.Sub(expectedBalance).Abs().GreaterThan(constants.ReconciliationTolerance) {
			return goserver.Response(400, "Cannot reconcile: account balance does not match bank balance"), nil
		}
		isManual = body.Balance.IsPositive()
	} else {
		// No-importer path: the user is confirming the app balance is correct.
		// Set expectedBalance = balance so the history record shows delta = 0.
		// IsManual is always true for no-importer accounts.
		expectedBalance = balance
	}

	rec, err := s.db.CreateReconciliation(userID, &goserver.ReconciliationNoId{
		AccountId:         id,
		CurrencyId:        body.CurrencyId,
		ReconciledBalance: balance,
		ExpectedBalance:   expectedBalance,
		IsManual:          isManual,
	})
	if err != nil {
		s.logger.With("error", err).Error("Failed to create reconciliation")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, rec), nil
}
```

- [ ] **Step 2: Format the file**

```bash
cd /Users/ek/work/GeekBudgetBE/backend && go tool mvdan.cc/gofumpt -w pkg/server/api/api_reconciliation.go
```

- [ ] **Step 3: Run the tests**

```bash
cd /Users/ek/work/GeekBudgetBE/backend && go tool github.com/onsi/ginkgo/v2/ginkgo -r ./pkg/server/api/
```

Expected: all 5 new tests pass, all existing tests pass.

- [ ] **Step 4: Run the full backend build + lint**

```bash
cd /Users/ek/work/GeekBudgetBE && make build && make lint
```

Expected: clean build, no lint errors.

- [ ] **Step 5: Commit**

```bash
cd /Users/ek/work/GeekBudgetBE
git add backend/pkg/server/api/api_reconciliation.go
git commit -m "fix: allow reconciliation for accounts without bank importers"
```

---

## Task 3: Frontend — Update Template (Disable Condition + Delta Color)

**Files:**
- Modify: `frontend/src/app/features/reconciliation/reconciliation.component.html`

- [ ] **Step 1: Update the `[disabled]` binding (lines 125–129)**

Replace:
```html
                        [disabled]="
                            (element.delta || 0) > 0.01 ||
                            (element.delta || 0) < -0.01 ||
                            element.hasUnprocessedTransactions
                        "
```

With:
```html
                        [disabled]="
                            element.hasUnprocessedTransactions ||
                            (element.hasBankImporter &&
                                ((element.delta || 0) > 0.01 || (element.delta || 0) < -0.01))
                        "
```

- [ ] **Step 2: Update the delta cell coloring (lines 69–71)**

The delta cell currently applies `text-danger` unconditionally for any large delta. For no-importer accounts a large delta is expected — guard with `hasBankImporter`:

Replace:
```html
                <span
                    [class.text-danger]="
                        (element.delta || 0) > 0.01 || (element.delta || 0) < -0.01
                    "
                    [class.text-success]="
                        (element.delta || 0) <= 0.01 && (element.delta || 0) >= -0.01
                    "
                >
```

With:
```html
                <span
                    [class.text-danger]="
                        element.hasBankImporter &&
                        ((element.delta || 0) > 0.01 || (element.delta || 0) < -0.01)
                    "
                    [class.text-success]="
                        (element.delta || 0) <= 0.01 && (element.delta || 0) >= -0.01
                    "
                >
```

- [ ] **Step 3: Build the frontend to catch any template errors**

```bash
cd /Users/ek/work/GeekBudgetBE/frontend && npm run build -- --no-progress 2>&1 | tail -20
```

Expected: build succeeds with no template errors.

- [ ] **Step 4: Commit**

```bash
cd /Users/ek/work/GeekBudgetBE
git add frontend/src/app/features/reconciliation/reconciliation.component.html
git commit -m "fix: remove delta gate and red coloring for no-importer accounts"
```

---

## Task 4: Frontend — Update TypeScript Component (Dialog + Tooltip + Row Color)

**Files:**
- Modify: `frontend/src/app/features/reconciliation/reconciliation.component.ts`

- [ ] **Step 1: Add `MatDialog` and `ConfirmationDialogComponent` imports**

In the import section at the top of the file, add:

```typescript
import { MatDialog } from '@angular/material/dialog';
import { ConfirmationDialogComponent, ConfirmationDialogData } from '../../shared/components/confirmation-dialog/confirmation-dialog.component';
```

- [ ] **Step 2: Inject `MatDialog` in the component**

In the `ReconciliationComponent` class body (after `private router = inject(Router);`), add:

```typescript
    private dialog = inject(MatDialog);
```

- [ ] **Step 3: Replace `getReconcileTooltip()`**

Replace the existing method (lines 88–97):

```typescript
    getReconcileTooltip(element: ReconciliationStatus): string {
        if (element.hasUnprocessedTransactions) {
            return 'Cannot reconcile while there are unprocessed transactions';
        }
        const delta = Math.abs(element.delta || 0);
        if (delta > RECONCILIATION_TOLERANCE) {
            if (!element.hasBankImporter) {
                return `Balance differs by ${delta.toFixed(2)} from last reconciliation — click to confirm`;
            }
            return `Delta is too large to reconcile (${delta.toFixed(2)})`;
        }
        return 'Mark as Reconciled';
    }
```

- [ ] **Step 4: Replace `reconcile()` to add confirmation dialog for large-delta no-importer accounts**

Replace the existing method (lines 99–118):

```typescript
    reconcile(status: ReconciliationStatus): void {
        if (!status.accountId || !status.currencyId) return;

        const needsConfirmation =
            !status.hasBankImporter &&
            Math.abs(status.delta ?? 0) > RECONCILIATION_TOLERANCE;

        if (needsConfirmation) {
            const delta = Math.abs(status.delta ?? 0).toFixed(2);
            const dialogRef = this.dialog.open(ConfirmationDialogComponent, {
                data: {
                    title: 'Confirm Balance',
                    message: `The current balance differs from the last reconciled balance by ${delta}. This exceeds the normal tolerance. Are you sure the current balance is correct?`,
                    confirmText: 'Confirm',
                    cancelText: 'Cancel',
                } as ConfirmationDialogData,
            });
            dialogRef.afterClosed().subscribe((confirmed: boolean) => {
                if (confirmed) {
                    this.doReconcile(status);
                }
            });
        } else {
            this.doReconcile(status);
        }
    }

    private doReconcile(status: ReconciliationStatus): void {
        this.reconciliationService
            .reconcile(status.accountId!, {
                currencyId: status.currencyId!,
                balance: 0,
            })
            .subscribe({
                next: () => {
                    this.snackBar.open('Account reconciled successfully', 'Close', {
                        duration: 3000,
                    });
                    this.loadStatuses();
                },
                error: (_err: any) => {
                    this.snackBar.open('Failed to reconcile account', 'Close', { duration: 3000 });
                },
            });
    }
```

- [ ] **Step 5: Replace `getStatusClass()`**

Replace the existing method (lines 182–186):

```typescript
    getStatusClass(status: ReconciliationStatus): string {
        if (status.hasUnprocessedTransactions) return 'status-yellow';
        // No-importer accounts are always in a valid state (large delta is expected and confirmable)
        if (!status.hasBankImporter) return 'status-green';
        if (Math.abs(status.delta || 0) > RECONCILIATION_TOLERANCE) return 'status-red';
        return 'status-green';
    }
```

- [ ] **Step 6: Build the frontend**

```bash
cd /Users/ek/work/GeekBudgetBE/frontend && npm run build -- --no-progress 2>&1 | tail -20
```

Expected: build succeeds with no TypeScript errors.

- [ ] **Step 7: Run frontend lint**

```bash
cd /Users/ek/work/GeekBudgetBE/frontend && npm run lint 2>&1 | tail -20
```

Expected: no lint errors.

- [ ] **Step 8: Commit**

```bash
cd /Users/ek/work/GeekBudgetBE
git add frontend/src/app/features/reconciliation/reconciliation.component.ts
git commit -m "feat: add confirmation dialog and update tooltip/row-color for no-importer reconciliation"
```

---

## Task 5: Frontend — Add Component Tests

**Files:**
- Create or modify: `frontend/src/app/features/reconciliation/reconciliation.component.spec.ts`

Check if the spec file already exists first. If it does, append the new `describe` block inside the existing one. If not, create it from scratch.

- [ ] **Step 1: Check for existing spec file**

```bash
ls /Users/ek/work/GeekBudgetBE/frontend/src/app/features/reconciliation/reconciliation.component.spec.ts 2>&1
```

If absent, proceed to Step 2. If present, read it first to understand its structure before modifying.

- [ ] **Step 2: Write the spec file (create or append)**

Create (or replace if empty) `frontend/src/app/features/reconciliation/reconciliation.component.spec.ts`:

```typescript
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReconciliationComponent } from './reconciliation.component';
import { ReconciliationService } from './services/reconciliation.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { of, Subject } from 'rxjs';
import { ReconciliationStatus } from '../../core/api/models/reconciliation-status';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('ReconciliationComponent', () => {
    let component: ReconciliationComponent;
    let fixture: ComponentFixture<ReconciliationComponent>;
    let mockReconciliationService: jasmine.SpyObj<ReconciliationService>;
    let mockSnackBar: jasmine.SpyObj<MatSnackBar>;
    let mockDialog: jasmine.SpyObj<MatDialog>;
    let mockRouter: jasmine.SpyObj<Router>;

    const baseStatus: ReconciliationStatus = {
        accountId: 'acc1',
        accountName: 'Cash',
        currencyId: 'USD',
        currencySymbol: '$',
        bankBalance: 100,
        appBalance: 200,
        delta: 100,
        hasUnprocessedTransactions: false,
        hasBankImporter: false,
        isManualReconciliationEnabled: true,
    };

    beforeEach(async () => {
        mockReconciliationService = jasmine.createSpyObj('ReconciliationService', [
            'loadStatuses',
            'reconcile',
            'enableManual',
            'getTransactionsSince',
            'analyzeDisbalance',
        ]);
        mockSnackBar = jasmine.createSpyObj('MatSnackBar', ['open']);
        mockDialog = jasmine.createSpyObj('MatDialog', ['open']);
        mockRouter = jasmine.createSpyObj('Router', ['navigate']);

        mockReconciliationService.loadStatuses.and.returnValue(of([]));

        await TestBed.configureTestingModule({
            imports: [ReconciliationComponent, NoopAnimationsModule],
            providers: [
                { provide: ReconciliationService, useValue: mockReconciliationService },
                { provide: MatSnackBar, useValue: mockSnackBar },
                { provide: MatDialog, useValue: mockDialog },
                { provide: Router, useValue: mockRouter },
            ],
        }).compileComponents();

        fixture = TestBed.createComponent(ReconciliationComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    describe('getReconcileTooltip', () => {
        it('returns unprocessed message when hasUnprocessedTransactions', () => {
            const s = { ...baseStatus, hasUnprocessedTransactions: true };
            expect(component.getReconcileTooltip(s)).toContain('unprocessed');
        });

        it('returns confirm message for no-importer account with large delta', () => {
            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            const tooltip = component.getReconcileTooltip(s);
            expect(tooltip).toContain('click to confirm');
            expect(tooltip).toContain('100.00');
        });

        it('returns too-large message for importer account with large delta', () => {
            const s = { ...baseStatus, hasBankImporter: true, delta: 100 };
            const tooltip = component.getReconcileTooltip(s);
            expect(tooltip).toContain('too large');
        });

        it('returns mark-as-reconciled for small delta', () => {
            const s = { ...baseStatus, delta: 0.005 };
            expect(component.getReconcileTooltip(s)).toBe('Mark as Reconciled');
        });
    });

    describe('getStatusClass', () => {
        it('returns status-yellow when hasUnprocessedTransactions', () => {
            expect(component.getStatusClass({ ...baseStatus, hasUnprocessedTransactions: true }))
                .toBe('status-yellow');
        });

        it('returns status-green for no-importer account even with large delta', () => {
            expect(component.getStatusClass({ ...baseStatus, hasBankImporter: false, delta: 100 }))
                .toBe('status-green');
        });

        it('returns status-red for importer account with large delta', () => {
            expect(component.getStatusClass({ ...baseStatus, hasBankImporter: true, delta: 100 }))
                .toBe('status-red');
        });

        it('returns status-green for importer account within tolerance', () => {
            expect(component.getStatusClass({ ...baseStatus, hasBankImporter: true, delta: 0.005 }))
                .toBe('status-green');
        });
    });

    describe('reconcile', () => {
        it('opens confirmation dialog when no-importer and large delta', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = { afterClosed: () => afterClosedSubject.asObservable() } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);

            expect(mockDialog.open).toHaveBeenCalled();
            expect(mockReconciliationService.reconcile).not.toHaveBeenCalled();
        });

        it('does not call API when user cancels confirmation dialog', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = { afterClosed: () => afterClosedSubject.asObservable() } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);
            afterClosedSubject.next(false); // user cancels

            expect(mockReconciliationService.reconcile).not.toHaveBeenCalled();
        });

        it('calls API when user confirms dialog', () => {
            const afterClosedSubject = new Subject<boolean>();
            const mockDialogRef = { afterClosed: () => afterClosedSubject.asObservable() } as MatDialogRef<any>;
            mockDialog.open.and.returnValue(mockDialogRef);
            mockReconciliationService.reconcile.and.returnValue(of({}));
            mockReconciliationService.loadStatuses.and.returnValue(of([]));

            const s = { ...baseStatus, hasBankImporter: false, delta: 100 };
            component.reconcile(s);
            afterClosedSubject.next(true); // user confirms

            expect(mockReconciliationService.reconcile).toHaveBeenCalledWith('acc1', {
                currencyId: 'USD',
                balance: 0,
            });
        });

        it('calls API directly (no dialog) when delta is within tolerance', () => {
            mockReconciliationService.reconcile.and.returnValue(of({}));
            mockReconciliationService.loadStatuses.and.returnValue(of([]));

            const s = { ...baseStatus, hasBankImporter: false, delta: 0.005 };
            component.reconcile(s);

            expect(mockDialog.open).not.toHaveBeenCalled();
            expect(mockReconciliationService.reconcile).toHaveBeenCalled();
        });

        it('disables reconcile button when hasUnprocessedTransactions regardless of importer status', () => {
            // Verify via getReconcileTooltip (template uses [disabled] which we test via tooltip)
            const s = { ...baseStatus, hasBankImporter: false, hasUnprocessedTransactions: true, delta: 100 };
            expect(component.getReconcileTooltip(s)).toContain('unprocessed');
        });
    });
});
```

- [ ] **Step 3: Run the frontend tests**

```bash
cd /Users/ek/work/GeekBudgetBE/frontend && npm run test -- --watch=false --browsers=ChromeHeadless 2>&1 | tail -40
```

Expected: all new tests pass. If any fail due to missing imports or provider issues, fix them before proceeding.

- [ ] **Step 4: Commit**

```bash
cd /Users/ek/work/GeekBudgetBE
git add frontend/src/app/features/reconciliation/reconciliation.component.spec.ts
git commit -m "test: add ReconciliationComponent unit tests for no-importer reconciliation"
```

---

## Task 6: Final Validation

- [ ] **Step 1: Run the full pipeline**

```bash
cd /Users/ek/work/GeekBudgetBE && make all
```

Expected: build, test, validate OpenAPI, lint — all pass.

- [ ] **Step 2: Commit if any auto-fixes were applied, then finish**

If `make all` applied any formatting changes:
```bash
cd /Users/ek/work/GeekBudgetBE
git add -p  # review and stage any auto-fixed files
git commit -m "chore: apply formatting fixes from make all"
```

- [ ] **Step 3: Use superpowers:finishing-a-development-branch**

Call the skill to decide on merge/PR options.
