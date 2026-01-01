import { Component, inject, OnInit, signal, computed, effect } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe, CurrencyPipe, CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { UserService } from '../../core/services/user.service'; // Added
import { CurrencyService } from '../currencies/services/currency.service'; // Added
import { BudgetItem } from '../../core/api/models/budget-item';
import { BudgetStatus } from '../../core/api/models/budget-status';
import { LayoutService } from '../../layout/services/layout.service';

import { BudgetMatrixEditComponent } from './components/budget-matrix-edit/budget-matrix-edit.component';

interface MatrixCell {
    month: string; // YYYY-MM-01
    amount: number; // Planned (Converted)
    rawAmount: number; // For editing
    spent: number; // (Converted)
    budgetItemId?: string; // If exists for editing
    calculatedAvailable?: number; // For styling
    isVirtual?: boolean;
    isPastMonth?: boolean;
}

interface MatrixRow {
    account: { id: string; name: string };
    cells: MatrixCell[];
    totalPlanned: number;
    totalSpent: number;
    averageSpent: number;
}

@Component({
    selector: 'app-budget-items',
    imports: [
        MatTableModule,
        MatSortModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
        DatePipe,
        CurrencyPipe,
        CommonModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatDialogModule,
        MatSlideToggleModule,
    ],
    template: `
        <div class="budget-items-container">
            @if (!sidenavOpened()) {
                <h1 class="page-title">Budget Items</h1>
            }

            @if (loading()) {
                <div class="loading-container">
                    <mat-spinner></mat-spinner>
                </div>
            } @else {
                <div class="matrix-container">
                    <div
                        class="matrix-header"
                        style="display: flex; flex-direction: column; gap: 16px; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid #eee;"
                    >
                        <div
                            class="title-row"
                            style="display: flex; align-items: center; gap: 16px;"
                        >
                            <h2>Budget Matrix</h2>
                            <span
                                class="currency-badge"
                                style="background: #e3f2fd; color: #1976d2; padding: 4px 8px; border-radius: 12px; font-size: 12px; font-weight: 500;"
                                *ngIf="targetCurrencySymbol() as curr"
                                >Currency: {{ curr }}</span
                            >
                        </div>

                        <div
                            class="controls-row"
                            style="display: flex; justify-content: space-between; align-items: center; width: 100%;"
                        >
                            <div
                                class="nav-controls"
                                style="display: flex; align-items: center; gap: 8px;"
                            >
                                <button mat-icon-button (click)="shiftMonths(-1)">
                                    <mat-icon>chevron_left</mat-icon>
                                </button>
                                <div
                                    class="month-display"
                                    style="display: flex; align-items: center; gap: 8px; font-weight: 500; font-size: 14px; min-width: 150px; justify-content: center;"
                                >
                                    <span>{{ months()[0] | date: 'MMM yyyy' }}</span>
                                    <span class="separator">-</span>
                                    <span>{{
                                        months()[months().length - 1] | date: 'MMM yyyy'
                                    }}</span>
                                </div>
                                <button mat-icon-button (click)="shiftMonths(1)">
                                    <mat-icon>chevron_right</mat-icon>
                                </button>
                            </div>

                            <div
                                class="density-controls"
                                style="display: flex; align-items: center; gap: 4px;"
                            >
                                <button
                                    mat-icon-button
                                    (click)="changeMonthCount(-1)"
                                    [disabled]="monthCount() <= 1"
                                >
                                    <mat-icon>remove</mat-icon>
                                </button>
                                <span style="font-size: 12px; color: #666; font-weight: 500;"
                                    >{{ monthCount() }} Months</span
                                >
                                <button
                                    mat-icon-button
                                    (click)="changeMonthCount(1)"
                                    [disabled]="monthCount() >= 12"
                                >
                                    <mat-icon>add</mat-icon>
                                </button>
                            </div>

                            <mat-slide-toggle
                                [checked]="includeHidden()"
                                (change)="includeHidden.set($event.checked)"
                                color="primary"
                                style="margin-left: 16px;"
                            >
                                <span style="font-size: 12px;">Hidden</span>
                            </mat-slide-toggle>
                        </div>
                    </div>

                    <div class="grid-wrapper">
                        <table class="matrix-table">
                            <thead>
                                <tr>
                                    <th class="sticky-col-header">Account</th>
                                    @for (month of months(); track month) {
                                        <th>{{ month | date: 'MMM yy' }}</th>
                                    }
                                    <th class="total-col-header">Total</th>
                                </tr>
                            </thead>
                            <tbody>
                                @for (row of matrixData(); track row.account.id) {
                                    <tr>
                                        <td class="sticky-col-header row-header">
                                            <div class="account-cell">
                                                <span class="account-name">{{
                                                    row.account.name
                                                }}</span>
                                                <span class="account-average"
                                                    >Avg:
                                                    {{
                                                        row.averageSpent
                                                            | currency
                                                                : preferredCurrency()?.name || ''
                                                    }}</span
                                                >
                                            </div>
                                        </td>
                                        @for (cell of row.cells; track cell.month) {
                                            <td
                                                class="cell-interactive"
                                                [class.status-ok]="
                                                    !!cell.budgetItemId &&
                                                    (cell.calculatedAvailable || 0) >= 0
                                                "
                                                [class.status-over]="
                                                    (cell.calculatedAvailable || 0) < 0
                                                "
                                                [class.status-virtual]="cell.isVirtual"
                                                (click)="editCell(row.account, cell)"
                                            >
                                                <div class="cell-content-modern">
                                                    <div class="cell-row top-row">
                                                        <span class="label">Budget</span>
                                                        <span class="value-primary">{{
                                                            cell.amount
                                                                | currency
                                                                    : preferredCurrency()?.name ||
                                                                          ''
                                                                    : 'symbol-narrow'
                                                                    : '1.0-0'
                                                        }}</span>
                                                    </div>
                                                    <div class="cell-row middle-row">
                                                        <span class="label">Spent</span>
                                                        <span class="value-secondary">{{
                                                            cell.spent
                                                                | currency
                                                                    : preferredCurrency()?.name ||
                                                                          ''
                                                                    : 'symbol-narrow'
                                                                    : '1.0-0'
                                                        }}</span>
                                                    </div>
                                                    <div class="progress-bar-container">
                                                        <div
                                                            class="progress-bar-fill"
                                                            [style.width.%]="
                                                                (cell.amount > 0
                                                                    ? (cell.spent / cell.amount) *
                                                                      100
                                                                    : cell.spent > 0
                                                                      ? 100
                                                                      : 0
                                                                ) | number: '1.0-0'
                                                            "
                                                            [class.over-budget]="
                                                                cell.spent > cell.amount
                                                            "
                                                            [class.unbudgeted-spent]="
                                                                cell.isVirtual
                                                            "
                                                        ></div>
                                                    </div>
                                                    <div class="percentage-label">
                                                        {{
                                                            cell.amount > 0
                                                                ? (cell.spent / cell.amount
                                                                  | percent: '1.0-0')
                                                                : '0%'
                                                        }}
                                                    </div>
                                                </div>
                                            </td>
                                        }
                                        <td
                                            class="total-cell"
                                            [class.status-over]="
                                                row.totalPlanned - row.totalSpent < 0
                                            "
                                            [class.status-ok]="
                                                row.totalPlanned - row.totalSpent >= 0
                                            "
                                        >
                                            <div class="cell-content">
                                                <span class="planned">{{
                                                    row.totalPlanned
                                                        | currency: preferredCurrency()?.name || ''
                                                }}</span>
                                                <span class="divider">/</span>
                                                <span class="spent">{{
                                                    row.totalSpent
                                                        | currency: preferredCurrency()?.name || ''
                                                }}</span>
                                            </div>
                                        </td>
                                    </tr>
                                }
                                <!-- Grand Total Row -->
                                <tr class="total-row">
                                    <td class="sticky-col-header row-header">Total</td>
                                    @for (col of columnTotals(); track col.month) {
                                        <td
                                            [class.status-over]="col.planned - col.spent < 0"
                                            [class.status-ok]="col.planned - col.spent >= 0"
                                        >
                                            <div class="cell-content">
                                                <span class="planned">{{
                                                    col.planned
                                                        | currency: preferredCurrency()?.name || ''
                                                }}</span>
                                                <span class="divider">/</span>
                                                <span class="spent">{{
                                                    col.spent
                                                        | currency: preferredCurrency()?.name || ''
                                                }}</span>
                                            </div>
                                        </td>
                                    }
                                    <td
                                        class="total-cell"
                                        [class.status-over]="
                                            grandTotal().planned - grandTotal().spent < 0
                                        "
                                        [class.status-ok]="
                                            grandTotal().planned - grandTotal().spent >= 0
                                        "
                                    >
                                        <div class="cell-content">
                                            <span class="planned">{{
                                                grandTotal().planned
                                                    | currency: preferredCurrency()?.name || ''
                                            }}</span>
                                            <span class="divider">/</span>
                                            <span class="spent">{{
                                                grandTotal().spent
                                                    | currency: preferredCurrency()?.name || ''
                                            }}</span>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            }
        </div>
    `,
    styleUrl: './budget-items.component.scss',
})
export class BudgetItemsComponent implements OnInit {
    private readonly budgetItemService = inject(BudgetItemService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);
    private readonly accountService = inject(AccountService);
    private readonly userService = inject(UserService);
    private readonly currencyService = inject(CurrencyService);
    private readonly dialog = inject(MatDialog);
    private readonly fb = inject(FormBuilder);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;
    protected readonly Creating = signal(false);

    protected readonly budgetItems = this.budgetItemService.budgetItems;
    protected readonly budgetStatus = this.budgetItemService.budgetStatus;
    protected readonly accounts = this.accountService.accounts;
    protected readonly loading = this.budgetItemService.loading;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly user = this.userService.user;

    // Configuration Signals
    protected readonly startDate = signal(new Date()); // Start of the view
    protected readonly monthCount = signal(3); // Number of months to show
    protected readonly includeHidden = signal(false);

    // Computed state
    protected readonly preferredCurrency = computed(() => {
        const user = this.user();
        const currencies = this.currencies();
        if (user?.favoriteCurrencyId) {
            return currencies.find((c) => c.id === user.favoriteCurrencyId);
        }
        return currencies.length > 0 ? currencies[0] : null;
    });

    protected readonly currencyMap = computed(() => {
        const map = new Map<string, string>(); // Id -> Symbol/Name
        this.currencies().forEach((c) => map.set(c.id, c.name));
        return map;
    });

    protected readonly targetCurrencySymbol = computed(() => {
        return this.preferredCurrency()?.name || '';
    });

    protected readonly months = computed(() => {
        const anchor = this.startDate(); // This is now the END of the view (inclusive)
        const count = this.monthCount();
        const result: string[] = [];

        // Use UTC to avoid timezone shifts when calling toISOString
        const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));

        for (let i = -(count - 1); i <= 0; i++) {
            const d = new Date(anchorMonth);
            d.setUTCMonth(d.getUTCMonth() + i);
            result.push(d.toISOString());
        }
        return result;
    });

    protected readonly matrixData = computed(() => {
        // ... existing matrixData logic ...
        const accs = this.accounts().filter((a) => a.type === 'expense');

        const items = this.budgetItems();
        const status = this.budgetStatus();
        const months = this.months();
        const averages = this.accountService.averages();

        // Maps
        const itemMap = new Map<string, BudgetItem>();
        items.forEach((i) => {
            const m = i.date ? i.date.substring(0, 7) : '';
            itemMap.set(`${i.accountId}_${m}`, i);
        });

        const statusMap = new Map<string, BudgetStatus>();
        status.forEach((s) => {
            const m = s.date.substring(0, 7);
            statusMap.set(`${s.accountId}_${m}`, s);
        });

        const avgMap = new Map(averages.map((a) => [a.accountId, a.averageSpent]));

        const now = new Date();
        const currentYear = now.getUTCFullYear();
        const currentMonth = now.getUTCMonth();

        const rows: MatrixRow[] = accs.map((acc) => {
            let rowTotalPlanned = 0;
            let rowTotalSpent = 0;

            const cells: MatrixCell[] = months.map((mStr) => {
                const mKey = mStr.substring(0, 7);
                const item = itemMap.get(`${acc.id}_${mKey}`);
                const stat = statusMap.get(`${acc.id}_${mKey}`);

                // Check if this cell is for the current month (Real World time)
                const cellDate = new Date(mStr);
                const isCurrentMonth =
                    cellDate.getUTCFullYear() === currentYear &&
                    cellDate.getUTCMonth() === currentMonth;

                // Check if past month (simple comparison since we iterate)
                // Actually easier to compare tokens or date objects
                // currentYear/Month is local time from new Date()
                // cellDate is UTC from string.
                // Let's rely on flexible comparison:
                const nowTotalMonths = currentYear * 12 + currentMonth;
                const cellTotalMonths = cellDate.getUTCFullYear() * 12 + cellDate.getUTCMonth();
                const isPastMonth = cellTotalMonths < nowTotalMonths;

                const spentDisplay = stat?.spent ?? 0;

                let amountDisplay = 0;
                const rawAmount = item?.amount ?? 0;
                let isVirtual = false;

                if (item) {
                    // Explicit budget exists
                    amountDisplay = stat?.budgeted ?? item.amount;
                } else {
                    // No explicit budget
                    if (isCurrentMonth) {
                        // Current month: No virtual budget. Strict.
                        amountDisplay = 0;
                    } else {
                        // Past/Future months: Virtual budget = Spent
                        amountDisplay = spentDisplay;
                        if (spentDisplay > 0) {
                            isVirtual = true;
                        }
                    }
                }

                const available = amountDisplay - spentDisplay;

                rowTotalPlanned += amountDisplay;
                rowTotalSpent += spentDisplay;

                return {
                    month: mStr,
                    amount: amountDisplay,
                    rawAmount: rawAmount,
                    spent: spentDisplay,
                    budgetItemId: item?.id,
                    calculatedAvailable: available,
                    isVirtual: isVirtual,
                    isPastMonth: isPastMonth,
                };
            });

            return {
                account: { id: acc.id, name: acc.name },
                cells,
                totalPlanned: rowTotalPlanned,
                totalSpent: rowTotalSpent,
                averageSpent: avgMap.get(acc.id) ?? 0,
            };
        });

        return rows;
    });

    protected readonly columnTotals = computed(() => {
        const rows = this.matrixData();
        const months = this.months();
        if (rows.length === 0) return [];

        // Init totals
        const totals = months.map((m) => ({ month: m, planned: 0, spent: 0 }));

        rows.forEach((row) => {
            row.cells.forEach((cell, idx) => {
                if (totals[idx]) {
                    totals[idx].planned += cell.amount;
                    totals[idx].spent += cell.spent;
                }
            });
        });
        return totals;
    });

    protected readonly grandTotal = computed(() => {
        const cols = this.columnTotals();
        return cols.reduce(
            (acc, curr) => ({
                planned: acc.planned + curr.planned,
                spent: acc.spent + curr.spent,
            }),
            { planned: 0, spent: 0 },
        );
    });

    // Effects
    constructor() {
        // Reload data when params change
        effect(() => {
            const anchor = this.startDate();
            const count = this.monthCount();
            const currency = this.preferredCurrency();
            const includeHidden = this.includeHidden();

            // Calculate query range
            const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));

            // From: Anchor - (count-1) months
            const fromDate = new Date(anchorMonth);
            fromDate.setUTCMonth(fromDate.getUTCMonth() - (count - 1));

            // To: Anchor + 1 month (exclusive)
            const toDate = new Date(anchorMonth);
            toDate.setUTCMonth(toDate.getUTCMonth() + 1);

            const from = fromDate.toISOString();
            const to = toDate.toISOString();

            const currencyId = currency?.id;

            // Load status with conversion
            if (currencyId) {
                this.budgetItemService
                    .loadBudgetStatus(from, to, currencyId, includeHidden)
                    .subscribe();
                this.accountService.loadYearlyExpenses(currencyId).subscribe();
            } else {
                this.budgetItemService
                    .loadBudgetStatus(from, to, undefined, includeHidden)
                    .subscribe();
                this.accountService.loadYearlyExpenses().subscribe();
            }
        });
    }

    ngOnInit(): void {
        this.budgetItemService.loadBudgetItems().subscribe();
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
        this.userService.loadUser().subscribe(); // Loads user -> triggers effect via preferredCurrency
    }

    // Actions
    protected shiftMonths(delta: number): void {
        const current = this.startDate();
        this.startDate.set(new Date(current.getFullYear(), current.getMonth() + delta, 1));
    }

    protected changeMonthCount(delta: number): void {
        const current = this.monthCount();
        const next = Math.max(1, Math.min(12, current + delta));
        this.monthCount.set(next);
    }

    protected editCell(account: { id: string; name: string }, cell: MatrixCell): void {
        // We pass the RAW amount for editing if we can find it?
        // Wait, cell.amount is CONVERTED. cell.rawAmount is RAW.
        // If we don't have a budget item yet, rawAmount is 0.
        // The dialog should behave effectively.

        // We need to tell the dialog what currency the Account checks?
        // The dialog just shows "Amount".
        // I added `rawAmount` to MatrixCell interface in computed above. I need to update interface def.

        const dialogRef = this.dialog.open(BudgetMatrixEditComponent, {
            data: {
                accountId: account.id,
                accountName: account.name,
                month: cell.month,
                currentAmount: cell.rawAmount ?? 0,
            },
            width: '300px',
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result !== undefined) {
                this.saveBudget(account.id, cell, Number(result));
            }
        });
    }

    private saveBudget(accountId: string, cell: MatrixCell, newAmount: number): void {
        if (cell.budgetItemId) {
            this.budgetItemService
                .update(cell.budgetItemId, {
                    accountId: accountId,
                    amount: newAmount,
                    date: cell.month,
                    description: 'Matrix Edit',
                })
                .subscribe({
                    next: () => this.refreshData(),
                    error: () =>
                        this.snackBar.open('Failed to update', 'Close', { duration: 3000 }),
                });
        } else {
            this.budgetItemService
                .create({
                    accountId: accountId,
                    amount: newAmount,
                    date: cell.month,
                    description: 'Matrix Create',
                })
                .subscribe({
                    next: () => this.refreshData(),
                    error: () =>
                        this.snackBar.open('Failed to create', 'Close', { duration: 3000 }),
                });
        }
    }

    private refreshData() {
        this.snackBar.open('Budget saved', 'Close', { duration: 2000 });
        this.budgetItemService.loadBudgetItems().subscribe();

        // We need to trigger re-fetch of status with current signals using correct view logic
        const anchor = this.startDate();
        const count = this.monthCount();
        const currency = this.preferredCurrency();
        const includeHidden = this.includeHidden();

        // Use UTC to avoid timezone shifts
        const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));

        // From: Anchor - (count-1) months
        const fromDate = new Date(anchorMonth);
        fromDate.setUTCMonth(fromDate.getUTCMonth() - (count - 1));

        // To: Anchor + 1 month (exclusive)
        const toDate = new Date(anchorMonth);
        toDate.setUTCMonth(toDate.getUTCMonth() + 1);

        const from = fromDate.toISOString();
        const to = toDate.toISOString();

        const currencyId = currency?.id;

        if (currencyId) {
            this.budgetItemService
                .loadBudgetStatus(from, to, currencyId, includeHidden)
                .subscribe();
            this.accountService.loadYearlyExpenses(currencyId).subscribe();
        } else {
            this.budgetItemService.loadBudgetStatus(from, to, undefined, includeHidden).subscribe();
            this.accountService.loadYearlyExpenses().subscribe();
        }
    }
}
