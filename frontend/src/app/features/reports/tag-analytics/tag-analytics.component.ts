import {
    Component,
    inject,
    OnInit,
    signal,
    computed,
    ChangeDetectionStrategy,
} from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { CommonModule, DecimalPipe } from '@angular/common';
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
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration } from 'chart.js';
import { map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { getExpenses } from '../../../core/api/fn/aggregations/get-expenses';
import { Aggregation } from '../../../core/api/models/aggregation';
import { Currency } from '../../../core/api/models/currency';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';
import { AccountService } from '../../accounts/services/account.service';
import { Account } from '../../../core/api/models/account';

interface TagSummary {
    tagName: string;
    total: number;
    monthlyAmounts: Map<string, number>;
}

interface CurrencyReport {
    currencyId: string;
    currencyName: string;
    tagSummaries: TagSummary[];
    grandTotal: number;
    tagChartData: ChartConfiguration['data'];
    monthlyChartData: ChartConfiguration['data'];
}

type SortColumn = 'tag' | 'amount';
type SortDirection = 'asc' | 'desc';

@Component({
    selector: 'app-tag-analytics',
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
        MatSelectModule,
        MatChipsModule,
        MatTooltipModule,
        DecimalPipe,
        BaseChartDirective,
    ],
    templateUrl: './tag-analytics.component.html',
    styleUrl: './tag-analytics.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TagAnalyticsComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly currencyService = inject(CurrencyService);
    private readonly userService = inject(UserService);
    private readonly accountService = inject(AccountService);

    protected readonly filterForm: FormGroup;
    protected readonly loading = signal(false);
    protected readonly tagsLoading = signal(true);
    protected readonly analyticsData = signal<Aggregation | null>(null);
    protected readonly sortColumn = signal<SortColumn>('amount');
    protected readonly sortDirection = signal<SortDirection>('desc');
    protected readonly selectedOutputCurrencyId = signal<string | null>(null);

    protected readonly availableTags = signal<{ id: string; name: string }[]>([]);
    protected readonly hiddenTags = signal<string[]>([]);
    protected readonly availableAccounts = computed(() => {
        return this.accountService.accounts().filter((a) => a.type === 'expense');
    });

    protected readonly currencyReports = computed<CurrencyReport[]>(() => {
        const data = this.analyticsData();
        if (!data) return [];

        const currencies = this.currencyService.currencies();
        const currenciesById = new Map<string, Currency>(
            currencies.map((currency: Currency) => [currency.id, currency]),
        );
        const reports: CurrencyReport[] = [];

        data.currencies.forEach((currencyAgg) => {
            const summaries: TagSummary[] = [];
            const currency = currenciesById.get(currencyAgg.currencyId);
            const currencyName = currency?.name ?? currencyAgg.currencyId;
            const hidden = this.hiddenTags();

            currencyAgg.accounts?.forEach((accountAgg) => {
                const tagName = accountAgg.accountId || 'Untagged';
                if (hidden.includes(tagName)) {
                    return;
                }

                const monthlyAmounts = new Map<string, number>();
                let total = 0;

                data.intervals.forEach((interval, intervalIndex) => {
                    let amount = 0;
                    if (accountAgg.amounts[intervalIndex] !== undefined) {
                        amount = accountAgg.amounts[intervalIndex];
                    }
                    monthlyAmounts.set(interval, amount);
                    total += amount;
                });

                if (total > 0 || Array.from(monthlyAmounts.values()).some((v) => v > 0)) {
                    summaries.push({
                        tagName,
                        total,
                        monthlyAmounts,
                    });
                }
            });

            if (!summaries.length) {
                return;
            }

            const sortedSummaries = this.sortSummaries(summaries);
            const grandTotal = sortedSummaries.reduce((sum, tag) => sum + tag.total, 0);

            const tagChartData: ChartConfiguration['data'] = {
                labels: sortedSummaries.map((s) => s.tagName),
                datasets: [
                    {
                        label: 'Total Expenses by Tag',
                        data: sortedSummaries.map((s) => s.total),
                        backgroundColor: 'rgba(33, 150, 243, 0.6)',
                        borderColor: 'rgba(33, 150, 243, 1)',
                        borderWidth: 1,
                    },
                ],
            };

            const colors = [
                'rgba(33, 150, 243, 0.6)',
                'rgba(76, 175, 80, 0.6)',
                'rgba(255, 152, 0, 0.6)',
                'rgba(156, 39, 176, 0.6)',
                'rgba(244, 67, 54, 0.6)',
                'rgba(0, 188, 212, 0.6)',
                'rgba(103, 58, 183, 0.6)',
            ];

            const monthlyChartData: ChartConfiguration['data'] = {
                labels: data.intervals.map((interval) => this.formatMonth(interval)),
                // Line for each tag summary, up to 10
                datasets: sortedSummaries.slice(0, 10).map((tagSum, index) => ({
                    label: tagSum.tagName,
                    data: data.intervals.map(
                        (interval) => tagSum.monthlyAmounts.get(interval) || 0,
                    ),
                    borderColor: colors[index % colors.length].replace('0.6', '1'),
                    backgroundColor: colors[index % colors.length],
                    fill: false,
                    tension: 0.1,
                })),
            };

            reports.push({
                currencyId: currencyAgg.currencyId,
                currencyName,
                tagSummaries: sortedSummaries,
                grandTotal,
                tagChartData,
                monthlyChartData,
            });
        });

        return reports;
    });

    protected readonly tagChartOptions: ChartConfiguration['options'] = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: 'top',
            },
            title: {
                display: true,
                text: 'Total Expenses by Tag',
            },
        },
    };

    protected readonly monthlyChartOptions: ChartConfiguration['options'] = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: 'top',
            },
            title: {
                display: true,
                text: 'Monthly Expense Trends by Tag',
            },
        },
    };

    constructor() {
        const now = new Date();
        const twelveMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 11, 1);

        this.filterForm = this.fb.group({
            startDate: [twelveMonthsAgo, [Validators.required]],
            endDate: [now, [Validators.required]],
            selectedAccounts: [[]],
            selectedTags: [['ALL'], [Validators.required]],
        });
    }

    ngOnInit(): void {
        this.currencyService.loadCurrencies().subscribe();
        this.accountService.loadAccounts().subscribe();

        // Load user data and set default currency
        this.userService.loadUser().subscribe({
            next: (user) => {
                if (user.favoriteCurrencyId && !this.selectedOutputCurrencyId()) {
                    this.selectedOutputCurrencyId.set(user.favoriteCurrencyId);
                }
                this.fetchAllTags();
            },
            error: () => {
                this.fetchAllTags();
            },
        });

        // Listen to account selection changes to reload available tags
        this.filterForm.get('selectedAccounts')?.valueChanges.subscribe(() => {
            this.fetchAllTags();
            if (this.filterForm.valid) {
                this.loadAnalyticsData();
            }
        });
    }

    private fetchAllTags(): void {
        this.tagsLoading.set(true);
        // Fetch all possible tags user has used by querying a large date range
        const from = '2000-01-01T00:00:00.000Z';
        const to = new Date().toISOString();

        const selectedAccounts = this.filterForm.value.selectedAccounts;
        const params: any = {
            from,
            to,
            groupBy: 'tag',
            granularity: 'year',
        };

        if (selectedAccounts && selectedAccounts.length > 0) {
            params.accounts = selectedAccounts;
        }

        getExpenses(this.http, this.apiConfig.rootUrl, params)
            .pipe(map((response) => response.body))
            .subscribe({
                next: (data) => {
                    const tagNames = new Set<string>();
                    data?.currencies?.forEach((curr) => {
                        curr.accounts?.forEach((acc) => {
                            if (acc.accountId && acc.accountId !== 'Untagged') {
                                // Split by comma just in case, though backend should return distinct ones
                                acc.accountId.split(',').forEach((tag) => {
                                    const cleaned = tag.trim();
                                    if (cleaned) tagNames.add(cleaned);
                                });
                            }
                        });
                    });

                    const tagObjects = Array.from(tagNames)
                        .sort()
                        .map((name) => ({ id: name, name }));

                    this.availableTags.set([
                        { id: 'ALL', name: 'All' },
                        ...tagObjects,
                        { id: 'UNTAGGED', name: 'No tags' },
                    ]);

                    if (this.filterForm.value.selectedTags.length === 0) {
                        this.filterForm.patchValue({ selectedTags: ['ALL'] }, { emitEvent: false });
                    }

                    // Re-trigger data load after tags are loaded for the first time
                    // or reloaded after accounts change
                    if (this.filterForm.valid) {
                        this.loadAnalyticsData();
                    }

                    this.tagsLoading.set(false);
                },
                error: (err) => {
                    console.error('Failed to load tags', err);
                    this.tagsLoading.set(false);
                },
            });
    }

    protected onFilterSubmit(): void {
        this.loadAnalyticsData();
    }

    protected onSortChange(column: SortColumn): void {
        if (this.sortColumn() === column) {
            this.sortDirection.set(this.sortDirection() === 'asc' ? 'desc' : 'asc');
        } else {
            this.sortColumn.set(column);
            this.sortDirection.set('asc');
        }
    }

    protected onOutputCurrencyToggle(currencyId: string): void {
        const current = this.selectedOutputCurrencyId();
        const next = current === currencyId ? null : currencyId;
        this.selectedOutputCurrencyId.set(next);
        this.loadAnalyticsData();
    }

    protected toggleTagVisibility(tagName: string): void {
        const current = this.hiddenTags();
        if (current.includes(tagName)) {
            this.hiddenTags.set(current.filter((t) => t !== tagName));
        } else {
            this.hiddenTags.set([...current, tagName]);
        }
    }

    protected unhideTag(tagName: string): void {
        this.hiddenTags.set(this.hiddenTags().filter((t) => t !== tagName));
    }

    private loadAnalyticsData(): void {
        this.loading.set(true);
        const formValue = this.filterForm.value;
        const outputCurrencyId = this.selectedOutputCurrencyId();

        const params: any = {
            from: formValue.startDate.toISOString(),
            to: formValue.endDate.toISOString(),
            includeHidden: false,
            groupBy: 'tag',
        };

        if (formValue.selectedAccounts && formValue.selectedAccounts.length > 0) {
            params.accounts = formValue.selectedAccounts;
        }

        if (!formValue.selectedTags.includes('ALL')) {
            params.tags = formValue.selectedTags.map((tag: string) =>
                tag === 'UNTAGGED' ? 'Untagged' : tag,
            );
        }

        if (outputCurrencyId) {
            params.outputCurrencyId = outputCurrencyId;
        }

        getExpenses(this.http, this.apiConfig.rootUrl, params)
            .pipe(map((response) => response.body))
            .subscribe({
                next: (data) => {
                    this.analyticsData.set(data);
                    this.loading.set(false);
                },
                error: (error) => {
                    console.error('Error loading tag analytics:', error);
                    this.loading.set(false);
                },
            });
    }

    private sortSummaries(summaries: TagSummary[]): TagSummary[] {
        const sorted = [...summaries];
        const column = this.sortColumn();
        const direction = this.sortDirection();

        sorted.sort((a, b) => {
            let comparison = 0;
            if (column === 'tag') {
                const aName = this.removeLeadingEmoji(a.tagName);
                const bName = this.removeLeadingEmoji(b.tagName);
                comparison = aName.localeCompare(bName);
            } else if (column === 'amount') {
                comparison = a.total - b.total;
            }
            return direction === 'asc' ? comparison : -comparison;
        });

        return sorted;
    }

    private removeLeadingEmoji(text: string): string {
        return text
            .replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '')
            .trim();
    }

    protected formatMonth(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
    }
}
