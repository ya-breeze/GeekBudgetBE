import { Component, inject, OnInit, signal, computed, effect } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe, CurrencyPipe, CommonModule } from '@angular/common';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { UserService } from '../../core/services/user.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { BudgetItem } from '../../core/api/models/budget-item';
import { BudgetStatus } from '../../core/api/models/budget-status';
import { LayoutService } from '../../layout/services/layout.service';

import { BudgetMatrixEditComponent } from './components/budget-matrix-edit/budget-matrix-edit.component';

interface MatrixCell {
    month: string; // YYYY-MM-01
    amount: number; // Planned (Converted)
    rawAmount: number; // For editing
    spent: number; // (Converted)
    budgetItemId?: string;
    calculatedAvailable?: number;
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

interface ListRow {
    account: { id: string; name: string };
    totalPlanned: number;
    totalSpent: number;
    averageSpent: number;
    percent: number;
    isOver: boolean;
    isVirtual: boolean;
    progressWidth: number;
}

@Component({
    selector: 'app-budget-items',
    imports: [
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
        DatePipe,
        CurrencyPipe,
        CommonModule,
        MatDialogModule,
        MatSlideToggleModule,
    ],
    template: `
        <div class="budget-items-container">
            @if (!sidenavOpened()) {
                <h1 class="page-title">Budget</h1>
            }

            @if (loading()) {
                <div class="loading-container">
                    <mat-spinner></mat-spinner>
                </div>
            } @else {
                <div class="budget-container">
                    <!-- Header -->
                    <div class="list-header">
                        <div class="nav-row">
                            <button mat-icon-button (click)="shiftMonths(-1)">
                                <mat-icon>chevron_left</mat-icon>
                            </button>
                            <span class="period-display">
                                {{ months()[0] | date: 'MMM yyyy' }}
                                @if (selectedPeriod() > 1) {
                                    &ndash; {{ months()[months().length - 1] | date: 'MMM yyyy' }}
                                }
                            </span>
                            <button mat-icon-button (click)="shiftMonths(1)">
                                <mat-icon>chevron_right</mat-icon>
                            </button>

                            <div class="period-selector">
                                @for (p of periodOptions; track p) {
                                    <button
                                        class="period-btn"
                                        [class.active]="selectedPeriod() === p"
                                        (click)="selectPeriod(p)"
                                    >
                                        {{ p }}m
                                    </button>
                                }
                            </div>

                            @if (targetCurrencySymbol(); as curr) {
                                <span class="currency-badge">{{ curr }}</span>
                            }

                            <mat-slide-toggle
                                [checked]="includeHidden()"
                                (change)="includeHidden.set($event.checked)"
                                color="primary"
                                class="hidden-toggle"
                            >
                                <span class="toggle-label">Hidden</span>
                            </mat-slide-toggle>
                        </div>

                        <!-- Summary progress bar -->
                        <div class="summary-bar">
                            <div class="summary-progress">
                                <div
                                    class="summary-fill"
                                    [style.width.%]="grandTotalPercent()"
                                    [class.fill-over]="grandTotal().spent > grandTotal().planned"
                                ></div>
                            </div>
                            <span class="summary-text">
                                {{
                                    grandTotal().spent
                                        | currency
                                            : preferredCurrency()?.name || ''
                                            : 'symbol-narrow'
                                            : '1.0-0'
                                }}
                                /
                                {{
                                    grandTotal().planned
                                        | currency
                                            : preferredCurrency()?.name || ''
                                            : 'symbol-narrow'
                                            : '1.0-0'
                                }}
                            </span>
                        </div>
                    </div>

                    <!-- Account list -->
                    <div class="budget-list">
                        @for (row of listData(); track row.account.id) {
                            <div
                                class="budget-row"
                                [class.over-budget]="row.isOver"
                                (click)="editCurrentMonth(row.account)"
                            >
                                <div class="row-info">
                                    <span class="account-name">{{ row.account.name }}</span>
                                    <span class="account-avg">
                                        Avg:
                                        {{
                                            row.averageSpent
                                                | currency
                                                    : preferredCurrency()?.name || ''
                                                    : 'symbol-narrow'
                                                    : '1.0-0'
                                        }}
                                    </span>
                                </div>
                                <div class="row-progress-area">
                                    <div class="progress-track">
                                        <div
                                            class="progress-fill"
                                            [style.width.%]="row.progressWidth"
                                            [class.fill-over]="row.isOver"
                                            [class.fill-virtual]="row.isVirtual && !row.isOver"
                                        ></div>
                                    </div>
                                    <div class="row-stats">
                                        <span class="pct" [class.text-over]="row.isOver">
                                            {{ row.percent | number: '1.0-0' }}%
                                        </span>
                                        @if (row.isOver) {
                                            <mat-icon class="warn-icon">warning</mat-icon>
                                        }
                                        <span class="amounts">
                                            {{
                                                row.totalSpent
                                                    | currency
                                                        : preferredCurrency()?.name || ''
                                                        : 'symbol-narrow'
                                                        : '1.0-0'
                                            }}
                                            /
                                            {{
                                                row.totalPlanned
                                                    | currency
                                                        : preferredCurrency()?.name || ''
                                                        : 'symbol-narrow'
                                                        : '1.0-0'
                                            }}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        }
                    </div>

                    <!-- Total row -->
                    <div
                        class="total-row"
                        [class.over-budget]="grandTotal().spent > grandTotal().planned"
                    >
                        <div class="row-info">
                            <span class="account-name">Total</span>
                        </div>
                        <div class="row-progress-area">
                            <div class="progress-track">
                                <div
                                    class="progress-fill"
                                    [style.width.%]="grandTotalPercent()"
                                    [class.fill-over]="grandTotal().spent > grandTotal().planned"
                                ></div>
                            </div>
                            <div class="row-stats">
                                <span
                                    class="pct"
                                    [class.text-over]="grandTotal().spent > grandTotal().planned"
                                >
                                    {{ grandTotalPercent() | number: '1.0-0' }}%
                                </span>
                                <span class="amounts">
                                    {{
                                        grandTotal().spent
                                            | currency
                                                : preferredCurrency()?.name || ''
                                                : 'symbol-narrow'
                                                : '1.0-0'
                                    }}
                                    /
                                    {{
                                        grandTotal().planned
                                            | currency
                                                : preferredCurrency()?.name || ''
                                                : 'symbol-narrow'
                                                : '1.0-0'
                                    }}
                                </span>
                            </div>
                        </div>
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

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;

    protected readonly budgetItems = this.budgetItemService.budgetItems;
    protected readonly budgetStatus = this.budgetItemService.budgetStatus;
    protected readonly accounts = this.accountService.accounts;
    protected readonly loading = this.budgetItemService.loading;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly user = this.userService.user;

    // Configuration
    protected readonly startDate = signal(new Date());
    protected readonly monthCount = signal(1);
    protected readonly includeHidden = signal(false);
    protected readonly selectedPeriod = signal<1 | 3 | 6 | 12>(1);
    protected readonly periodOptions = [1, 3, 6, 12] as const;

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
        const map = new Map<string, string>();
        this.currencies().forEach((c) => map.set(c.id, c.name));
        return map;
    });

    protected readonly targetCurrencySymbol = computed(() => {
        return this.preferredCurrency()?.name || '';
    });

    protected readonly months = computed(() => {
        const anchor = this.startDate();
        const count = this.monthCount();
        const result: string[] = [];
        const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));
        for (let i = -(count - 1); i <= 0; i++) {
            const d = new Date(anchorMonth);
            d.setUTCMonth(d.getUTCMonth() + i);
            result.push(d.toISOString());
        }
        return result;
    });

    protected readonly matrixData = computed((): MatrixRow[] => {
        const accs = this.accounts().filter((a) => a.type === 'expense');
        const items = this.budgetItems();
        const status = this.budgetStatus();
        const months = this.months();
        const averages = this.accountService.averages();

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

        return accs.map((acc) => {
            let rowTotalPlanned = 0;
            let rowTotalSpent = 0;

            const cells: MatrixCell[] = months.map((mStr) => {
                const mKey = mStr.substring(0, 7);
                const item = itemMap.get(`${acc.id}_${mKey}`);
                const stat = statusMap.get(`${acc.id}_${mKey}`);

                const cellDate = new Date(mStr);
                const isCurrentMonth =
                    cellDate.getUTCFullYear() === currentYear &&
                    cellDate.getUTCMonth() === currentMonth;

                const nowTotalMonths = currentYear * 12 + currentMonth;
                const cellTotalMonths = cellDate.getUTCFullYear() * 12 + cellDate.getUTCMonth();
                const isPastMonth = cellTotalMonths < nowTotalMonths;

                const spentDisplay = stat?.spent ?? 0;

                let amountDisplay = 0;
                const rawAmount = item?.amount ?? 0;
                let isVirtual = false;

                if (item) {
                    amountDisplay = stat?.budgeted ?? item.amount;
                } else {
                    if (isCurrentMonth) {
                        amountDisplay = 0;
                    } else {
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
                    rawAmount,
                    spent: spentDisplay,
                    budgetItemId: item?.id,
                    calculatedAvailable: available,
                    isVirtual,
                    isPastMonth,
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
    });

    protected readonly columnTotals = computed(() => {
        const rows = this.matrixData();
        const months = this.months();
        if (rows.length === 0) return [];

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
        return this.columnTotals().reduce(
            (acc, curr) => ({
                planned: acc.planned + curr.planned,
                spent: acc.spent + curr.spent,
            }),
            { planned: 0, spent: 0 },
        );
    });

    protected readonly listData = computed((): ListRow[] => {
        return this.matrixData().map((row) => {
            const { totalPlanned, totalSpent } = row;
            const pct =
                totalPlanned > 0 ? (totalSpent / totalPlanned) * 100 : totalSpent > 0 ? 100 : 0;
            const isVirtual = row.cells.length > 0 && row.cells.every((c) => c.isVirtual);
            return {
                account: row.account,
                totalPlanned,
                totalSpent,
                averageSpent: row.averageSpent,
                percent: pct,
                isOver: totalSpent > totalPlanned && totalPlanned > 0,
                isVirtual,
                progressWidth: Math.min(pct, 100),
            };
        });
    });

    protected readonly grandTotalPercent = computed((): number => {
        const gt = this.grandTotal();
        if (gt.planned <= 0) return gt.spent > 0 ? 100 : 0;
        return Math.min((gt.spent / gt.planned) * 100, 100);
    });

    constructor() {
        effect(() => {
            const anchor = this.startDate();
            const count = this.monthCount();
            const currency = this.preferredCurrency();
            const includeHidden = this.includeHidden();

            const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));

            const fromDate = new Date(anchorMonth);
            fromDate.setUTCMonth(fromDate.getUTCMonth() - (count - 1));

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
        this.userService.loadUser().subscribe();
    }

    protected selectPeriod(p: 1 | 3 | 6 | 12): void {
        this.selectedPeriod.set(p);
        this.monthCount.set(p);
    }

    protected shiftMonths(delta: number): void {
        const current = this.startDate();
        const step = this.selectedPeriod();
        this.startDate.set(new Date(current.getFullYear(), current.getMonth() + delta * step, 1));
    }

    protected editCurrentMonth(account: { id: string; name: string }): void {
        const months = this.months();
        const lastMonth = months[months.length - 1];
        const row = this.matrixData().find((r) => r.account.id === account.id);
        const lastCell = row?.cells[row.cells.length - 1];
        const cell: MatrixCell = lastCell ?? {
            month: lastMonth,
            amount: 0,
            rawAmount: 0,
            spent: 0,
        };
        this.editCell(account, cell);
    }

    private editCell(account: { id: string; name: string }, cell: MatrixCell): void {
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
                    accountId,
                    amount: newAmount,
                    date: cell.month,
                    description: 'Budget Edit',
                })
                .subscribe({
                    next: () => this.refreshData(),
                    error: () =>
                        this.snackBar.open('Failed to update', 'Close', { duration: 3000 }),
                });
        } else {
            this.budgetItemService
                .create({
                    accountId,
                    amount: newAmount,
                    date: cell.month,
                    description: 'Budget Create',
                })
                .subscribe({
                    next: () => this.refreshData(),
                    error: () =>
                        this.snackBar.open('Failed to create', 'Close', { duration: 3000 }),
                });
        }
    }

    private refreshData(): void {
        this.snackBar.open('Budget saved', 'Close', { duration: 2000 });
        this.budgetItemService.loadBudgetItems().subscribe();

        const anchor = this.startDate();
        const count = this.monthCount();
        const currency = this.preferredCurrency();
        const includeHidden = this.includeHidden();

        const anchorMonth = new Date(Date.UTC(anchor.getFullYear(), anchor.getMonth(), 1));

        const fromDate = new Date(anchorMonth);
        fromDate.setUTCMonth(fromDate.getUTCMonth() - (count - 1));

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
