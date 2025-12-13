import { Component, inject, OnInit, signal, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { FormsModule } from '@angular/forms';
import { DatePipe, CurrencyPipe } from '@angular/common';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';
import { MatcherAndTransaction } from '../../../core/api/models/matcher-and-transaction';
import { Transaction } from '../../../core/api/models/transaction';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';

export type UnprocessedTransactionDialogResult =
    | { action: 'convert'; match: MatcherAndTransaction }
    | { action: 'delete'; duplicateOf: Transaction }
    | { action: 'manual'; accountId: string };

@Component({
    selector: 'app-unprocessed-transaction-dialog',
    standalone: true,
    imports: [
        MatDialogModule,
        MatButtonModule,
        MatIconModule,
        MatListModule,
        MatCardModule,
        MatFormFieldModule,
        MatSelectModule,
        MatProgressSpinnerModule,
        FormsModule,
        DatePipe,
        CurrencyPipe,
        CommonModule
    ],
    templateUrl: './unprocessed-transaction-dialog.component.html',
    styleUrl: './unprocessed-transaction-dialog.component.scss'
})
export class UnprocessedTransactionDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<UnprocessedTransactionDialogComponent>);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly initialData = inject<UnprocessedTransaction>(MAT_DIALOG_DATA);

    protected readonly accounts = this.accountService.accounts;
    protected readonly currencies = this.currencyService.currencies;

    // Reactive state for the current transaction being processed
    readonly transactionData = signal<UnprocessedTransaction>(this.initialData);
    readonly loading = signal<boolean>(false);

    // Output event for parent component
    readonly action = new EventEmitter<UnprocessedTransactionDialogResult>();

    protected selectedAccountId = signal<string | null>(null);

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
    }

    updateTransaction(transaction: UnprocessedTransaction) {
        this.loading.set(false);
        this.selectedAccountId.set(null);
        this.transactionData.set(transaction);
    }

    setLoading(isLoading: boolean) {
        this.loading.set(isLoading);
    }

    getAccountName(accountId: string | undefined): string {
        if (!accountId) return 'Unknown Account';
        const account = this.accounts().find(a => a.id === accountId);
        return account ? account.name : accountId;
    }

    getCurrencyName(currencyId: string): string {
        const currency = this.currencies().find(c => c.id === currencyId);
        return currency ? currency.name : currencyId;
    }

    get shouldShowEffectiveAmount(): boolean {
        return (this.transactionData().transaction.movements || []).length === 2;
    }

    get effectiveAmount(): { amount: number, currencyName: string, knownAccountId?: string } | null {
        if (!this.shouldShowEffectiveAmount) return null;

        const movements = this.transactionData().transaction.movements || [];
        const unknownAccountMovement = movements.find(m => !m.accountId);
        const knownAccountMovement = movements.find(m => !!m.accountId);

        if (unknownAccountMovement) {
            return {
                amount: unknownAccountMovement.amount,
                currencyName: this.getCurrencyName(unknownAccountMovement.currencyId),
                knownAccountId: knownAccountMovement?.accountId
            };
        }

        const m = movements[0];
        return { amount: m.amount, currencyName: this.getCurrencyName(m.currencyId), knownAccountId: movements.find(mov => mov !== m)?.accountId };
    }

    get knownAccountName(): string | null {
        const amountData = this.effectiveAmount;
        if (amountData?.knownAccountId) {
            return this.getAccountName(amountData.knownAccountId);
        }
        return null;
    }

    applyMatch(match: MatcherAndTransaction): void {
        this.loading.set(true);
        this.action.emit({ action: 'convert', match });
    }

    applyDuplicate(duplicate: Transaction): void {
        this.loading.set(true);
        this.action.emit({ action: 'delete', duplicateOf: duplicate });
    }

    processManual(): void {
        const accountId = this.selectedAccountId();
        if (accountId) {
            this.loading.set(true);
            this.action.emit({ action: 'manual', accountId });
        }
    }

    close(): void {
        this.dialogRef.close();
    }
}
