import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { DatePipe } from '@angular/common';
import { TransactionService } from './services/transaction.service';
import { Transaction } from '../../core/api/models/transaction';
import { TransactionFormDialogComponent } from './transaction-form-dialog/transaction-form-dialog.component';
import { AccountService } from '../accounts/services/account.service';
import { Account } from '../../core/api/models/account';
import { forkJoin, Observable } from 'rxjs';

@Component({
  selector: 'app-transactions',
  imports: [
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatSnackBarModule,
    MatChipsModule,
    DatePipe,
  ],
  templateUrl: './transactions.component.html',
  styleUrl: './transactions.component.scss',
})
export class TransactionsComponent implements OnInit {
  private readonly transactionService = inject(TransactionService);
  private readonly accountService = inject(AccountService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);

  protected readonly transactions = this.transactionService.transactions;
  protected readonly loading = signal(false); // Combined loading for transactions and accounts
  protected readonly displayedColumns = signal([
    'date',
    'description',
    'movements',
    'tags',
    'actions',
  ]);

  protected readonly accounts = this.accountService.accounts;
  protected readonly accountMap = computed<Map<string, Account>>(() => {
    const map = new Map<string, Account>();
    this.accounts().forEach((account) => map.set(account.id!, account));
    return map;
  });

  // Month navigation state
  protected readonly currentMonth = signal(new Date().getMonth());
  protected readonly currentYear = signal(new Date().getFullYear());

  // Computed property for displaying current month/year
  protected readonly currentMonthDisplay = computed(() => {
    const date = new Date(this.currentYear(), this.currentMonth(), 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  });

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.loading.set(true);
    forkJoin([this.accountService.loadAccounts(), this.loadTransactions()]).subscribe({
      next: () => this.loading.set(false),
      error: () => this.loading.set(false),
    });
  }

  loadTransactions(): Observable<any> {
    const startOfMonth = new Date(this.currentYear(), this.currentMonth(), 1);
    const endOfMonth = new Date(this.currentYear(), this.currentMonth() + 1, 0, 23, 59, 59, 999);

    const params = {
      dateFrom: startOfMonth.toISOString(),
      dateTo: endOfMonth.toISOString(),
    };

    return this.transactionService.loadTransactions(params);
  }

  goToPreviousMonth(): void {
    if (this.currentMonth() === 0) {
      this.currentMonth.set(11);
      this.currentYear.set(this.currentYear() - 1);
    } else {
      this.currentMonth.set(this.currentMonth() - 1);
    }
    this.loadData();
  }

  goToNextMonth(): void {
    if (this.currentMonth() === 11) {
      this.currentMonth.set(0);
      this.currentYear.set(this.currentYear() + 1);
    } else {
      this.currentMonth.set(this.currentMonth() + 1);
    }
    this.loadData();
  }

  openCreateDialog(): void {
    const dialogRef = this.dialog.open(TransactionFormDialogComponent, {
      width: '90vw',
      maxWidth: '1200px',
      height: 'auto',
      maxHeight: '90vh',
      data: { mode: 'create' },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.transactionService.create(result).subscribe({
          next: () => {
            this.snackBar.open('Transaction created successfully', 'Close', { duration: 3000 });
            this.loadData(); // Reload data after creation
          },
          error: () => {
            this.snackBar.open('Failed to create transaction', 'Close', { duration: 3000 });
          },
        });
      }
    });
  }

  openEditDialog(transaction: Transaction): void {
    const dialogRef = this.dialog.open(TransactionFormDialogComponent, {
      width: '90vw',
      maxWidth: '1200px',
      height: 'auto',
      maxHeight: '90vh',
      data: { mode: 'edit', transaction },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result && transaction.id) {
        this.transactionService.update(transaction.id, result).subscribe({
          next: () => {
            this.snackBar.open('Transaction updated successfully', 'Close', { duration: 3000 });
            this.loadData(); // Reload data after update
          },
          error: () => {
            this.snackBar.open('Failed to update transaction', 'Close', { duration: 3000 });
          },
        });
      }
    });
  }

  deleteTransaction(transaction: Transaction): void {
    if (confirm(`Are you sure you want to delete this transaction?`)) {
      if (transaction.id) {
        this.transactionService.delete(transaction.id).subscribe({
          next: () => {
            this.snackBar.open('Transaction deleted successfully', 'Close', { duration: 3000 });
            this.loadData(); // Reload data after deletion
          },
          error: () => {
            this.snackBar.open('Failed to delete transaction', 'Close', { duration: 3000 });
          },
        });
      }
    }
  }

  formatMovements(transaction: Transaction): string {
    if (!transaction.movements || transaction.movements.length === 0) {
      return 'No movements';
    }

    const accountMap = this.accountMap();

    // Get input movements (sources of money - negative amounts)
    const inputMovements = transaction.movements.filter((movement) => movement.amount < 0);
    const inputAccountNames = inputMovements.map((movement) => {
      if (!movement.accountId) {
        return 'Undefined';
      }
      return accountMap.get(movement.accountId)?.name || movement.accountId;
    });

    // Get output movements (destinations of money - positive amounts)
    const outputMovements = transaction.movements.filter((movement) => movement.amount > 0);
    const outputAccountNames = outputMovements.map((movement) => {
      if (!movement.accountId) {
        return 'Undefined';
      }
      return accountMap.get(movement.accountId)?.name || movement.accountId;
    });

    // Format as "input => output"
    if (inputAccountNames.length === 0 && outputAccountNames.length === 0) {
      return 'No valid movements';
    } else if (inputAccountNames.length === 0) {
      return outputAccountNames.join(', ');
    } else if (outputAccountNames.length === 0) {
      return inputAccountNames.join(', ');
    } else {
      const inputPart = inputAccountNames.join(', ');
      const outputPart = outputAccountNames.join(', ');
      return `${inputPart} => ${outputPart}`;
    }
  }
}
