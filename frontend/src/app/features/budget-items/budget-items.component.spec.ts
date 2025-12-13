import { ComponentFixture, TestBed } from '@angular/core/testing';
import { BudgetItemsComponent } from './budget-items.component';
import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { LayoutService } from '../../layout/services/layout.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { of } from 'rxjs';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { signal } from '@angular/core';

describe('BudgetItemsComponent', () => {
    let component: BudgetItemsComponent;
    let fixture: ComponentFixture<BudgetItemsComponent>;
    let mockBudgetItemService: any;
    let mockAccountService: any;
    let mockLayoutService: any;
    let mockSnackBar: any;

    beforeEach(async () => {
        mockBudgetItemService = {
            budgetItems: signal([]),
            budgetStatus: signal([]),
            loading: signal(false),
            loadBudgetItems: jasmine.createSpy('loadBudgetItems').and.returnValue(of([])),
            loadBudgetStatus: jasmine.createSpy('loadBudgetStatus').and.returnValue(of([])),
            create: jasmine.createSpy('create').and.returnValue(of({})),
        };

        mockAccountService = {
            accounts: signal([]),
            loadAccounts: jasmine.createSpy('loadAccounts').and.returnValue(of([])),
        };

        mockLayoutService = {
            sidenavOpened: signal(true),
        };

        mockSnackBar = {
            open: jasmine.createSpy('open'),
        };

        await TestBed.configureTestingModule({
            imports: [BudgetItemsComponent, BrowserAnimationsModule, HttpClientTestingModule],
            providers: [
                { provide: BudgetItemService, useValue: mockBudgetItemService },
                { provide: AccountService, useValue: mockAccountService },
                { provide: LayoutService, useValue: mockLayoutService },
            ],
        })
            .overrideComponent(BudgetItemsComponent, {
                add: {
                    providers: [{ provide: MatSnackBar, useValue: mockSnackBar }]
                }
            })
            .compileComponents();

        fixture = TestBed.createComponent(BudgetItemsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
        const injectedSnackBar = (component as any).snackBar;
        expect(injectedSnackBar).toBe(mockSnackBar, 'Injected SnackBar should be the mock');
    });

    it('should initialize form', () => {
        // Access property using any/protected access workaround for test
        const form = (component as any).addBudgetForm;
        expect(form).toBeTruthy();
        expect(form.get('accountId')).toBeTruthy();
        expect(form.get('amount')).toBeTruthy();
        expect(form.get('date')).toBeTruthy();
    });

    it('should call create service on submit', () => {
        const form = (component as any).addBudgetForm;
        form.patchValue({
            accountId: 'acc-1',
            amount: 100,
            date: new Date(),
        });

        // Determine expected methods based on button click or direct method call
        // Calling protected method directly for unit test
        (component as any).addBudgetItem();

        expect(mockBudgetItemService.create).toHaveBeenCalled();
        expect(mockSnackBar.open).toHaveBeenCalled();
    });

    it('should not call create service if form is invalid', () => {
        const form = (component as any).addBudgetForm;
        form.patchValue({
            accountId: '', // Invalid
            amount: 100,
            date: new Date(),
        });

        (component as any).addBudgetItem();

        expect(mockBudgetItemService.create).not.toHaveBeenCalled();
    });
});
