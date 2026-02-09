import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { ReconciliationStatus } from '../../../core/api/models/reconciliation-status';
import { Reconciliation } from '../../../core/api/models/reconciliation';
import { Transaction } from '../../../core/api/models/transaction';
import { getReconciliationStatus } from '../../../core/api/fn/reconciliation/get-reconciliation-status';
import { reconcileAccount } from '../../../core/api/fn/reconciliation/reconcile-account';
import { getTransactionsSinceReconciliation } from '../../../core/api/fn/reconciliation/get-transactions-since-reconciliation';
import { enableAccountReconciliation } from '../../../core/api/fn/reconciliation/enable-account-reconciliation';
import { ReconcileAccountRequest } from '../../../core/api/models/reconcile-account-request';
import { EnableReconciliationRequest } from '../../../core/api/models/enable-reconciliation-request';

@Injectable({
    providedIn: 'root',
})
export class ReconciliationService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly statuses = signal<ReconciliationStatus[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadStatuses(): Observable<ReconciliationStatus[]> {
        this.loading.set(true);
        this.error.set(null);

        return getReconciliationStatus(this.http, this.apiConfig.rootUrl).pipe(
            map((response) => response.body),
            tap({
                next: (statuses) => {
                    this.statuses.set(statuses);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load reconciliation status');
                    this.loading.set(false);
                },
            }),
        );
    }

    reconcile(id: string, body: ReconcileAccountRequest): Observable<Reconciliation> {
        return reconcileAccount(this.http, this.apiConfig.rootUrl, { id, body }).pipe(
            map((response) => response.body),
        );
    }

    getTransactionsSince(id: string, currencyId: string): Observable<Transaction[]> {
        return getTransactionsSinceReconciliation(this.http, this.apiConfig.rootUrl, {
            id,
            currencyId,
        }).pipe(map((response) => response.body));
    }

    enableManual(id: string, body: EnableReconciliationRequest): Observable<Reconciliation> {
        return enableAccountReconciliation(this.http, this.apiConfig.rootUrl, { id, body }).pipe(
            map((response) => response.body),
        );
    }
}
