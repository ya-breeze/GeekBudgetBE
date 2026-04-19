import { Component, inject, OnInit, signal, computed, effect } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe, CurrencyPipe, CommonModule } from '@angular/common';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { RouterLink } from '@angular/router';

import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { UserService } from '../../core/services/user.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { LayoutService } from '../../layout/services/layout.service';
import { BudgetMatrixEditComponent } from './components/budget-matrix-edit/budget-matrix-edit.component';
import { BudgetItem } from '../../core/api/models/budget-item';
import { BudgetStatus } from '../../core/api/models/budget-status';

export type ViewMode = 'split' | 'focus' | 'year';

interface MatrixCell {
    month: string;
    amount: number;
    rawAmount: number;
    spent: number;
    rollover: number;
    available: number;
    budgetItemId?: string;
    isVirtual?: boolean;
    isPastMonth?: boolean;
    isFutureMonth?: boolean;
}

interface MatrixRow {
    account: { id: string; name: string };
    hue: number;
    cells: MatrixCell[];
    totalPlanned: number;
    totalSpent: number;
    averageSpent: number;
}

interface FocusedRow {
    account: { id: string; name: string };
    hue: number;
    averageSpent: number;
    cell: MatrixCell;
    status: string;
    pct: number;
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
        RouterLink,
    ],
    templateUrl: './budget-items.component.html',
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
    protected readonly Math = Math;

    protected readonly startDate = signal(new Date());
    protected readonly viewMode = signal<ViewMode>('split');
    protected readonly density = signal<'comfortable' | 'compact'>('comfortable');
    protected readonly selectedAccountId = signal<string | null>(null);
    protected readonly focusedMonthIdx = signal(5);
    protected readonly includeHidden = signal(false);

    protected readonly monthCount = computed(() => (this.viewMode() === 'year' ? 12 : 6));

    protected readonly preferredCurrency = computed(() => {
        const user = this.user();
        const currencies = this.currencies();
        if (user?.favoriteCurrencyId) {
            return currencies.find((c) => c.id === user.favoriteCurrencyId);
        }
        return currencies.length > 0 ? currencies[0] : null;
    });

    protected readonly currencyName = computed(() => this.preferredCurrency()?.name || '');

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

    protected readonly safeFocusedIdx = computed(() =>
        Math.min(this.focusedMonthIdx(), this.months().length - 1),
    );

    protected readonly focusedMonthLabel = computed(() => {
        const months = this.months();
        if (months.length === 0) return '';
        const d = new Date(months[this.safeFocusedIdx()]);
        return d.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    });

    protected readonly prevMonthLabel = computed(() => {
        const months = this.months();
        const idx = this.safeFocusedIdx();
        if (idx === 0 || months.length === 0) return 'previous month';
        const d = new Date(months[idx - 1]);
        return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
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
            const hue = this.hueFromId(acc.id);

            const cells: MatrixCell[] = months.map((mStr) => {
                const mKey = mStr.substring(0, 7);
                const item = itemMap.get(`${acc.id}_${mKey}`);
                const stat = statusMap.get(`${acc.id}_${mKey}`);

                const cellDate = new Date(mStr);
                const isCurrentMonth =
                    cellDate.getUTCFullYear() === currentYear &&
                    cellDate.getUTCMonth() === currentMonth;

                const nowTotalMonths = currentYear * 12 + currentMonth;
                const cellTotalMonths =
                    cellDate.getUTCFullYear() * 12 + cellDate.getUTCMonth();
                const isPastMonth = cellTotalMonths < nowTotalMonths;
                const isFutureMonth = cellTotalMonths > nowTotalMonths;

                const spentDisplay = stat?.spent ?? 0;
                const rollover = stat?.rollover ?? 0;

                let amountDisplay = 0;
                const rawAmount = item?.amount ?? 0;
                let isVirtual = false;

                if (item) {
                    amountDisplay = stat?.budgeted ?? item.amount;
                } else {
                    if (!isCurrentMonth) {
                        amountDisplay = spentDisplay;
                        if (spentDisplay > 0) isVirtual = true;
                    }
                }

                const available = stat?.available ?? amountDisplay + rollover - spentDisplay;

                rowTotalPlanned += amountDisplay;
                rowTotalSpent += spentDisplay;

                return {
                    month: mStr,
                    amount: amountDisplay,
                    rawAmount,
                    spent: spentDisplay,
                    rollover,
                    available,
                    budgetItemId: item?.id,
                    isVirtual,
                    isPastMonth,
                    isFutureMonth,
                };
            });

            return {
                account: { id: acc.id, name: acc.name },
                hue,
                cells,
                totalPlanned: rowTotalPlanned,
                totalSpent: rowTotalSpent,
                averageSpent: avgMap.get(acc.id) ?? 0,
            };
        });
    });

    protected readonly focusedMonthRows = computed((): FocusedRow[] => {
        const idx = this.safeFocusedIdx();
        return this.matrixData().map((row) => {
            const cell = row.cells[idx] ?? {
                month: '',
                amount: 0,
                rawAmount: 0,
                spent: 0,
                rollover: 0,
                available: 0,
            };
            const pct =
                cell.amount > 0 ? cell.spent / cell.amount : cell.spent > 0 ? 1 : 0;
            return {
                account: row.account,
                hue: row.hue,
                averageSpent: row.averageSpent,
                cell,
                status: this.statusFor(cell.amount, cell.spent),
                pct,
                progressWidth: Math.min(100, pct * 100),
            };
        });
    });

    protected readonly focusedTotals = computed(() => {
        const rows = this.focusedMonthRows();
        const assigned = rows.reduce((s, r) => s + r.cell.amount, 0);
        const spent = rows.reduce((s, r) => s + r.cell.spent, 0);
        const available = rows.reduce((s, r) => s + r.cell.available, 0);
        return { assigned, spent, available, pct: assigned > 0 ? spent / assigned : 0 };
    });

    protected readonly maxCellSpent = computed(() => {
        let max = 0;
        this.matrixData().forEach((row) =>
            row.cells.forEach((c) => {
                if (c.spent > max) max = c.spent;
            }),
        );
        return max || 1;
    });

    protected readonly selectedAccount = computed((): MatrixRow | null => {
        const id = this.selectedAccountId();
        if (!id) return null;
        return this.matrixData().find((r) => r.account.id === id) ?? null;
    });

    protected readonly selectedSpentHistory = computed((): number[] => {
        return this.selectedAccount()?.cells.map((c) => c.spent) ?? [];
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

    protected shiftMonths(delta: number): void {
        const current = this.startDate();
        const step = this.monthCount();
        this.startDate.set(new Date(current.getFullYear(), current.getMonth() + delta * step, 1));
        this.focusedMonthIdx.set(step - 1);
    }

    protected setFocusedMonth(idx: number): void {
        this.focusedMonthIdx.set(idx);
    }

    protected setViewMode(mode: ViewMode): void {
        this.viewMode.set(mode);
        const newCount = mode === 'year' ? 12 : 6;
        this.focusedMonthIdx.set(newCount - 1);
    }

    protected editAssigned(account: { id: string; name: string }, cell: MatrixCell): void {
        this.openEditDialog(account, cell, cell.rawAmount);
    }

    protected editYearCell(account: { id: string; name: string }, cell: MatrixCell): void {
        this.openEditDialog(account, cell, cell.rawAmount);
    }

    protected drawerAssignAvg(): void {
        const sa = this.selectedAccount();
        if (!sa) return;
        const cell = sa.cells[this.safeFocusedIdx()];
        if (!cell) return;
        this.openEditDialog(sa.account, cell, Math.round(sa.averageSpent));
    }

    protected drawerCopyFromPrev(): void {
        const sa = this.selectedAccount();
        if (!sa) return;
        const idx = this.safeFocusedIdx();
        if (idx === 0) return;
        const cell = sa.cells[idx];
        const prevCell = sa.cells[idx - 1];
        if (!cell) return;
        this.openEditDialog(sa.account, cell, prevCell?.rawAmount ?? 0);
    }

    protected drawerCoverOverspend(): void {
        const sa = this.selectedAccount();
        if (!sa) return;
        const cell = sa.cells[this.safeFocusedIdx()];
        if (!cell || cell.available >= 0) return;
        this.openEditDialog(sa.account, cell, cell.amount + Math.abs(cell.available));
    }

    private openEditDialog(
        account: { id: string; name: string },
        cell: MatrixCell,
        prefillAmount: number,
    ): void {
        const dialogRef = this.dialog.open(BudgetMatrixEditComponent, {
            data: {
                accountId: account.id,
                accountName: account.name,
                month: cell.month,
                currentAmount: prefillAmount,
            },
            width: '300px',
        });
        dialogRef.afterClosed().subscribe((result) => {
            if (result !== undefined) this.saveBudget(account.id, cell, Number(result));
        });
    }

    private saveBudget(accountId: string, cell: MatrixCell, newAmount: number): void {
        const obs = cell.budgetItemId
            ? this.budgetItemService.update(cell.budgetItemId, {
                  accountId,
                  amount: newAmount,
                  date: cell.month,
                  description: 'Budget Edit',
              })
            : this.budgetItemService.create({
                  accountId,
                  amount: newAmount,
                  date: cell.month,
                  description: 'Budget Create',
              });

        obs.subscribe({
            next: () => this.refreshData(),
            error: () => this.snackBar.open('Failed to save', 'Close', { duration: 3000 }),
        });
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
        const currencyId = currency?.id;

        if (currencyId) {
            this.budgetItemService
                .loadBudgetStatus(fromDate.toISOString(), toDate.toISOString(), currencyId, includeHidden)
                .subscribe();
            this.accountService.loadYearlyExpenses(currencyId).subscribe();
        } else {
            this.budgetItemService
                .loadBudgetStatus(fromDate.toISOString(), toDate.toISOString(), undefined, includeHidden)
                .subscribe();
            this.accountService.loadYearlyExpenses().subscribe();
        }
    }

    protected monthTotals(idx: number): { assigned: number; spent: number; pct: number } {
        const rows = this.matrixData();
        if (rows.length === 0) return { assigned: 0, spent: 0, pct: 0 };
        const assigned = rows.reduce((s, r) => s + (r.cells[idx]?.amount ?? 0), 0);
        const spent = rows.reduce((s, r) => s + (r.cells[idx]?.spent ?? 0), 0);
        return { assigned, spent, pct: assigned > 0 ? spent / assigned : spent > 0 ? 1 : 0 };
    }

    protected statusFor(assigned: number, spent: number): string {
        if (assigned === 0 && spent === 0) return '';
        if (spent > assigned) return 'over';
        const pct = assigned > 0 ? spent / assigned : 0;
        if (spent === assigned || spent > assigned * 0.98) return 'done';
        if (pct > 0.85) return 'warn';
        return 'ok';
    }

    protected fmtK(n: number): string {
        const abs = Math.abs(n);
        if (abs >= 1e6) return (n / 1e6).toFixed(1) + 'M';
        if (abs >= 1e3) return (n / 1e3).toFixed(1) + 'k';
        return n.toFixed(0);
    }

    protected hueFromId(id: string): number {
        let hash = 5381;
        for (let i = 0; i < id.length; i++) {
            hash = ((hash << 5) + hash) ^ id.charCodeAt(i);
        }
        return Math.abs(hash) % 360;
    }

    protected catColor(hue: number, l = 0.65, c = 0.14): string {
        return `oklch(${l} ${c} ${hue})`;
    }

    protected heatBg(spent: number, maxSpent: number): string {
        const pct = maxSpent > 0 ? spent / maxSpent : 0;
        if (pct <= 0) return 'transparent';
        if (pct < 0.2) return 'oklch(0.975 0.012 255)';
        if (pct < 0.4) return 'oklch(0.93 0.04 255)';
        if (pct < 0.6) return 'oklch(0.87 0.08 255)';
        if (pct < 0.8) return 'oklch(0.78 0.12 255)';
        return 'oklch(0.65 0.17 255)';
    }

    protected sparklinePath(values: number[], width = 340, height = 60): string {
        if (values.length < 2) return '';
        const maxAbs = Math.max(...values.map((v) => Math.abs(v))) || 1;
        const min = Math.min(...values, 0);
        const range = maxAbs - min || 1;
        const step = width / (values.length - 1);
        const pts = values.map((v, i) => {
            const x = i * step;
            const y = height - ((v - min) / range) * (height - 2) - 1;
            return `${x.toFixed(1)},${y.toFixed(1)}`;
        });
        return 'M' + pts.join(' L');
    }

    protected sparklineAreaPath(values: number[], width = 340, height = 60): string {
        const line = this.sparklinePath(values, width, height);
        if (!line) return '';
        return `${line} L${width},${height} L0,${height} Z`;
    }
}
