import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Transaction } from '../../../core/api/models/transaction';
import { TransactionNoId } from '../../../core/api/models/transaction-no-id';
import {
    getTransactions,
    GetTransactions$Params,
} from '../../../core/api/fn/transactions/get-transactions';
import { createTransaction } from '../../../core/api/fn/transactions/create-transaction';
import { updateTransaction } from '../../../core/api/fn/transactions/update-transaction';
import { deleteTransaction } from '../../../core/api/fn/transactions/delete-transaction';

import { mergeTransactions } from '../../../core/api/fn/transactions/merge-transactions';

@Injectable({
    providedIn: 'root',
})
export class TransactionService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly transactions = signal<Transaction[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadTransactions(params?: GetTransactions$Params): Observable<Transaction[]> {
        this.loading.set(true);
        this.error.set(null);

        return getTransactions(this.http, this.apiConfig.rootUrl, params).pipe(
            map((response) => response.body),
            tap({
                next: (transactions) => {
                    this.transactions.set(transactions);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load transactions');
                    this.loading.set(false);
                },
            }),
        );
    }

    create(transaction: TransactionNoId): Observable<Transaction> {
        this.loading.set(true);
        this.error.set(null);

        return createTransaction(this.http, this.apiConfig.rootUrl, { body: transaction }).pipe(
            map((response) => response.body),
            tap({
                next: (transaction) => {
                    this.transactions.update((transactions) => [...transactions, transaction]);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to create transaction');
                    this.loading.set(false);
                },
            }),
        );
    }

    update(id: string, transaction: TransactionNoId): Observable<Transaction> {
        this.loading.set(true);
        this.error.set(null);

        return updateTransaction(this.http, this.apiConfig.rootUrl, { id, body: transaction }).pipe(
            map((response) => response.body),
            tap({
                next: (updatedTransaction) => {
                    this.transactions.update((transactions) =>
                        transactions.map((t) => (t.id === id ? updatedTransaction : t)),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to update transaction');
                    this.loading.set(false);
                },
            }),
        );
    }

    delete(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return deleteTransaction(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.transactions.update((transactions) =>
                        transactions.filter((t) => t.id !== id),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to delete transaction');
                    this.loading.set(false);
                },
            }),
        );
    }

    merge(keepId: string, mergeId: string): Observable<Transaction> {
        this.loading.set(true);
        this.error.set(null);

        return mergeTransactions(this.http, this.apiConfig.rootUrl, {
            body: { keepId, mergeId },
        }).pipe(
            map((response) => response.body),
            tap({
                next: (updatedTransaction) => {
                    this.transactions.update((transactions) => {
                        // Remove the merged transaction
                        const filtered = transactions.filter((t) => t.id !== mergeId);
                        // Update the kept transaction
                        return filtered.map((t) => (t.id === keepId ? updatedTransaction : t));
                    });
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to merge transactions');
                    this.loading.set(false);
                },
            }),
        );
    }
}
