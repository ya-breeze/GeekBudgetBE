import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { BudgetItemService } from './budget-item.service';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { BudgetItem } from '../../../core/api/models/budget-item';

describe('BudgetItemService', () => {
  let service: BudgetItemService;
  let httpMock: HttpTestingController;
  let apiConfig: ApiConfiguration;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [BudgetItemService, ApiConfiguration],
    });

    service = TestBed.inject(BudgetItemService);
    httpMock = TestBed.inject(HttpTestingController);
    apiConfig = TestBed.inject(ApiConfiguration);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('loadBudgetItems', () => {
    it('should load all budget items', (done) => {
      const mockBudgetItems: BudgetItem[] = [
        {
          id: '1',
          date: '2024-01-01',
          accountId: 'acc1',
          amount: 1000,
          description: 'Monthly budget',
        },
      ];

      service.loadBudgetItems().subscribe({
        next: (items) => {
          expect(items).toEqual(mockBudgetItems);
          expect(service.budgetItems()).toEqual(mockBudgetItems);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/budgetItems`);
      expect(req.request.method).toBe('GET');
      req.flush(mockBudgetItems);
    });
  });

  describe('create', () => {
    it('should create a budget item', (done) => {
      const newBudgetItem = {
        date: '2024-02-01',
        accountId: 'acc1',
        amount: 1500,
        description: 'February budget',
      };
      const createdBudgetItem: BudgetItem = { id: '2', ...newBudgetItem };

      service.create(newBudgetItem).subscribe({
        next: (item) => {
          expect(item).toEqual(createdBudgetItem);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/budgetItems`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(newBudgetItem);
      req.flush(createdBudgetItem);
    });
  });

  describe('update', () => {
    it('should update a budget item', (done) => {
      const updatedBudgetItem: BudgetItem = {
        id: '1',
        date: '2024-01-01',
        accountId: 'acc1',
        amount: 1200,
        description: 'Updated budget',
      };

      service.update('1', updatedBudgetItem).subscribe({
        next: (item) => {
          expect(item).toEqual(updatedBudgetItem);
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/budgetItems/1`);
      expect(req.request.method).toBe('PUT');
      req.flush(updatedBudgetItem);
    });
  });

  describe('delete', () => {
    it('should delete a budget item', (done) => {
      service.delete('1').subscribe({
        next: () => {
          done();
        },
      });

      const req = httpMock.expectOne(`${apiConfig.rootUrl}/v1/budgetItems/1`);
      expect(req.request.method).toBe('DELETE');
      req.flush(null);
    });
  });
});
