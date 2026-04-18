import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DecimalPipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { map, forkJoin } from 'rxjs';
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
import { AccountFormDialogComponent } from '../accounts/account-form-dialog/account-form-dialog.component';
import { AccountDisplayComponent } from '../../shared/components/account-display/account-display.component';
import { ReconciliationService } from '../reconciliation/services/reconciliation.service';
import { ReconciliationStatus } from '../../core/api/models/reconciliation-status';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AssetCard, AssetTotal } from './models/dashboard.models';
import { BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration, ChartOptions } from 'chart.js';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

interface ExpenseTableCell {
    value: number;
    heatClass: string | null;
}

interface ExpenseTableRow {
    accountId: string;
    accountName: string;
    monthCells: Map<string, ExpenseTableCell>;
    total: number;
    averageSpent: number;
    accountImage?: string;
    sparkValues: number[];
}

interface CurrencyTable {
    currencyId: string;
    currencyName: string;
    rows: ExpenseTableRow[];
    totalRow: ExpenseTableRow;
}

interface KpiData {
    curTotal: number;
    prevTotal: number | null;
    avgPerMonth: number;
    delta: number | null;
    topName: string;
    topTotal: number;
    topPct: number;
    curMonthLabel: string;
    prevMonthLabel: string;
}

interface RankedCategory {
    accountId: string;
    accountName: string;
    total: number;
    barPct: number;
    color: string;
    rank: number;
}

type RangeId = '3m' | '6m' | '12m' | 'ytd';

// Golden-angle hue stepping: each consecutive color is ~137.5° away, maximising distance
const HUES = [0, 138, 275, 53, 190, 328, 105, 243, 20, 158, 295, 73, 210, 348, 125];

@Component({
    selector: 'app-dashboard',
    imports: [
        MatProgressSpinnerModule,
        MatIconModule,
        MatButtonModule,
        MatSnackBarModule,
        MatDialogModule,
        MatSlideToggleModule,
        MatTooltipModule,
        DecimalPipe,
        RouterLink,
        AccountDisplayComponent,
        BaseChartDirective,
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
    private readonly reconciliationService = inject(ReconciliationService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;
    protected readonly loading = signal(true);
    protected readonly expenseData = signal<Aggregation | null>(null);
    protected readonly accounts = this.accountService.accounts;
    protected readonly selectedOutputCurrencyId = signal<string | null>(null);
    protected readonly includeHidden = signal(false);

    // Range selector
    protected readonly selectedRange = signal<RangeId>('6m');
    protected readonly rangeOptions: { id: RangeId; label: string }[] = [
        { id: '3m', label: '3M' },
        { id: '6m', label: '6M' },
        { id: '12m', label: '12M' },
        { id: 'ytd', label: 'YTD' },
    ];

    // Expansion state
    protected readonly expandedRowId = signal<string | null>(null);

    // Search
    protected readonly searchQuery = signal('');

    // Sort state
    protected readonly sortColumn = signal<string>('total');
    protected readonly sortDirection = signal<'asc' | 'desc'>('desc');

    // Asset stuff
    protected readonly assetData = signal<Aggregation | null>(null);
    protected readonly reconciliationStatuses = signal<ReconciliationStatus[]>([]);
    protected readonly showAssetDetails = signal(false);

    // ---- Computed: expense accounts ----
    protected readonly accountColumns = computed(() =>
        this.accounts().filter((acc) => acc.type === 'expense'),
    );

    // ---- Computed: color map ----
    private readonly accountColorMap = computed(() => {
        const map = new Map<string, string>();
        this.accountColumns().forEach((acc, i) => {
            map.set(acc.id!, `oklch(0.62 0.14 ${HUES[i % HUES.length]})`);
        });
        return map;
    });

    // ---- Computed: visible months based on range ----
    protected readonly monthColumns = computed((): string[] => {
        const data = this.expenseData();
        if (!data?.intervals?.length) return [];
        const all = data.intervals;
        const range = this.selectedRange();
        if (range === 'ytd') {
            const year = new Date().getFullYear();
            return all.filter((m) => new Date(m).getFullYear() === year);
        }
        const count = range === '3m' ? 3 : range === '6m' ? 6 : 12;
        return all.slice(-count);
    });

    // ---- KPIs ----
    protected readonly kpis = computed((): KpiData | null => {
        const data = this.expenseData();
        const months = this.monthColumns();
        if (!data || months.length < 1) return null;
        const cur = data.currencies[0];
        if (!cur) return null;

        const allIntervals = data.intervals;
        const idxs = months.map((m) => allIntervals.indexOf(m)).filter((i) => i >= 0);
        if (!idxs.length) return null;

        const totalsPerMonth = idxs.map((mi) =>
            cur.accounts.reduce((s, a) => s + Math.max(0, a.amounts[mi] ?? 0), 0),
        );
        const curTotal = totalsPerMonth[totalsPerMonth.length - 1];
        const prevTotal = totalsPerMonth.length >= 2 ? totalsPerMonth[totalsPerMonth.length - 2] : null;
        const avgPerMonth = totalsPerMonth.reduce((a, b) => a + b, 0) / totalsPerMonth.length;
        const delta = prevTotal ? ((curTotal - prevTotal) / prevTotal) * 100 : null;

        const catTotals = cur.accounts.map((a) => ({
            id: a.accountId,
            name: this.accounts().find((acc) => acc.id === a.accountId)?.name ?? a.accountId,
            total: idxs.reduce((s, mi) => s + Math.max(0, a.amounts[mi] ?? 0), 0),
        }));
        catTotals.sort((a, b) => b.total - a.total);
        const top = catTotals[0];
        const totalAll = catTotals.reduce((s, a) => s + a.total, 0);

        return {
            curTotal,
            prevTotal,
            avgPerMonth,
            delta,
            topName: top?.name ?? '',
            topTotal: top?.total ?? 0,
            topPct: totalAll > 0 ? ((top?.total ?? 0) / totalAll) * 100 : 0,
            curMonthLabel: this.formatMonth(months[months.length - 1]),
            prevMonthLabel: months.length >= 2 ? this.formatMonth(months[months.length - 2]) : '',
        };
    });

    // ---- Stacked bar chart ----
    protected readonly stackedBarChartData = computed((): ChartConfiguration<'bar'>['data'] => {
        const data = this.expenseData();
        const months = this.monthColumns();
        if (!data || !months.length) return { labels: [], datasets: [] };

        const cur = data.currencies[0];
        if (!cur) return { labels: [], datasets: [] };

        const allIntervals = data.intervals;
        const idxs = months.map((m) => allIntervals.indexOf(m));
        const colorMap = this.accountColorMap();

        const datasets = this.accountColumns().map((acc) => {
            const color = colorMap.get(acc.id!) ?? '#888';
            const accAgg = cur.accounts.find((a) => a.accountId === acc.id);
            return {
                label: acc.name,
                data: idxs.map((mi) => (mi >= 0 && accAgg ? Math.max(0, accAgg.amounts[mi] ?? 0) : 0)),
                backgroundColor: color,
                borderRadius: 2,
                borderSkipped: false as const,
            };
        });

        return {
            labels: months.map((m) => this.formatMonthShort(m)),
            datasets,
        };
    });

    protected readonly stackedBarOptions: ChartOptions<'bar'> = {
        responsive: true,
        maintainAspectRatio: false,
        animation: { duration: 400 },
        scales: {
            x: {
                stacked: true,
                grid: { display: false },
                ticks: { font: { family: 'monospace', size: 10 }, color: '#999' },
            },
            y: {
                stacked: true,
                grid: { color: 'rgba(0,0,0,0.05)' },
                ticks: {
                    font: { family: 'monospace', size: 10 },
                    color: '#999',
                    callback: (v) => {
                        const n = +v;
                        if (n >= 1e6) return (n / 1e6).toFixed(1) + 'M';
                        if (n >= 1e3) return (n / 1e3).toFixed(1) + 'k';
                        return n.toFixed(0);
                    },
                },
            },
        },
        plugins: {
            legend: { display: false },
            tooltip: {
                mode: 'index',
                itemSort: (a, b) => (b.raw as number) - (a.raw as number),
                callbacks: {
                    label: (ctx) => {
                        const v = ctx.raw as number;
                        return v > 0 ? `${ctx.dataset.label}: ${this.fmtK(v)}` : '';
                    },
                },
                filter: (item) => (item.raw as number) > 0,
            },
        },
    };

    // ---- Ranked categories ----
    protected readonly rankedCategories = computed((): RankedCategory[] => {
        const data = this.expenseData();
        const months = this.monthColumns();
        if (!data || !months.length) return [];
        const cur = data.currencies[0];
        if (!cur) return [];

        const allIntervals = data.intervals;
        const idxs = months.map((m) => allIntervals.indexOf(m)).filter((i) => i >= 0);
        const colorMap = this.accountColorMap();

        const cats = this.accountColumns().map((acc) => {
            const accAgg = cur.accounts.find((a) => a.accountId === acc.id);
            const total = idxs.reduce((s, mi) => s + Math.max(0, accAgg?.amounts[mi] ?? 0), 0);
            return { accountId: acc.id!, accountName: acc.name, total, color: colorMap.get(acc.id!) ?? '#888' };
        });
        cats.sort((a, b) => b.total - a.total);
        const max = cats[0]?.total || 1;
        return cats
            .filter((c) => c.total > 0)
            .slice(0, 10)
            .map((c, i) => ({ ...c, barPct: Math.max(2, (c.total / max) * 100), rank: i + 1 }));
    });

    // ---- Expense table ----
    protected readonly currencyTables = computed<CurrencyTable[]>(() => {
        const data = this.expenseData();
        if (!data?.intervals?.length) return [];

        const expenseAccounts = this.accountColumns();
        const currencies = this.currencyService.currencies();
        const currenciesById = new Map<string, Currency>(currencies.map((c: Currency) => [c.id, c]));

        const visibleMonths = new Set(this.monthColumns());
        const colorMap = this.accountColorMap();
        const allIntervals = data.intervals;
        const tables: CurrencyTable[] = [];

        data.currencies.forEach((currencyAgg) => {
            const rowsData: {
                accountId: string;
                accountName: string;
                accountImage?: string;
                monthValues: Map<string, number>;
                rowTotal: number;
                sparkValues: number[];
            }[] = [];

            expenseAccounts.forEach((account) => {
                const accountData = currencyAgg.accounts.find((a) => a.accountId === account.id);
                const monthValues = new Map<string, number>();
                let rowTotal = 0;
                const sparkValues: number[] = [];

                allIntervals.forEach((interval, idx) => {
                    const v = accountData?.amounts[idx] ?? 0;
                    monthValues.set(interval, v);
                    if (visibleMonths.has(interval)) rowTotal += v;
                    sparkValues.push(v);
                });

                if (rowTotal > 0 || Array.from(monthValues.values()).some((v) => v > 0)) {
                    rowsData.push({
                        accountId: account.id!,
                        accountName: account.name,
                        accountImage: account.image,
                        monthValues,
                        rowTotal,
                        sparkValues,
                    });
                }
            });

            if (!rowsData.length) return;

            const monthTotals = new Map<string, number>();
            let grandTotal = 0;
            allIntervals.forEach((interval) => {
                const t = rowsData.reduce((s, r) => s + (r.monthValues.get(interval) ?? 0), 0);
                monthTotals.set(interval, t);
                if (visibleMonths.has(interval)) grandTotal += t;
            });

            const rows: ExpenseTableRow[] = rowsData.map((rowData) => {
                const vals = Array.from(visibleMonths).map((m) => rowData.monthValues.get(m) ?? 0);
                const rowMax = Math.max(...vals.map((v) => Math.max(0, v)));
                const monthCells = new Map<string, ExpenseTableCell>();
                rowData.monthValues.forEach((value, interval) => {
                    monthCells.set(interval, { value, heatClass: this.heatClass(value, rowMax) });
                });
                const nonZero = vals.filter((v) => v > 0);
                const avg = nonZero.length ? nonZero.reduce((a, b) => a + b, 0) / nonZero.length : 0;
                return {
                    accountId: rowData.accountId,
                    accountName: rowData.accountName,
                    monthCells,
                    total: rowData.rowTotal,
                    averageSpent:
                        this.accountService.averages().find((a) => a.accountId === rowData.accountId)
                            ?.averageSpent ?? avg,
                    accountImage: rowData.accountImage,
                    sparkValues: rowData.sparkValues,
                };
            });

            const sorted = this.sortRows(rows);

            const totalRowCells = new Map<string, ExpenseTableCell>();
            const totVals = Array.from(visibleMonths).map((m) => monthTotals.get(m) ?? 0);
            const totMax = Math.max(...totVals.map((v) => Math.max(0, v)));
            monthTotals.forEach((value, interval) => {
                totalRowCells.set(interval, { value, heatClass: this.heatClass(value, totMax) });
            });

            const totalRow: ExpenseTableRow = {
                accountId: 'total',
                accountName: 'Total',
                monthCells: totalRowCells,
                total: grandTotal,
                averageSpent: 0,
                sparkValues: allIntervals.map((m) => monthTotals.get(m) ?? 0),
            };

            const currencyMeta = currenciesById.get(currencyAgg.currencyId);
            tables.push({
                currencyId: currencyAgg.currencyId,
                currencyName: currencyMeta?.name ?? currencyAgg.currencyId,
                rows: sorted,
                totalRow,
            });
        });

        return tables;
    });

    // ---- Filtered table rows ----
    protected filteredRows(rows: ExpenseTableRow[]): ExpenseTableRow[] {
        const q = this.searchQuery().trim().toLowerCase();
        if (!q) return rows;
        return rows.filter((r) => r.accountName.toLowerCase().includes(q));
    }

    // ---- Drill-down data ----
    protected getDrillSparklineData(row: ExpenseTableRow): ChartConfiguration<'line'>['data'] {
        const months = this.monthColumns();
        const allIntervals = this.expenseData()?.intervals ?? [];
        const idxs = months.map((m) => allIntervals.indexOf(m)).filter((i) => i >= 0);
        const values = idxs.map((i) => row.sparkValues[i] ?? 0);
        const color = this.accountColorMap().get(row.accountId) ?? '#888';
        return {
            labels: months.map((m) => this.formatMonthShort(m)),
            datasets: [
                {
                    data: values,
                    borderColor: color,
                    backgroundColor: color + '1a',
                    fill: true,
                    tension: 0.3,
                    borderWidth: 2,
                    pointRadius: 3,
                    pointHoverRadius: 5,
                },
            ],
        };
    }

    protected readonly drillSparklineOptions: ChartOptions<'line'> = {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
            x: {
                display: true,
                grid: { display: false },
                ticks: { font: { family: 'monospace', size: 10 }, color: '#999' },
            },
            y: {
                display: true,
                grid: { color: 'rgba(0,0,0,0.05)' },
                ticks: {
                    font: { family: 'monospace', size: 10 },
                    color: '#999',
                    callback: (v) => {
                        const n = +v;
                        if (n >= 1e6) return (n / 1e6).toFixed(1) + 'M';
                        if (n >= 1e3) return (n / 1e3).toFixed(1) + 'k';
                        return n.toFixed(0);
                    },
                },
            },
        },
        plugins: {
            legend: { display: false },
            tooltip: { mode: 'index', intersect: false },
        },
    };

    protected getDrillStats(row: ExpenseTableRow) {
        const months = this.monthColumns();
        const allIntervals = this.expenseData()?.intervals ?? [];
        const idxs = months.map((m) => allIntervals.indexOf(m)).filter((i) => i >= 0);
        const values = idxs.map((i) => row.sparkValues[i] ?? 0);
        const positive = values.filter((v) => v > 0);
        const total = values.reduce((a, b) => a + b, 0);
        const avg = positive.length ? positive.reduce((a, b) => a + b, 0) / positive.length : 0;
        const peak = Math.max(0, ...values);
        const peakIdx = values.indexOf(peak);
        const peakMonth = peakIdx >= 0 ? this.formatMonthShort(months[peakIdx]) : '';
        const mean = avg;
        const stdev = Math.sqrt(
            positive.reduce((s, v) => s + Math.pow(v - mean, 2), 0) / (positive.length || 1),
        );
        return { total, avg, peak, peakMonth, stdev };
    }

    // ---- Row sparkline (inline in category cell) ----
    protected readonly inlineSparklineOptions: ChartOptions<'line'> = {
        responsive: true,
        maintainAspectRatio: false,
        elements: { point: { radius: 0 }, line: { tension: 0.3, borderWidth: 1.5 } },
        scales: { x: { display: false }, y: { display: false } },
        plugins: { legend: { display: false }, tooltip: { enabled: false } },
    };

    protected getInlineSparklineData(row: ExpenseTableRow): ChartConfiguration<'line'>['data'] {
        const months = this.monthColumns();
        const allIntervals = this.expenseData()?.intervals ?? [];
        const idxs = months.map((m) => allIntervals.indexOf(m)).filter((i) => i >= 0);
        const values = idxs.map((i) => row.sparkValues[i] ?? 0);
        const color = this.accountColorMap().get(row.accountId) ?? '#888';
        return {
            labels: values.map((_, i) => i.toString()),
            datasets: [{ data: values, borderColor: color, fill: false }],
        };
    }

    // ---- Assets ----
    protected readonly assetTotals = computed<AssetTotal[]>(() => {
        const data = this.assetData();
        const accounts = this.accounts();
        const includeHidden = this.includeHidden();
        if (!data || !accounts.length) return [];
        const assetAccounts = accounts.filter(
            (acc) => acc.type === 'asset' && (includeHidden || acc.showInDashboardSummary !== false),
        );
        const assetAccountIds = new Set(assetAccounts.map((a) => a.id));
        return data.currencies
            .map((currencyAgg) => {
                let totalBalance = 0;
                let prevTotalBalance = 0;
                const history: number[] = new Array(data.intervals.length).fill(0);
                currencyAgg.accounts.forEach((accAgg) => {
                    if (!assetAccountIds.has(accAgg.accountId)) return;
                    const bal = accAgg.total || 0;
                    totalBalance += bal;
                    prevTotalBalance += bal / (1 + (accAgg.changePercent || 0) / 100);
                    accAgg.amounts.forEach((amt, idx) => { history[idx] += amt; });
                });
                let trendPercent = 0;
                let trendDirection: 'up' | 'down' | 'neutral' = 'neutral';
                if (prevTotalBalance > 0) {
                    trendPercent = ((totalBalance - prevTotalBalance) / prevTotalBalance) * 100;
                    if (trendPercent > 0.01) trendDirection = 'up';
                    else if (trendPercent < -0.01) trendDirection = 'down';
                }
                const currency = this.currencyService.currencies().find((c) => c.id === currencyAgg.currencyId);
                return { currencyId: currencyAgg.currencyId, currencyName: currency?.name || currencyAgg.currencyId, totalBalance, trendPercent: Math.abs(trendPercent), trendDirection, history };
            })
            .filter((t) => t.totalBalance !== 0 || t.history.some((v) => v !== 0));
    });

    protected readonly assetCards = computed<AssetCard[]>(() => {
        const data = this.assetData();
        const accounts = this.accounts();
        const statuses = this.reconciliationStatuses() || [];
        if (!data || !accounts.length) return [];
        const includeHidden = this.includeHidden();
        const assetAccounts = accounts.filter(
            (acc) => acc.type === 'asset' && (includeHidden || acc.showInDashboardSummary !== false),
        );
        if (!assetAccounts.length) return [];
        const statusByAccountId = new Map<string, ReconciliationStatus>(statuses.map((s) => [s.accountId, s]));
        const cards: AssetCard[] = [];
        data.currencies.forEach((currencyAgg) => {
            currencyAgg.accounts.forEach((accAgg) => {
                const account = assetAccounts.find((a) => a.id === accAgg.accountId);
                if (!account) return;
                const totalBalance = accAgg.total || 0;
                const trendPercent = accAgg.changePercent || 0;
                let trendDirection: 'up' | 'down' | 'neutral' = 'neutral';
                if (trendPercent > 0.01) trendDirection = 'up';
                else if (trendPercent < -0.01) trendDirection = 'down';
                const currency = this.currencyService.currencies().find((c) => c.id === currencyAgg.currencyId);
                const reconciliationStatus = statusByAccountId.get(account.id!);
                const { state: reconciliationState, tooltip: reconciliationTooltip } = this.getReconciliationState(reconciliationStatus);
                cards.push({ accountId: account.id, accountName: account.name, balance: totalBalance, currencyId: currencyAgg.currencyId, currencyName: currency?.name || currencyAgg.currencyId, trendPercent: Math.abs(trendPercent), trendDirection, accountImage: account.image, reconciliationState, reconciliationTooltip });
            });
        });
        return cards;
    });

    protected readonly sparklineOptions: ChartOptions<'line'> = {
        responsive: true, maintainAspectRatio: false,
        elements: { point: { radius: 0 }, line: { tension: 0.3, borderWidth: 2 } },
        scales: { x: { display: false }, y: { display: false } },
        plugins: { legend: { display: false }, tooltip: { enabled: false } },
    };

    protected getSparklineData(history: number[]): ChartConfiguration<'line'>['data'] {
        return {
            labels: history.map((_, i) => i.toString()),
            datasets: [{ data: history, borderColor: '#1967d2', fill: false }],
        };
    }

    // ---- Lifecycle ----
    ngOnInit(): void {
        this.currencyService.loadCurrencies().subscribe();
        this.userService.loadUser().subscribe({
            next: (user) => {
                if (user.favoriteCurrencyId && !this.selectedOutputCurrencyId()) {
                    this.selectedOutputCurrencyId.set(user.favoriteCurrencyId);
                }
                this.loadDashboardData();
            },
            error: () => this.loadDashboardData(),
        });
    }

    // ---- Actions ----
    protected toggleAssetDetails(): void {
        this.showAssetDetails.update((v) => !v);
    }

    protected onToggleExpand(id: string): void {
        this.expandedRowId.update((cur) => (cur === id ? null : id));
    }

    protected onHideAccount(accountId: string): void {
        const account = this.accounts().find((a) => a.id === accountId);
        if (!account?.id) return;
        this.dialog.open(ConfirmationDialogComponent, {
            data: { title: 'Hide Account', message: `Hide "${account.name}" from the dashboard?`, confirmText: 'Hide', cancelText: 'Cancel' },
        }).afterClosed().subscribe((result) => {
            if (!result) return;
            const update: AccountNoId = { name: account.name!, type: account.type, description: account.description, bankInfo: account.bankInfo, showInDashboardSummary: false };
            this.accountService.update(account.id!, update).subscribe({
                next: () => this.snackBar.open('Account hidden', 'Undo', { duration: 3000 }).onAction().subscribe(() => this.accountService.update(account.id!, { ...update, showInDashboardSummary: true }).subscribe()),
                error: () => this.snackBar.open('Failed to hide account', 'Close', { duration: 3000 }),
            });
        });
    }

    protected onOutputCurrencyToggle(currencyId: string): void {
        const current = this.selectedOutputCurrencyId();
        this.selectedOutputCurrencyId.set(current === currencyId ? null : currencyId);
        this.loadDashboardData();
    }

    protected onIncludeHiddenToggle(): void {
        this.includeHidden.update((c) => !c);
        this.loadDashboardData();
    }

    protected onSettingsClick(event: Event, accountId: string): void {
        event.stopPropagation();
        const account = this.accounts().find((a) => a.id === accountId);
        if (!account) return;
        this.dialog.open(AccountFormDialogComponent, { width: '600px', data: { mode: 'edit', account } })
            .afterClosed().subscribe((result) => {
                if (result && account.id) {
                    this.accountService.handleAccountDialogResult(account, result, this.snackBar);
                }
            });
    }

    protected onBalanceClick(accountId: string): void {
        const now = new Date();
        const oneYearAgo = new Date(now.getFullYear(), now.getMonth() - 11, now.getDate());
        this.router.navigate(['/reports/balance'], { queryParams: { accountId, from: oneYearAgo.toISOString(), to: now.toISOString() } });
    }

    protected onColumnClick(column: string): void {
        if (this.sortColumn() === column) {
            this.sortDirection.set(this.sortDirection() === 'asc' ? 'desc' : 'asc');
        } else {
            this.sortColumn.set(column);
            this.sortDirection.set('desc');
        }
    }

    protected onCellClick(accountId: string, monthDateString: string): void {
        const d = new Date(monthDateString);
        this.router.navigate(['/transactions'], { queryParams: { accountId, month: d.getMonth(), year: d.getFullYear() } });
    }

    protected navigateToTransactions(accountId: string): void {
        const months = this.monthColumns();
        if (!months.length) { this.router.navigate(['/transactions'], { queryParams: { accountId } }); return; }
        const from = new Date(months[0]);
        const to = new Date(months[months.length - 1]);
        to.setMonth(to.getMonth() + 1);
        this.router.navigate(['/transactions'], { queryParams: { accountId, from: from.toISOString(), to: to.toISOString() } });
    }

    protected getAccountColor(accountId: string): string {
        return this.accountColorMap().get(accountId) ?? '#888';
    }

    protected getAccountGlyph(accountName: string): string {
        return accountName.replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '').trim().slice(0, 2).toUpperCase();
    }

    // ---- Helpers ----
    protected formatMonth(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
    }

    protected formatMonthShort(dateStr: string): string {
        const d = new Date(dateStr);
        return d.toLocaleDateString('en-US', { month: 'short' }) + " '" + String(d.getFullYear()).slice(2);
    }

    protected fmtK(n: number): string {
        const abs = Math.abs(n);
        if (abs >= 1e6) return (n / 1e6).toFixed(1) + 'M';
        if (abs >= 1e3) return (n / 1e3).toFixed(1) + 'k';
        return n.toFixed(0);
    }

    protected fmt(n: number, dp = 0): string {
        if (n === 0) return '0';
        return n.toLocaleString('en-US', { minimumFractionDigits: dp, maximumFractionDigits: dp });
    }

    // ---- Private ----
    private heatClass(value: number, rowMax: number): string | null {
        if (value <= 0) return null;
        const t = Math.min(1, value / (rowMax || 1));
        if (t < 0.15) return 'heat-1';
        if (t < 0.35) return 'heat-2';
        if (t < 0.60) return 'heat-3';
        if (t < 0.85) return 'heat-4';
        return 'heat-5';
    }

    private readonly reconciliationTolerance = 0.01;

    private getReconciliationState(status: ReconciliationStatus | undefined): { state: 'reconciled' | 'warning' | 'unreconciled' | 'none'; tooltip: string } {
        if (!status || (status.lastReconciledAt == null && status.bankBalance == null)) return { state: 'none', tooltip: '' };
        const delta = Math.abs(status.delta ?? 0);
        const isBalanced = delta <= this.reconciliationTolerance;
        if (status.hasUnprocessedTransactions) return { state: 'warning', tooltip: 'Has unprocessed transactions' };
        if (!isBalanced) return { state: 'unreconciled', tooltip: `Balance mismatch: ${(status.delta ?? 0) > 0 ? '+' : ''}${status.delta?.toFixed(2)} ${status.currencySymbol ?? ''}` };
        if (status.hasTransactionsAfterBankBalance) return { state: 'warning', tooltip: 'Reconciled, but newer transactions exist' };
        const date = status.lastReconciledAt ? new Date(status.lastReconciledAt).toLocaleDateString() : '';
        return { state: 'reconciled', tooltip: date ? `Reconciled on ${date}` : 'Reconciled' };
    }

    private loadDashboardData(): void {
        this.loading.set(true);
        const now = new Date();
        const twelveMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 11, 1);
        const outputCurrencyId = this.selectedOutputCurrencyId();
        const includeHidden = this.includeHidden();
        const params: { from: string; to: string; outputCurrencyId?: string; includeHidden?: boolean } = {
            from: twelveMonthsAgo.toISOString(),
            to: now.toISOString(),
            includeHidden,
        };
        if (outputCurrencyId) params.outputCurrencyId = outputCurrencyId;

        forkJoin({
            accounts: this.accountService.loadAccounts(),
            expenseData: getExpenses(this.http, this.apiConfig.rootUrl, params).pipe(map((r) => r.body)),
            assetData: getBalances(this.http, this.apiConfig.rootUrl, params).pipe(map((r) => r.body)),
            averages: this.accountService.loadYearlyExpenses(outputCurrencyId ?? undefined),
            reconciliationStatuses: this.reconciliationService.loadStatuses(),
        }).subscribe({
            next: ({ expenseData, assetData, reconciliationStatuses }) => {
                this.expenseData.set(expenseData);
                this.assetData.set(assetData);
                this.reconciliationStatuses.set(reconciliationStatuses ?? []);
                this.loading.set(false);
            },
            error: () => this.loading.set(false),
        });
    }

    private sortRows(rows: ExpenseTableRow[]): ExpenseTableRow[] {
        const column = this.sortColumn();
        const dir = this.sortDirection() === 'asc' ? 1 : -1;
        return [...rows].sort((a, b) => {
            if (column === 'accountName') return a.accountName.localeCompare(b.accountName) * dir;
            if (column === 'total') return (a.total - b.total) * dir;
            const va = a.monthCells.get(column)?.value ?? 0;
            const vb = b.monthCells.get(column)?.value ?? 0;
            return (va - vb) * dir;
        });
    }
}
