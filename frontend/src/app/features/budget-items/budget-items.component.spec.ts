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
            update: jasmine.createSpy('update').and.returnValue(of({}))
        };

        mockAccountService = {
            accounts: signal([]),
            loadAccounts: jasmine.createSpy('loadAccounts').and.returnValue(of([]))
        };

        mockUserService = {
            user: signal({ favoriteCurrencyId: 'curr1' }),
            loadUser: jasmine.createSpy('loadUser').and.returnValue(of({ favoriteCurrencyId: 'curr1' }))
        };

        mockCurrencyService = {
            currencies: signal([{ id: 'curr1', name: 'USD' }, { id: 'curr2', name: 'EUR' }]),
            loadCurrencies: jasmine.createSpy('loadCurrencies').and.returnValue(of([]))
        };

        mockLayoutService = {
            sidenavOpened: signal(false)
        };

        mockSnackBar = {
            open: jasmine.createSpy('open')
        };

        mockDialog = {
            open: jasmine.createSpy('open').and.returnValue({
                afterClosed: () => of(100) // Simulate user entering 100
            })
        };

        await TestBed.configureTestingModule({
            imports: [BudgetItemsComponent, NoopAnimationsModule],
            providers: [
                { provide: BudgetItemService, useValue: mockBudgetItemService },
                { provide: AccountService, useValue: mockAccountService },
                { provide: UserService, useValue: mockUserService },
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: LayoutService, useValue: mockLayoutService }
            ],
        })
            .overrideComponent(BudgetItemsComponent, {
                add: {
                    providers: [
                        { provide: MatSnackBar, useValue: mockSnackBar },
                        { provide: MatDialog, useValue: mockDialog }
                    ]
                }
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
        mockAccountService.accounts.set([{ id: 'acc1', name: 'Test Account' }]);
        mockBudgetItemService.budgetItems.set([
            { id: 'item1', accountId: 'acc1', amount: 500, date: new Date().toISOString() }
        ]);
        // Status with converted amounts
        mockBudgetItemService.budgetStatus.set([
            { accountId: 'acc1', date: new Date().toISOString(), spent: 200, available: 300, budgeted: 500 }
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
});
