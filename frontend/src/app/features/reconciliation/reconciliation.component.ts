import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReconciliationService } from './services/reconciliation.service';
import { ReconciliationStatus } from '../../core/api/models/reconciliation-status';
import { Transaction } from '../../core/api/models/transaction';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { animate, state, style, transition, trigger } from '@angular/animations';

const RECONCILIATION_TOLERANCE = 0.01;
@Component({
    selector: 'app-reconciliation',
    standalone: true,
    imports: [
        CommonModule,
        MatTableModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        MatDialogModule,
        MatSnackBarModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './reconciliation.component.html',
    styleUrls: ['./reconciliation.component.scss'],
    animations: [
        trigger('detailExpand', [
            state('collapsed,void', style({ height: '0px', minHeight: '0', visibility: 'hidden' })),
            state('expanded', style({ height: '*', visibility: 'visible' })),
            transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
            transition('expanded <=> void', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
        ]),
    ],
})
export class ReconciliationComponent implements OnInit {
    private reconciliationService = inject(ReconciliationService);
    private snackBar = inject(MatSnackBar);

    statuses: ReconciliationStatus[] = [];
    loading = true;
    displayedColumns: string[] = [
        'accountName',
        'bankBalance',
        'bankBalanceAt',
        'appBalance',
        'delta',
        'lastReconciledAt',
        'actions',
    ];

    expandedElement: ReconciliationStatus | null = null;
    transactionsSinceRec: Transaction[] = [];
    loadingTransactions = false;

    ngOnInit(): void {
        this.loadStatuses();
    }

    loadStatuses(): void {
        this.loading = true;
        this.reconciliationService.loadStatuses().subscribe({
            next: (statuses: ReconciliationStatus[]) => {
                this.statuses = statuses;
                this.loading = false;
            },
            error: (_err: any) => {
                this.snackBar.open('Failed to load reconciliation status', 'Close', {
                    duration: 3000,
                });
                this.loading = false;
            },
        });
    }

    getReconcileTooltip(element: ReconciliationStatus): string {
        if (element.hasUnprocessedTransactions) {
            return 'Cannot reconcile while there are unprocessed transactions';
        }
        const delta = Math.abs(element.delta || 0);
        if (delta > RECONCILIATION_TOLERANCE) {
            return `Delta is too large to reconcile (${delta.toFixed(2)})`;
        }
        return 'Mark as Reconciled';
    }

    reconcile(status: ReconciliationStatus): void {
        if (!status.accountId || !status.currencyId) return;

        this.reconciliationService
            .reconcile(status.accountId, {
                currencyId: status.currencyId,
                balance: 0, // Use account balance
            })
            .subscribe({
                next: () => {
                    this.snackBar.open('Account reconciled successfully', 'Close', {
                        duration: 3000,
                    });
                    this.loadStatuses();
                },
                error: (_err: any) => {
                    this.snackBar.open('Failed to reconcile account', 'Close', { duration: 3000 });
                },
            });
    }

    toggleExpand(status: ReconciliationStatus): void {
        if (this.expandedElement === status) {
            this.expandedElement = null;
        } else {
            this.expandedElement = status;
            this.loadTransactionsSince(status);
        }
    }

    loadTransactionsSince(status: ReconciliationStatus): void {
        if (!status.accountId || !status.currencyId) return;

        this.loadingTransactions = true;
        this.transactionsSinceRec = [];
        this.reconciliationService
            .getTransactionsSince(status.accountId, status.currencyId)
            .subscribe({
                next: (txs: Transaction[]) => {
                    this.transactionsSinceRec = txs;
                    this.loadingTransactions = false;
                },
                error: (_err: any) => {
                    this.snackBar.open(
                        'Failed to load transactions since reconciliation',
                        'Close',
                        { duration: 3000 },
                    );
                    this.loadingTransactions = false;
                },
            });
    }

    getStatusClass(status: ReconciliationStatus): string {
        if (status.hasUnprocessedTransactions) return 'status-yellow';
        if (Math.abs(status.delta || 0) > RECONCILIATION_TOLERANCE) return 'status-red';
        return 'status-green';
    }

    enableManual(status: ReconciliationStatus): void {
        if (!status.accountId || !status.currencyId) return;

        if (
            confirm(
                `Enable manual reconciliation starting with balance ${status.appBalance} ${status.currencySymbol}?`
            )
        ) {
            this.reconciliationService
                .enableManual(status.accountId, {
                    currencyId: status.currencyId,
                    initialBalance: status.appBalance,
                })
                .subscribe({
                    next: () => {
                        this.snackBar.open('Manual reconciliation enabled', 'Close', {
                            duration: 3000,
                        });
                        this.loadStatuses();
                    },
                    error: (_err: any) => {
                        this.snackBar.open('Failed to enable manual reconciliation', 'Close', {
                            duration: 3000,
                        });
                    },
                });
        }
    }
}
