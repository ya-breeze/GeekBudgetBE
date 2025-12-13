import { Component, inject, OnInit, signal, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef, MatDialog } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';
import { UnprocessedTransactionService } from '../../unprocessed-transactions/services/unprocessed-transaction.service';
import { MatcherEditDialogComponent } from '../../matchers/matcher-edit-dialog/matcher-edit-dialog.component';
import { MatcherAndTransaction } from '../../../core/api/models/matcher-and-transaction';
import { Transaction } from '../../../core/api/models/transaction';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { BudgetItemService } from '../../budget-items/services/budget-item.service';
import { BudgetItem } from '../../../core/api/models/budget-item';
import { Movement } from '../../../core/api/models/movement';
import { MatInputModule } from '@angular/material/input';

export type UnprocessedTransactionDialogResult =
    | { action: 'convert'; match: MatcherAndTransaction }
    | { action: 'delete'; duplicateOf: Transaction }
    | { action: 'manual'; accountId: string };

@Component({
    selector: 'app-unprocessed-transaction-dialog',
    standalone: true,
    imports: [
        CommonModule,
        MatButtonModule,
        MatDialogModule,
        MatIconModule,
        MatListModule,
        MatCardModule,
        MatProgressSpinnerModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        FormsModule,
        ReactiveFormsModule,
        DatePipe
    ],
    templateUrl: './unprocessed-transaction-dialog.component.html',
    styleUrl: './unprocessed-transaction-dialog.component.scss'
})
export class UnprocessedTransactionDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<UnprocessedTransactionDialogComponent>);
    private readonly dialog = inject(MatDialog);
    private readonly service = inject(UnprocessedTransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly budgetItemService = inject(BudgetItemService);

    // Inject data but treat it as initial state
    readonly initialData = inject<UnprocessedTransaction>(MAT_DIALOG_DATA);

    // Signal for the current state of the transaction (can be refreshed)
    protected readonly transaction = signal<UnprocessedTransaction>(this.initialData);

    protected readonly loading = signal(false);
    protected readonly accounts = this.accountService.accounts;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly budgetItems = signal<BudgetItem[]>([]);
    protected readonly manualProcessing = signal(false);

    // Output event for parent component
    readonly action = new EventEmitter<UnprocessedTransactionDialogResult>();

    protected selectedAccountId = signal<string | null>(null);

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
        this.budgetItemService.loadBudgetItems().subscribe(items => {
            this.budgetItems.set(items);
        });
    }

    updateTransaction(transaction: UnprocessedTransaction) {
        this.loading.set(false);
        this.selectedAccountId.set(null);
        this.transaction.set(transaction);
    }

    setLoading(isLoading: boolean) {
        this.loading.set(isLoading);
    }

    getAccountName(accountId: string | undefined): string {
        if (!accountId) return 'Unknown Account';
        const account = this.accounts().find(a => a.id === accountId);
        return account ? account.name : accountId;
    }

    createMatcher(): void {
        const dialogRef = this.dialog.open(MatcherEditDialogComponent, {
            data: { transaction: this.transaction().transaction },
            width: '98%',
            maxWidth: '98vw',
            height: '95%',
            maxHeight: '98vh',
            disableClose: true
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.refresh();
            }
        });
    }

    private refresh(): void {
        this.loading.set(true);
        this.service.getUnprocessedTransaction(this.transaction().transaction.id).subscribe({
            next: (updatedTransaction) => {
                this.transaction.set(updatedTransaction);
                this.loading.set(false);
            },
            error: () => {
                this.loading.set(false);
                // Handle error? Maybe close dialog or show message
            }
        });
    }

    getCurrencyName(currencyId: string): string {
        const currency = this.currencies().find(c => c.id === currencyId);
        return currency ? currency.name : currencyId;
    }

    get shouldShowEffectiveAmount(): boolean {
        return (this.transaction().transaction.movements || []).length === 2;
    }

    get effectiveAmount(): { amount: number, currencyName: string, knownAccountId?: string } | null {
        if (!this.shouldShowEffectiveAmount) return null;

        const movements = this.transaction().transaction.movements || [];
        const unknownAccountMovement = movements.find((m: Movement) => !m.accountId);
        const knownAccountMovement = movements.find((m: Movement) => !!m.accountId);

        if (unknownAccountMovement) {
            return {
                amount: unknownAccountMovement.amount,
                currencyName: this.getCurrencyName(unknownAccountMovement.currencyId),
                knownAccountId: knownAccountMovement?.accountId
            };
        }

        const m = movements[0];
        return { amount: m.amount, currencyName: this.getCurrencyName(m.currencyId), knownAccountId: movements.find((mov: Movement) => mov !== m)?.accountId };
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
