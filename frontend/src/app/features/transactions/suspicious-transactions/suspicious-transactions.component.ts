import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
import { DatePipe } from '@angular/common';
import { TransactionService } from '../services/transaction.service';
import { Transaction } from '../../../core/api/models/transaction';
import { TransactionFormDialogComponent } from '../transaction-form-dialog/transaction-form-dialog.component';
import { AccountService } from '../../accounts/services/account.service';
import { Account } from '../../../core/api/models/account';
import { CurrencyService } from '../../currencies/services/currency.service';
import { Currency } from '../../../core/api/models/currency';
import { forkJoin, Observable } from 'rxjs';
import { TransactionUtils } from '../utils/transaction.utils';
import { AccountDisplayComponent } from '../../../shared/components/account-display/account-display.component';

@Component({
    selector: 'app-suspicious-transactions',
    standalone: true,
    imports: [
        MatTableModule,
        MatSortModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatDialogModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatChipsModule,
        DatePipe,
        AccountDisplayComponent,
    ],
    templateUrl: './suspicious-transactions.component.html',
    styleUrl: './suspicious-transactions.component.scss',
})
export class SuspiciousTransactionsComponent implements OnInit {
    private readonly transactionService = inject(TransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly dialog = inject(MatDialog);
    private readonly snackBar = inject(MatSnackBar);

    protected readonly transactions = this.transactionService.transactions;
    protected readonly loading = signal(false);
    protected readonly displayedColumns = signal([
        'date',
        'movements',
        'description',
        'effectiveAmount',
        'suspiciousReasons',
        'actions',
    ]);

    protected readonly accounts = this.accountService.accounts;
    protected readonly accountMap = computed<Map<string, Account>>(() => {
        const map = new Map<string, Account>();
        this.accounts().forEach((account) => map.set(account.id!, account));
        return map;
    });

    protected readonly currencies = this.currencyService.currencies;
    protected readonly currencyMap = computed<Map<string, Currency>>(() => {
        const map = new Map<string, Currency>();
        this.currencies().forEach((currency) => map.set(currency.id!, currency));
        return map;
    });

    // Sorting state
    protected readonly sortActive = signal<string | null>('date');
    protected readonly sortDirection = signal<'asc' | 'desc'>('desc');

    ngOnInit(): void {
        this.loadData();
    }

    loadData(): void {
        this.loading.set(true);
        forkJoin([
            this.accountService.loadAccounts(),
            this.currencyService.loadCurrencies(),
            this.loadTransactions(),
        ]).subscribe({
            next: () => this.loading.set(false),
            error: () => this.loading.set(false),
        });
    }

    loadTransactions(): Observable<any> {
        const params = {
            onlySuspicious: true,
        };
        return this.transactionService.loadTransactions(params);
    }

    protected onSortChange(sort: Sort): void {
        this.sortActive.set(sort.active);
        this.sortDirection.set(sort.direction || 'desc');
    }

    protected readonly sortedTransactions = computed(() => {
        const data = this.transactions();
        const active = this.sortActive();
        const direction = this.sortDirection();

        if (!active || !direction) return data;

        return [...data].sort((a, b) => {
            const valueA = this.getSortValue(a, active);
            const valueB = this.getSortValue(b, active);

            if (valueA === valueB) return 0;
            if (valueA === undefined || valueA === null) return 1;
            if (valueB === undefined || valueB === null) return -1;

            const factor = direction === 'asc' ? 1 : -1;
            return valueA < valueB ? -factor : factor;
        });
    });

    private getSortValue(transaction: Transaction, active: string): any {
        switch (active) {
            case 'date':
                return transaction.date;
            case 'description':
                return transaction.description?.toLowerCase();
            case 'effectiveAmount': {
                const ea = TransactionUtils.getEffectiveAmounts(transaction);
                return ea.length > 0 ? ea[0].amount : 0;
            }
            default:
                return null;
        }
    }

    openEditDialog(transaction: Transaction): void {
        const dialogRef = this.dialog.open(TransactionFormDialogComponent, {
            width: '90vw',
            maxWidth: '1200px',
            data: { mode: 'edit', transaction },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result && transaction.id) {
                this.transactionService.update(transaction.id, result).subscribe({
                    next: () => {
                        this.snackBar.open('Transaction updated successfully', 'Close', {
                            duration: 3000,
                        });
                        this.loadData();
                    },
                    error: (err) => {
                        this.snackBar.open('Failed to update transaction', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        });
    }

    deleteTransaction(transaction: Transaction): void {
        if (confirm(`Are you sure you want to delete this transaction?`)) {
            if (transaction.id) {
                this.transactionService.delete(transaction.id).subscribe({
                    next: () => {
                        this.snackBar.open('Transaction deleted successfully', 'Close', {
                            duration: 3000,
                        });
                        this.loadData();
                    },
                    error: (err) => {
                        this.snackBar.open('Failed to delete transaction', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        }
    }

    getTargetAccountList(transaction: Transaction): Account[] {
        if (!transaction.movements || transaction.movements.length === 0) return [];
        const accountMap = this.accountMap();
        return transaction.movements
            .filter((m) => m.amount > 0)
            .map((m) => (m.accountId ? accountMap.get(m.accountId) : null))
            .filter((a): a is Account => !!a);
    }

    formatEffectiveAmounts(transaction: Transaction): string {
        const effectiveAmounts = TransactionUtils.getEffectiveAmounts(transaction);
        if (effectiveAmounts.length === 0) return 'N/A';
        const currencyMap = this.currencyMap();
        return effectiveAmounts
            .map((ea) => {
                const currency = currencyMap.get(ea.currencyId);
                return `${ea.amount.toFixed(2)} ${currency?.name || ea.currencyId}`;
            })
            .join(', ');
    }
}
