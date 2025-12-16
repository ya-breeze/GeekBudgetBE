import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Account } from '../../../core/api/models/account';
import { Aggregation } from '../../../core/api/models/aggregation';
import { AccountNoId } from '../../../core/api/models/account-no-id';
import { getAccounts } from '../../../core/api/fn/accounts/get-accounts';
import { createAccount } from '../../../core/api/fn/accounts/create-account';
import { updateAccount } from '../../../core/api/fn/accounts/update-account';
import { deleteAccount } from '../../../core/api/fn/accounts/delete-account';
import { getExpenses } from '../../../core/api/fn/aggregations/get-expenses';
// Assuming AppStateService is also in a core module, adjust path if needed
// import { AppStateService } from '../../../core/state/app-state.service'; // Added import for AppStateService (commented out as path is unknown)

@Injectable({
  providedIn: 'root',
})
export class AccountService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);
  // private readonly state = inject(AppStateService); // Uncomment and adjust path if AppStateService is needed

  readonly accounts = signal<Account[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadAccounts(): Observable<Account[]> {
    this.loading.set(true);
    this.error.set(null);

    return getAccounts(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (accounts) => {
          this.accounts.set(accounts);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load accounts');
          this.loading.set(false);
        },
      })
    );
  }

  create(account: AccountNoId): Observable<Account> {
    this.loading.set(true);
    this.error.set(null);

    return createAccount(this.http, this.apiConfig.rootUrl, { body: account }).pipe(
      map((response) => response.body),
      tap({
        next: (account) => {
          this.accounts.update((accounts) => [...accounts, account]);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to create account');
          this.loading.set(false);
        },
      })
    );
  }

  update(id: string, account: AccountNoId): Observable<Account> {
    this.loading.set(true);
    this.error.set(null);

    return updateAccount(this.http, this.apiConfig.rootUrl, { id, body: account }).pipe(
      map((response) => response.body),
      tap({
        next: (updatedAccount) => {
          this.accounts.update((accounts) =>
            accounts.map((a) => (a.id === id ? updatedAccount : a))
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to update account');
          this.loading.set(false);
        },
      })
    );
  }

  delete(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return deleteAccount(this.http, this.apiConfig.rootUrl, { id }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.accounts.update((accounts) => accounts.filter((a) => a.id !== id));
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete account');
          this.loading.set(false);
        },
      })
    );
  }

  readonly averages = signal<AccountAverage[]>([]);

  loadYearlyExpenses(currencyId?: string): Observable<Aggregation> {
    this.loading.set(true);
    this.error.set(null);
    // Calculate last 12 months range
    const now = new Date();
    const to = now.toISOString();
    const fromDate = new Date();
    fromDate.setFullYear(fromDate.getFullYear() - 1);
    const from = fromDate.toISOString();

    return getExpenses(this.http, this.apiConfig.rootUrl, { from, to, outputCurrencyId: currencyId, granularity: 'year' }).pipe(
      map((response) => response.body),
      tap({
        next: (aggregation) => {
          console.log('Yearly Expenses Loaded:', aggregation);
          const avgs: AccountAverage[] = [];

          // Parse aggregation to extract account totals
          // Aggregation -> Currencies -> Accounts -> Amounts
          if (aggregation.currencies && aggregation.currencies.length > 0) {
            // We expect one currency if currencyId provided, or mixed if not. 
            // Logic: For each account, find its total in the aggregation.
            // Since we request 'year' granularity for 1 year, we expect amounts[0] to be the yearly total.

            aggregation.currencies.forEach(curr => {
              curr.accounts?.forEach(acc => {
                const totalSpent = acc.amounts?.[0] ?? 0;
                avgs.push({
                  accountId: acc.accountId!,
                  averageSpent: totalSpent / 12,
                  averageBudgeted: 0 // We don't have budgeted from expenses endpoint usually? Expenses = Spent.
                });
              });
            });
          }

          this.averages.set(avgs);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load yearly expenses');
          console.error('Failed to load yearly expenses', err);
          this.loading.set(false);
        }
      })
    );
  }
}

export interface AccountAverage {
  accountId: string;
  averageSpent: number;
  averageBudgeted: number;
}
