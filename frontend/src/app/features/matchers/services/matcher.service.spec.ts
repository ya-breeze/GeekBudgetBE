import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { MatcherService } from './matcher.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Matcher } from '../../../core/api/models/matcher';

describe('MatcherService', () => {
    let service: MatcherService;
    let httpMock: HttpTestingController;
    let apiConfig: ApiConfiguration;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [MatcherService, ApiConfiguration],
        });

        service = TestBed.inject(MatcherService);
        httpMock = TestBed.inject(HttpTestingController);
        apiConfig = TestBed.inject(ApiConfiguration);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    describe('loadMatchers', () => {
        it('should load all matchers', (done) => {
            const mockMatchers: Matcher[] = [
                {
                    id: '1',
                    descriptionRegExp: 'GROCERY.*',
                    outputAccountId: 'acc1',
                    outputDescription: 'Groceries',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
            ];

            service.loadMatchers().subscribe({
                next: (matchers) => {
                    expect(matchers).toEqual(mockMatchers);
                    expect(service.matchers()).toEqual(mockMatchers);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers`);
            expect(req.request.method).toBe('GET');
            req.flush(mockMatchers);
        });

        it('should validate regex patterns', (done) => {
            const mockMatchers: Matcher[] = [
                {
                    id: '1',
                    descriptionRegExp: '^GROCERY.*$',
                    outputAccountId: 'acc1',
                    outputDescription: 'Groceries',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
                {
                    id: '2',
                    descriptionRegExp: 'SALARY',
                    outputAccountId: 'acc2',
                    outputDescription: 'Salary',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
            ];

            service.loadMatchers().subscribe({
                next: (matchers) => {
                    expect(matchers[0].descriptionRegExp).toBe('^GROCERY.*$');
                    expect(matchers[1].descriptionRegExp).toBe('SALARY');
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers`);
            req.flush(mockMatchers);
        });

        it('should handle multiple regex patterns', (done) => {
            const mockMatchers: Matcher[] = [
                {
                    id: '1',
                    descriptionRegExp: 'PATTERN1',
                    partnerNameRegExp: 'PARTNER1',
                    currencyRegExp: 'USD',
                    outputAccountId: 'acc1',
                    outputDescription: 'Output 1',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
                {
                    id: '2',
                    descriptionRegExp: 'PATTERN2',
                    partnerAccountNumberRegExp: '123456',
                    outputAccountId: 'acc2',
                    outputDescription: 'Output 2',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
                {
                    id: '3',
                    extraRegExp: 'EXTRA',
                    outputAccountId: 'acc3',
                    outputDescription: 'Output 3',
                    confirmationsCount: 0,
                    confirmationsTotal: 0,
                },
            ];

            service.loadMatchers().subscribe({
                next: (matchers) => {
                    expect(matchers[0].descriptionRegExp).toBe('PATTERN1');
                    expect(matchers[1].partnerAccountNumberRegExp).toBe('123456');
                    expect(matchers[2].extraRegExp).toBe('EXTRA');
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers`);
            req.flush(mockMatchers);
        });
    });

    describe('create', () => {
        it('should create a matcher', (done) => {
            const newMatcher = {
                descriptionRegExp: 'NEW_PATTERN',
                outputAccountId: 'acc1',
                outputDescription: 'New Output',
            };
            const createdMatcher: Matcher = {
                id: '4',
                ...newMatcher,
                confirmationsCount: 0,
                confirmationsTotal: 0,
            };

            service.create(newMatcher).subscribe({
                next: (matcher) => {
                    expect(matcher).toEqual(createdMatcher);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers`);
            expect(req.request.method).toBe('POST');
            expect(req.request.body).toEqual(newMatcher);
            req.flush(createdMatcher);
        });
    });

    describe('update', () => {
        it('should update a matcher', (done) => {
            const updatedMatcher: Matcher = {
                id: '1',
                descriptionRegExp: 'UPDATED_PATTERN',
                outputAccountId: 'acc2',
                outputDescription: 'Updated Output',
                confirmationsCount: 0,
                confirmationsTotal: 0,
            };

            service.update('1', updatedMatcher).subscribe({
                next: (matcher) => {
                    expect(matcher).toEqual(updatedMatcher);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers/1`);
            expect(req.request.method).toBe('PUT');
            req.flush({ matcher: updatedMatcher, autoProcessedIds: [] });
        });
    });

    describe('delete', () => {
        it('should delete a matcher', (done) => {
            service.delete('1').subscribe({
                next: () => {
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/matchers/1`);
            expect(req.request.method).toBe('DELETE');
            req.flush(null);
        });
    });
});
