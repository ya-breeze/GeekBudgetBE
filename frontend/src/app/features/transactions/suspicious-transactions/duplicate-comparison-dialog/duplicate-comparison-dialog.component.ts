import { Component, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { Transaction } from '../../../../core/api/models/transaction';
import { AppDatePipe } from '../../../../shared/pipes/app-date.pipe';
import { AccountDisplayComponent } from '../../../../shared/components/account-display/account-display.component';
import { TransactionService } from '../../services/transaction.service';
import { AccountService } from '../../../accounts/services/account.service';
import { CurrencyService } from '../../../currencies/services/currency.service';
import { TransactionUtils } from '../../utils/transaction.utils';

export interface DuplicateComparisonData {
    transaction1: Transaction;
    transaction2: Transaction;
}

@Component({
    selector: 'app-duplicate-comparison-dialog',
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        MatDividerModule,
        MatSnackBarModule,
        MatProgressSpinnerModule,
        AppDatePipe,
        AccountDisplayComponent,
    ],
    templateUrl: './duplicate-comparison-dialog.component.html',
    styleUrl: './duplicate-comparison-dialog.component.scss',
})
export class DuplicateComparisonDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<DuplicateComparisonDialogComponent>);
    private readonly data = inject<DuplicateComparisonData>(MAT_DIALOG_DATA);
    private readonly transactionService = inject(TransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly snackBar = inject(MatSnackBar);

    protected readonly t1 = this.data.transaction1;
    protected readonly t2 = this.data.transaction2;
    protected readonly loading = signal(false);

    protected readonly accountMap = computed(() => {
        const map = new Map();
        this.accountService.accounts().forEach((a) => map.set(a.id, a));
        return map;
    });

    protected readonly currencyMap = computed(() => {
        const map = new Map();
        this.currencyService.currencies().forEach((c) => map.set(c.id, c));
        return map;
    });

    getAccount(id: string | undefined) {
        return id ? this.accountMap().get(id) : null;
    }

    getCurrency(id: string) {
        return this.currencyMap().get(id);
    }

    formatAmount(amount: number, currencyId: string) {
        const currency = this.getCurrency(currencyId);
        return `${amount.toFixed(2)} ${currency?.name || currencyId}`;
    }

    getEffectiveAmounts(t: Transaction) {
        return TransactionUtils.getEffectiveAmounts(t);
    }

    onDismiss() {
        this.loading.set(true);
        // Dismissing means marking both as NO LONGER duplicates
        // In our backend, setting DuplicateDismissed = true on any of them
        // clears the link. For safety, we can update the one the user is looking at.
        // We MUST NOT send 'id' or other Entity fields to the backend UpdateTransaction
        // as it expects TransactionNoId and fails with "json: unknown field id".
        const update = {
            date: this.t1.date,
            description: this.t1.description,
            movements: this.t1.movements,
            partnerName: this.t1.partnerName,
            partnerAccount: this.t1.partnerAccount,
            place: this.t1.place,
            tags: this.t1.tags,
            externalIds: this.t1.externalIds,
            duplicateDismissed: true,
            suspiciousReasons: [], // Clear suspicious reasons
        };

        if (this.t1.id) {
            this.transactionService.update(this.t1.id, update).subscribe({
                next: () => {
                    this.snackBar.open('Duplicate dismissed', 'Close', { duration: 3000 });
                    this.dialogRef.close(true);
                },
                error: (err) => {
                    this.loading.set(false);
                    // Log more info to help debug if it fails again
                    console.error('Dismiss failed:', err);
                    this.snackBar.open('Failed to dismiss duplicate', 'Close', { duration: 3000 });
                },
            });
        }
    }

    onMerge(keep: 't1' | 't2') {
        const keepId = keep === 't1' ? this.t1.id : this.t2.id;
        const mergeId = keep === 't1' ? this.t2.id : this.t1.id;

        if (!keepId || !mergeId) return;

        if (
            confirm(
                'Are you sure you want to merge these transactions? The merged transaction will be soft-deleted.',
            )
        ) {
            this.loading.set(true);
            this.transactionService.merge(keepId, mergeId).subscribe({
                next: () => {
                    this.snackBar.open('Transactions merged successfully', 'Close', {
                        duration: 3000,
                    });
                    this.dialogRef.close(true);
                },
                error: () => {
                    this.loading.set(false);
                    this.snackBar.open('Failed to merge transactions', 'Close', { duration: 3000 });
                },
            });
        }
    }

    onCancel() {
        this.dialogRef.close();
    }
}
