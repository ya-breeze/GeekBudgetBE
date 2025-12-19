import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { CurrencyService } from './currency.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Currency } from '../../../core/api/models/currency';

describe('CurrencyService', () => {
    let service: CurrencyService;
    let httpMock: HttpTestingController;
    let apiConfig: ApiConfiguration;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [CurrencyService, ApiConfiguration],
        });

        service = TestBed.inject(CurrencyService);
        httpMock = TestBed.inject(HttpTestingController);
        apiConfig = TestBed.inject(ApiConfiguration);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    describe('loadCurrencies', () => {
        it('should load all currencies', (done) => {
            const mockCurrencies: Currency[] = [
                { id: '1', name: 'USD', description: 'US Dollar' },
                { id: '2', name: 'EUR', description: 'Euro' },
            ];

            service.loadCurrencies().subscribe({
                next: (currencies) => {
                    expect(currencies).toEqual(mockCurrencies);
                    expect(service.currencies()).toEqual(mockCurrencies);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies`);
            expect(req.request.method).toBe('GET');
            req.flush(mockCurrencies);
        });

        it('should update signal state on successful load', (done) => {
            const mockCurrencies: Currency[] = [{ id: '1', name: 'USD', description: 'US Dollar' }];

            expect(service.loading()).toBe(false);

            service.loadCurrencies().subscribe({
                next: () => {
                    expect(service.currencies()).toEqual(mockCurrencies);
                    expect(service.loading()).toBe(false);
                    expect(service.error()).toBeNull();
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies`);
            req.flush(mockCurrencies);
        });

        it('should set error state on failed load', (done) => {
            const mockError = { status: 500, statusText: 'Server Error' };

            service.loadCurrencies().subscribe({
                error: () => {
                    expect(service.error()).toBeTruthy();
                    expect(service.loading()).toBe(false);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies`);
            req.flush('Server Error', mockError);
        });

        it('should set loading state during API calls', () => {
            service.loadCurrencies().subscribe();

            // Loading should be set before request completes
            expect(service.loading()).toBe(true);

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies`);
            req.flush({ body: [] });

            // Loading should be false after request completes
            expect(service.loading()).toBe(false);
        });
    });

    describe('create', () => {
        it('should create a new currency', (done) => {
            const newCurrency = { name: 'GBP', description: 'British Pound' };
            const createdCurrency: Currency = { id: '3', ...newCurrency };

            service.create(newCurrency).subscribe({
                next: (currency) => {
                    expect(currency).toEqual(createdCurrency);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies`);
            expect(req.request.method).toBe('POST');
            expect(req.request.body).toEqual(newCurrency);
            req.flush(createdCurrency);
        });
    });

    describe('update', () => {
        it('should update an existing currency', (done) => {
            const updatedCurrency: Currency = {
                id: '1',
                name: 'USD',
                description: 'United States Dollar',
            };

            service.update('1', updatedCurrency).subscribe({
                next: (currency) => {
                    expect(currency).toEqual(updatedCurrency);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies/1`);
            expect(req.request.method).toBe('PUT');
            expect(req.request.body).toEqual(updatedCurrency);
            req.flush(updatedCurrency);
        });
    });

    describe('delete', () => {
        it('should delete a currency', (done) => {
            service.delete('1').subscribe({
                next: () => {
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/currencies/1`);
            expect(req.request.method).toBe('DELETE');
            req.flush(null);
        });
    });
});
