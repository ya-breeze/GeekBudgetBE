import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { BankImporterService } from './bank-importer.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { BankImporter } from '../../../core/api/models/bank-importer';

describe('BankImporterService', () => {
    let service: BankImporterService;
    let httpMock: HttpTestingController;
    let apiConfig: ApiConfiguration;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [BankImporterService, ApiConfiguration],
        });

        service = TestBed.inject(BankImporterService);
        httpMock = TestBed.inject(HttpTestingController);
        apiConfig = TestBed.inject(ApiConfiguration);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    describe('loadBankImporters', () => {
        it('should load all bank importers', (done) => {
            const mockImporters: BankImporter[] = [
                {
                    id: '1',
                    name: 'FIO Bank',
                    type: 'fio',
                    accountId: 'acc1',
                    extra: 'api-token-123',
                },
            ];

            service.loadBankImporters().subscribe({
                next: (importers) => {
                    expect(importers).toEqual(mockImporters);
                    expect(service.bankImporters()).toEqual(mockImporters);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters`);
            expect(req.request.method).toBe('GET');
            req.flush(mockImporters);
        });

        it('should handle different bank types', (done) => {
            const mockImporters: BankImporter[] = [
                { id: '1', name: 'FIO', type: 'fio', accountId: 'acc1', extra: 'token1' },
                { id: '2', name: 'KB', type: 'kb', accountId: 'acc2', extra: 'token2' },
                { id: '3', name: 'Revolut', type: 'revolut', accountId: 'acc3', extra: 'token3' },
            ];

            service.loadBankImporters().subscribe({
                next: (importers) => {
                    expect(importers.length).toBe(3);
                    expect(importers[0].type).toBe('fio');
                    expect(importers[1].type).toBe('kb');
                    expect(importers[2].type).toBe('revolut');
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters`);
            req.flush(mockImporters);
        });

        it('should validate extra data field', (done) => {
            const mockImporters: BankImporter[] = [
                {
                    id: '1',
                    name: 'FIO Bank',
                    type: 'fio',
                    accountId: 'acc1',
                    extra: 'api-token-123',
                },
            ];

            service.loadBankImporters().subscribe({
                next: (importers) => {
                    expect(importers[0].extra).toBeTruthy();
                    expect(importers[0].extra).toBe('api-token-123');
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters`);
            req.flush(mockImporters);
        });
    });

    describe('create', () => {
        it('should create a bank importer', (done) => {
            const newImporter = {
                name: 'My FIO Bank',
                type: 'fio' as const,
                accountId: 'acc1',
                extra: 'my-api-token',
            };
            const createdImporter: BankImporter = { id: '2', ...newImporter };

            service.create(newImporter).subscribe({
                next: (importer) => {
                    expect(importer).toEqual(createdImporter);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters`);
            expect(req.request.method).toBe('POST');
            expect(req.request.body).toEqual(newImporter);
            req.flush(createdImporter);
        });
    });

    describe('update', () => {
        it('should update a bank importer', (done) => {
            const updatedImporter: BankImporter = {
                id: '1',
                name: 'Updated FIO',
                type: 'fio',
                accountId: 'acc1',
                extra: 'new-token',
            };

            service.update('1', updatedImporter).subscribe({
                next: (importer) => {
                    expect(importer).toEqual(updatedImporter);
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters/1`);
            expect(req.request.method).toBe('PUT');
            req.flush(updatedImporter);
        });
    });

    describe('delete', () => {
        it('should delete a bank importer', (done) => {
            service.delete('1').subscribe({
                next: () => {
                    done();
                },
            });

            const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/bankImporters/1`);
            expect(req.request.method).toBe('DELETE');
            req.flush(null);
        });
    });
});
