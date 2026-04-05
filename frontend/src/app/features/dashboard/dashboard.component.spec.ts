import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DashboardComponent } from './dashboard.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { ApiConfiguration } from '../../core/api/api-configuration';
import { AccountService } from '../accounts/services/account.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { UserService } from '../../core/services/user.service';
import { LayoutService } from '../../layout/services/layout.service';
import { ReconciliationService } from '../reconciliation/services/reconciliation.service';
import { of } from 'rxjs';
import { signal } from '@angular/core';
import { provideRouter } from '@angular/router';

describe('DashboardComponent', () => {
    let component: DashboardComponent;
    let fixture: ComponentFixture<DashboardComponent>;
    let httpClient: jasmine.SpyObj<HttpClient>;
    let accountService: jasmine.SpyObj<AccountService>;

    let layoutService: jasmine.SpyObj<LayoutService>;
    let dialogSpy: jasmine.SpyObj<MatDialog>;
    let accountsSignal: any;

    // Golden Mock Data satisfying all tests
    const goldenMockData = {
        from: '2024-01-01',
        to: '2024-12-01',
        granularity: 'month',
        intervals: [
            '2024-01-01',
            '2024-02-01',
            '2024-03-01',
            '2024-04-01',
            '2024-05-01',
            '2024-06-01',
            '2024-07-01',
            '2024-08-01',
            '2024-09-01',
            '2024-10-01',
            '2024-11-01',
            '2024-12-01',
        ],
        currencies: [
            {
                currencyId: 'usd',
                accounts: [
                    {
                        accountId: 'asset1',
                        amounts: [1000, 200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], // 1200 total
                        total: 1200,
                        changePercent: 20,
                    },
                    {
                        accountId: 'acc1', // expense account
                        amounts: [50, 60, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                        total: 110,
                        changePercent: 20,
                    },
                    // For filtering test
                    {
                        accountId: 'visible',
                        amounts: [100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                        total: 100,
                        changePercent: 0,
                    },
                    {
                        accountId: 'hidden',
                        amounts: [200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                        total: 200,
                        changePercent: 0,
                    },
                    {
                        accountId: 'default',
                        amounts: [300, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                        total: 300,
                        changePercent: 0,
                    },
                ],
            },
        ],
    };

    beforeEach(async () => {
        const httpClientSpy = jasmine.createSpyObj('HttpClient', ['request']);
        httpClientSpy.request.and.returnValue(
            of(new HttpResponse({ body: goldenMockData as any })),
        );

        // Create writable signals
        accountsSignal = signal([]);
        const averagesSignal = signal([]);
        const currenciesSignal = signal([]);
        const userSignal = signal(null);

        // Create spy object with properties
        const accountServiceSpy = jasmine.createSpyObj(
            'AccountService',
            ['loadAccounts', 'update', 'loadYearlyExpenses'],
            {
                accounts: accountsSignal,
                averages: averagesSignal,
            },
        );

        accountServiceSpy.loadAccounts.and.returnValue(of([]));
        accountServiceSpy.loadYearlyExpenses.and.returnValue(of([]));
        accountServiceSpy.update.and.returnValue(of({} as any));

        const currencyServiceSpy = jasmine.createSpyObj('CurrencyService', ['loadCurrencies'], {
            currencies: currenciesSignal,
        });
        currencyServiceSpy.loadCurrencies.and.returnValue(
            of([{ id: 'usd', name: 'US Dollar', description: '' }] as any),
        );

        const userServiceSpy = jasmine.createSpyObj('UserService', ['loadUser'], {
            user: userSignal,
        });
        userServiceSpy.loadUser.and.returnValue(of({ favoriteCurrencyId: null } as any));

        const sidenavOpenedSignal = signal(true);
        const layoutServiceSpy = jasmine.createSpyObj('LayoutService', ['toggleSidenav']);
        layoutServiceSpy.sidenavOpened = sidenavOpenedSignal;
        layoutServiceSpy.sidenavWidth = 250;

        const apiConfigMock = { rootUrl: 'http://localhost:8080/api/v1' };

        dialogSpy = jasmine.createSpyObj('MatDialog', ['open']);

        const reconciliationStatusesSignal = signal([]);
        const reconciliationServiceSpy = jasmine.createSpyObj(
            'ReconciliationService',
            ['loadStatuses'],
            {
                statuses: reconciliationStatusesSignal,
            },
        );
        reconciliationServiceSpy.loadStatuses.and.returnValue(of([]));

        await TestBed.configureTestingModule({
            imports: [DashboardComponent],
            providers: [
                provideRouter([]),
                { provide: HttpClient, useValue: httpClientSpy },
                { provide: ApiConfiguration, useValue: apiConfigMock },
                { provide: AccountService, useValue: accountServiceSpy },
                { provide: CurrencyService, useValue: currencyServiceSpy },
                { provide: UserService, useValue: userServiceSpy },
                { provide: LayoutService, useValue: layoutServiceSpy },
                { provide: MatDialog, useValue: dialogSpy },
                { provide: ReconciliationService, useValue: reconciliationServiceSpy },
            ],
        })
            .overrideComponent(DashboardComponent, {
                remove: { imports: [MatDialogModule] },
                add: { providers: [{ provide: MatDialog, useValue: dialogSpy }] },
            })
            .compileComponents();

        httpClient = TestBed.inject(HttpClient) as jasmine.SpyObj<HttpClient>;
        accountService = TestBed.inject(AccountService) as jasmine.SpyObj<AccountService>;

        layoutService = TestBed.inject(LayoutService) as jasmine.SpyObj<LayoutService>;

        fixture = TestBed.createComponent(DashboardComponent);
        component = fixture.componentInstance;
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should load expense data on init', (done) => {
        component.ngOnInit();
        setTimeout(() => {
            expect(httpClient.request).toHaveBeenCalledTimes(2);
            expect(component['expenseData']()).toBeTruthy();
            done();
        }, 100);
    });

    it('should display loading spinner while loading', () => {
        fixture.detectChanges();
        component['loading'].set(true);
        fixture.detectChanges();
        const compiled = fixture.nativeElement;
        expect(compiled.querySelector('mat-spinner')).toBeTruthy();
    });

    it('should show all 12 months on large screens', () => {
        component['isSmallScreen'].set(false);
        component['expenseData'].set(goldenMockData as any);
        const monthColumns = component['monthColumns']();
        expect(monthColumns.length).toBe(12);
        expect(monthColumns).toEqual(goldenMockData.intervals);
    });

    it('should show only 6 months on small screens', () => {
        component['isSmallScreen'].set(true);
        component['expenseData'].set(goldenMockData as any);
        const monthColumns = component['monthColumns']();
        expect(monthColumns.length).toBe(6);
        expect(monthColumns).toEqual(goldenMockData.intervals.slice(-6));
    });

    it('should update screen size based on window width and sidenav state', () => {
        component['windowWidth'].set(2000);
        fixture.detectChanges();
        expect(component['isSmallScreen']()).toBe(false);

        component['windowWidth'].set(1600);
        fixture.detectChanges();
        expect(component['isSmallScreen']()).toBe(true);

        layoutService.sidenavOpened.set(false);
        fixture.detectChanges();
        expect(component['isSmallScreen']()).toBe(false);
    });

    it('should load asset data and create asset cards', (done) => {
        const mockAccounts = [{ id: 'asset1', name: 'My Asset', type: 'asset' }];
        accountsSignal.set(mockAccounts);
        accountService.loadAccounts.and.returnValue(of(mockAccounts as any));

        component.ngOnInit();

        setTimeout(() => {
            const cards = component['assetCards']();
            expect(cards.length).toBe(1);
            if (cards.length > 0) {
                expect(cards[0].accountName).toBe('My Asset');
                expect(cards[0].balance).toBe(1200);
            }
            done();
        }, 100);
    });

    it('should filter out asset accounts when showInDashboardSummary is false', (done) => {
        const mockAccounts = [
            { id: 'visible', name: 'Visible Asset', type: 'asset', showInDashboardSummary: true },
            { id: 'hidden', name: 'Hidden Asset', type: 'asset', showInDashboardSummary: false },
            { id: 'default', name: 'Default Asset', type: 'asset' },
        ];

        accountsSignal.set(mockAccounts);
        accountService.loadAccounts.and.returnValue(of(mockAccounts as any));

        component.ngOnInit();

        setTimeout(() => {
            const cards = component['assetCards']();
            expect(cards.length).toBe(2);
            expect(cards.find((c: any) => c.accountId === 'visible')).toBeTruthy();
            expect(cards.find((c: any) => c.accountId === 'default')).toBeTruthy();
            expect(cards.find((c: any) => c.accountId === 'hidden')).toBeFalsy();
            done();
        }, 100);
    });

    it('should call accountService.update when onHideAccount is called and dialog is confirmed', () => {
        const dialogRefSpy = jasmine.createSpyObj({ afterClosed: of(true) });
        dialogSpy.open.and.returnValue(dialogRefSpy);

        const mockAccounts = [
            { id: 'acc1', name: 'Test Asset', type: 'asset', showInDashboardSummary: true },
        ];

        accountsSignal.set(mockAccounts);
        accountService.update.and.returnValue(of({} as any));

        component['onHideAccount']('acc1');

        expect(dialogSpy.open).toHaveBeenCalled();
        expect(accountService.update).toHaveBeenCalledWith(
            'acc1',
            jasmine.objectContaining({
                showInDashboardSummary: false,
            }),
        );
    });

    describe('assetTotals', () => {
        it('should return empty array when assetData or accounts are empty', () => {
            component['assetData'].set(null);
            accountsSignal.set([]);
            expect(component['assetTotals']()).toEqual([]);

            component['assetData'].set({ intervals: [], currencies: [] } as any);
            expect(component['assetTotals']()).toEqual([]);
        });

        it('should aggregate asset totals across multiple accounts for the same currency', () => {
            const mockAccounts = [
                { id: 'asset1', name: 'Asset 1', type: 'asset' },
                { id: 'asset2', name: 'Asset 2', type: 'asset' },
            ];
            accountsSignal.set(mockAccounts);

            const mockAssetData = {
                intervals: ['2024-01-01', '2024-02-01'],
                currencies: [
                    {
                        currencyId: 'usd',
                        accounts: [
                            {
                                accountId: 'asset1',
                                amounts: [1000, 1100],
                                total: 1100,
                                changePercent: 10,
                            },
                            {
                                accountId: 'asset2',
                                amounts: [500, 600],
                                total: 600,
                                changePercent: 20,
                            },
                        ],
                    },
                ],
            };
            component['assetData'].set(mockAssetData as any);

            const totals = component['assetTotals']();
            expect(totals.length).toBe(1);
            expect(totals[0].currencyId).toBe('usd');
            expect(totals[0].totalBalance).toBe(1700); // 1100 + 600

            // Previous balances calculation:
            // asset1: 1100 / (1 + 10/100) = 1100 / 1.1 = 1000
            // asset2: 600 / (1 + 20/100) = 600 / 1.2 = 500
            // total prev: 1500
            // trend: (1700 - 1500) / 1500 * 100 = 200 / 1500 * 100 = 13.333%
            expect(totals[0].trendPercent).toBeCloseTo(13.333, 3);
            expect(totals[0].trendDirection).toBe('up');
            expect(totals[0].history).toEqual([1500, 1700]); // 1000+500, 1100+600
        });

        it('should calculate correct trend (e.g. 25% for 100/80)', () => {
            const mockAccounts = [{ id: 'asset1', name: 'Asset 1', type: 'asset' }];
            accountsSignal.set(mockAccounts);

            const mockAssetData = {
                intervals: ['2024-01-01', '2024-02-01'],
                currencies: [
                    {
                        currencyId: 'usd',
                        accounts: [
                            {
                                accountId: 'asset1',
                                amounts: [80, 100],
                                total: 100,
                                changePercent: 25, // (100-80)/80 * 100 = 25%
                            },
                        ],
                    },
                ],
            };
            component['assetData'].set(mockAssetData as any);

            const totals = component['assetTotals']();
            expect(totals[0].totalBalance).toBe(100);
            expect(totals[0].trendPercent).toBe(25);
            expect(totals[0].trendDirection).toBe('up');
        });

        it('should handle multiple currencies', () => {
            const mockAccounts = [
                { id: 'usd-acc', name: 'USD Acc', type: 'asset' },
                { id: 'eur-acc', name: 'EUR Acc', type: 'asset' },
            ];
            accountsSignal.set(mockAccounts);

            const mockAssetData = {
                intervals: ['2024-01-01'],
                currencies: [
                    {
                        currencyId: 'usd',
                        accounts: [{ accountId: 'usd-acc', amounts: [100], total: 100 }],
                    },
                    {
                        currencyId: 'eur',
                        accounts: [{ accountId: 'eur-acc', amounts: [200], total: 200 }],
                    },
                ],
            };
            component['assetData'].set(mockAssetData as any);

            const totals = component['assetTotals']();
            expect(totals.length).toBe(2);
            expect(totals.find((t) => t.currencyId === 'usd')?.totalBalance).toBe(100);
            expect(totals.find((t) => t.currencyId === 'eur')?.totalBalance).toBe(200);
        });

        it('should respect showInDashboardSummary and includeHidden', () => {
            const mockAccounts = [
                { id: 'visible', name: 'Visible', type: 'asset', showInDashboardSummary: true },
                { id: 'hidden', name: 'Hidden', type: 'asset', showInDashboardSummary: false },
            ];
            accountsSignal.set(mockAccounts);

            const mockAssetData = {
                intervals: ['2024-01-01'],
                currencies: [
                    {
                        currencyId: 'usd',
                        accounts: [
                            { accountId: 'visible', amounts: [100], total: 100 },
                            { accountId: 'hidden', amounts: [50], total: 50 },
                        ],
                    },
                ],
            };
            component['assetData'].set(mockAssetData as any);

            // By default, should only include visible
            expect(component['assetTotals']()[0].totalBalance).toBe(100);

            // When includeHidden is true, should include both
            component['includeHidden'].set(true);
            expect(component['assetTotals']()[0].totalBalance).toBe(150);
        });
    });
});
