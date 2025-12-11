import { Component, inject, OnInit, signal } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { DatePipe } from '@angular/common';
import { UnprocessedTransactionService } from './services/unprocessed-transaction.service';
import { UnprocessedTransaction } from '../../core/api/models/unprocessed-transaction';
import { LayoutService } from '../../layout/services/layout.service';

@Component({
  selector: 'app-unprocessed-transactions',
  imports: [
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatChipsModule,
    DatePipe,
  ],
  templateUrl: './unprocessed-transactions.component.html',
  styleUrl: './unprocessed-transactions.component.scss',
})
export class UnprocessedTransactionsComponent implements OnInit {
  private readonly unprocessedTransactionService = inject(UnprocessedTransactionService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly layoutService = inject(LayoutService);

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  protected readonly unprocessedTransactions =
    this.unprocessedTransactionService.unprocessedTransactions;
  protected readonly loading = this.unprocessedTransactionService.loading;
  protected readonly displayedColumns = signal([
    'date',
    'description',
    'matches',
    'duplicates',
    'actions',
  ]);

  ngOnInit(): void {
    this.loadUnprocessedTransactions();
  }

  loadUnprocessedTransactions(): void {
    this.unprocessedTransactionService.loadUnprocessedTransactions().subscribe();
  }

  convertTransaction(transaction: UnprocessedTransaction): void {
    if (transaction.transaction.id) {
      this.unprocessedTransactionService
        .convert(transaction.transaction.id, transaction)
        .subscribe({
          next: () => {
            this.snackBar.open('Transaction converted successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to convert transaction', 'Close', { duration: 3000 });
          },
        });
    }
  }

  deleteTransaction(transaction: UnprocessedTransaction): void {
    if (confirm('Are you sure you want to delete this transaction?')) {
      if (transaction.transaction.id) {
        this.unprocessedTransactionService.delete(transaction.transaction.id).subscribe({
          next: () => {
            this.snackBar.open('Transaction deleted successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to delete transaction', 'Close', { duration: 3000 });
          },
        });
      }
    }
  }
}
