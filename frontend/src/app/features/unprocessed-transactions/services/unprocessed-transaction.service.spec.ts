import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { UnprocessedTransactionService } from './unprocessed-transaction.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { UnprocessedTransaction } from '../../../core/api/models/unprocessed-transaction';

describe('UnprocessedTransactionService', () => {
  let service: UnprocessedTransactionService;
  let httpMock: HttpTestingController;
  let apiConfig: ApiConfiguration;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [UnprocessedTransactionService, ApiConfiguration],
    });

    service = TestBed.inject(UnprocessedTransactionService);
    httpMock = TestBed.inject(HttpTestingController);
    apiConfig = TestBed.inject(ApiConfiguration);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('loadUnprocessedTransactions', () => {
    it('should load unprocessed transactions', (done) => {
      const mockTransactions: UnprocessedTransaction[] = [
        {
          transaction: {
            id: '1',
            date: '2024-01-01',
            description: 'Imported transaction',
            movements: [
              { accountId: 'acc1', amount: -100, currencyId: 'usd' },
              { accountId: 'acc2', amount: 100, currencyId: 'usd' },
            ],
          },
          matched: [],
          duplicates: [],
        },
      ];

      service.loadUnprocessedTransactions().subscribe({
        next: (transactions) => {
          expect(transactions).toEqual(mockTransactions);
          expect(service.unprocessedTransactions()).toEqual(mockTransactions);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/unprocessedTransactions`);
      expect(req.request.method).toBe('GET');
      req.flush(mockTransactions);
    });

    it('should handle duplicate detection', (done) => {
      const mockTransactions: UnprocessedTransaction[] = [
        {
          transaction: {
            id: '1',
            date: '2024-01-01',
            description: 'Transaction with duplicates',
            movements: [
              { accountId: 'acc1', amount: -100, currencyId: 'usd' },
              { accountId: 'acc2', amount: 100, currencyId: 'usd' },
            ],
          },
          matched: [],
          duplicates: [
            {
              id: 'dup1',
              date: '2024-01-01',
              description: 'Duplicate 1',
              movements: [
                { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                { accountId: 'acc2', amount: 100, currencyId: 'usd' },
              ],
            },
            {
              id: 'dup2',
              date: '2024-01-01',
              description: 'Duplicate 2',
              movements: [
                { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                { accountId: 'acc2', amount: 100, currencyId: 'usd' },
              ],
            },
          ],
        },
      ];

      service.loadUnprocessedTransactions().subscribe({
        next: (transactions) => {
          expect(transactions[0].duplicates.length).toBe(2);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/unprocessedTransactions`);
      req.flush(mockTransactions);
    });

    it('should handle matcher suggestions', (done) => {
      const mockTransactions: UnprocessedTransaction[] = [
        {
          transaction: {
            id: '1',
            date: '2024-01-01',
            description: 'Transaction with matches',
            movements: [
              { accountId: 'acc1', amount: -100, currencyId: 'usd' },
              { accountId: 'acc2', amount: 100, currencyId: 'usd' },
            ],
          },
          matched: [
            {
              matcherId: 'matcher1',
              transaction: {
                date: '2024-01-01',
                description: 'Matched transaction 1',
                movements: [
                  { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                  { accountId: 'acc2', amount: 100, currencyId: 'usd' },
                ],
              },
            },
            {
              matcherId: 'matcher2',
              transaction: {
                date: '2024-01-01',
                description: 'Matched transaction 2',
                movements: [
                  { accountId: 'acc1', amount: -100, currencyId: 'usd' },
                  { accountId: 'acc2', amount: 100, currencyId: 'usd' },
                ],
              },
            },
          ],
          duplicates: [],
        },
      ];

      service.loadUnprocessedTransactions().subscribe({
        next: (transactions) => {
          expect(transactions[0].matched.length).toBe(2);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/unprocessedTransactions`);
      req.flush(mockTransactions);
    });
  });

  describe('convert', () => {
    it('should convert unprocessed transaction to transaction', (done) => {
      const unprocessedTransaction: UnprocessedTransaction = {
        transaction: {
          id: '1',
          date: '2024-01-01',
          description: 'Converted transaction',
          movements: [
            { accountId: 'acc1', amount: -100, currencyId: 'usd' },
            { accountId: 'acc2', amount: 100, currencyId: 'usd' },
          ],
        },
        matched: [],
        duplicates: [],
      };

      // Set initial unprocessed transactions
      service.unprocessedTransactions.set([unprocessedTransaction]);

      service.convert('1', unprocessedTransaction).subscribe({
        next: () => {
          // Verify transaction was removed from the list
          expect(service.unprocessedTransactions().length).toBe(0);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/unprocessedTransactions/1/convert`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual({
        date: '2024-01-01',
        description: 'Converted transaction',
        movements: [
          { accountId: 'acc1', amount: -100, currencyId: 'usd' },
          { accountId: 'acc2', amount: 100, currencyId: 'usd' },
        ],
        tags: undefined,
        partnerName: undefined,
        partnerAccount: undefined,
        partnerInternalId: undefined,
        place: undefined,
        extra: undefined,
        externalIds: undefined,
        unprocessedSources: undefined,
      });
      req.flush(unprocessedTransaction.transaction);
    });
  });

  describe('delete', () => {
    it('should delete unprocessed transaction', (done) => {
      service.delete('1').subscribe({
        next: () => {
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/unprocessedTransactions/1`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });
});
