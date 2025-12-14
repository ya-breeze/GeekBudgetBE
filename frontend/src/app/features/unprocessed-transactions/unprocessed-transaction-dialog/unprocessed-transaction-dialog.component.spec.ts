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
            movements: []
        },
        matched: [],
        duplicates: []
    };

    beforeEach(async () => {
        const dialogSpy = jasmine.createSpyObj('MatDialog', ['open']);
        const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['close', 'afterClosed']);
        const unprocessedServiceSpy = jasmine.createSpyObj('UnprocessedTransactionService', ['getUnprocessedTransaction']);
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
                NoopAnimationsModule
            ],
            providers: [
                { provide: MAT_DIALOG_DATA, useValue: mockTransaction },
                { provide: MatDialogRef, useValue: dialogRefSpy },
                // { provide: MatDialog, useValue: dialogSpy }, // redundant if overriding?
                { provide: UnprocessedTransactionService, useValue: unprocessedServiceSpy },
                { provide: AccountService, useValue: accountServiceSpy },
                { provide: CurrencyService, useValue: currencyServiceSpy },
                { provide: BudgetItemService, useValue: budgetItemServiceSpy },
                { provide: MatcherService, useValue: matcherServiceSpy }
            ]
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

    it('should toggle edit popover', () => {
        expect(component['showEditPopover']()).toBeFalse();
        component.toggleEditPopover();
        expect(component['showEditPopover']()).toBeTrue();
        component.toggleEditPopover();
        expect(component['showEditPopover']()).toBeFalse();
    });

    it('should filter matchers', () => {
        const matchers: Matcher[] = [
            { id: '1', outputDescription: 'Netflix', descriptionRegExp: 'netflix' } as any,
            { id: '2', outputDescription: 'Uber', descriptionRegExp: 'uber' } as any
        ];

        // Update the signal
        matcherServiceSpy.matchers.set(matchers);

        // Test empty search
        component['matcherSearchControl'].setValue('');
        expect(component['filteredMatchers']().length).toBe(2);

        // Test search "net"
        component['matcherSearchControl'].setValue('net');
        expect(component['filteredMatchers']().length).toBe(1);
        expect(component['filteredMatchers']()[0].outputDescription).toBe('Netflix');

        // Test search "Uber" (case insensitive check implied by implementation?)
        // Implementation uses .toLowerCase()
        component['matcherSearchControl'].setValue('UBER');
        expect(component['filteredMatchers']().length).toBe(1);
        expect(component['filteredMatchers']()[0].outputDescription).toBe('Uber');
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
});
