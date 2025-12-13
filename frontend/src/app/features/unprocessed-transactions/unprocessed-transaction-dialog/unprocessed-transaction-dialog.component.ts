import { Component, inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { DatePipe, CurrencyPipe } from '@angular/common';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';
import { MatcherAndTransaction } from '../../../core/api/models/matcher-and-transaction';
import { Transaction } from '../../../core/api/models/transaction';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';

export type UnprocessedTransactionDialogResult =
    | { action: 'convert'; match: MatcherAndTransaction }
    | { action: 'delete'; duplicateOf: Transaction };

@Component({
    selector: 'app-unprocessed-transaction-dialog',
    standalone: true,
    imports: [
        MatDialogModule,
        MatButtonModule,
        MatIconModule,
        MatListModule,
        MatCardModule,
        DatePipe,
        CurrencyPipe
    ],
    templateUrl: './unprocessed-transaction-dialog.component.html',
    styleUrl: './unprocessed-transaction-dialog.component.scss'
})
export class UnprocessedTransactionDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<UnprocessedTransactionDialogComponent>);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    readonly data = inject<UnprocessedTransaction>(MAT_DIALOG_DATA);

    protected readonly accounts = this.accountService.accounts;
    protected readonly currencies = this.currencyService.currencies;

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
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
        return (this.data.transaction.movements || []).length === 2;
    }

    get effectiveAmount(): { amount: number, currencyName: string } | null {
        if (!this.shouldShowEffectiveAmount) return null;
        // Typically in unprocessed transaction logic, we might want to show the 'unknown' amount or just the positive one?
        // Let's assume we want to show the non-primary amount? Or just the first one?
        // Requirement just says "effective amount should be shown".
        // In the table component, it calculates effective amount by summing unknown account movements.
        // If we have 2 movements, and one is 'known' (imported), likely the other is the effective impact on the user's budget.
        // Let's look for the movement that DOES NOT have accountId matching the bank importer?
        // But the unprocessed transaction doesn't explicitly link to bank importer here easily without context.
        // However, usually one movement has accountId (the bank account) and one is empty (the one to be assigned).
        // Let's check movements without accountId.

        const movements = this.data.transaction.movements || [];
        const unknownAccountMovement = movements.find(m => !m.accountId);

        if (unknownAccountMovement) {
            return {
                amount: unknownAccountMovement.amount,
                currencyName: this.getCurrencyName(unknownAccountMovement.currencyId)
            };
        }

        // If both have account IDs (unlikely for unprocessed) or none, just return the first one?
        // Let's fallback to the first movement.
        const m = movements[0];
        return { amount: m.amount, currencyName: this.getCurrencyName(m.currencyId) };
    }

    applyMatch(match: MatcherAndTransaction): void {
        this.dialogRef.close({ action: 'convert', match });
    }

    applyDuplicate(duplicate: Transaction): void {
        this.dialogRef.close({ action: 'delete', duplicateOf: duplicate });
    }

    close(): void {
        this.dialogRef.close();
    }
}
