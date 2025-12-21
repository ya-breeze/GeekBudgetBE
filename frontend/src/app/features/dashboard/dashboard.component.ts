import { Component, inject, OnInit, signal, computed, effect } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DecimalPipe, JsonPipe } from '@angular/common';
import { map, forkJoin, fromEvent } from 'rxjs';
import { ApiConfiguration } from '../../core/api/api-configuration';
import { getExpenses } from '../../core/api/fn/aggregations/get-expenses';
import { Aggregation } from '../../core/api/models/aggregation';
import { Currency } from '../../core/api/models/currency';
import { AccountService } from '../accounts/services/account.service';
import { getBalances } from '../../core/api/fn/aggregations/get-balances';
import { CurrencyService } from '../currencies/services/currency.service';
import { UserService } from '../../core/services/user.service';
import { LayoutService } from '../../layout/services/layout.service';
import { AccountNoId } from '../../core/api/models/account-no-id';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { ConfirmationDialogComponent } from '../../shared/components/confirmation-dialog/confirmation-dialog.component';

interface ExpenseTableCell {
    value: number;
    color: string;
}

interface ExpenseTableRow {
    accountId: string;
    accountName: string;
    monthCells: Map<string, ExpenseTableCell>;
    total: ExpenseTableCell;
    averageSpent: number; // Added averageSpent
}

interface CurrencyTable {
    currencyId: string;
    currencyName: string;
    rows: ExpenseTableRow[];
    totalRow: ExpenseTableRow;
}

import { MatSlideToggleModule } from '@angular/material/slide-toggle';

@Component({
    selector: 'app-dashboard',
    imports: [
        MatCardModule,
        MatProgressSpinnerModule,
        MatIconModule,
        MatTableModule,
        MatButtonToggleModule,
        MatButtonModule,
        MatSnackBarModule,
        MatDialogModule,
        MatSlideToggleModule,
        DecimalPipe,
        JsonPipe,
    ],
    templateUrl: './dashboard.component.html',
    styleUrl: './dashboard.component.scss',
})
export class DashboardComponent implements OnInit {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly userService = inject(UserService);
    private readonly layoutService = inject(LayoutService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly dialog = inject(MatDialog);
    private readonly router = inject(Router);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;
    protected readonly loading = signal(true);
    protected readonly expenseData = signal<Aggregation | null>(null);
    protected readonly accounts = this.accountService.accounts;
    protected readonly selectedOutputCurrencyId = signal<string | null>(null);
    protected readonly isSmallScreen = signal(false);
    protected readonly includeHidden = signal(false);
    private readonly windowWidth = signal(typeof window !== 'undefined' ? window.innerWidth : 1024);

    constructor() {
        // Use effect to react to sidenav state changes
        effect(() => {
            const sidenavWidth = this.layoutService.sidenavOpened()
                ? this.layoutService.sidenavWidth
                : 0;
            const effectiveWidth = this.windowWidth() - sidenavWidth;
            console.log('Effective width:', effectiveWidth);
            this.isSmallScreen.set(effectiveWidth <= 1500);
        });
    }

    // Sorting state
    protected readonly sortColumn = signal<string>('accountName');
    protected readonly sortDirection = signal<'asc' | 'desc'>('asc');

    // Computed values for the expense table
    protected readonly accountColumns = computed(() => {
        return this.accounts().filter((acc) => acc.type === 'expense');
    });

    protected readonly monthColumns = computed(() => {
        const data = this.expenseData();
        const allIntervals = data?.intervals || [];

        // Show different number of months based on screen size
        if (this.isSmallScreen()) {
            // On small screens (< 768px), show only the last 6 months
            return allIntervals.slice(-6);
        } else {
            // On larger screens (>= 768px), show all 12 months
            return allIntervals;
        }
    });

    protected readonly currencyTables = computed<CurrencyTable[]>(() => {
        const data = this.expenseData();

        if (!data || !data.intervals || data.intervals.length === 0) {
            return [];
        }

        const expenseAccounts = this.accountColumns();
        const currencies = this.currencyService.currencies();
        const currenciesById = new Map<string, Currency>(
            currencies.map((currency: Currency) => [currency.id, currency]),
        );

        const visibleMonths = new Set(this.monthColumns());
        const tables: CurrencyTable[] = [];

        data.currencies.forEach((currencyAgg) => {
            const allValues: number[] = [];

            // PHASE 1: Collect all values first (month cells, row totals, month totals, grand total)
            const rowsData: {
                accountId: string;
                accountName: string;
                monthValues: Map<string, number>;
                rowTotal: number;
            }[] = [];

            // Build rows for each expense account for this currency
            expenseAccounts.forEach((account) => {
                const monthValues = new Map<string, number>();
                let rowTotal = 0;

                const accountData = currencyAgg.accounts.find(
                    (acc) => acc.accountId === account.id,
                );

                data.intervals.forEach((interval, intervalIndex) => {
                    let cellValue = 0;

                    if (accountData && accountData.amounts[intervalIndex] !== undefined) {
                        cellValue = accountData.amounts[intervalIndex];
                    }

                    monthValues.set(interval, cellValue);

                    // Only sum visible months for total
                    if (visibleMonths.has(interval)) {
                        allValues.push(cellValue);
                        rowTotal += cellValue;
                    }
                });

                if (rowTotal > 0 || Array.from(monthValues.values()).some((v) => v > 0)) {
                    allValues.push(rowTotal);

                    rowsData.push({
                        accountId: account.id!,
                        accountName: account.name,
                        monthValues,
                        rowTotal,
                    });
                }
            });

            if (!rowsData.length) {
                return;
            }

            // Calculate total row values and add to allValues
            const monthTotals = new Map<string, number>();
            let grandTotal = 0;

            data.intervals.forEach((interval) => {
                let monthTotal = 0;

                rowsData.forEach((rowData) => {
                    const cellValue = rowData.monthValues.get(interval) ?? 0;
                    monthTotal += cellValue;
                });

                monthTotals.set(interval, monthTotal);

                if (visibleMonths.has(interval)) {
                    allValues.push(monthTotal);
                    grandTotal += monthTotal;
                }
            });

            // Note: grandTotal is NOT added to allValues - it will always be white

            // PHASE 2: Now create the actual rows with colors calculated from the complete allValues array
            const rows: ExpenseTableRow[] = rowsData.map((rowData) => {
                const monthCells = new Map<string, ExpenseTableCell>();

                rowData.monthValues.forEach((value, interval) => {
                    monthCells.set(interval, {
                        value,
                        color: this.calculateColor(value, allValues),
                    });
                });

                return {
                    accountId: rowData.accountId,
                    accountName: rowData.accountName,
                    monthCells,
                    total: {
                        value: rowData.rowTotal,
                        color: this.calculateColor(rowData.rowTotal, allValues),
                    },
                    averageSpent:
                        this.accountService
                            .averages()
                            .find((a) => a.accountId === rowData.accountId)?.averageSpent ?? 0,
                };
            });

            const sortedRows = this.sortRows(rows);

            // Build total row with colors
            const totalRowMonthCells = new Map<string, ExpenseTableCell>();

            monthTotals.forEach((value, interval) => {
                totalRowMonthCells.set(interval, {
                    value,
                    color: this.calculateColor(value, allValues),
                });
            });

            const totalRow: ExpenseTableRow = {
                accountId: 'total',
                accountName: 'Total',
                monthCells: totalRowMonthCells,
                total: {
                    value: grandTotal,
                    color: 'rgb(255, 255, 255)', // Grand total is always white
                },
                averageSpent: 0,
            };

            const currencyMeta = currenciesById.get(currencyAgg.currencyId);
            const currencyName = currencyMeta?.name ?? currencyAgg.currencyId;

            tables.push({
                currencyId: currencyAgg.currencyId,
                currencyName,
                rows: sortedRows,
                totalRow,
            });
        });

        return tables;
    });

    protected readonly assetData = signal<Aggregation | null>(null);

    protected readonly assetCards = computed(() => {
        const data = this.assetData();
        const accounts = this.accounts();
        if (!data || !accounts.length) {
            return [];
        }

        const assetAccounts = accounts.filter(
            (acc) => acc.type === 'asset' && acc.showInDashboardSummary !== false,
        );
        if (!assetAccounts.length) {
            return [];
        }

        // Map currency ID to symbol/name if needed (or just use code)
        const cards: any[] = [];

        // Group data by currency, but we need to extract account data
        data.currencies.forEach((currencyAgg) => {
            currencyAgg.accounts.forEach((accAgg) => {
                const account = assetAccounts.find((a) => a.id === accAgg.accountId);
                if (!account) return;

                // Calculate balances
                // Amounts are net changes per interval. Sum them up for total balance.
                const totalBalance = accAgg.amounts.reduce((sum, val) => sum + val, 0);

                // Calculate last month balance (total minus last interval)
                // Check if we have at least 1 interval
                let trendPercent = 0;
                let trendDirection: 'up' | 'down' | 'neutral' = 'neutral';

                if (accAgg.amounts.length > 0) {
                    // Assuming the last interval is the current partially complete month or just the last month
                    const lastAmount = accAgg.amounts[accAgg.amounts.length - 1];
                    const previousBalance = totalBalance - lastAmount;

                    if (previousBalance !== 0) {
                        trendPercent = (lastAmount / Math.abs(previousBalance)) * 100;
                    } else if (lastAmount !== 0) {
                        trendPercent = 100; // From 0 to something
                    }

                    if (trendPercent > 0) trendDirection = 'up';
                    else if (trendPercent < 0) trendDirection = 'down';
                }

                const currency = this.currencyService
                    .currencies()
                    .find((c) => c.id === currencyAgg.currencyId);

                cards.push({
                    accountId: account.id,
                    accountName: account.name,
                    balance: totalBalance,
                    currencyId: currencyAgg.currencyId,
                    currencyName: currency?.name || currencyAgg.currencyId,
                    trendPercent: Math.abs(trendPercent),
                    trendDirection,
                });
            });
        });

        return cards;
    });

    ngOnInit(): void {
        this.currencyService.loadCurrencies().subscribe();

        // Listen to window resize events
        fromEvent(window, 'resize').subscribe(() => {
            this.windowWidth.set(window.innerWidth);
        });

        // Load user data and set default currency
        this.userService.loadUser().subscribe({
            next: (user) => {
                if (user.favoriteCurrencyId && !this.selectedOutputCurrencyId()) {
                    this.selectedOutputCurrencyId.set(user.favoriteCurrencyId);
                }
                this.loadDashboardData();
            },
            error: () => {
                // If user loading fails, still load dashboard data
                this.loadDashboardData();
            },
        });
    }

    protected onHideAccount(accountId: string): void {
        const account = this.accounts().find((a) => a.id === accountId);
        if (!account || !account.id) return;

        this.dialog
            .open(ConfirmationDialogComponent, {
                data: {
                    title: 'Hide Account',
                    message: `Are you sure you want to hide "${account.name}" from the dashboard summary?`,
                    confirmText: 'Hide',
                    cancelText: 'Cancel',
                },
            })
            .afterClosed()
            .subscribe((result) => {
                if (result) {
                    const update: AccountNoId = {
                        name: account.name!,
                        type: account.type,
                        description: account.description,
                        bankInfo: account.bankInfo,
                        showInDashboardSummary: false,
                    };

                    this.accountService.update(account.id!, update).subscribe({
                        next: () => {
                            this.snackBar
                                .open('Account hidden from dashboard', 'Undo', { duration: 3000 })
                                .onAction()
                                .subscribe(() => {
                                    // Undo action
                                    this.accountService
                                        .update(account.id!, {
                                            ...update,
                                            showInDashboardSummary: true,
                                        })
                                        .subscribe();
                                });
                        },
                        error: () => {
                            this.snackBar.open('Failed to hide account', 'Close', {
                                duration: 3000,
                            });
                        },
                    });
                }
            });
    }

    protected onOutputCurrencyToggle(currencyId: string): void {
        const current = this.selectedOutputCurrencyId();
        const next = current === currencyId ? null : currencyId;
        this.selectedOutputCurrencyId.set(next);
        this.loadDashboardData();
    }

    protected onIncludeHiddenToggle(): void {
        this.includeHidden.update((current) => !current);
        this.loadDashboardData();
    }

    private loadDashboardData(): void {
        this.loading.set(true);

        const now = new Date();
        // Get data for the last 12 months for expenses
        const twelveMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 11, 1);
        const outputCurrencyId = this.selectedOutputCurrencyId();
        const includeHidden = this.includeHidden();

        const expenseParams: {
            from: string;
            to: string;
            outputCurrencyId?: string;
            includeHidden?: boolean;
        } = {
            from: twelveMonthsAgo.toISOString(),
            to: now.toISOString(),
            includeHidden: includeHidden,
        };

        const balanceParams: { to: string; outputCurrencyId?: string; includeHidden?: boolean } = {
            to: now.toISOString(),
            includeHidden: includeHidden,
        };

        if (outputCurrencyId) {
            expenseParams.outputCurrencyId = outputCurrencyId;
            balanceParams.outputCurrencyId = outputCurrencyId;
        }

        forkJoin({
            accounts: this.accountService.loadAccounts(),
            expenseData: getExpenses(this.http, this.apiConfig.rootUrl, expenseParams).pipe(
                map((response) => response.body),
            ),
            assetData: getBalances(this.http, this.apiConfig.rootUrl, balanceParams).pipe(
                map((response) => response.body),
            ),
            averages: this.accountService.loadYearlyExpenses(outputCurrencyId ?? undefined), // Load averages
        }).subscribe({
            next: ({ expenseData, assetData }) => {
                console.log('Dashboard data loaded:', { expenseData, assetData });
                this.expenseData.set(expenseData);
                this.assetData.set(assetData);
                this.loading.set(false);
            },
            error: (error) => {
                console.error('Error loading dashboard data:', error);
                this.loading.set(false);
            },
        });
    }

    private calculateColor(value: number, allValues: number[]): string {
        if (allValues.length === 0 || value <= 0) {
            return 'rgb(255, 255, 255)';
        }

        const min = Math.min(...allValues);

        // Use 99th percentile instead of max to avoid outliers skewing the color scale
        const sortedValues = [...allValues].sort((a, b) => a - b);
        const percentile99Index = Math.floor(sortedValues.length * 0.99);
        const max = sortedValues[percentile99Index];

        if (min === max) {
            return 'rgb(255, 255, 200)';
        }

        // Normalize value between 0 and 1
        const normalized = (value - min) / (max - min);

        // Create gradient from green (low) to red (high)
        // Green: rgb(200, 255, 200)
        // Yellow: rgb(255, 255, 200)
        // Red: rgb(255, 200, 200)

        let r: number, g: number, b: number, a: number;

        if (normalized < 0.5) {
            // Green to Yellow
            const t = normalized * 2;
            r = Math.round(200 + 55 * t);
            g = 255;
            b = 200;
            a = 0.3;
        } else {
            // Yellow to Red
            const t = (normalized - 0.5) * 2;
            r = 255;
            g = Math.round(255 - 55 * t);
            b = 200;
            a = 0.9;
        }

        return `rgb(${r}, ${g}, ${b}, ${a})`;
    }

    protected formatMonth(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
    }

    protected onColumnClick(column: string): void {
        if (this.sortColumn() === column) {
            // Toggle direction if clicking the same column
            this.sortDirection.set(this.sortDirection() === 'asc' ? 'desc' : 'asc');
        } else {
            // Set new column and default to ascending
            this.sortColumn.set(column);
            this.sortDirection.set('asc');
        }
    }

    private sortRows(rows: ExpenseTableRow[]): ExpenseTableRow[] {
        const column = this.sortColumn();
        const direction = this.sortDirection();
        const factor = direction === 'asc' ? 1 : -1;

        return [...rows].sort((a, b) => {
            let valueA: number | string;
            let valueB: number | string;

            if (column === 'accountName') {
                valueA = this.removeLeadingEmoji(a.accountName);
                valueB = this.removeLeadingEmoji(b.accountName);
            } else if (column === 'total') {
                valueA = a.total.value;
                valueB = b.total.value;
            } else {
                // Sorting by a specific month
                valueA = a.monthCells.get(column)?.value ?? 0;
                valueB = b.monthCells.get(column)?.value ?? 0;
            }

            if (typeof valueA === 'string' && typeof valueB === 'string') {
                return valueA.localeCompare(valueB) * factor;
            }

            return ((valueA as number) - (valueB as number)) * factor;
        });
    }

    private removeLeadingEmoji(text: string): string {
        // Remove leading emoji characters for sorting purposes
        // This regex matches emoji at the start of the string and removes them
        return text
            .replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '')
            .trim();
    }

    protected onCellClick(accountId: string, monthDateString: string): void {
        const d = new Date(monthDateString);
        this.router.navigate(['/transactions'], {
            queryParams: {
                accountId,
                month: d.getMonth(),
                year: d.getFullYear(),
            },
        });
    }
}
