import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { TransactionService } from './transaction.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Transaction } from '../../../core/api/models/transaction';

describe('TransactionService', () => {
  let service: TransactionService;
  let httpMock: HttpTestingController;
  let apiConfig: ApiConfiguration;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [TransactionService, ApiConfiguration],
    });

    service = TestBed.inject(TransactionService);
    httpMock = TestBed.inject(HttpTestingController);
    apiConfig = TestBed.inject(ApiConfiguration);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('loadTransactions', () => {
    it('should load all transactions', (done) => {
      const mockTransactions: Transaction[] = [
        {
          id: '1',
          date: '2024-01-01',
          description: 'Grocery shopping',
          movements: [
            { accountId: 'acc1', amount: -50, currencyId: 'usd' },
            { accountId: 'acc2', amount: 50, currencyId: 'usd' },
          ],
        },
      ];

      service.loadTransactions().subscribe({
        next: (transactions) => {
          expect(transactions).toEqual(mockTransactions);
          expect(service.transactions()).toEqual(mockTransactions);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/transactions`);
      expect(req.request.method).toBe('GET');
      req.flush(mockTransactions);
    });

    it('should handle multiple movements correctly', (done) => {
      const mockTransactions: Transaction[] = [
        {
          id: '1',
          date: '2024-01-01',
          description: 'Split payment',
          movements: [
            { accountId: 'acc1', amount: -100, currencyId: 'usd' },
            { accountId: 'acc2', amount: 50, currencyId: 'usd' },
            { accountId: 'acc3', amount: 50, currencyId: 'usd' },
          ],
        },
      ];

      service.loadTransactions().subscribe({
        next: (transactions) => {
          expect(transactions[0].movements.length).toBe(3);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/transactions`);
      req.flush(mockTransactions);
    });
  });

  describe('create', () => {
    it('should create a transaction with movements', (done) => {
      const newTransaction = {
        date: '2024-01-15',
        description: 'Test transaction',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
        ],
      };
      const createdTransaction: Transaction = { id: '2', ...newTransaction };

      service.create(newTransaction).subscribe({
        next: (transaction) => {
          expect(transaction).toEqual(createdTransaction);
          expect(transaction.movements.length).toBe(2);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/transactions`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(newTransaction);
      req.flush(createdTransaction);
    });
  });

  describe('update', () => {
    it('should update a transaction', (done) => {
      const updatedTransaction: Transaction = {
        id: '1',
        date: '2024-01-15',
        description: 'Updated transaction',
        movements: [
          { accountId: 'acc1', amount: -150, currencyId: 'usd' },
          { accountId: 'acc2', amount: 150, currencyId: 'usd' },
        ],
      };

      service.update('1', updatedTransaction).subscribe({
        next: (transaction) => {
          expect(transaction).toEqual(updatedTransaction);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/transactions/1`);
      expect(req.request.method).toBe('PUT');
      req.flush(updatedTransaction);
    });
  });

  describe('delete', () => {
    it('should delete a transaction', (done) => {
      service.delete('1').subscribe({
        next: () => {
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/transactions/1`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });
});
