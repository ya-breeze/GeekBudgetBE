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
import { BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration } from 'chart.js';
import { map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { getExpenses } from '../../../core/api/fn/aggregations/get-expenses';
import { Aggregation } from '../../../core/api/models/aggregation';
import { Currency } from '../../../core/api/models/currency';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';

interface CategorySummary {
  categoryId: string;
  categoryName: string;
  total: number;
  monthlyAmounts: Map<string, number>;
}

interface CurrencyReport {
  currencyId: string;
  currencyName: string;
  categorySummaries: CategorySummary[];
  grandTotal: number;
  categoryChartData: ChartConfiguration['data'];
  monthlyChartData: ChartConfiguration['data'];
}

type SortColumn = 'category' | 'amount';
type SortDirection = 'asc' | 'desc';

@Component({
  selector: 'app-expense-report',
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
    DecimalPipe,
    BaseChartDirective,
  ],
  templateUrl: './expense-report.component.html',
  styleUrl: './expense-report.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ExpenseReportComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);
  private readonly accountService = inject(AccountService);
  private readonly currencyService = inject(CurrencyService);
  private readonly userService = inject(UserService);

  protected readonly filterForm: FormGroup;
  protected readonly loading = signal(false);
  protected readonly expenseData = signal<Aggregation | null>(null);
  protected readonly sortColumn = signal<SortColumn>('amount');
  protected readonly sortDirection = signal<SortDirection>('desc');
  protected readonly selectedOutputCurrencyId = signal<string | null>(null);

  protected readonly currencyReports = computed<CurrencyReport[]>(() => {
    const data = this.expenseData();
    if (!data) return [];

    const accounts = this.accountService.accounts();
    const expenseAccounts = accounts.filter((acc) => acc.type === 'expense');
    const currencies = this.currencyService.currencies();
    const currenciesById = new Map<string, Currency>(
      currencies.map((currency: Currency) => [currency.id, currency])
    );
    const reports: CurrencyReport[] = [];

    data.currencies.forEach((currencyAgg) => {
      const summaries: CategorySummary[] = [];

      const currency = currenciesById.get(currencyAgg.currencyId);
      const currencyName = currency?.name ?? currencyAgg.currencyId;

      expenseAccounts.forEach((account) => {
        const monthlyAmounts = new Map<string, number>();
        let total = 0;

        data.intervals.forEach((interval, intervalIndex) => {
          let amount = 0;
          const accountData = currencyAgg.accounts.find((acc) => acc.accountId === account.id);
          if (accountData && accountData.amounts[intervalIndex] !== undefined) {
            amount = accountData.amounts[intervalIndex];
          }
          monthlyAmounts.set(interval, amount);
          total += amount;
        });

        if (total > 0) {
          summaries.push({
            categoryId: account.id!,
            categoryName: account.name,
            total,
            monthlyAmounts,
          });
        }
      });

      if (!summaries.length) {
        return;
      }

      const sortedSummaries = this.sortSummaries(summaries);
      const grandTotal = sortedSummaries.reduce((sum, cat) => sum + cat.total, 0);

      const categoryChartData: ChartConfiguration['data'] = {
        labels: sortedSummaries.map((s) => s.categoryName),
        datasets: [
          {
            label: 'Total Expenses by Category',
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
      ];

      const topCategories = sortedSummaries.slice(0, 5);
      const monthlyChartData: ChartConfiguration['data'] = {
        labels: data.intervals.map((interval) => this.formatMonth(interval)),
        datasets: topCategories.map((category, index) => ({
          label: category.categoryName,
          data: data.intervals.map((interval) => category.monthlyAmounts.get(interval) || 0),
          borderColor: colors[index % colors.length].replace('0.6', '1'),
          backgroundColor: colors[index % colors.length],
          fill: false,
          tension: 0.1,
        })),
      };

      reports.push({
        currencyId: currencyAgg.currencyId,
        currencyName,
        categorySummaries: sortedSummaries,
        grandTotal,
        categoryChartData,
        monthlyChartData,
      });
    });

    return reports;
  });

  protected readonly categoryChartOptions: ChartConfiguration['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'top',
      },
      title: {
        display: true,
        text: 'Expenses by Category',
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
        text: 'Monthly Expense Trends (Top 5 Categories)',
      },
    },
  };

  constructor() {
    const now = new Date();
    const twelveMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 11, 1); // First day of month 12 months ago

    this.filterForm = this.fb.group({
      startDate: [twelveMonthsAgo, [Validators.required]],
      endDate: [now, [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.accountService.loadAccounts().subscribe();
    this.currencyService.loadCurrencies().subscribe();

    // Load user data and set default currency
    this.userService.loadUser().subscribe({
      next: (user) => {
        if (user.favoriteCurrencyId && !this.selectedOutputCurrencyId()) {
          this.selectedOutputCurrencyId.set(user.favoriteCurrencyId);
        }
        this.loadExpenseData();
      },
      error: () => {
        // If user loading fails, still load expense data
        this.loadExpenseData();
      },
    });
  }

  protected onFilterSubmit(): void {
    if (this.filterForm.valid) {
      this.loadExpenseData();
    }
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
    this.loadExpenseData();
  }

  private loadExpenseData(): void {
    this.loading.set(true);

    const formValue = this.filterForm.value;
    const outputCurrencyId = this.selectedOutputCurrencyId();
    const params: { from: string; to: string; outputCurrencyId?: string } = {
      from: formValue.startDate.toISOString(),
      to: formValue.endDate.toISOString(),
    };

    if (outputCurrencyId) {
      params.outputCurrencyId = outputCurrencyId;
    }

    getExpenses(this.http, this.apiConfig.rootUrl, params)
      .pipe(map((response) => response.body))
      .subscribe({
        next: (data) => {
          this.expenseData.set(data);
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error loading expense data:', error);
          this.loading.set(false);
        },
      });
  }

  private sortSummaries(summaries: CategorySummary[]): CategorySummary[] {
    const sorted = [...summaries];
    const column = this.sortColumn();
    const direction = this.sortDirection();

    sorted.sort((a, b) => {
      let comparison = 0;

      if (column === 'category') {
        const aName = this.removeLeadingEmoji(a.categoryName);
        const bName = this.removeLeadingEmoji(b.categoryName);
        comparison = aName.localeCompare(bName);
      } else if (column === 'amount') {
        comparison = a.total - b.total;
      }

      return direction === 'asc' ? comparison : -comparison;
    });

    return sorted;
  }

  private removeLeadingEmoji(text: string): string {
    // Remove leading emoji characters for sorting
    return text.replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '').trim();
  }

  protected formatMonth(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
  }
}
