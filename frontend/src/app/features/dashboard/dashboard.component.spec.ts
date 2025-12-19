import { MatDialog, MatDialogModule } from '@angular/material/dialog';
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
                    },
                    {
                        accountId: 'acc1', // expense account
                        amounts: [50, 60, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
                    },
                    // For filtering test
                    { accountId: 'visible', amounts: [100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0] },
                    { accountId: 'hidden', amounts: [200, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0] },
                    { accountId: 'default', amounts: [300, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0] },
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

        await TestBed.configureTestingModule({
            imports: [DashboardComponent],
            providers: [
                { provide: HttpClient, useValue: httpClientSpy },
                { provide: ApiConfiguration, useValue: apiConfigMock },
                { provide: AccountService, useValue: accountServiceSpy },
                { provide: CurrencyService, useValue: currencyServiceSpy },
                { provide: UserService, useValue: userServiceSpy },
                { provide: LayoutService, useValue: layoutServiceSpy },
                { provide: MatDialog, useValue: dialogSpy },
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
});
