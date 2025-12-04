import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Currency } from '../../../core/api/models/currency';
import { CurrencyNoId } from '../../../core/api/models/currency-no-id';
import { getCurrencies } from '../../../core/api/fn/currencies/get-currencies';
import { createCurrency } from '../../../core/api/fn/currencies/create-currency';
import { updateCurrency } from '../../../core/api/fn/currencies/update-currency';
import { deleteCurrency } from '../../../core/api/fn/currencies/delete-currency';

@Injectable({
  providedIn: 'root',
})
export class CurrencyService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly currencies = signal<Currency[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadCurrencies(): Observable<Currency[]> {
    this.loading.set(true);
    this.error.set(null);

    return getCurrencies(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (currencies) => {
          this.currencies.set(currencies);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load currencies');
          this.loading.set(false);
        },
      })
    );
  }

  create(currency: CurrencyNoId): Observable<Currency> {
    this.loading.set(true);
    this.error.set(null);

    return createCurrency(this.http, this.apiConfig.rootUrl, { body: currency }).pipe(
      map((response) => response.body),
      tap({
        next: (currency) => {
          this.currencies.update((currencies) => [...currencies, currency]);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to create currency');
          this.loading.set(false);
        },
      })
    );
  }

  update(id: string, currency: CurrencyNoId): Observable<Currency> {
    this.loading.set(true);
    this.error.set(null);

    return updateCurrency(this.http, this.apiConfig.rootUrl, { id, body: currency }).pipe(
      map((response) => response.body),
      tap({
        next: (updatedCurrency) => {
          this.currencies.update((currencies) =>
            currencies.map((c) => (c.id === id ? updatedCurrency : c))
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to update currency');
          this.loading.set(false);
        },
      })
    );
  }

  delete(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return deleteCurrency(this.http, this.apiConfig.rootUrl, { id }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.currencies.update((currencies) => currencies.filter((c) => c.id !== id));
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete currency');
          this.loading.set(false);
        },
      })
    );
  }
}
