import { Component, inject, OnInit, signal, ViewChild } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, MatSort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { DatePipe } from '@angular/common';
import { UnprocessedTransactionService } from './services/unprocessed-transaction.service';
import { UnprocessedTransaction } from '../../core/api/models/unprocessed-transaction';
import { LayoutService } from '../../layout/services/layout.service';
import { CurrencyService } from '../currencies/services/currency.service';

@Component({
  selector: 'app-unprocessed-transactions',
  imports: [
    MatTableModule,
    MatSortModule,
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
  private readonly currencyService = inject(CurrencyService);

  @ViewChild(MatSort) sort!: MatSort;

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  protected readonly unprocessedTransactions =
    this.unprocessedTransactionService.unprocessedTransactions;
  protected readonly loading = this.unprocessedTransactionService.loading;
  protected readonly currencies = this.currencyService.currencies;
  protected readonly displayedColumns = signal([
    'date',
    'description',
    'effectiveAmount',
    'matches',
    'duplicates',
    'actions',
  ]);
  protected readonly sortedTransactions = signal<UnprocessedTransaction[]>([]);

  ngOnInit(): void {
    this.loadUnprocessedTransactions();
    this.currencyService.loadCurrencies().subscribe();
  }

  loadUnprocessedTransactions(): void {
    this.unprocessedTransactionService.loadUnprocessedTransactions().subscribe();
  }

  getCurrencyName(currencyId: string): string {
    const currency = this.currencies().find((c) => c.id === currencyId);
    return currency ? currency.name : currencyId;
  }

  getUnknownAccountMovements(transaction: any) {
    return transaction.movements?.filter((m: any) => !m.accountId) || [];
  }

  getSortedData(): UnprocessedTransaction[] {
    const data = this.unprocessedTransactions();
    if (!this.sort || !this.sort.active || this.sort.direction === '') {
      return data;
    }

    return data.slice().sort((a, b) => {
      const isAsc = this.sort.direction === 'asc';
      switch (this.sort.active) {
        case 'date':
          return this.compare(new Date(a.transaction.date), new Date(b.transaction.date), isAsc);
        case 'description':
          return this.compare(a.transaction.description || '', b.transaction.description || '', isAsc);
        case 'effectiveAmount':
          return this.compareAmounts(a, b, isAsc);
        case 'matches':
          return this.compare(a.matched.length, b.matched.length, isAsc);
        case 'duplicates':
          return this.compare(a.duplicates.length, b.duplicates.length, isAsc);
        default:
          return 0;
      }
    });
  }

  private compare(a: any, b: any, isAsc: boolean): number {
    if (a < b) {
      return isAsc ? -1 : 1;
    }
    if (a > b) {
      return isAsc ? 1 : -1;
    }
    return 0;
  }

  private compareAmounts(a: UnprocessedTransaction, b: UnprocessedTransaction, isAsc: boolean): number {
    const aAmount = this.getTotalUnknownAccountAmount(a.transaction);
    const bAmount = this.getTotalUnknownAccountAmount(b.transaction);
    return this.compare(aAmount, bAmount, isAsc);
  }

  private getTotalUnknownAccountAmount(transaction: any): number {
    return this.getUnknownAccountMovements(transaction).reduce((sum: number, m: any) => sum + m.amount, 0);
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
