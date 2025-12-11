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

describe('TransactionsComponent', () => {
  let component: TransactionsComponent;
  let fixture: ComponentFixture<TransactionsComponent>;
  let mockTransactionService: jasmine.SpyObj<TransactionService>;
  let mockAccountService: jasmine.SpyObj<AccountService>;
  let mockCurrencyService: jasmine.SpyObj<CurrencyService>;
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
    mockTransactionService = jasmine.createSpyObj('TransactionService', ['loadTransactions', 'create', 'update', 'delete'], {
      transactions: signal<Transaction[]>([]),
      loading: signal(false),
      error: signal<string | null>(null),
    });

    mockAccountService = jasmine.createSpyObj('AccountService', ['loadAccounts'], {
      accounts: signal<Account[]>([]),
    });

    mockCurrencyService = jasmine.createSpyObj('CurrencyService', ['loadCurrencies'], {
      currencies: signal<Currency[]>([]),
    });

    mockDialog = jasmine.createSpyObj('MatDialog', ['open']);
    mockSnackBar = jasmine.createSpyObj('MatSnackBar', ['open']);

    mockTransactionService.loadTransactions.and.returnValue(of(mockTransactions));
    mockAccountService.loadAccounts.and.returnValue(of(mockAccounts));
    mockCurrencyService.loadCurrencies.and.returnValue(of(mockCurrencies));

    await TestBed.configureTestingModule({
      imports: [TransactionsComponent, NoopAnimationsModule],
      providers: [
        { provide: TransactionService, useValue: mockTransactionService },
        { provide: AccountService, useValue: mockAccountService },
        { provide: CurrencyService, useValue: mockCurrencyService },
        { provide: MatDialog, useValue: mockDialog },
        { provide: MatSnackBar, useValue: mockSnackBar },
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
      expect(component.filteredTransactions().map(t => t.id)).toEqual(['1', '2', '4']);
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
      expect(component.filteredTransactions().map(t => t.id)).toEqual(['1', '3']);
    });

    it('should filter transactions by multiple tags', () => {
      component.selectedTags.set(['food', 'transport']);
      expect(component.filteredTransactions().length).toBe(3);
      expect(component.filteredTransactions().map(t => t.id)).toEqual(['1', '2', '3']);
    });

    it('should exclude transactions without tags when tag filter is applied', () => {
      component.selectedTags.set(['food']);
      const filtered = component.filteredTransactions();
      expect(filtered.every(t => t.tags && t.tags.length > 0)).toBe(true);
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
      expect(component.filteredTransactions().map(t => t.id)).toEqual(['1', '4']);
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
      expect(component.filteredTransactions().map(t => t.id)).toEqual(['1', '4']);
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
});


