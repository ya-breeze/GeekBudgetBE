import { ComponentFixture, TestBed } from '@angular/core/testing';
import { TagAnalyticsComponent } from './tag-analytics.component';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';
import { AccountService } from '../../accounts/services/account.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { of } from 'rxjs';

describe('TagAnalyticsComponent', () => {
    let component: TagAnalyticsComponent;
    let fixture: ComponentFixture<TagAnalyticsComponent>;

    beforeEach(async () => {
        const mockCurrencyService = {
            currencies: () => [],
            loadCurrencies: () => of(null),
        };
        const mockUserService = {
            loadUser: () => of({ favoriteCurrencyId: 'USD' }),
        };
        const mockAccountService = {
            accounts: () => [],
            loadAccounts: () => of(null),
        };

        await TestBed.configureTestingModule({
            imports: [TagAnalyticsComponent, NoopAnimationsModule],
            providers: [
                provideHttpClient(),
                provideHttpClientTesting(),
                ApiConfiguration,
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: UserService, useValue: mockUserService },
                { provide: AccountService, useValue: mockAccountService },
            ],
        }).compileComponents();

        fixture = TestBed.createComponent(TagAnalyticsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
