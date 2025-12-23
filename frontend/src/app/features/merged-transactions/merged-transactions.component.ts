import { Component, inject, OnInit, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MergedTransactionService } from './services/merged-transaction.service';
import { MergedTransaction } from '../../core/api/models/merged-transaction';
import { LayoutService } from '../../layout/services/layout.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { AccountService } from '../accounts/services/account.service';
import { TransactionUtils } from '../transactions/utils/transaction.utils';
import { AccountDisplayComponent } from '../../shared/components/account-display/account-display.component';

@Component({
    selector: 'app-merged-transactions',
    standalone: true,
    imports: [
        CommonModule,
        MatTableModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
        MatTooltipModule,
        DatePipe,
        AccountDisplayComponent,
    ],
    templateUrl: './merged-transactions.component.html',
    styleUrl: './merged-transactions.component.scss',
})
export class MergedTransactionsComponent implements OnInit {
    private readonly mergedTransactionService = inject(MergedTransactionService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);
    private readonly currencyService = inject(CurrencyService);
    private readonly accountService = inject(AccountService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;
    protected readonly mergedTransactions = this.mergedTransactionService.mergedTransactions;
    protected readonly loading = this.mergedTransactionService.loading;

    protected readonly currencies = this.currencyService.currencies;
    protected readonly accounts = this.accountService.accounts;

    protected readonly currencyMap = computed(() => {
        const map = new Map<string, any>();
        this.currencies().forEach((c) => map.set(c.id!, c));
        return map;
    });

    protected readonly accountMap = computed(() => {
        const map = new Map<string, any>();
        this.accounts().forEach((a) => map.set(a.id!, a));
        return map;
    });

    protected readonly displayedColumns = [
        'date',
        'mergedDescription',
        'mergedAmount',
        'mergedAt',
        'keptDescription',
        'keptAmount',
        'actions',
    ];

    ngOnInit(): void {
        this.loadData();
    }

    loadData(): void {
        this.mergedTransactionService.loadMergedTransactions().subscribe();
        this.currencyService.loadCurrencies().subscribe();
        this.accountService.loadAccounts().subscribe();
    }

    unmerge(transaction: MergedTransaction): void {
        if (!transaction.transaction.id) return;

        this.mergedTransactionService.unmerge(transaction.transaction.id).subscribe({
            next: () => {
                this.snackBar.open('Transaction unmerged successfully', 'Close', {
                    duration: 3000,
                });
            },
            error: () => {
                this.snackBar.open('Failed to unmerge transaction', 'Close', {
                    duration: 3000,
                });
            },
        });
    }

    formatEffectiveAmounts(transaction: any): string {
        const effectiveAmounts = TransactionUtils.getEffectiveAmounts(transaction);
        if (effectiveAmounts.length === 0) return 'N/A';

        const map = this.currencyMap();
        return effectiveAmounts
            .map((ea) => {
                const currency = map.get(ea.currencyId);
                const name = currency?.name || ea.currencyId;
                return `${ea.amount.toFixed(2)} ${name}`;
            })
            .join(', ');
    }

    getTargetAccountList(transaction: any): any[] {
        if (!transaction.movements) return [];
        const map = this.accountMap();
        return transaction.movements
            .filter((m: any) => m.amount > 0 && m.accountId)
            .map((m: any) => map.get(m.accountId))
            .filter((a: any) => !!a);
    }
}
