import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DashboardComponent } from './dashboard.component';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { ApiConfiguration } from '../../core/api/api-configuration';
import { AccountService } from '../accounts/services/account.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { of } from 'rxjs';

describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let httpClient: jasmine.SpyObj<HttpClient>;
  let accountService: jasmine.SpyObj<AccountService>;
  let currencyService: jasmine.SpyObj<CurrencyService>;

  beforeEach(async () => {
    const httpClientSpy = jasmine.createSpyObj('HttpClient', ['request']);
    httpClientSpy.request.and.returnValue(
      of(new HttpResponse({ body: { intervals: [], currencies: [] } as any }))
    );

    const accountServiceSpy = jasmine.createSpyObj('AccountService', ['loadAccounts'], {
      accounts: jasmine.createSpy('accounts').and.returnValue([]),
    });
    accountServiceSpy.loadAccounts.and.returnValue(of([]));

    const currencyServiceSpy = jasmine.createSpyObj('CurrencyService', ['loadCurrencies'], {
      currencies: jasmine.createSpy('currencies').and.returnValue([]),
    });
    currencyServiceSpy.loadCurrencies.and.returnValue(of([]));

    const apiConfigMock = { rootUrl: 'http://localhost:8080/api/v1' };

    await TestBed.configureTestingModule({
      imports: [DashboardComponent],
      providers: [
        { provide: HttpClient, useValue: httpClientSpy },
        { provide: ApiConfiguration, useValue: apiConfigMock },
        { provide: AccountService, useValue: accountServiceSpy },
        { provide: CurrencyService, useValue: currencyServiceSpy },
      ],
    }).compileComponents();

    httpClient = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
    accountService = TestBed.inject(AccountService) as jasmine.SpyObj<AccountService>;
    currencyService = TestBed.inject(CurrencyService) as jasmine.SpyObj<CurrencyService>;

    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load expense data on init', (done) => {
    const mockExpensesWithIntervals = {
      body: {
        from: '2024-01-01',
        to: '2024-06-01',
        granularity: 'month',
        intervals: ['2024-01-01', '2024-02-01'],
        currencies: [
          {
            currencyId: 'usd',
            accounts: [
              {
                accountId: 'acc1',
                amounts: [50, 60],
              },
            ],
          },
        ],
      },
    };

    accountService.loadAccounts.and.returnValue(of([]));
    currencyService.loadCurrencies.and.returnValue(of([]));
    httpClient.request.and.returnValue(of(new HttpResponse(mockExpensesWithIntervals as any)));

    component.ngOnInit();

    // Wait for async operations to complete
    setTimeout(() => {
      expect(httpClient.request).toHaveBeenCalledTimes(1);
      expect(component['expenseData']()).toBeTruthy();
      done();
    }, 100);
  });

  it('should display loading spinner while loading', () => {
    // Run initial change detection to trigger ngOnInit
    fixture.detectChanges();

    // Manually set loading to true and check that spinner is shown
    component['loading'].set(true);
    fixture.detectChanges();

    const compiled = fixture.nativeElement;
    const spinner = compiled.querySelector('mat-spinner');
    expect(spinner).toBeTruthy();
  });
});
