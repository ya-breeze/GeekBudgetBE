import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { MergedTransaction } from '../../../core/api/models/merged-transaction';
import { getMergedTransactions } from '../../../core/api/fn/merged-transactions/get-merged-transactions';
import { getMergedTransaction } from '../../../core/api/fn/merged-transactions/get-merged-transaction';
import { unmergeMergedTransaction } from '../../../core/api/fn/merged-transactions/unmerge-merged-transaction';

@Injectable({
    providedIn: 'root',
})
export class MergedTransactionService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly mergedTransactions = signal<MergedTransaction[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadMergedTransactions(): Observable<MergedTransaction[]> {
        this.loading.set(true);
        this.error.set(null);

        return getMergedTransactions(this.http, this.apiConfig.rootUrl).pipe(
            map((response) => response.body),
            tap({
                next: (transactions) => {
                    this.mergedTransactions.set(transactions);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load merged transactions');
                    this.loading.set(false);
                },
            }),
        );
    }

    getMergedTransaction(id: string): Observable<MergedTransaction> {
        return getMergedTransaction(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map((response) => response.body),
        );
    }

    unmerge(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return unmergeMergedTransaction(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.mergedTransactions.update((transactions) =>
                        transactions.filter((t) => t.transaction.id !== id),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to unmerge transaction');
                    this.loading.set(false);
                },
            }),
        );
    }
}
