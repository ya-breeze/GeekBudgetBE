import {
    Component,
    inject,
    OnInit,
    signal,
    computed,
    ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule, DecimalPipe } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration, ChartType } from 'chart.js';
import { map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { getBalances } from '../../../core/api/fn/aggregations/get-balances';
import { Aggregation } from '../../../core/api/models/aggregation';
import { Currency } from '../../../core/api/models/currency';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';
import { AccountDisplayComponent } from '../../../shared/components/account-display/account-display.component';

interface AccountSummary {
    accountId: string;
    accountName: string;
    accountImage?: string;
    history: number[];
}

interface CurrencyReport {
    currencyId: string;
    currencyName: string;
    accountSummaries: AccountSummary[];
    intervals: string[];
    chartData: ChartConfiguration['data'];
    tableDataSource: any[];
    tableColumns: string[];
}

@Component({
    selector: 'app-balance-report',
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatButtonModule,
        MatButtonToggleModule,
        MatProgressSpinnerModule,
        MatTableModule,
        MatIconModule,
        MatSlideToggleModule,
        DecimalPipe,
        BaseChartDirective,
        AccountDisplayComponent,
    ],
    templateUrl: './balance-report.component.html',
    styleUrl: './balance-report.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BalanceReportComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly userService = inject(UserService);

    protected readonly filterForm: FormGroup;
    protected readonly loading = signal(false);
    protected readonly balanceData = signal<Aggregation | null>(null);
    protected readonly selectedOutputCurrencyId = signal<string | null>(null);
    protected readonly includeHidden = signal(false);

    protected readonly chartOptions: ChartConfiguration['options'] = {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
            x: {
                stacked: true,
            },
            y: {
                stacked: true,
                beginAtZero: false,
            },
        },
        plugins: {
            legend: {
                display: true,
                position: 'bottom',
            },
            tooltip: {
                mode: 'index',
                intersect: false,
            },
        },
    };

    protected readonly currencyReports = computed<CurrencyReport[]>(() => {
        const data = this.balanceData();
        if (!data) return [];

        const accounts = this.accountService.accounts();
        const currencies = this.currencyService.currencies();
        const currenciesById = new Map<string, Currency>(currencies.map((c) => [c.id, c]));

        const reports: CurrencyReport[] = [];

        data.currencies.forEach((currencyAgg) => {
            const currency = currenciesById.get(currencyAgg.currencyId);
            const currencyName = currency?.name ?? currencyAgg.currencyId;

            const accountSummaries: AccountSummary[] = [];

            currencyAgg.accounts.forEach((accAgg) => {
                const account = accounts.find((a) => a.id === accAgg.accountId);
                if (!account) return;

                accountSummaries.push({
                    accountId: accAgg.accountId,
                    accountName: account.name,
                    accountImage: account.image,
                    history: accAgg.amounts,
                });
            });

            if (!accountSummaries.length) return;

            // Chart Data (Stacked Area)
            const colors = [
                'rgba(46, 204, 113, 0.5)', // Emerald
                'rgba(52, 152, 219, 0.5)', // Blue
                'rgba(155, 89, 182, 0.5)', // Amethyst
                'rgba(241, 194, 50, 0.5)', // Sun Flower
                'rgba(231, 76, 60, 0.5)', // Alizarin
            ];

            const chartData: ChartConfiguration['data'] = {
                labels: data.intervals.map((i) => this.formatMonth(i)),
                datasets: accountSummaries.map((acc, index) => ({
                    label: acc.accountName,
                    data: acc.history,
                    fill: true,
                    backgroundColor: colors[index % colors.length],
                    borderColor: colors[index % colors.length].replace('0.5', '1'),
                    pointRadius: 2,
                    tension: 0.3,
                })),
            };

            // Table Data (Rotated: Accounts as Rows, Months as Columns)
            const intervalKeys = data.intervals;
            const tableColumns = ['account', ...intervalKeys, 'total'];

            // 1. Create rows for each account
            const tableDataSource: any[] = accountSummaries.map((acc) => {
                const row: any = { account: acc.accountName, accountImage: acc.accountImage };
                const accountTotal = 0;
                acc.history.forEach((val, i) => {
                    const key = intervalKeys[i];
                    row[key] = val;
                });
                // The total is the last value in the history (since it's cumulative)
                // Wait, if it's cumulative balance, the "Total" at the end of the year
                // is just the last month's balance?
                // Or should it be the average?
                // Usually for balance reports, the last point is the "current" balance.
                // Let's use the last available value as the "Total/Current" balance.
                row.total = acc.history.length > 0 ? acc.history[acc.history.length - 1] : 0;
                return row;
            });

            // 2. Create a "Total" row for the bottom
            const totalRow: any = { account: 'Total Assets' };
            const grandTotal = 0;
            intervalKeys.forEach((key, i) => {
                let intervalSum = 0;
                accountSummaries.forEach((acc) => {
                    intervalSum += acc.history[i] || 0;
                });
                totalRow[key] = intervalSum;
            });
            totalRow.total = accountSummaries.reduce((sum, acc) => {
                return sum + (acc.history.length > 0 ? acc.history[acc.history.length - 1] : 0);
            }, 0);

            tableDataSource.push(totalRow);

            reports.push({
                currencyId: currencyAgg.currencyId,
                currencyName,
                accountSummaries,
                intervals: data.intervals,
                chartData,
                tableDataSource,
                tableColumns,
            });
        });

        return reports;
    });

    constructor() {
        const now = new Date();
        const sixMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 5, 1);

        this.filterForm = this.fb.group({
            startDate: [sixMonthsAgo, [Validators.required]],
            endDate: [now, [Validators.required]],
        });
    }

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();

        this.userService.loadUser().subscribe({
            next: (user) => {
                if (user.favoriteCurrencyId && !this.selectedOutputCurrencyId()) {
                    this.selectedOutputCurrencyId.set(user.favoriteCurrencyId);
                }
                this.loadData();
            },
            error: () => this.loadData(),
        });
    }

    protected onFilterSubmit(): void {
        if (this.filterForm.valid) {
            this.loadData();
        }
    }

    protected onOutputCurrencyToggle(currencyId: string): void {
        const current = this.selectedOutputCurrencyId();
        const next = current === currencyId ? null : currencyId;
        this.selectedOutputCurrencyId.set(next);
        this.loadData();
    }

    protected onIncludeHiddenToggle(): void {
        this.includeHidden.update((c) => !c);
        this.loadData();
    }

    private loadData(): void {
        this.loading.set(true);
        const formValue = this.filterForm.value;

        const params: any = {
            from: formValue.startDate.toISOString(),
            to: formValue.endDate.toISOString(),
            includeHidden: this.includeHidden(),
        };

        if (this.selectedOutputCurrencyId()) {
            params.outputCurrencyId = this.selectedOutputCurrencyId();
        }

        getBalances(this.http, this.apiConfig.rootUrl, params)
            .pipe(map((r) => r.body))
            .subscribe({
                next: (data) => {
                    this.balanceData.set(data);
                    this.loading.set(false);
                },
                error: (err) => {
                    console.error('Error loading Balance Report:', err);
                    this.loading.set(false);
                },
            });
    }

    protected formatMonth(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
    }
}
