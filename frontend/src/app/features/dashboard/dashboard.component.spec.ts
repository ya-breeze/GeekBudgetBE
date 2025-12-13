import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DashboardComponent } from './dashboard.component';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { ApiConfiguration } from '../../core/api/api-configuration';
import { AccountService } from '../accounts/services/account.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { UserService } from '../../core/services/user.service';
import { LayoutService } from '../../layout/services/layout.service';
import { of } from 'rxjs';
import { signal } from '@angular/core';

describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let httpClient: jasmine.SpyObj<HttpClient>;
  let accountService: jasmine.SpyObj<AccountService>;
  let currencyService: jasmine.SpyObj<CurrencyService>;
  let userService: jasmine.SpyObj<UserService>;
  let layoutService: jasmine.SpyObj<LayoutService>;

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

    const userServiceSpy = jasmine.createSpyObj('UserService', ['loadUser'], {
      user: jasmine.createSpy('user').and.returnValue(null),
    });
    userServiceSpy.loadUser.and.returnValue(of({ favoriteCurrencyId: null } as any));

    const sidenavOpenedSignal = signal(true);
    const layoutServiceSpy = jasmine.createSpyObj('LayoutService', ['toggleSidenav']);
    layoutServiceSpy.sidenavOpened = sidenavOpenedSignal;
    layoutServiceSpy.sidenavWidth = 250;

    const apiConfigMock = { rootUrl: 'http://localhost:8080/api/v1' };

    await TestBed.configureTestingModule({
      imports: [DashboardComponent],
      providers: [
        { provide: HttpClient, useValue: httpClientSpy },
        { provide: ApiConfiguration, useValue: apiConfigMock },
        { provide: AccountService, useValue: accountServiceSpy },
        { provide: CurrencyService, useValue: currencyServiceSpy },
        { provide: UserService, useValue: userServiceSpy },
        { provide: LayoutService, useValue: layoutServiceSpy },
      ],
    }).compileComponents();

    httpClient = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
    accountService = TestBed.inject(AccountService) as jasmine.SpyObj<AccountService>;
    currencyService = TestBed.inject(CurrencyService) as jasmine.SpyObj<CurrencyService>;
    userService = TestBed.inject(UserService) as jasmine.SpyObj<UserService>;
    layoutService = TestBed.inject(LayoutService) as jasmine.SpyObj<LayoutService>;

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
      expect(httpClient.request).toHaveBeenCalledTimes(2);
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

  it('should show all 12 months on large screens', () => {
    // Mock 12 months of data
    const mockData = {
      intervals: [
        '2024-01-01', '2024-02-01', '2024-03-01', '2024-04-01',
        '2024-05-01', '2024-06-01', '2024-07-01', '2024-08-01',
        '2024-09-01', '2024-10-01', '2024-11-01', '2024-12-01'
      ],
      currencies: []
    };

    // Set up large screen (not small screen)
    component['isSmallScreen'].set(false);
    component['expenseData'].set(mockData as any);

    const monthColumns = component['monthColumns']();
    expect(monthColumns.length).toBe(12);
    expect(monthColumns).toEqual(mockData.intervals);
  });

  it('should show only 6 months on small screens', () => {
    // Mock 12 months of data
    const mockData = {
      intervals: [
        '2024-01-01', '2024-02-01', '2024-03-01', '2024-04-01',
        '2024-05-01', '2024-06-01', '2024-07-01', '2024-08-01',
        '2024-09-01', '2024-10-01', '2024-11-01', '2024-12-01'
      ],
      currencies: []
    };

    // Set up small screen
    component['isSmallScreen'].set(true);
    component['expenseData'].set(mockData as any);

    const monthColumns = component['monthColumns']();
    expect(monthColumns.length).toBe(6);
    // Should show the last 6 months
    expect(monthColumns).toEqual([
      '2024-07-01', '2024-08-01', '2024-09-01',
      '2024-10-01', '2024-11-01', '2024-12-01'
    ]);
  });

  it('should update screen size based on window width and sidenav state', () => {
    // Set initial window width to 2000
    component['windowWidth'].set(2000);
    fixture.detectChanges();

    // With sidenav open (250px), effective width = 2000 - 250 = 1750 > 1500 (not small)
    expect(component['isSmallScreen']()).toBe(false);

    // Change window width to 1600px
    // With sidenav open (250px), effective width = 1600 - 250 = 1350 <= 1500 (small)
    component['windowWidth'].set(1600);
    fixture.detectChanges();
    expect(component['isSmallScreen']()).toBe(true);

    // Test sidenav toggle: close sidenav
    // With sidenav closed (0px), effective width = 1600 - 0 = 1600 > 1500 (not small)
    layoutService.sidenavOpened.set(false);
    fixture.detectChanges();
    expect(component['isSmallScreen']()).toBe(false);
  });
  it('should load asset data and create asset cards', (done) => {
    const mockExpenses = {
      intervals: ['2023-01-01'],
      currencies: []
    };

    const mockAssets = {
      from: '2023-01-01',
      to: '2023-02-01',
      granularity: 'month',
      intervals: ['2023-01-01', '2023-02-01'],
      currencies: [
        {
          currencyId: 'usd',
          accounts: [
            {
              accountId: 'asset1',
              amounts: [1000, 200] // Balance: 1200. Prev: 1000. Trend: +200 (+20%)
            }
          ]
        }
      ]
    };

    const mockAccounts = [
      { id: 'asset1', name: 'My Asset', type: 'asset' },
      { id: 'expense1', name: 'My Expense', type: 'expense' }
    ];

    accountService.loadAccounts.and.returnValue(of(mockAccounts as any));

    // We need to differentiate the calls. 
    // Since we can't easily differentiate by arguments in this simple spy setup without more complex logic,
    // we can assume the first call is expenses and second is balances or vice versa depending on forkJoin order.
    // However, forkJoin usually subscribes concurrently.
    // Let's improve the spy to return based on URL or parameters if possible, 
    // but the spy is on 'request' method of HttpClient.
    // We can use `.and.callFake` to inspect arguments.

    httpClient.request.and.callFake((first: any, second?: any, third?: any) => {
      // Handle overloads
      let url = '';
      if (typeof first === 'string') {
        url = second;
      } else if (first && first.url) {
        url = first.url;
      }

      // getExpenses calls /v1/expenses, getBalances calls /v1/balances
      if (url && url.includes('/balances')) {
        return of(new HttpResponse({ body: mockAssets })) as any;
      }
      return of(new HttpResponse({ body: mockExpenses })) as any;
    });

    currencyService.loadCurrencies.and.returnValue(of([
      { id: 'usd', name: 'US Dollar', description: '' }
    ] as any));

    // Mock the accounts signal to return our mock accounts
    (accountService.accounts as unknown as jasmine.Spy).and.returnValue(mockAccounts);

    component.ngOnInit();

    setTimeout(() => {
      expect(component['assetData']()).toEqual(mockAssets as any);

      const cards = component['assetCards']();
      expect(cards.length).toBe(1);
      expect(cards[0].accountName).toBe('My Asset');
      expect(cards[0].balance).toBe(1200); // 1000 + 200
      expect(cards[0].trendDirection).toBe('up');
      expect(cards[0].trendPercent).toBe(20);

      done();
    }, 100);
  });
});
