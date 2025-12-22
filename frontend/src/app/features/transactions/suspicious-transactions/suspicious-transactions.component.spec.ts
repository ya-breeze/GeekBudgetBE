import { ComponentFixture, TestBed } from '@angular/core/testing';
import { SuspiciousTransactionsComponent } from './suspicious-transactions.component';
import { TransactionService } from '../services/transaction.service';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { signal } from '@angular/core';
import { Transaction } from '../../../core/api/models/transaction';
import { of } from 'rxjs';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';

describe('SuspiciousTransactionsComponent', () => {
    let component: SuspiciousTransactionsComponent;
    let fixture: ComponentFixture<SuspiciousTransactionsComponent>;
    let mockTransactionService: jasmine.SpyObj<TransactionService>;
    let mockAccountService: jasmine.SpyObj<AccountService>;
    let mockCurrencyService: jasmine.SpyObj<CurrencyService>;

    beforeEach(async () => {
        mockTransactionService = jasmine.createSpyObj('TransactionService', ['loadTransactions'], {
            transactions: signal<Transaction[]>([]),
        });
        mockAccountService = jasmine.createSpyObj('AccountService', ['loadAccounts'], {
            accounts: signal([]),
        });
        mockCurrencyService = jasmine.createSpyObj('CurrencyService', ['loadCurrencies'], {
            currencies: signal([]),
        });

        mockTransactionService.loadTransactions.and.returnValue(of([]));
        mockAccountService.loadAccounts.and.returnValue(of([]));
        mockCurrencyService.loadCurrencies.and.returnValue(of([]));

        await TestBed.configureTestingModule({
            imports: [SuspiciousTransactionsComponent, NoopAnimationsModule],
            providers: [
                { provide: TransactionService, useValue: mockTransactionService },
                { provide: AccountService, useValue: mockAccountService },
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: MatDialog, useValue: {} },
                { provide: MatSnackBar, useValue: {} },
                provideRouter([]),
            ],
        }).compileComponents();

        fixture = TestBed.createComponent(SuspiciousTransactionsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should load transactions with onlySuspicious filter on init', () => {
        expect(mockTransactionService.loadTransactions).toHaveBeenCalledWith({
            onlySuspicious: true,
        });
    });
});
