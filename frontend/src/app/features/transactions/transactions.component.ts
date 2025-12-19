import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TransactionService } from './services/transaction.service';
import { Transaction } from '../../core/api/models/transaction';
import { TransactionFormDialogComponent } from './transaction-form-dialog/transaction-form-dialog.component';
import { AccountService } from '../accounts/services/account.service';
import { Account } from '../../core/api/models/account';
import { CurrencyService } from '../currencies/services/currency.service';
import { Currency } from '../../core/api/models/currency';
import { forkJoin, Observable } from 'rxjs';
import { TransactionUtils } from './utils/transaction.utils';
import { LayoutService } from '../../layout/services/layout.service';
import { MatcherService } from '../matchers/services/matcher.service';
import { Matcher } from '../../core/api/models/matcher';

@Component({
    selector: 'app-transactions',
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
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        DatePipe,
        FormsModule,
    ],
    templateUrl: './transactions.component.html',
    styleUrl: './transactions.component.scss',
})
export class TransactionsComponent implements OnInit {
    private readonly transactionService = inject(TransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly matcherService = inject(MatcherService);
    private readonly dialog = inject(MatDialog);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);

    private readonly route = inject(ActivatedRoute);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;

    protected readonly transactions = this.transactionService.transactions;
    protected readonly loading = signal(false); // Combined loading for transactions, accounts, and currencies
    protected readonly displayedColumns = signal([
        'date',
        'movements',
        'description',
        'effectiveAmount',
        'tags',
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

    protected readonly matchers = this.matcherService.matchers;
    protected readonly matcherMap = computed<Map<string, Matcher>>(() => {
        const map = new Map<string, Matcher>();
        this.matchers().forEach((matcher) => map.set(matcher.id!, matcher));
        return map;
    });

    // Month navigation state
    protected readonly currentMonth = signal(new Date().getMonth());
    protected readonly currentYear = signal(new Date().getFullYear());

    // Computed property for displaying current month/year
    protected readonly currentMonthDisplay = computed(() => {
        const date = new Date(this.currentYear(), this.currentMonth(), 1);
        return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    });

    // Sorting state
    protected readonly sortActive = signal<string | null>(null);
    protected readonly sortDirection = signal<'asc' | 'desc'>('asc');

    // Filter state - public for testing
    readonly selectedAccountIds = signal<string[]>([]);
    readonly selectedTags = signal<string[]>([]);
    readonly descriptionFilter = signal<string>('');

    // Computed property for all unique tags from transactions - public for testing
    readonly availableTags = computed(() => {
        const tagsSet = new Set<string>();
        this.transactions().forEach((transaction) => {
            if (transaction.tags) {
                transaction.tags.forEach((tag) => tagsSet.add(tag));
            }
        });
        return Array.from(tagsSet).sort();
    });

    // Computed property for filtered transactions - public for testing
    readonly filteredTransactions = computed(() => {
        let filtered = this.transactions();

        // Filter by accounts
        const accountIds = this.selectedAccountIds();
        if (accountIds.length > 0) {
            filtered = filtered.filter((transaction) => {
                if (!transaction.movements || transaction.movements.length === 0) {
                    return false;
                }
                // Check if any movement's accountId is in the selected accounts
                return transaction.movements.some(
                    (movement) => movement.accountId && accountIds.includes(movement.accountId),
                );
            });
        }

        // Filter by tags
        const tags = this.selectedTags();
        if (tags.length > 0) {
            filtered = filtered.filter((transaction) => {
                if (!transaction.tags || transaction.tags.length === 0) {
                    return false;
                }
                // Check if transaction has any of the selected tags
                return transaction.tags.some((tag) => tags.includes(tag));
            });
        }

        // Filter by description
        const descFilter = this.descriptionFilter().toLowerCase().trim();
        if (descFilter) {
            filtered = filtered.filter((transaction) => {
                const description = transaction.description?.toLowerCase() || '';
                return description.includes(descFilter);
            });
        }

        return filtered;
    });

    // Computed property for sorted and filtered transactions
    protected readonly sortedAndFilteredTransactions = computed(() => {
        const data = this.filteredTransactions();
        const columns = this.displayedColumns();

        if (!columns.length) {
            return data;
        }

        const active = this.sortActive() ?? columns[0];
        const direction = this.sortDirection();

        return [...data].sort((a, b) => this.compareTransactions(a, b, active, direction));
    });

    // Computed property to check if any filters are active - public for testing
    readonly hasActiveFilters = computed(() => {
        return (
            this.selectedAccountIds().length > 0 ||
            this.selectedTags().length > 0 ||
            this.descriptionFilter().trim() !== ''
        );
    });

    ngOnInit(): void {
        // Read query parameters
        this.route.queryParams.subscribe((params) => {
            let dataReloadNeeded = false;

            if (params['month'] && params['year']) {
                const month = parseInt(params['month'], 10);
                const year = parseInt(params['year'], 10);
                if (!isNaN(month) && !isNaN(year)) {
                    this.currentMonth.set(month);
                    this.currentYear.set(year);
                    dataReloadNeeded = true;
                }
            }

            if (params['accountId']) {
                this.selectedAccountIds.set([params['accountId']]);
            }

            this.loadData();
        });
    }

    loadData(): void {
        this.loading.set(true);
        forkJoin([
            this.accountService.loadAccounts(),
            this.currencyService.loadCurrencies(),
            this.matcherService.loadMatchers(),
            this.loadTransactions(),
        ]).subscribe({
            next: () => this.loading.set(false),
            error: () => this.loading.set(false),
        });
    }

    clearFilters(): void {
        this.selectedAccountIds.set([]);
        this.selectedTags.set([]);
        this.descriptionFilter.set('');
    }

    loadTransactions(): Observable<any> {
        const startOfMonth = new Date(this.currentYear(), this.currentMonth(), 1);
        const endOfMonth = new Date(
            this.currentYear(),
            this.currentMonth() + 1,
            0,
            23,
            59,
            59,
            999,
        );

        const params = {
            dateFrom: startOfMonth.toISOString(),
            dateTo: endOfMonth.toISOString(),
        };

        return this.transactionService.loadTransactions(params);
    }

    goToPreviousMonth(): void {
        if (this.currentMonth() === 0) {
            this.currentMonth.set(11);
            this.currentYear.set(this.currentYear() - 1);
        } else {
            this.currentMonth.set(this.currentMonth() - 1);
        }
        this.loadData();
    }

    goToNextMonth(): void {
        if (this.currentMonth() === 11) {
            this.currentMonth.set(0);
            this.currentYear.set(this.currentYear() + 1);
        } else {
            this.currentMonth.set(this.currentMonth() + 1);
        }
        this.loadData();
    }

    openCreateDialog(): void {
        const dialogRef = this.dialog.open(TransactionFormDialogComponent, {
            width: '90vw',
            maxWidth: '1200px',
            height: 'auto',
            maxHeight: '90vh',
            data: { mode: 'create' },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result) {
                this.transactionService.create(result).subscribe({
                    next: () => {
                        this.snackBar.open('Transaction created successfully', 'Close', {
                            duration: 3000,
                        });
                        this.loadData(); // Reload data after creation
                    },
                    error: () => {
                        this.snackBar.open('Failed to create transaction', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        });
    }

    openEditDialog(transaction: Transaction): void {
        const dialogRef = this.dialog.open(TransactionFormDialogComponent, {
            width: '90vw',
            maxWidth: '1200px',
            height: 'auto',
            maxHeight: '90vh',
            data: { mode: 'edit', transaction },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result && transaction.id) {
                this.transactionService.update(transaction.id, result).subscribe({
                    next: () => {
                        this.snackBar.open('Transaction updated successfully', 'Close', {
                            duration: 3000,
                        });
                        this.loadData(); // Reload data after update
                    },
                    error: () => {
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
                        this.loadData(); // Reload data after deletion
                    },
                    error: () => {
                        this.snackBar.open('Failed to delete transaction', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        }
    }

    formatMovements(transaction: Transaction): string {
        if (!transaction.movements || transaction.movements.length === 0) {
            return 'No movements';
        }

        const accountMap = this.accountMap();

        // Get input movements (sources of money - negative amounts)
        const inputMovements = transaction.movements.filter((movement) => movement.amount < 0);
        const inputAccountNames = inputMovements.map((movement) => {
            if (!movement.accountId) {
                return 'Undefined';
            }
            return accountMap.get(movement.accountId)?.name || movement.accountId;
        });

        // Get output movements (destinations of money - positive amounts)
        const outputMovements = transaction.movements.filter((movement) => movement.amount > 0);
        const outputAccountNames = outputMovements.map((movement) => {
            if (!movement.accountId) {
                return 'Undefined';
            }
            return accountMap.get(movement.accountId)?.name || movement.accountId;
        });

        // Format as "input => output"
        if (inputAccountNames.length === 0 && outputAccountNames.length === 0) {
            return 'No valid movements';
        } else if (inputAccountNames.length === 0) {
            return outputAccountNames.join(', ');
        } else if (outputAccountNames.length === 0) {
            return inputAccountNames.join(', ');
        } else {
            const inputPart = inputAccountNames.join(', ');
            const outputPart = outputAccountNames.join(', ');
            return `${inputPart} => ${outputPart}`;
        }
    }

    /**
     * Get target accounts (positive amounts) for display in the Accounts column
     * @param transaction The transaction to get target accounts from
     * @returns A formatted string showing only target accounts
     */
    getTargetAccounts(transaction: Transaction): string {
        if (!transaction.movements || transaction.movements.length === 0) {
            return '-';
        }

        const accountMap = this.accountMap();

        // Get output movements (destinations of money - positive amounts)
        const outputMovements = transaction.movements.filter((movement) => movement.amount > 0);
        const outputAccountNames = outputMovements.map((movement) => {
            if (!movement.accountId) {
                return 'Undefined';
            }
            return accountMap.get(movement.accountId)?.name || movement.accountId;
        });

        if (outputAccountNames.length === 0) {
            return '-';
        }

        return outputAccountNames.join(', ');
    }

    /**
     * Get the name of a matcher by its ID
     * @param matcherId The ID of the matcher
     * @returns The output description of the matcher, or the ID if not found
     */
    getMatcherName(matcherId: string): string {
        const matcher = this.matcherMap().get(matcherId);
        return matcher?.outputDescription || matcher?.descriptionRegExp || 'Unknown Matcher';
    }

    protected onSortChange(sort: Sort): void {
        if (!sort.direction) {
            this.sortActive.set(null);
            this.sortDirection.set('asc');
            return;
        }

        this.sortActive.set(sort.active);
        this.sortDirection.set(sort.direction);
    }

    /**
     * Format the effective amounts for a transaction
     * @param transaction The transaction to format
     * @returns A formatted string showing effective amounts per currency
     */
    formatEffectiveAmounts(transaction: Transaction): string {
        const effectiveAmounts = TransactionUtils.getEffectiveAmounts(transaction);

        if (effectiveAmounts.length === 0) {
            return 'N/A';
        }

        const currencyMap = this.currencyMap();

        return effectiveAmounts
            .map((ea) => {
                const currency = currencyMap.get(ea.currencyId);
                const currencyName = currency?.name || ea.currencyId;
                return `${ea.amount.toFixed(2)} ${currencyName}`;
            })
            .join(', ');
    }

    private compareTransactions(
        a: Transaction,
        b: Transaction,
        active: string,
        direction: 'asc' | 'desc',
    ): number {
        const valueA = this.getTransactionSortValue(a, active);
        const valueB = this.getTransactionSortValue(b, active);
        return this.comparePrimitiveValues(valueA, valueB, direction);
    }

    private getTransactionSortValue(
        transaction: Transaction,
        active: string,
    ): string | number | null {
        switch (active) {
            case 'date':
                return transaction.date || '';
            case 'description':
                return transaction.description?.toLowerCase() || '';
            case 'effectiveAmount': {
                // Sort by the first effective amount value
                const effectiveAmounts = TransactionUtils.getEffectiveAmounts(transaction);
                return effectiveAmounts.length > 0 ? effectiveAmounts[0].amount : 0;
            }
            case 'movements': {
                // Sort by formatted movements string
                return this.formatMovements(transaction).toLowerCase();
            }
            case 'tags': {
                // Sort by first tag alphabetically
                if (transaction.tags && transaction.tags.length > 0) {
                    return transaction.tags.sort()[0].toLowerCase();
                }
                return '';
            }
            default:
                return null;
        }
    }

    private comparePrimitiveValues(
        a: string | number | null,
        b: string | number | null,
        direction: 'asc' | 'desc',
    ): number {
        const factor = direction === 'asc' ? 1 : -1;

        if (a === null && b === null) return 0;
        if (a === null) return factor;
        if (b === null) return -factor;

        if (typeof a === 'number' && typeof b === 'number') {
            return (a - b) * factor;
        }

        const strA = String(a);
        const strB = String(b);
        return strA.localeCompare(strB) * factor;
    }
}
