import { ComponentFixture, TestBed } from '@angular/core/testing';
import { UnprocessedTransactionDialogComponent } from './unprocessed-transaction-dialog.component';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialog } from '@angular/material/dialog';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { UnprocessedTransactionService } from '../../unprocessed-transactions/services/unprocessed-transaction.service';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { BudgetItemService } from '../../budget-items/services/budget-item.service';
import { MatcherService } from '../../matchers/services/matcher.service';
import { of } from 'rxjs';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { Matcher } from '../../../core/api/models/matcher';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';
import { signal } from '@angular/core';

describe('UnprocessedTransactionDialogComponent', () => {
    let component: UnprocessedTransactionDialogComponent;
    let fixture: ComponentFixture<UnprocessedTransactionDialogComponent>;
    let matcherServiceSpy: any; // Use any to allow property assignment

    const mockTransaction = {
        transaction: {
            id: 't1',
            date: '2023-01-01',
            description: 'Test',
            amount: 100,
            movements: [],
        },
        matched: [
            {
                matcherId: 'm1',
                transaction: {
                    description: 'Match Desc',
                    movements: [{ accountId: 'acc-1', amount: 100, currencyId: 'CZK' }],
                    tags: ['tag1'],
                },
            } as any,
        ],
        duplicates: [{ id: 'd1', description: 'Dup Desc', date: '2023-01-01' } as any],
    };

    beforeEach(async () => {
        const dialogSpy = jasmine.createSpyObj('MatDialog', ['open']);
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['close', 'afterClosed']);
        const unprocessedServiceSpy = jasmine.createSpyObj('UnprocessedTransactionService', [
            'getUnprocessedTransaction',
        ]);
        const accountServiceSpy = jasmine.createSpyObj('AccountService', ['loadAccounts']);
        accountServiceSpy.accounts = signal([]);
        const currencyServiceSpy = jasmine.createSpyObj('CurrencyService', ['loadCurrencies']);
        currencyServiceSpy.currencies = signal([]);
        const budgetItemServiceSpy = jasmine.createSpyObj('BudgetItemService', ['loadBudgetItems']);

        matcherServiceSpy = jasmine.createSpyObj('MatcherService', ['loadMatchers']);
        matcherServiceSpy.matchers = signal([]);

        // Setup spy returns immediately
        accountServiceSpy.loadAccounts.and.returnValue(of([]));
        currencyServiceSpy.loadCurrencies.and.returnValue(of([]));
        matcherServiceSpy.loadMatchers.and.returnValue(of([]));
        budgetItemServiceSpy.loadBudgetItems.and.returnValue(of([]));
        unprocessedServiceSpy.getUnprocessedTransaction.and.returnValue(of(mockTransaction));

        await TestBed.configureTestingModule({
            imports: [
                UnprocessedTransactionDialogComponent,
                HttpClientTestingModule,
                NoopAnimationsModule,
            ],
            providers: [
                { provide: MAT_DIALOG_DATA, useValue: mockTransaction },
                { provide: MatDialogRef, useValue: dialogRefSpy },
                // { provide: MatDialog, useValue: dialogSpy }, // redundant if overriding?
                { provide: UnprocessedTransactionService, useValue: unprocessedServiceSpy },
                { provide: AccountService, useValue: accountServiceSpy },
                { provide: CurrencyService, useValue: currencyServiceSpy },
                { provide: BudgetItemService, useValue: budgetItemServiceSpy },
                { provide: MatcherService, useValue: matcherServiceSpy },
            ],
        })
            .overrideProvider(MatDialog, { useValue: dialogSpy })
            .compileComponents();

        fixture = TestBed.createComponent(UnprocessedTransactionDialogComponent);
        component = fixture.componentInstance;

        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should handle transactions with 3 movements correctly (repro)', () => {
        // Setup transaction with 3 movements: Source, Fee, Target(Unknown)
        const threeMovTransaction: UnprocessedTransaction = {
            transaction: {
                id: 't3',
                description: 'Complex Transaction',
                date: '2023-01-01',
                movements: [
                    { accountId: 'acc-1', amount: -100, currencyId: 'USD' }, // Source
                    { accountId: 'acc-1', amount: -1, currencyId: 'USD' },   // Fee
                    { accountId: undefined, amount: 101, currencyId: 'USD' } // Target (Unknown)
                ],
                tags: []
            } as any,
            matched: [],
            duplicates: []
        };

        // Re-create component with new data
        // Note: In a real app we'd likely use a different setup or override the injection, 
        // but here we can just update the signal if it were writable or re-instantiate.
        // Actually, the component's `transaction` signal is `protected readonly transaction = signal...`. We can try casting to any to set it.

        (component as any).transaction.set(threeMovTransaction);
        fixture.detectChanges();

        // Initialize manual state (since we updated the transaction signal manually in the test setup above, 
        // we might also need to trigger the initialization logic or call updateTransaction if the component doesn't auto-detect signal changes deeply for side effects in ngOnInit, 
        // but we called ngOnInit once. 
        // Ideally we should call `updateTransaction` or similar. 
        // Let's call the private/public update method if accessible or just check if manualMovements updated if it's using an effect (it's not, it's init logic).
        // Since we are hacking the signal set in the test: (component as any).transaction.set(...)
        // The `initializeManualState` was called in ngOnInit with the OLD data.
        // We need to re-initialize it.
        (component as any).initializeManualState();

        const movements = component['manualMovements']();
        expect(movements.length).toBe(3);
        expect(movements[2].amount).toBe(101);
    });

    it('should toggle edit popover', () => {
        expect(component['showEditPopover']()).toBeFalse();
        component.toggleEditPopover();
        expect(component['showEditPopover']()).toBeTrue();
        component.toggleEditPopover();
        expect(component['showEditPopover']()).toBeFalse();
    });

    it('should filter matchers and sort them', () => {
        const matchers: Matcher[] = [
            { id: '1', outputDescription: 'Zebra Service', outputAccountId: 'acc1' } as any,
            { id: '2', outputDescription: 'Banana Store', outputAccountId: 'acc2' } as any,
            { id: '3', outputDescription: 'Apple Store', outputAccountId: 'acc1' } as any,
        ];

        // Mock accounts
        const accounts = [
            { id: 'acc1', name: 'Expenses: Groceries' },
            { id: 'acc2', name: 'Expenses: Entertainment' },
        ];
        // We can't easily mock the protected accounts signal directly if it's derived from service in the component constructor/field init
        // But we mocked keys in LoadAccounts.
        // Actually the component calls `this.accountService.accounts` which is a signal (mocked as spy signal above?)
        // The spy set up `accountServiceSpy.accounts = signal([])`. Let's update that signal.

        const accountServiceSpy = TestBed.inject(AccountService) as any;
        accountServiceSpy.accounts.set(accounts);

        matcherServiceSpy.matchers.set(matchers);
        fixture.detectChanges(); // Update computed

        // 1. Test Sorting (Expected: "Expenses: Entertainment: Banana...", "Expenses: Groceries: Apple...", "Expenses: Groceries: Zebra...")
        // Wait, "Expenses: Entertainment" comes before "Expenses: Groceries"
        const filtered = component['filteredMatchers']();
        expect(filtered.length).toBe(3);
        expect(filtered[0].id).toBe('2'); // Ent: Banana
        expect(filtered[1].id).toBe('3'); // Groc: Apple
        expect(filtered[2].id).toBe('1'); // Groc: Zebra

        // 2. Test Search by Account Name ("Entertain")
        component['matcherSearchControl'].setValue('Entertain');
        fixture.detectChanges();
        const searched = component['filteredMatchers']();
        expect(searched.length).toBe(1);
        expect(searched[0].id).toBe('2');
    });
    it('should open edit dialog when matcher selected', () => {
        const dialogSpy = TestBed.inject(MatDialog) as jasmine.SpyObj<MatDialog>;
        dialogSpy.open.and.returnValue({ afterClosed: () => of(true) } as any);

        const matcher = { id: '1', outputDescription: 'Test' } as Matcher;
        component.onMatcherSelected({ option: { value: matcher } });

        expect(dialogSpy.open).toHaveBeenCalled();
        // Verify showEditPopover is set to false
        expect(component['showEditPopover']()).toBeFalse();
    });

    it('should toggle match expansion and initialize edit state', () => {
        const match = mockTransaction.matched[0];
        // Toggle On
        component.toggleMatch(match.matcherId, match);
        expect(component['expandedMatchId']()).toBe(match.matcherId);
        expect(component['editState']()).toEqual({
            description: match.transaction.description!,
            tags: match.transaction.tags!,
            movements: match.transaction.movements!,
        });

        // Toggle Off
        component.toggleMatch(match.matcherId, match);
        expect(component['expandedMatchId']()).toBeNull();
        expect(component['editState']()).toBeNull();
    });

    it('should toggle duplicate expansion exclusive of match', () => {
        // Need to ensure signals are updated or check initial state
        const match = mockTransaction.matched[0];
        const duplicate = mockTransaction.duplicates[0];

        component.toggleMatch(match.matcherId, match);
        expect(component['expandedMatchId']()).toBe(match.matcherId);

        component.toggleDuplicate(duplicate.id);
        expect(component['expandedDuplicateId']()).toBe(duplicate.id);
        expect(component['expandedMatchId']()).toBeNull(); // Should close match
    });

    it('should emit edited values when applying an edited match', () => {
        spyOn(component.action, 'emit');
        const match = mockTransaction.matched[0];

        // 1. Expand Match
        component.toggleMatch(match.matcherId, match);

        // 2. Modify State
        component.addMatchTag({ value: 'new-tag', chipInput: { clear: () => { } } });
        const currentState = component['editState']()!;
        expect(currentState.tags).toContain('new-tag');

        // 3. Apply
        component.applyMatch(match);

        // 4. Verify Emit
        expect(component.action.emit).toHaveBeenCalledWith(
            jasmine.objectContaining({
                action: 'convert',
                match: jasmine.objectContaining({
                    matcherId: match.matcherId,
                    transaction: jasmine.objectContaining({
                        tags: jasmine.arrayContaining(['new-tag']),
                        description: match.transaction.description, // unchanged
                    }),
                }),
            }),
        );
    });
});
