import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Account } from '../../../core/api/models/account';
import { AccountNoId } from '../../../core/api/models/account-no-id';
import { getAccounts } from '../../../core/api/fn/accounts/get-accounts';
import { createAccount } from '../../../core/api/fn/accounts/create-account';
import { updateAccount } from '../../../core/api/fn/accounts/update-account';
import { deleteAccount } from '../../../core/api/fn/accounts/delete-account';

@Injectable({
  providedIn: 'root',
})
export class AccountService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

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
}

