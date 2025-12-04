import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AccountService } from './account.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Account } from '../../../core/api/models/account';

describe('AccountService', () => {
  let service: AccountService;
  let httpMock: HttpTestingController;
  let apiConfig: ApiConfiguration;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AccountService, ApiConfiguration],
    });

    service = TestBed.inject(AccountService);
    httpMock = TestBed.inject(HttpTestingController);
    apiConfig = TestBed.inject(ApiConfiguration);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('loadAccounts', () => {
    it('should load all accounts', (done) => {
      const mockAccounts: Account[] = [
        { id: '1', name: 'Checking', type: 'asset' },
        { id: '2', name: 'Groceries', type: 'expense' },
      ];

      service.loadAccounts().subscribe({
        next: (accounts) => {
          expect(accounts).toEqual(mockAccounts);
          expect(service.accounts()).toEqual(mockAccounts);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts`);
      expect(req.request.method).toBe('GET');
      req.flush(mockAccounts);
    });

    it('should handle different account types', (done) => {
      const mockAccounts: Account[] = [
        { id: '1', name: 'Checking', type: 'asset' },
        { id: '2', name: 'Groceries', type: 'expense' },
        { id: '3', name: 'Salary', type: 'income' },
      ];

      service.loadAccounts().subscribe({
        next: (accounts) => {
          expect(accounts.length).toBe(3);
          expect(accounts[0].type).toBe('asset');
          expect(accounts[1].type).toBe('expense');
          expect(accounts[2].type).toBe('income');
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts`);
      req.flush(mockAccounts);
    });

    it('should update signal state correctly', (done) => {
      const mockAccounts: Account[] = [{ id: '1', name: 'Test', type: 'asset' }];

      service.loadAccounts().subscribe({
        next: () => {
          expect(service.accounts()).toEqual(mockAccounts);
          expect(service.loading()).toBe(false);
          expect(service.error()).toBeNull();
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts`);
      req.flush(mockAccounts);
    });

    it('should handle API errors gracefully', (done) => {
      const mockError = { status: 500, statusText: 'Server Error' };

      service.loadAccounts().subscribe({
        error: () => {
          expect(service.error()).toBeTruthy();
          expect(service.loading()).toBe(false);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts`);
      req.flush('Server Error', mockError);
    });
  });

  describe('create', () => {
    it('should create a new account', (done) => {
      const newAccount = { name: 'Savings', type: 'asset' as const };
      const createdAccount: Account = { id: '3', ...newAccount };

      service.create(newAccount).subscribe({
        next: (account) => {
          expect(account).toEqual(createdAccount);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(newAccount);
      req.flush(createdAccount);
    });
  });

  describe('update', () => {
    it('should update an existing account', (done) => {
      const updatedAccount: Account = {
        id: '1',
        name: 'Updated Checking',
        type: 'asset',
      };

      service.update('1', updatedAccount).subscribe({
        next: (account) => {
          expect(account).toEqual(updatedAccount);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts/1`);
      expect(req.request.method).toBe('PUT');
      req.flush(updatedAccount);
    });
  });

  describe('delete', () => {
    it('should delete an account', (done) => {

      service.delete('1').subscribe({
        next: () => {
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/accounts/1`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });
});
