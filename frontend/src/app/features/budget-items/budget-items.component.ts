import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe, DecimalPipe } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { BudgetItemService } from './services/budget-item.service';
import { AccountService } from '../accounts/services/account.service';
import { BudgetItem } from '../../core/api/models/budget-item';
import { BudgetStatus } from '../../core/api/models/budget-status';
import { LayoutService } from '../../layout/services/layout.service';

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
    DecimalPipe,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
  ],
  template: `
    <div class="budget-items-container">
      @if (!sidenavOpened()) {
        <h1 class="page-title">Budget Items</h1>
      }

      @if (loading()) {
      <div class="loading-container">
        <mat-spinner></mat-spinner>
      </div>
      } @else {
      <div class="form-container">
        <h2>Add Budget Item</h2>
        <form [formGroup]="addBudgetForm" (ngSubmit)="addBudgetItem()" class="add-budget-form">
          <mat-form-field appearance="outline">
            <mat-label>Account</mat-label>
            <mat-select formControlName="accountId">
              @for (account of accounts(); track account.id) {
                <mat-option [value]="account.id">{{ account.name }}</mat-option>
              }
            </mat-select>
          </mat-form-field>

          <mat-form-field appearance="outline">
            <mat-label>Amount</mat-label>
            <input matInput type="number" formControlName="amount" min="0">
          </mat-form-field>

          <mat-form-field appearance="outline">
            <mat-label>Date</mat-label>
            <input matInput [matDatepicker]="picker" formControlName="date">
            <mat-datepicker-toggle matIconSuffix [for]="picker"></mat-datepicker-toggle>
            <mat-datepicker #picker></mat-datepicker>
          </mat-form-field>

          <button mat-flat-button color="primary" type="submit" [disabled]="addBudgetForm.invalid || creating()">
            @if (creating()) {
              <mat-spinner diameter="20"></mat-spinner>
            } @else {
              Add
            }
          </button>
        </form>
      </div>

      <div class="table-container">
        <h2>Budget Configuration</h2>
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

      <div class="table-container" style="margin-top: 20px;">
        <h2>Budget Status (vs Actuals)</h2>
        <table
          mat-table
          [dataSource]="budgetStatus()"
          class="budget-items-table"
        >
          <ng-container matColumnDef="date">
            <th mat-header-cell *matHeaderCellDef>Month</th>
            <td mat-cell *matCellDef="let item">{{ item.date | date : 'MM/yyyy' }}</td>
          </ng-container>

          <ng-container matColumnDef="account">
            <th mat-header-cell *matHeaderCellDef>Account</th>
            <td mat-cell *matCellDef="let item">{{ item.accountId }}</td>
          </ng-container>

          <ng-container matColumnDef="budgeted">
            <th mat-header-cell *matHeaderCellDef>Budgeted</th>
            <td mat-cell *matCellDef="let item">{{ item.budgeted | number:'1.2-2' }}</td>
          </ng-container>

          <ng-container matColumnDef="spent">
            <th mat-header-cell *matHeaderCellDef>Spent</th>
            <td mat-cell *matCellDef="let item">{{ item.spent | number:'1.2-2' }}</td>
          </ng-container>

          <ng-container matColumnDef="rollover">
            <th mat-header-cell *matHeaderCellDef>Rollover</th>
            <td mat-cell *matCellDef="let item">{{ item.rollover | number:'1.2-2' }}</td>
          </ng-container>

          <ng-container matColumnDef="available">
            <th mat-header-cell *matHeaderCellDef>Available</th>
            <td mat-cell *matCellDef="let item" [style.color]="(item.available || 0) < 0 ? 'red' : 'green'">
              {{ item.available | number:'1.2-2' }}
            </td>
          </ng-container>

          <tr mat-header-row *matHeaderRowDef="statusColumns"></tr>
          <tr mat-row *matRowDef="let row; columns: statusColumns"></tr>
          
           <tr class="mat-row" *matNoDataRow>
            <td class="mat-cell" [attr.colspan]="statusColumns.length">
              No budget status data.
            </td>
          </tr>
        </table>
      </div>
      }
    </div>
  `,
  styles: `
    .budget-items-container {
      padding: 0;
    }
    .page-title {
      margin: 0 0 16px 0;
      font-size: 24px;
      font-weight: 500;
    }
    .loading-container {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 300px;
    }
    .form-container {
      background: white;
      border-radius: 4px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      padding: 16px;
      margin-bottom: 24px;

      h2 {
        margin-top: 0;
        margin-bottom: 16px;
        font-size: 18px;
        font-weight: 500;
      }
    }
    .add-budget-form {
      display: flex;
      gap: 16px;
      align-items: flex-start;
      flex-wrap: wrap;

      mat-form-field {
        flex: 1;
        min-width: 200px;
      }

      button {
        height: 56px;
        min-width: 100px;
      }
    }
  `,
})
export class BudgetItemsComponent implements OnInit {
  private readonly budgetItemService = inject(BudgetItemService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly layoutService = inject(LayoutService);
  private readonly accountService = inject(AccountService);
  private readonly fb = inject(FormBuilder);

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;
  protected readonly Creating = signal(false);

  protected readonly budgetItems = this.budgetItemService.budgetItems;
  protected readonly budgetStatus = this.budgetItemService.budgetStatus;
  protected readonly accounts = this.accountService.accounts;
  protected readonly loading = this.budgetItemService.loading;
  protected readonly creating = this.budgetItemService.loading; // Re-use loading or create separate state? Re-using for now implicitly or create new one.
  // Actually, let's use a local signal for creating state to avoid global loading spinner overlap if desired, 
  // but budgetItemService has global loading. Ideally we separate list loading from create loading.
  // inspecting service: 'create' sets 'loading' to true. So global spinner will show. 
  // If we want inline spinner on button, we might want to check if data is empty. 
  // But global loading hides ALL content when 'loading()' is true. This is UX issue for "Add".
  // We should probably rely on 'loading' but maybe change template logic later.
  // For now, let's just use it.

  protected readonly displayedColumns = signal(['date', 'account', 'amount', 'description']);

  protected readonly addBudgetForm = this.fb.group({
    accountId: ['', Validators.required],
    amount: ['', [Validators.required, Validators.min(0)]],
    date: [new Date(), Validators.required],
  });
  protected readonly statusColumns = ['date', 'account', 'budgeted', 'spent', 'rollover', 'available'];

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
    this.accountService.loadAccounts().subscribe();
    // Load status for current year/month range ideally, but let's just load defaults (all)
    // We pass empty strings to let backend decide defaults or load all if supported
    this.budgetItemService.loadBudgetStatus(new Date(new Date().getFullYear(), 0, 1).toISOString(), new Date().toISOString()).subscribe();
  }

  protected addBudgetItem(): void {
    if (this.addBudgetForm.invalid) {
      return;
    }

    const { accountId, amount, date } = this.addBudgetForm.value;

    this.budgetItemService.create({
      accountId: accountId!,
      amount: Number(amount),
      date: new Date(date!).toISOString(),
      description: 'Budget Item', // Optional, maybe add field later
    }).subscribe({
      next: () => {
        this.snackBar.open('Budget item added', 'Close', { duration: 3000 });
        // Refresh status
        this.budgetItemService.loadBudgetStatus(new Date(new Date().getFullYear(), 0, 1).toISOString(), new Date().toISOString()).subscribe();
      },
      error: (err) => {
        this.snackBar.open('Failed to add budget item', 'Close', { duration: 3000 });
      }
    });
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
