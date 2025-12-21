import { Component, inject, OnInit, signal, EventEmitter, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import {
    MAT_DIALOG_DATA,
    MatDialogModule,
    MatDialogRef,
    MatDialog,
} from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
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
import { MatcherService } from '../../matchers/services/matcher.service';
import { MatInputModule } from '@angular/material/input';
import { HttpClient } from '@angular/common/http';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { getMatcher } from '../../../core/api/fn/matchers/get-matcher';
import { Matcher } from '../../../core/api/models/matcher';
import { startWith, map } from 'rxjs';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';
import { toSignal } from '@angular/core/rxjs-interop';
import { OverlayModule } from '@angular/cdk/overlay';
import { COMMA, ENTER } from '@angular/cdk/keycodes';

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
        AccountSelectComponent,
        OverlayModule,
        MatChipsModule,
    ],
    templateUrl: './unprocessed-transaction-dialog.component.html',
    styleUrl: './unprocessed-transaction-dialog.component.scss',
    styles: [
        `
            .badge-loading {
                opacity: 0.7;
            }
        `,
    ],
})
export class UnprocessedTransactionDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<UnprocessedTransactionDialogComponent>);
    private readonly dialog = inject(MatDialog);
    private readonly service = inject(UnprocessedTransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly budgetItemService = inject(BudgetItemService);
    private readonly matcherService = inject(MatcherService);
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    // Inject data but treat it as initial state
    readonly initialData = inject<UnprocessedTransaction>(MAT_DIALOG_DATA);
    readonly separatorKeysCodes = [ENTER, COMMA] as const;

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

    protected readonly showEditPopover = signal(false);
    protected readonly matcherSearchControl = new FormControl<string | Matcher>('');
    protected readonly matcherSearchQuery = toSignal(
        this.matcherSearchControl.valueChanges.pipe(
            startWith(''),
            map((v) => (typeof v === 'string' ? v : v?.outputDescription || '')),
        ),
        { initialValue: '' },
    );

    protected readonly allMatchers = this.matcherService.matchers; // Signal from service
    protected readonly filteredMatchers = computed(() => {
        const query = this.matcherSearchQuery();
        const matchers = this.allMatchers();

        let filtered = matchers;

        if (query) {
            const lowerQuery = query.toLowerCase();
            filtered = matchers.filter((m) => {
                const accountName = this.getAccountName(m.outputAccountId).toLowerCase();
                return (
                    (m.descriptionRegExp &&
                        m.descriptionRegExp.toLowerCase().includes(lowerQuery)) ||
                    (m.outputDescription &&
                        m.outputDescription.toLowerCase().includes(lowerQuery)) ||
                    accountName.includes(lowerQuery)
                );
            });
        }

        // Return sorted
        return filtered.sort((a, b) => {
            const nameA =
                `${this.getAccountName(a.outputAccountId)}: ${a.outputDescription || ''}`.toLowerCase();
            const nameB =
                `${this.getAccountName(b.outputAccountId)}: ${b.outputDescription || ''}`.toLowerCase();
            return nameA.localeCompare(nameB);
        });
    });

    // Output event for parent component
    readonly action = new EventEmitter<UnprocessedTransactionDialogResult>();

    protected readonly accountControl = new FormControl<string | null>(null);
    protected readonly descriptionControl = new FormControl<string>('');

    protected readonly matchersMap = signal<Map<string, Matcher>>(new Map());

    // Details / Edit State
    protected readonly expandedMatchId = signal<string | null>(null);
    protected readonly expandedDuplicateId = signal<string | null>(null);

    // Initial edit state for the expanded match
    protected readonly editState = signal<{
        description: string;
        tags: string[];
        movements: Movement[];
    } | null>(null);

    // Component-level editing for manual processing or transaction details
    protected readonly tagsControl = new FormControl<string[]>([]);
    protected readonly tagInputControl = new FormControl('');

    // Manual Processing State
    protected readonly manualMovements = signal<Movement[]>([]);

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();
        this.budgetItemService.loadBudgetItems().subscribe((items) => {
            this.budgetItems.set(items);
        });
        this.matcherService.loadMatchers().subscribe(); // Load all matchers for the search
        this.loadMatchers(); // Load specific matchers for this transaction matches (if any missing from global? logic slightly duplicative but okay)

        this.initializeManualState();
    }

    private initializeManualState() {
        const t = this.transaction().transaction;
        this.tagsControl.setValue(t.tags || []);

        // Initialize manual movements with a deep copy
        this.manualMovements.set(t.movements ? t.movements.map((m) => ({ ...m })) : []);
    }

    private loadMatchers() {
        const t = this.transaction();
        const missingIds = t.matched
            .map((m) => m.matcherId)
            .filter((id) => !this.matchersMap().has(id));

        new Set(missingIds).forEach((id) => {
            getMatcher(this.http, this.apiConfig.rootUrl, { id: id }).subscribe({
                next: (response) => {
                    this.matchersMap.update((map) => {
                        const newMap = new Map(map);
                        newMap.set(id, response.body);
                        return newMap;
                    });
                },
                error: (err) => console.error('Failed to load matcher', id, err),
            });
        });
    }

    updateTransaction(transaction: UnprocessedTransaction) {
        this.loading.set(false);
        this.accountControl.reset();
        this.transaction.set(transaction);

        this.descriptionControl.setValue(transaction.transaction.description || '');
        this.initializeManualState();
    }

    setLoading(isLoading: boolean) {
        this.loading.set(isLoading);
    }

    addTag(event: any, control: FormControl<string[] | null>): void {
        const value = (event.value || '').trim();
        const currentTags = control.value || [];

        if (value) {
            control.setValue([...currentTags, value]);
        }

        event.chipInput!.clear();
        this.tagInputControl.setValue('');
    }

    removeTag(tag: string, control: FormControl<string[] | null>): void {
        const currentTags = control.value || [];
        const index = currentTags.indexOf(tag);

        if (index >= 0) {
            const newTags = [...currentTags];
            newTags.splice(index, 1);
            control.setValue(newTags);
        }
    }

    // For Match Edit State Tags
    addMatchTag(event: any): void {
        const value = (event.value || '').trim();
        const state = this.editState();
        if (value && state) {
            this.editState.set({
                ...state,
                tags: [...(state.tags || []), value],
            });
        }
        event.chipInput!.clear();
    }

    removeMatchTag(tag: string): void {
        const state = this.editState();
        if (state && state.tags) {
            this.editState.set({
                ...state,
                tags: state.tags.filter((t) => t !== tag),
            });
        }
    }

    updateMatchMovementAccount(index: number, accountId: string) {
        const state = this.editState();
        if (!state) return;

        const newMovements = [...state.movements];
        newMovements[index] = { ...newMovements[index], accountId };

        this.editState.set({ ...state, movements: newMovements });
    }

    updateManualMovementAccount(index: number, accountId: string) {
        const current = this.manualMovements();
        const updated = [...current];
        updated[index] = { ...updated[index], accountId };
        this.manualMovements.set(updated);
    }

    isOriginalAccountDefined(index: number): boolean {
        const originalMovements = this.transaction().transaction.movements;
        return !!originalMovements && !!originalMovements[index]?.accountId;
    }

    updateMatchDescription(description: string) {
        const state = this.editState();
        if (state) {
            this.editState.set({ ...state, description });
        }
    }

    toggleMatch(matchId: string, matchData: MatcherAndTransaction) {
        if (this.expandedMatchId() === matchId) {
            this.expandedMatchId.set(null);
            this.editState.set(null);
        } else {
            this.expandedMatchId.set(matchId);
            this.expandedDuplicateId.set(null); // exclusive

            // Initialize edit state
            this.editState.set({
                description: matchData.transaction.description || '',
                tags: matchData.transaction.tags ? [...matchData.transaction.tags] : [],
                movements: matchData.transaction.movements
                    ? matchData.transaction.movements.map((m) => ({ ...m }))
                    : [],
            });
        }
    }

    toggleDuplicate(id: string) {
        if (this.expandedDuplicateId() === id) {
            this.expandedDuplicateId.set(null);
        } else {
            this.expandedDuplicateId.set(id);
            this.expandedMatchId.set(null); // exclusive
            this.editState.set(null);
        }
    }

    getAccountName(accountId: string | undefined): string {
        if (!accountId) return 'Unknown Account';
        const account = this.accounts().find((a) => a.id === accountId);
        return account ? account.name : accountId;
    }

    createMatcher(): void {
        this.openMatcherDialog();
    }

    editMatcher(matcher: Matcher): void {
        this.openMatcherDialog(matcher);
    }

    private openMatcherDialog(matcher?: Matcher): void {
        const dialogRef = this.dialog.open(MatcherEditDialogComponent, {
            data: {
                transaction: this.transaction().transaction,
                matcher: matcher,
            },
            width: '98%',
            maxWidth: '98vw',
            height: '95%',
            maxHeight: '98vh',
            disableClose: true,
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result) {
                // If we edited a matcher, we should refresh the transaction to see if it now matches
                // Or if we created one.
                this.refresh();
                // Also reload matchers list to update the dropdown
                this.matcherService.loadMatchers().subscribe();
            }
        });
    }

    toggleEditPopover(): void {
        this.showEditPopover.update((v) => !v);
        if (this.showEditPopover()) {
            this.matcherSearchControl.setValue('');
        }
    }

    onMatcherSelected(event: any): void {
        const matcher: Matcher = event.option.value;
        this.editMatcher(matcher);
        this.showEditPopover.set(false);
    }

    displayMatcherFn = (matcher: Matcher): string => {
        if (!matcher) return '';
        const accountName = this.getAccountName(matcher.outputAccountId);
        return `${accountName}: ${matcher.outputDescription || ''}`;
    };

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
            },
        });
    }

    getCurrencyName(currencyId: string): string {
        const currency = this.currencies().find((c) => c.id === currencyId);
        return currency ? currency.name : currencyId;
    }

    applyMatch(match: MatcherAndTransaction): void {
        // If this match is currently being edited, use the edit state
        if (this.expandedMatchId() === match.matcherId && this.editState()) {
            const state = this.editState()!;
            const updatedMatch: MatcherAndTransaction = {
                ...match,
                transaction: {
                    ...match.transaction,
                    description: state.description,
                    tags: state.tags,
                    movements: state.movements,
                },
            };
            this.loading.set(true);
            this.action.emit({ action: 'convert', match: updatedMatch });
            return;
        }

        this.loading.set(true);
        this.action.emit({ action: 'convert', match });
    }

    applyDuplicate(duplicate: Transaction): void {
        this.loading.set(true);
        this.action.emit({ action: 'delete', duplicateOf: duplicate });
    }

    processManual(): void {
        const movements = this.manualMovements();
        // optionally validate that all movements have accounts?
        // Logic: at least one side should be known, or typically all should be assigned for a 'processed' transaction.
        // For now, let's allow saving if at least one account is selected (which the UI might enforce).

        const manualTransaction: Transaction = {
            ...this.transaction().transaction,
            description:
                this.descriptionControl.value || this.transaction().transaction.description,
            tags: this.tagsControl.value || [],
            movements: movements,
        };

        // Create a dummy matcher structure for the 'convert' action
        // effectively treating this as a "Manual Match"
        const manualMatch: MatcherAndTransaction = {
            matcherId: 'manual', // specific ID or empty
            transaction: manualTransaction,
        };

        this.loading.set(true);
        this.action.emit({ action: 'convert', match: manualMatch });
    }

    skip(): void {
        this.action.emit({ action: 'skip' });
    }

    getConfidenceBadge(match: MatcherAndTransaction): {
        text: string;
        class: string;
        tooltip: string;
    } {
        const matcher = this.matchersMap().get(match.matcherId);

        if (!matcher) {
            return {
                text: '...',
                class: 'badge-secondary badge-loading',
                tooltip: 'Loading confirmation history...',
            };
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
            tooltip: `${count} successful confirmations out of ${total} attempts (${percentage.toFixed(0)}%)`,
        };
    }

    close(): void {
        this.dialogRef.close();
    }
}
