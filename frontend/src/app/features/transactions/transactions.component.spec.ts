import { ComponentFixture, TestBed } from '@angular/core/testing';
import { TransactionsComponent } from './transactions.component';
import { TransactionService } from './services/transaction.service';
import { AccountService } from '../accounts/services/account.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { signal } from '@angular/core';
import { Transaction } from '../../core/api/models/transaction';
import { Account } from '../../core/api/models/account';
import { Currency } from '../../core/api/models/currency';
import { of } from 'rxjs';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';

import { MatcherService } from '../matchers/services/matcher.service';
import { Matcher } from '../../core/api/models/matcher';

describe('TransactionsComponent', () => {
    let component: TransactionsComponent;
    let fixture: ComponentFixture<TransactionsComponent>;
    let mockTransactionService: jasmine.SpyObj<TransactionService>;
    let mockAccountService: jasmine.SpyObj<AccountService>;
    let mockCurrencyService: jasmine.SpyObj<CurrencyService>;
    let mockMatcherService: jasmine.SpyObj<MatcherService>;
    let mockDialog: jasmine.SpyObj<MatDialog>;
    let mockSnackBar: jasmine.SpyObj<MatSnackBar>;

    const mockTransactions: Transaction[] = [
        {
            id: '1',
            date: '2024-01-15',
            description: 'Grocery shopping',
            tags: ['food', 'groceries'],
            movements: [
                { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                { accountId: 'acc2', amount: 100, currencyId: 'usd' },
            ],
        },
        {
            id: '2',
            date: '2024-01-16',
            description: 'Gas station',
            tags: ['transport', 'fuel'],
            movements: [
                { accountId: 'acc1', amount: -50, currencyId: 'usd' },
                { accountId: 'acc3', amount: 50, currencyId: 'usd' },
            ],
        },
        {
            id: '3',
            date: '2024-01-17',
            description: 'Restaurant dinner',
            tags: ['food', 'dining'],
            movements: [
                { accountId: 'acc2', amount: -75, currencyId: 'usd' },
                { accountId: 'acc3', amount: 75, currencyId: 'usd' },
            ],
        },
        {
            id: '4',
            date: '2024-01-18',
            description: 'Online shopping',
            movements: [
                { accountId: 'acc1', amount: -200, currencyId: 'usd' },
                { accountId: 'acc2', amount: 200, currencyId: 'usd' },
            ],
        },
    ];

    const mockAccounts: Account[] = [
        { id: 'acc1', name: 'Checking Account', type: 'asset' },
        { id: 'acc2', name: 'Savings Account', type: 'asset' },
        { id: 'acc3', name: 'Expense Account', type: 'expense' },
    ];

    const mockCurrencies: Currency[] = [
        { id: 'usd', name: 'USD', description: 'US Dollar' },
        { id: 'eur', name: 'EUR', description: 'Euro' },
        { id: 'gbp', name: 'GBP', description: 'British Pound' },
    ];

    beforeEach(async () => {
        mockTransactionService = jasmine.createSpyObj(
            'TransactionService',
            ['loadTransactions', 'create', 'update', 'delete'],
            {
                transactions: signal<Transaction[]>([]),
                loading: signal(false),
                error: signal<string | null>(null),
            },
        );

        mockAccountService = jasmine.createSpyObj('AccountService', ['loadAccounts'], {
            accounts: signal<Account[]>([]),
        });

        mockCurrencyService = jasmine.createSpyObj('CurrencyService', ['loadCurrencies'], {
            currencies: signal<Currency[]>([]),
        });

        mockMatcherService = jasmine.createSpyObj('MatcherService', ['loadMatchers'], {
            matchers: signal<Matcher[]>([]),
        });

        mockDialog = jasmine.createSpyObj('MatDialog', ['open']);
        mockSnackBar = jasmine.createSpyObj('MatSnackBar', ['open']);

        mockTransactionService.loadTransactions.and.returnValue(of(mockTransactions));
        mockAccountService.loadAccounts.and.returnValue(of(mockAccounts));
        mockCurrencyService.loadCurrencies.and.returnValue(of(mockCurrencies));
        mockMatcherService.loadMatchers.and.returnValue(of([]));

        await TestBed.configureTestingModule({
            imports: [TransactionsComponent, NoopAnimationsModule],
            providers: [
                { provide: TransactionService, useValue: mockTransactionService },
                { provide: AccountService, useValue: mockAccountService },
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: MatcherService, useValue: mockMatcherService },
                { provide: MatDialog, useValue: mockDialog },
                { provide: MatSnackBar, useValue: mockSnackBar },
                provideRouter([]),
            ],
        }).compileComponents();

        fixture = TestBed.createComponent(TransactionsComponent);
        component = fixture.componentInstance;

        // Set up initial data
        mockTransactionService.transactions.set(mockTransactions);
        mockAccountService.accounts.set(mockAccounts);
        mockCurrencyService.currencies.set(mockCurrencies);

        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    describe('Filtering by Account', () => {
        it('should show all transactions when no account filter is applied', () => {
            expect(component.filteredTransactions().length).toBe(4);
        });

        it('should filter transactions by single account', () => {
            component.selectedAccountIds.set(['acc1']);
            expect(component.filteredTransactions().length).toBe(3);
            expect(component.filteredTransactions().map((t) => t.id)).toEqual(['1', '2', '4']);
        });

        it('should filter transactions by multiple accounts', () => {
            component.selectedAccountIds.set(['acc1', 'acc3']);
            expect(component.filteredTransactions().length).toBe(4);
        });

        it('should return empty array when filtering by non-existent account', () => {
            component.selectedAccountIds.set(['non-existent']);
            expect(component.filteredTransactions().length).toBe(0);
        });
    });

    describe('Filtering by Tags', () => {
        it('should show all transactions when no tag filter is applied', () => {
            expect(component.filteredTransactions().length).toBe(4);
        });

        it('should filter transactions by single tag', () => {
            component.selectedTags.set(['food']);
            expect(component.filteredTransactions().length).toBe(2);
            expect(component.filteredTransactions().map((t) => t.id)).toEqual(['1', '3']);
        });

        it('should filter transactions by multiple tags', () => {
            component.selectedTags.set(['food', 'transport']);
            expect(component.filteredTransactions().length).toBe(3);
            expect(component.filteredTransactions().map((t) => t.id)).toEqual(['1', '2', '3']);
        });

        it('should exclude transactions without tags when tag filter is applied', () => {
            component.selectedTags.set(['food']);
            const filtered = component.filteredTransactions();
            expect(filtered.every((t) => t.tags && t.tags.length > 0)).toBe(true);
        });

        it('should return empty array when filtering by non-existent tag', () => {
            component.selectedTags.set(['non-existent-tag']);
            expect(component.filteredTransactions().length).toBe(0);
        });
    });

    describe('Filtering by Description', () => {
        it('should show all transactions when no description filter is applied', () => {
            expect(component.filteredTransactions().length).toBe(4);
        });

        it('should filter transactions by description (case-insensitive)', () => {
            component.descriptionFilter.set('shopping');
            expect(component.filteredTransactions().length).toBe(2);
            expect(component.filteredTransactions().map((t) => t.id)).toEqual(['1', '4']);
        });

        it('should filter transactions by partial description match', () => {
            component.descriptionFilter.set('gas');
            expect(component.filteredTransactions().length).toBe(1);
            expect(component.filteredTransactions()[0].id).toBe('2');
        });

        it('should handle uppercase search term', () => {
            component.descriptionFilter.set('RESTAURANT');
            expect(component.filteredTransactions().length).toBe(1);
            expect(component.filteredTransactions()[0].id).toBe('3');
        });

        it('should return empty array when no description matches', () => {
            component.descriptionFilter.set('xyz123');
            expect(component.filteredTransactions().length).toBe(0);
        });

        it('should ignore leading and trailing whitespace', () => {
            component.descriptionFilter.set('  shopping  ');
            expect(component.filteredTransactions().length).toBe(2);
        });
    });

    describe('Combined Filters', () => {
        it('should apply account and tag filters together', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['food']);
            expect(component.filteredTransactions().length).toBe(1);
            expect(component.filteredTransactions()[0].id).toBe('1');
        });

        it('should apply account and description filters together', () => {
            component.selectedAccountIds.set(['acc1']);
            component.descriptionFilter.set('shopping');
            expect(component.filteredTransactions().length).toBe(2);
            expect(component.filteredTransactions().map((t) => t.id)).toEqual(['1', '4']);
        });

        it('should apply tag and description filters together', () => {
            component.selectedTags.set(['food']);
            component.descriptionFilter.set('grocery');
            expect(component.filteredTransactions().length).toBe(1);
            expect(component.filteredTransactions()[0].id).toBe('1');
        });

        it('should apply all three filters together', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['food']);
            component.descriptionFilter.set('grocery');
            expect(component.filteredTransactions().length).toBe(1);
            expect(component.filteredTransactions()[0].id).toBe('1');
        });

        it('should return empty array when combined filters match nothing', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['transport']);
            component.descriptionFilter.set('restaurant');
            expect(component.filteredTransactions().length).toBe(0);
        });
    });

    describe('Available Tags', () => {
        it('should extract all unique tags from transactions', () => {
            const tags = component.availableTags();
            expect(tags.length).toBe(5);
            expect(tags).toContain('food');
            expect(tags).toContain('groceries');
            expect(tags).toContain('transport');
            expect(tags).toContain('fuel');
            expect(tags).toContain('dining');
        });

        it('should return sorted tags', () => {
            const tags = component.availableTags();
            const sortedTags = [...tags].sort();
            expect(tags).toEqual(sortedTags);
        });

        it('should return empty array when no transactions have tags', () => {
            mockTransactionService.transactions.set([
                {
                    id: '1',
                    date: '2024-01-15',
                    description: 'Test',
                    movements: [],
                },
            ]);
            expect(component.availableTags().length).toBe(0);
        });
    });

    describe('Clear Filters', () => {
        it('should clear all filters', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['food']);
            component.descriptionFilter.set('test');

            component.clearFilters();

            expect(component.selectedAccountIds()).toEqual([]);
            expect(component.selectedTags()).toEqual([]);
            expect(component.descriptionFilter()).toBe('');
        });

        it('should show all transactions after clearing filters', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['food']);
            component.descriptionFilter.set('grocery');

            expect(component.filteredTransactions().length).toBe(1);

            component.clearFilters();

            expect(component.filteredTransactions().length).toBe(4);
        });
    });

    describe('Has Active Filters', () => {
        it('should return false when no filters are active', () => {
            expect(component.hasActiveFilters()).toBe(false);
        });

        it('should return true when account filter is active', () => {
            component.selectedAccountIds.set(['acc1']);
            expect(component.hasActiveFilters()).toBe(true);
        });

        it('should return true when tag filter is active', () => {
            component.selectedTags.set(['food']);
            expect(component.hasActiveFilters()).toBe(true);
        });

        it('should return true when description filter is active', () => {
            component.descriptionFilter.set('test');
            expect(component.hasActiveFilters()).toBe(true);
        });

        it('should return false when description filter contains only whitespace', () => {
            component.descriptionFilter.set('   ');
            expect(component.hasActiveFilters()).toBe(false);
        });

        it('should return true when multiple filters are active', () => {
            component.selectedAccountIds.set(['acc1']);
            component.selectedTags.set(['food']);
            component.descriptionFilter.set('test');
            expect(component.hasActiveFilters()).toBe(true);
        });
    });

    describe('Format Effective Amounts', () => {
        it('should format effective amount for single currency transaction', () => {
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: 'acc2', amount: 100, currencyId: 'usd' },
                ],
            };

            const result = component.formatEffectiveAmounts(transaction);
            expect(result).toBe('100.00 USD');
        });

        it('should format effective amounts for multi-currency transaction', () => {
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: 'acc2', amount: 100, currencyId: 'usd' },
                    { accountId: 'acc3', amount: -50, currencyId: 'eur' },
                    { accountId: 'acc4', amount: 50, currencyId: 'eur' },
                ],
            };

            const result = component.formatEffectiveAmounts(transaction);
            // Result should contain both currencies
            expect(result).toContain('USD');
            expect(result).toContain('EUR');
            expect(result).toContain('100.00');
            expect(result).toContain('50.00');
        });

        it('should return N/A for transaction with no movements', () => {
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [],
            };

            const result = component.formatEffectiveAmounts(transaction);
            expect(result).toBe('N/A');
        });

        it('should use currency ID when currency not found in map', () => {
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'unknown-currency' },
                    { accountId: 'acc2', amount: 100, currencyId: 'unknown-currency' },
                ],
            };

            const result = component.formatEffectiveAmounts(transaction);
            expect(result).toBe('100.00 unknown-currency');
        });

        it('should format decimal amounts with 2 decimal places', () => {
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [
                    { accountId: 'acc1', amount: -123.456, currencyId: 'usd' },
                    { accountId: 'acc2', amount: 123.456, currencyId: 'usd' },
                ],
            };

            const result = component.formatEffectiveAmounts(transaction);
            expect(result).toBe('123.46 USD');
        });

        it('should calculate effective amount correctly for complex transaction', () => {
            // +100 USD, -50 USD, -30 USD => max(100, 80) = 100
            const transaction: Transaction = {
                id: '1',
                date: '2024-01-01',
                movements: [
                    { accountId: 'acc1', amount: 100, currencyId: 'usd' },
                    { accountId: 'acc2', amount: -50, currencyId: 'usd' },
                    { accountId: 'acc3', amount: -30, currencyId: 'usd' },
                ],
            };

            const result = component.formatEffectiveAmounts(transaction);
            expect(result).toBe('100.00 USD');
        });
    });

    describe('Get Target Accounts', () => {
        it('should return target account for simple transaction', () => {
            const transaction = mockTransactions[0]; // acc1 => acc2
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('Savings Account');
        });

        it('should return multiple target accounts joined by comma', () => {
            const transaction: Transaction = {
                id: '5',
                date: '2024-01-19',
                description: 'Multi-target transaction',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: 'acc2', amount: 50, currencyId: 'usd' },
                    { accountId: 'acc3', amount: 50, currencyId: 'usd' },
                ],
            };
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('Savings Account, Expense Account');
        });

        it('should return dash when no movements', () => {
            const transaction: Transaction = {
                id: '6',
                date: '2024-01-20',
                description: 'No movements',
                movements: [],
            };
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('-');
        });

        it('should return dash when no positive movements (no target accounts)', () => {
            const transaction: Transaction = {
                id: '7',
                date: '2024-01-21',
                description: 'Only negative movements',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: 'acc2', amount: -50, currencyId: 'usd' },
                ],
            };
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('-');
        });

        it('should use account ID when account not found in map', () => {
            const transaction: Transaction = {
                id: '8',
                date: '2024-01-22',
                description: 'Unknown account',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: 'unknown-acc', amount: 100, currencyId: 'usd' },
                ],
            };
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('unknown-acc');
        });

        it('should handle undefined accountId', () => {
            const transaction: Transaction = {
                id: '9',
                date: '2024-01-23',
                description: 'Undefined account',
                movements: [
                    { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                    { accountId: undefined, amount: 100, currencyId: 'usd' },
                ],
            };
            const result = component.getTargetAccounts(transaction);
            expect(result).toBe('Undefined');
        });
    });

    describe('Sorting', () => {
        it('should initialize with no active sort', () => {
            expect(component['sortActive']()).toBeNull();
            expect(component['sortDirection']()).toBe('asc');
        });

        it('should update sort state when onSortChange is called', () => {
            component['onSortChange']({ active: 'date', direction: 'desc' });
            expect(component['sortActive']()).toBe('date');
            expect(component['sortDirection']()).toBe('desc');
        });

        it('should reset sort when direction is empty', () => {
            component['onSortChange']({ active: 'date', direction: 'desc' });
            component['onSortChange']({ active: 'date', direction: '' });
            expect(component['sortActive']()).toBeNull();
            expect(component['sortDirection']()).toBe('asc');
        });

        it('should sort transactions by date ascending', () => {
            component['onSortChange']({ active: 'date', direction: 'asc' });
            const sorted = component['sortedAndFilteredTransactions']();
            expect(sorted[0].id).toBe('1'); // 2024-01-15
            expect(sorted[1].id).toBe('2'); // 2024-01-16
            expect(sorted[2].id).toBe('3'); // 2024-01-17
            expect(sorted[3].id).toBe('4'); // 2024-01-18
        });

        it('should sort transactions by date descending', () => {
            component['onSortChange']({ active: 'date', direction: 'desc' });
            const sorted = component['sortedAndFilteredTransactions']();
            expect(sorted[0].id).toBe('4'); // 2024-01-18
            expect(sorted[1].id).toBe('3'); // 2024-01-17
            expect(sorted[2].id).toBe('2'); // 2024-01-16
            expect(sorted[3].id).toBe('1'); // 2024-01-15
        });

        it('should sort transactions by description ascending', () => {
            component['onSortChange']({ active: 'description', direction: 'asc' });
            const sorted = component['sortedAndFilteredTransactions']();
            expect(sorted[0].id).toBe('2'); // Gas station
            expect(sorted[1].id).toBe('1'); // Grocery shopping
            expect(sorted[2].id).toBe('4'); // Online shopping
            expect(sorted[3].id).toBe('3'); // Restaurant dinner
        });

        it('should sort transactions by description descending', () => {
            component['onSortChange']({ active: 'description', direction: 'desc' });
            const sorted = component['sortedAndFilteredTransactions']();
            expect(sorted[0].id).toBe('3'); // Restaurant dinner
            expect(sorted[1].id).toBe('4'); // Online shopping
            expect(sorted[2].id).toBe('1'); // Grocery shopping
            expect(sorted[3].id).toBe('2'); // Gas station
        });

        it('should sort transactions by effective amount', () => {
            component['onSortChange']({ active: 'effectiveAmount', direction: 'asc' });
            const sorted = component['sortedAndFilteredTransactions']();
            // All have same amount (50, 75, 100, 200), so check ascending order
            expect(sorted[0].id).toBe('2'); // 50
            expect(sorted[1].id).toBe('3'); // 75
            expect(sorted[2].id).toBe('1'); // 100
            expect(sorted[3].id).toBe('4'); // 200
        });

        it('should apply sorting after filtering', () => {
            component.selectedAccountIds.set(['acc1']);
            component['onSortChange']({ active: 'date', direction: 'desc' });
            const sorted = component['sortedAndFilteredTransactions']();
            // Should have transactions 1, 2, 4 (all with acc1), sorted by date desc
            expect(sorted.length).toBe(3);
            expect(sorted[0].id).toBe('4'); // 2024-01-18
            expect(sorted[1].id).toBe('2'); // 2024-01-16
            expect(sorted[2].id).toBe('1'); // 2024-01-15
        });
    });

    describe('Get Matcher Name', () => {
        const mockMatchers: Matcher[] = [
            {
                id: 'm1',
                outputDescription: 'Netflix',
                outputAccountId: 'acc3',
                confirmationsCount: 1,
                confirmationsTotal: 1,
                outputTags: [],
                currencyRegExp: '',
                partnerNameRegExp: '',
                partnerAccountNumberRegExp: '',
                descriptionRegExp: '',
                extraRegExp: '',
                confirmationHistory: [],
            },
            {
                id: 'm2',
                outputDescription: '',
                descriptionRegExp: 'Spotify',
                outputAccountId: 'acc3',
                confirmationsCount: 1,
                confirmationsTotal: 1,
                outputTags: [],
                currencyRegExp: '',
                partnerNameRegExp: '',
                partnerAccountNumberRegExp: '',
                extraRegExp: '',
                confirmationHistory: [],
            },
        ];

        beforeEach(() => {
            mockMatcherService.matchers.set(mockMatchers);
        });

        it('should return outputDescription if available', () => {
            expect(component.getMatcherName('m1')).toBe('Netflix');
        });

        it('should return descriptionRegExp if outputDescription is missing', () => {
            expect(component.getMatcherName('m2')).toBe('Spotify');
        });

        it('should return "Unknown Matcher" if matcher not found', () => {
            expect(component.getMatcherName('unknown')).toBe('Unknown Matcher');
        });
    });

    describe('Template Rendering', () => {
        it('should display suspicious icon when transaction has suspicious reasons', () => {
            const suspiciousTx: Transaction = {
                id: 'suspicious1',
                date: '2024-01-01',
                description: 'Suspicious Tx',
                movements: [],
                suspiciousReasons: ['Not present in import'],
            };
            mockTransactionService.transactions.set([suspiciousTx]);
            fixture.detectChanges();

            const icon = fixture.nativeElement.querySelector('.suspicious-icon');
            expect(icon).toBeTruthy();
        });

        it('should NOT display suspicious icon when transaction has NO suspicious reasons', () => {
            const cleanTx: Transaction = {
                id: 'clean1',
                date: '2024-01-01',
                description: 'Clean Tx',
                movements: [],
            };
            mockTransactionService.transactions.set([cleanTx]);
            fixture.detectChanges();

            const icon = fixture.nativeElement.querySelector('.suspicious-icon');
            expect(icon).toBeFalsy();
        });
    });
});
