import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';
import { getUnprocessedTransactions } from '../../../core/api/fn/unprocessed-transactions/get-unprocessed-transactions';
import { deleteUnprocessedTransaction } from '../../../core/api/fn/unprocessed-transactions/delete-unprocessed-transaction';
import { convertUnprocessedTransaction } from '../../../core/api/fn/unprocessed-transactions/convert-unprocessed-transaction';
import { getUnprocessedTransaction } from '../../../core/api/fn/unprocessed-transactions/get-unprocessed-transaction';

@Injectable({
  providedIn: 'root',
})
export class UnprocessedTransactionService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly unprocessedTransactions = signal<UnprocessedTransaction[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadUnprocessedTransactions(): Observable<UnprocessedTransaction[]> {
    this.loading.set(true);
    this.error.set(null);

    return getUnprocessedTransactions(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (transactions) => {
          this.unprocessedTransactions.set(transactions);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load unprocessed transactions');
          this.loading.set(false);
        },
      })
    );
  }

  convert(id: string, transaction: UnprocessedTransaction, matcherId?: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    const body = {
      date: transaction.transaction.date,
      description: transaction.transaction.description,
      movements: transaction.transaction.movements,
      tags: transaction.transaction.tags,
      partnerName: transaction.transaction.partnerName,
      partnerAccount: transaction.transaction.partnerAccount,
      partnerInternalId: transaction.transaction.partnerInternalId,
      place: transaction.transaction.place,
      extra: transaction.transaction.extra,
      externalIds: transaction.transaction.externalIds,
      unprocessedSources: transaction.transaction.unprocessedSources,
    };

    return convertUnprocessedTransaction(this.http, this.apiConfig.rootUrl, { id, body, matcherId }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.unprocessedTransactions.update((transactions) =>
            transactions.filter((t) => t.transaction.id !== id)
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to convert transaction');
          this.loading.set(false);
        },
      })
    );
  }

  delete(id: string, duplicateOf?: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return deleteUnprocessedTransaction(this.http, this.apiConfig.rootUrl, { id, duplicateOf }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.unprocessedTransactions.update((transactions) =>
            transactions.filter((t) => t.transaction.id !== id)
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete transaction');
          this.loading.set(false);
        },
      })
    );
  }

  getUnprocessedTransaction(id: string): Observable<UnprocessedTransaction> {
    return getUnprocessedTransaction(this.http, this.apiConfig.rootUrl, { id }).pipe(
      map((response) => response.body)
    );
  }
}
