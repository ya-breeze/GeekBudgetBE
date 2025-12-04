import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { DecimalPipe, JsonPipe } from '@angular/common';
import { map, forkJoin } from 'rxjs';
import { ApiConfiguration } from '../../core/api/api-configuration';
import { getExpenses } from '../../core/api/fn/aggregations/get-expenses';
import { Aggregation } from '../../core/api/models/aggregation';
import { Currency } from '../../core/api/models/currency';
import { AccountService } from '../accounts/services/account.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { UserService } from '../../core/services/user.service';

interface ExpenseTableCell {
  value: number;
  color: string;
}

interface ExpenseTableRow {
  accountId: string;
  accountName: string;
  monthCells: Map<string, ExpenseTableCell>;
  total: ExpenseTableCell;
}

interface CurrencyTable {
  currencyId: string;
  currencyName: string;
  rows: ExpenseTableRow[];
  totalRow: ExpenseTableRow;
}

@Component({
  selector: 'app-dashboard',
  imports: [
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatTableModule,
    MatButtonToggleModule,
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

  protected readonly loading = signal(true);
  protected readonly expenseData = signal<Aggregation | null>(null);
  protected readonly accounts = this.accountService.accounts;
  protected readonly selectedOutputCurrencyId = signal<string | null>(null);

  // Sorting state
  protected readonly sortColumn = signal<string>('accountName');
  protected readonly sortDirection = signal<'asc' | 'desc'>('asc');

  // Computed values for the expense table
  protected readonly accountColumns = computed(() => {
    return this.accounts().filter((acc) => acc.type === 'expense');
  });

  protected readonly monthColumns = computed(() => {
    const data = this.expenseData();
    return data?.intervals || [];
  });

  protected readonly currencyTables = computed<CurrencyTable[]>(() => {
    const data = this.expenseData();

    if (!data || !data.intervals || data.intervals.length === 0) {
      return [];
    }

    const expenseAccounts = this.accountColumns();
    const currencies = this.currencyService.currencies();
    const currenciesById = new Map<string, Currency>(
      currencies.map((currency: Currency) => [currency.id, currency])
    );

    const tables: CurrencyTable[] = [];

    data.currencies.forEach((currencyAgg) => {
      const allValues: number[] = [];

      // Collect all values for color calculation for this currency
      currencyAgg.accounts.forEach((account) => {
        account.amounts.forEach((amount) => {
          if (amount !== undefined) {
            allValues.push(amount);
          }
        });
      });

      const rows: ExpenseTableRow[] = [];

      // Build rows for each expense account for this currency
      expenseAccounts.forEach((account) => {
        const monthCells = new Map<string, ExpenseTableCell>();
        let rowTotal = 0;

        const accountData = currencyAgg.accounts.find((acc) => acc.accountId === account.id);

        data.intervals.forEach((interval, intervalIndex) => {
          let cellValue = 0;

          if (accountData && accountData.amounts[intervalIndex] !== undefined) {
            cellValue = accountData.amounts[intervalIndex];
          }

          monthCells.set(interval, {
            value: cellValue,
            color: this.calculateColor(cellValue, allValues),
          });

          rowTotal += cellValue;
        });

        if (rowTotal > 0) {
          allValues.push(rowTotal);

          rows.push({
            accountId: account.id!,
            accountName: account.name,
            monthCells,
            total: {
              value: rowTotal,
              color: this.calculateColor(rowTotal, allValues),
            },
          });
        }
      });

      if (!rows.length) {
        return;
      }

      const sortedRows = this.sortRows(rows);

      // Build total row for this currency
      const monthCells = new Map<string, ExpenseTableCell>();
      let grandTotal = 0;
      const totalValues: number[] = [];

      sortedRows.forEach((row) => {
        row.monthCells.forEach((cell) => totalValues.push(cell.value));
        totalValues.push(row.total.value);
      });

      data.intervals.forEach((interval) => {
        let monthTotal = 0;

        sortedRows.forEach((row) => {
          const cell = row.monthCells.get(interval);
          if (cell) {
            monthTotal += cell.value;
          }
        });

        monthCells.set(interval, {
          value: monthTotal,
          color: this.calculateColor(monthTotal, totalValues),
        });

        grandTotal += monthTotal;
      });

      totalValues.push(grandTotal);

      const totalRow: ExpenseTableRow = {
        accountId: 'total',
        accountName: 'Total',
        monthCells,
        total: {
          value: grandTotal,
          color: this.calculateColor(grandTotal, totalValues),
        },
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

  ngOnInit(): void {
    this.currencyService.loadCurrencies().subscribe();

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

  protected onOutputCurrencyToggle(currencyId: string): void {
    const current = this.selectedOutputCurrencyId();
    const next = current === currencyId ? null : currencyId;
    this.selectedOutputCurrencyId.set(next);
    this.loadDashboardData();
  }

  private loadDashboardData(): void {
    this.loading.set(true);

    const now = new Date();
    // Get data for the last 12 months
    const twelveMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 11, 1);
    const outputCurrencyId = this.selectedOutputCurrencyId();
    const params: { from: string; to: string; outputCurrencyId?: string } = {
      from: twelveMonthsAgo.toISOString(),
      to: now.toISOString(),
    };

    if (outputCurrencyId) {
      params.outputCurrencyId = outputCurrencyId;
    }

    forkJoin({
      accounts: this.accountService.loadAccounts(),
      expenseData: getExpenses(this.http, this.apiConfig.rootUrl, params).pipe(
        map((response) => response.body)
      ),
    }).subscribe({
      next: ({ expenseData }) => {
        console.log('Dashboard data loaded:', { expenseData });
        this.expenseData.set(expenseData);
        this.loading.set(false);
      },
      error: (error) => {
        console.error('Error loading dashboard data:', error);
        this.loading.set(false);
      },
    });
  }

  private calculateColor(value: number, allValues: number[]): string {
    if (allValues.length === 0 || value === 0) {
      return 'rgb(255, 255, 255)';
    }

    const min = Math.min(...allValues);
    const max = Math.max(...allValues);

    if (min === max) {
      return 'rgb(255, 255, 200)';
    }

    // Normalize value between 0 and 1
    const normalized = (value - min) / (max - min);

    // Create gradient from green (low) to red (high)
    // Green: rgb(200, 255, 200)
    // Yellow: rgb(255, 255, 200)
    // Red: rgb(255, 200, 200)

    let r: number, g: number, b: number;

    if (normalized < 0.5) {
      // Green to Yellow
      const t = normalized * 2;
      r = Math.round(200 + 55 * t);
      g = 255;
      b = 200;
    } else {
      // Yellow to Red
      const t = (normalized - 0.5) * 2;
      r = 255;
      g = Math.round(255 - 55 * t);
      b = 200;
    }

    return `rgb(${r}, ${g}, ${b})`;
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
    return text.replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '').trim();
  }
}
