import { Component, inject, OnInit, signal, EventEmitter, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef, MatDialog } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatTooltipModule } from '@angular/material/tooltip';
import { FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
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
import { Account } from '../../../core/api/models/account';
import { Movement } from '../../../core/api/models/movement';
import { MatInputModule } from '@angular/material/input';
import { HttpClient } from '@angular/common/http';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { getMatcher } from '../../../core/api/fn/matchers/get-matcher';
import { Matcher } from '../../../core/api/models/matcher';
import { startWith, map, combineLatest } from 'rxjs';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';
import { toSignal } from '@angular/core/rxjs-interop';

export interface AccountGroup {
    name: string;
    accounts: Account[];
}

export type UnprocessedTransactionDialogResult =
    | { action: 'convert'; match: MatcherAndTransaction }
    | { action: 'delete'; duplicateOf: Transaction }
    | { action: 'manual'; accountId: string; description?: string }
    | { action: 'skip' };

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
        MatAutocompleteModule,
        MatTooltipModule,
        FormsModule,
        ReactiveFormsModule,
        DatePipe,
        AccountSelectComponent
    ],
    templateUrl: './unprocessed-transaction-dialog.component.html',
    styleUrl: './unprocessed-transaction-dialog.component.scss',
    styles: [`
        .badge-loading { opacity: 0.7; }
    `]
})
export class UnprocessedTransactionDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<UnprocessedTransactionDialogComponent>);
    private readonly dialog = inject(MatDialog);
    private readonly service = inject(UnprocessedTransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly budgetItemService = inject(BudgetItemService);
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    // Inject data but treat it as initial state
    readonly initialData = inject<UnprocessedTransaction>(MAT_DIALOG_DATA);

    constructor() {
        this.descriptionControl.setValue(this.initialData.transaction.description || '');
    }

    // Signal for the current state of the transaction (can be refreshed)
    protected readonly transaction = signal<UnprocessedTransaction>(this.initialData);

    protected readonly loading = signal(false);
    protected readonly accounts = this.accountService.accounts;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly budgetItems = signal<BudgetItem[]>([]);
    protected readonly manualProcessing = signal(false);

    // Output event for parent component
    readonly action = new EventEmitter<UnprocessedTransactionDialogResult>();

    protected readonly accountControl = new FormControl<string | null>(null);
    protected readonly descriptionControl = new FormControl<string>('');

    protected readonly matchersMap = signal<Map<string, Matcher>>(new Map());


    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
        this.budgetItemService.loadBudgetItems().subscribe(items => {
            this.budgetItems.set(items);
        });
        this.loadMatchers();
    }

    // I will replace ngOnInit completely to include the effect/computed logic correctly


    private loadMatchers() {
        const t = this.transaction();
        const missingIds = t.matched
            .map(m => m.matcherId)
            .filter(id => !this.matchersMap().has(id));

        new Set(missingIds).forEach(id => {
            getMatcher(this.http, this.apiConfig.rootUrl, { id: id }).subscribe({
                next: (response) => {
                    this.matchersMap.update(map => {
                        const newMap = new Map(map);
                        newMap.set(id, response.body);
                        return newMap;
                    });
                },
                error: (err) => console.error('Failed to load matcher', id, err)
            });
        });
    }

    updateTransaction(transaction: UnprocessedTransaction) {
        this.loading.set(false);
        this.accountControl.reset();
        this.transaction.set(transaction);

        this.descriptionControl.setValue(transaction.transaction.description || '');
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
        const accountId = this.accountControl.value;

        if (accountId) {
            this.loading.set(true);
            this.action.emit({
                action: 'manual',
                accountId,
                description: this.descriptionControl.value || undefined
            });
        }
    }

    skip(): void {
        this.action.emit({ action: 'skip' });
    }



    getConfidenceBadge(match: MatcherAndTransaction): { text: string; class: string; tooltip: string } {
        const matcher = this.matchersMap().get(match.matcherId);

        if (!matcher) {
            return { text: '...', class: 'badge-secondary badge-loading', tooltip: 'Loading confirmation history...' };
        }

        const count = matcher.confirmationsCount || 0;
        const total = matcher.confirmationsTotal || 0;

        if (total === 0) {
            return { text: 'New', class: 'badge-secondary', tooltip: 'No confirmation history' };
        }

        const percentage = (count / total) * 100;
        const isPerfect = percentage === 100;
        const isLargeSample = count >= 10;

        let badgeClass = 'badge-danger'; // <40%
        if (isPerfect && isLargeSample) {
            badgeClass = 'badge-perfect';
        } else if (percentage >= 70) {
            badgeClass = 'badge-success';
        } else if (percentage >= 40) {
            badgeClass = 'badge-warning';
        }

        return {
            text: `${count}/${total}`,
            class: badgeClass,
            tooltip: `${count} successful confirmations out of ${total} attempts (${percentage.toFixed(0)}%)`
        };
    }

    close(): void {
        this.dialogRef.close();
    }
}
