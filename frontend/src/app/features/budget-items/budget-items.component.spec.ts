import { ComponentFixture, TestBed } from '@angular/core/testing';
import { BudgetItemsComponent } from './budget-items.component';
import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { UserService } from '../../core/services/user.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { LayoutService } from '../../layout/services/layout.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { of } from 'rxjs';
import { signal } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('BudgetItemsComponent', () => {
    let component: BudgetItemsComponent;
    let fixture: ComponentFixture<BudgetItemsComponent>;
    let mockBudgetItemService: any;
    let mockAccountService: any;
    let mockUserService: any;
    let mockCurrencyService: any;
    let mockLayoutService: any;
    let mockSnackBar: any;
    let mockDialog: any;

    beforeEach(async () => {
        mockBudgetItemService = {
            budgetItems: signal([]),
            budgetStatus: signal([]),
            loading: signal(false),
            loadBudgetItems: jasmine.createSpy('loadBudgetItems').and.returnValue(of([])),
            loadBudgetStatus: jasmine.createSpy('loadBudgetStatus').and.returnValue(of([])),
            create: jasmine.createSpy('create').and.returnValue(of({})),
            update: jasmine.createSpy('update').and.returnValue(of({})),
        };

        mockAccountService = {
            accounts: signal([]),
            averages: signal([]),
            loadAccounts: jasmine.createSpy('loadAccounts').and.returnValue(of([])),
            loadYearlyExpenses: jasmine.createSpy('loadYearlyExpenses').and.returnValue(of([])),
        };

        mockUserService = {
            user: signal({ favoriteCurrencyId: 'curr1' }),
            loadUser: jasmine
                .createSpy('loadUser')
                .and.returnValue(of({ favoriteCurrencyId: 'curr1' })),
        };

        mockCurrencyService = {
            currencies: signal([
                { id: 'curr1', name: 'USD' },
                { id: 'curr2', name: 'EUR' },
            ]),
            loadCurrencies: jasmine.createSpy('loadCurrencies').and.returnValue(of([])),
        };

        mockLayoutService = {
            sidenavOpened: signal(false),
        };

        mockSnackBar = {
            open: jasmine.createSpy('open'),
        };

        mockDialog = {
            open: jasmine.createSpy('open').and.returnValue({
                afterClosed: () => of(100), // Simulate user entering 100
            }),
        };

        await TestBed.configureTestingModule({
            imports: [BudgetItemsComponent, NoopAnimationsModule],
            providers: [
                { provide: BudgetItemService, useValue: mockBudgetItemService },
                { provide: AccountService, useValue: mockAccountService },
                { provide: UserService, useValue: mockUserService },
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: LayoutService, useValue: mockLayoutService },
            ],
        })
            .overrideComponent(BudgetItemsComponent, {
                add: {
                    providers: [
                        { provide: MatSnackBar, useValue: mockSnackBar },
                        { provide: MatDialog, useValue: mockDialog },
                    ],
                },
            })
            .compileComponents();

        fixture = TestBed.createComponent(BudgetItemsComponent);
        component = fixture.componentInstance;
        // Allow effects to run (Angular 17+ needs explicit handling usually, but detectChanges triggers init)
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should load initial data and trigger effect', () => {
        expect(mockBudgetItemService.loadBudgetItems).toHaveBeenCalled();
        expect(mockAccountService.loadAccounts).toHaveBeenCalled();
        expect(mockCurrencyService.loadCurrencies).toHaveBeenCalled();
        expect(mockUserService.loadUser).toHaveBeenCalled();
        // Effect should trigger status load
        expect(mockBudgetItemService.loadBudgetStatus).toHaveBeenCalled();
    });

    it('should load budget status with preferred currency', () => {
        // Assert called with currency 'curr1'
        const calls = mockBudgetItemService.loadBudgetStatus.calls.all();
        // Check finding one with arguments including 'curr1'
        // Args: from, to, outputCurrencyId
        const match = calls.some((c: any) => c.args[2] === 'curr1');
        expect(match).toBeTrue();
    });

    it('should calculate matrix data correctly with totals', () => {
        // Setup mock data
        mockAccountService.accounts.set([{ id: 'acc1', name: 'Test Account', type: 'expense' }]);
        mockBudgetItemService.budgetItems.set([
            { id: 'item1', accountId: 'acc1', amount: 500, date: new Date().toISOString() },
        ]);
        // Status with converted amounts
        mockBudgetItemService.budgetStatus.set([
            {
                accountId: 'acc1',
                date: new Date().toISOString(),
                spent: 200,
                available: 300,
                budgeted: 500,
            },
        ]);

        // Force re-computation
        fixture.detectChanges();

        const matrix = (component as any).matrixData();
        expect(matrix.length).toBe(1);

        // Check Totals
        expect(matrix[0].totalPlanned).toBeGreaterThanOrEqual(0);
        // Current month should have 500
        const currentMonthStr = new Date().toISOString().substring(0, 7);
        const cell = matrix[0].cells.find((c: any) => c.month.startsWith(currentMonthStr));
        expect(cell).toBeDefined();
        if (cell) {
            expect(cell.amount).toBe(500);
            expect(cell.spent).toBe(200);
        }
    });

    it('should change month count', () => {
        const initialCount = (component as any).monthCount();
        (component as any).changeMonthCount(1);
        expect((component as any).monthCount()).toBe(initialCount + 1);

        // Validate bounds
        (component as any).monthCount.set(12);
        (component as any).changeMonthCount(1);
        expect((component as any).monthCount()).toBe(12); // Max 12

        (component as any).monthCount.set(1);
        (component as any).changeMonthCount(-1);
        expect((component as any).monthCount()).toBe(1); // Min 1
    });

    it('should shift months', () => {
        const initialDate = (component as any).startDate();
        (component as any).shiftMonths(1);
        const newDate = (component as any).startDate();
        expect(newDate.getMonth()).not.toBe(initialDate.getMonth());
    });

    it('should correctly mark virtual and over-budget status properties', () => {
        // Setup scenarios
        // 1. Budgeted < Spent (Should be Red -> over-budget class)
        // 2. Unbudgeted & Spent > 0 & Current Month (Should be Red -> over-budget class)
        // 3. Unbudgeted & Spent > 0 & Past Month (Should be Yellow -> unbudgeted-spent class)

        const currentMonthStr = '2025-12';
        const pastMonthStr = '2025-11';

        mockAccountService.accounts.set([{ id: 'acc1', name: 'Test Account', type: 'expense' }]);

        // Mock Items (Budget Set)
        mockBudgetItemService.budgetItems.set([
            // Scenario 1: Budget Set
            { id: 'item1', accountId: 'acc1', amount: 100, date: currentMonthStr + '-01' },
        ]);

        // Mock Status
        mockBudgetItemService.budgetStatus.set([
            // Scenario 1: Budgeted 100, Spent 200
            {
                accountId: 'acc1',
                date: currentMonthStr + '-01',
                spent: 200,
                available: -100,
                budgeted: 100,
            },
            // Scenario 3: Past Month, Unbudgeted, Spent 50
            {
                accountId: 'acc1',
                date: pastMonthStr + '-01',
                spent: 50,
                available: -50,
                budgeted: 0,
            },
        ]);

        // Scenario 2: Current Month, Unbudgeted... wait, we need another item or account for this difference?
        // Let's use a different account for the unbudgeted current month scenario to avoid collision with Scenario 1
        mockAccountService.accounts.set([
            { id: 'acc1', name: 'Budgeted Account', type: 'expense' },
            { id: 'acc2', name: 'Unbudgeted Account', type: 'expense' },
        ]);

        // Status updates
        mockBudgetItemService.budgetStatus.set([
            // Scenario 1 (acc1, current): Budgeted 100, Spent 200
            {
                accountId: 'acc1',
                date: currentMonthStr + '-01',
                spent: 200,
                available: -100,
                budgeted: 100,
            },
            // Scenario 2 (acc2, current): Unbudgeted, Spent 50
            {
                accountId: 'acc2',
                date: currentMonthStr + '-01',
                spent: 50,
                available: -50,
                budgeted: 0,
            },
            // Scenario 3 (acc2, past): Unbudgeted, Spent 50
            {
                accountId: 'acc2',
                date: pastMonthStr + '-01',
                spent: 50,
                available: -50,
                budgeted: 0,
            },
        ]);

        // Ensure month count covers past month
        (component as any).monthCount.set(2);

        fixture.detectChanges();

        const matrix = (component as any).matrixData();

        // Check Scenario 1: Budgeted < Spent
        const row1 = matrix.find((r: any) => r.account.id === 'acc1');
        const cell1 = row1.cells.find((c: any) => c.month.startsWith(currentMonthStr));
        // Expect: spent > amount. Logic: 200 > 100.
        expect(cell1.spent).toBe(200);
        expect(cell1.amount).toBe(100);
        expect(cell1.spent > cell1.amount).toBeTrue();
        expect(cell1.isVirtual).toBeFalse();

        // Check Scenario 2: Unbudgeted & Spent > 0 & Current Month
        const row2 = matrix.find((r: any) => r.account.id === 'acc2');
        const cell2 = row2.cells.find((c: any) => c.month.startsWith(currentMonthStr));
        // Expect: Amount 0 (Strict), Spent 50.
        expect(cell2.amount).toBe(0);
        expect(cell2.spent).toBe(50);
        expect(cell2.spent > cell2.amount).toBeTrue(); // Triggers Red
        expect(cell2.isVirtual).toBeFalse();

        // Check Scenario 3: Unbudgeted & Spent > 0 & Past Month
        const cell3 = row2.cells.find((c: any) => c.month.startsWith(pastMonthStr));
        // Expect: Amount = Spent = 50 (Virtual).
        expect(cell3.amount).toBe(50);
        expect(cell3.spent).toBe(50);
        expect(cell3.spent > cell3.amount).toBeFalse(); // Not Red
        expect(cell3.isVirtual).toBeTrue(); // Triggers Yellow
    });
});
