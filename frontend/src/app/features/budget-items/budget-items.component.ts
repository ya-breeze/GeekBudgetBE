import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe } from '@angular/common';
import { BudgetItemService } from './services/budget-item.service';
import { BudgetItem } from '../../core/api/models/budget-item';

@Component({
  selector: 'app-budget-items',
  imports: [
    MatTableModule,
    MatSortModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    DatePipe,
  ],
  template: `
    <div class="budget-items-container">
      <div class="header">
        <h1>Budget Items</h1>
        <p class="subtitle">Planned budget allocations</p>
      </div>

      @if (loading()) {
      <div class="loading-container">
        <mat-spinner></mat-spinner>
      </div>
      } @else {
      <div class="table-container">
        <table
          mat-table
          [dataSource]="sortedBudgetItems()"
          matSort
          (matSortChange)="onSortChange($event)"
          [matSortActive]="sortActive() ?? displayedColumns()[0]"
          [matSortDirection]="sortDirection()"
          class="budget-items-table"
        >
          <ng-container matColumnDef="date">
            <th mat-header-cell *matHeaderCellDef mat-sort-header>Date</th>
            <td mat-cell *matCellDef="let item">{{ item.date | date : 'short' }}</td>
          </ng-container>

          <ng-container matColumnDef="account">
            <th mat-header-cell *matHeaderCellDef mat-sort-header>Account</th>
            <td mat-cell *matCellDef="let item">{{ item.accountId }}</td>
          </ng-container>

          <ng-container matColumnDef="amount">
            <th mat-header-cell *matHeaderCellDef mat-sort-header>Amount</th>
            <td mat-cell *matCellDef="let item">{{ item.amount }}</td>
          </ng-container>

          <ng-container matColumnDef="description">
            <th mat-header-cell *matHeaderCellDef mat-sort-header>Description</th>
            <td mat-cell *matCellDef="let item">{{ item.description || '-' }}</td>
          </ng-container>

          <tr mat-header-row *matHeaderRowDef="displayedColumns()"></tr>
          <tr mat-row *matRowDef="let row; columns: displayedColumns()"></tr>

          <tr class="mat-row" *matNoDataRow>
            <td class="mat-cell" [attr.colspan]="displayedColumns().length">
              No budget items found.
            </td>
          </tr>
        </table>
      </div>
      }
    </div>
  `,
  styles: `
    .budget-items-container {
      padding: 24px;
    }
    .header {
      margin-bottom: 24px;
      h1 {
        margin: 0 0 8px 0;
      }
      .subtitle {
        margin: 0;
        color: rgba(0, 0, 0, 0.6);
      }
    }
    .loading-container {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 300px;
    }
    .table-container {
      background: white;
      border-radius: 4px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      overflow: hidden;
    }
    .budget-items-table {
      width: 100%;
      th {
        font-weight: 600;
      }
      .mat-row:hover {
        background-color: rgba(0, 0, 0, 0.04);
      }
    }
  `,
})
export class BudgetItemsComponent implements OnInit {
  private readonly budgetItemService = inject(BudgetItemService);
  private readonly snackBar = inject(MatSnackBar);

  protected readonly budgetItems = this.budgetItemService.budgetItems;
  protected readonly loading = this.budgetItemService.loading;
  protected readonly displayedColumns = signal(['date', 'account', 'amount', 'description']);

  protected readonly sortActive = signal<string | null>(null);
  protected readonly sortDirection = signal<'asc' | 'desc'>('asc');
  protected readonly sortedBudgetItems = computed(() => {
    const data = this.budgetItems();
    const columns = this.displayedColumns();

    if (!columns.length) {
      return data;
    }

    const active = this.sortActive() ?? columns[0];
    const direction = this.sortDirection();

    return [...data].sort((a, b) => this.compareBudgetItems(a, b, active, direction));
  });

  ngOnInit(): void {
    this.budgetItemService.loadBudgetItems().subscribe();
  }

  protected onSortChange(sort: Sort): void {
    if (!sort.direction) {
      this.sortActive.set(null);
      this.sortDirection.set('asc');
      return;
    }

    this.sortActive.set(sort.active);
    this.sortDirection.set(sort.direction);
  }

  private compareBudgetItems(
    a: BudgetItem,
    b: BudgetItem,
    active: string,
    direction: 'asc' | 'desc'
  ): number {
    const valueA = this.getBudgetItemSortValue(a, active);
    const valueB = this.getBudgetItemSortValue(b, active);
    return this.comparePrimitiveValues(valueA, valueB, direction);
  }

  private getBudgetItemSortValue(item: BudgetItem, active: string): string | number | Date | null {
    switch (active) {
      case 'date':
        return item.date ? new Date(item.date) : null;
      case 'account':
        return item.accountId ?? '';
      case 'amount':
        return item.amount ?? 0;
      case 'description':
        return item.description ?? '';
      default:
        return null;
    }
  }

  private comparePrimitiveValues(
    a: string | number | Date | null | undefined,
    b: string | number | Date | null | undefined,
    direction: 'asc' | 'desc'
  ): number {
    const factor = direction === 'asc' ? 1 : -1;

    if (a == null && b == null) return 0;
    if (a == null) return 1 * factor;
    if (b == null) return -1 * factor;

    if (typeof a === 'string' && typeof b === 'string') {
      return a.localeCompare(b) * factor;
    }

    if (typeof a === 'number' && typeof b === 'number') {
      return (a - b) * factor;
    }

    if (a instanceof Date && b instanceof Date) {
      return (a.getTime() - b.getTime()) * factor;
    }

    return `${a}`.localeCompare(`${b}`) * factor;
  }
}
