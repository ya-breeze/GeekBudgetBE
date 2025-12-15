import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { BudgetItem } from '../../../core/api/models/budget-item';
import { BudgetItemNoId } from '../../../core/api/models/budget-item-no-id';
import { getBudgetItems } from '../../../core/api/fn/budget-items/get-budget-items';
import { createBudgetItem } from '../../../core/api/fn/budget-items/create-budget-item';
import { updateBudgetItem } from '../../../core/api/fn/budget-items/update-budget-item';

import { deleteBudgetItem } from '../../../core/api/fn/budget-items/delete-budget-item';
import { getBudgetStatus } from '../../../core/api/fn/budget-items/get-budget-status';
import { BudgetStatus } from '../../../core/api/models/budget-status';

@Injectable({
  providedIn: 'root',
})
export class BudgetItemService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly budgetItems = signal<BudgetItem[]>([]);
  readonly budgetStatus = signal<BudgetStatus[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadBudgetItems(): Observable<BudgetItem[]> {
    this.loading.set(true);
    this.error.set(null);

    return getBudgetItems(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (budgetItems) => {
          this.budgetItems.set(budgetItems);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load budget items');
          this.loading.set(false);
        },
      })
    );
  }

  create(budgetItem: BudgetItemNoId): Observable<BudgetItem> {
    this.loading.set(true);
    this.error.set(null);

    return createBudgetItem(this.http, this.apiConfig.rootUrl, { body: budgetItem }).pipe(
      map((response) => response.body),
      tap({
        next: (budgetItem) => {
          this.budgetItems.update((items) => [...items, budgetItem]);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to create budget item');
          this.loading.set(false);
        },
      })
    );
  }

  update(id: string, budgetItem: BudgetItemNoId): Observable<BudgetItem> {
    this.loading.set(true);
    this.error.set(null);

    return updateBudgetItem(this.http, this.apiConfig.rootUrl, { id, body: budgetItem }).pipe(
      map((response) => response.body),
      tap({
        next: (updatedItem) => {
          this.budgetItems.update((items) =>
            items.map((item) => (item.id === id ? updatedItem : item))
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to update budget item');
          this.loading.set(false);
        },
      })
    );
  }

  delete(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return deleteBudgetItem(this.http, this.apiConfig.rootUrl, { id }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.budgetItems.update((items) => items.filter((item) => item.id !== id));
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete budget item');
          this.loading.set(false);
        },
      })
    );
  }

  loadBudgetStatus(from?: string, to?: string, outputCurrencyId?: string): Observable<BudgetStatus[]> {
    this.loading.set(true);
    this.error.set(null);

    return getBudgetStatus(this.http, this.apiConfig.rootUrl, { from, to, outputCurrencyId }).pipe(
      map((response) => response.body),
      tap({
        next: (status) => {
          this.budgetStatus.set(status);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load budget status');
          this.loading.set(false);
        },
      })
    );
  }
}

