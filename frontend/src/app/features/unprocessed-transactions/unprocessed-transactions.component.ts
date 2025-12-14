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
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { UnprocessedTransactionDialogComponent, UnprocessedTransactionDialogResult } from './unprocessed-transaction-dialog/unprocessed-transaction-dialog.component';

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
    MatDialogModule,
  ],
  templateUrl: './unprocessed-transactions.component.html',
  styleUrl: './unprocessed-transactions.component.scss',
})
export class UnprocessedTransactionsComponent implements OnInit {
  private readonly unprocessedTransactionService = inject(UnprocessedTransactionService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly layoutService = inject(LayoutService);
  private readonly currencyService = inject(CurrencyService);
  private readonly dialog = inject(MatDialog);

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

  openProcessDialog(transaction: UnprocessedTransaction): void {
    const dialogRef = this.dialog.open(UnprocessedTransactionDialogComponent, {
      data: transaction,
      width: '95vw',
      maxWidth: '1200px',
      height: '90vh',
      autoFocus: false,
    });

    const componentInstance = dialogRef.componentInstance;

    // Handle actions from the dialog
    componentInstance.action.subscribe((result: UnprocessedTransactionDialogResult) => {
      // Access the signal directly (using any cast to bypass protected visibility if needed)
      const currentTransaction = (componentInstance as any).transaction();

      if (result.action === 'convert') {
        this.processMatch(currentTransaction, result.match, dialogRef);
      } else if (result.action === 'delete') {
        this.deleteTransaction(currentTransaction, result.duplicateOf.id, dialogRef);
      } else if (result.action === 'manual') {
        this.processManual(currentTransaction, result.accountId, dialogRef);
      }
    });

    dialogRef.afterClosed().subscribe(() => {
      this.loadUnprocessedTransactions();
    });
  }

  private processMatch(original: UnprocessedTransaction, match: any, dialogRef?: any) {
    const transactionToConvert = { ...original };
    transactionToConvert.transaction = {
      ...original.transaction,
      ...match.transaction
    };

    this.unprocessedTransactionService.convert(original.transaction.id!, transactionToConvert, match.matcherId).subscribe({
      next: () => {
        this.snackBar.open('Transaction processed (match applied)', 'Close', { duration: 3000 });
        if (dialogRef) {
          this.handleSuccess(original.transaction.id!, dialogRef);
        }
      },
      error: () => {
        this.snackBar.open('Failed to process transaction', 'Close', { duration: 3000 });
        if (dialogRef) dialogRef.componentInstance.setLoading(false);
      }
    });
  }

  private processManual(original: UnprocessedTransaction, accountId: string, dialogRef?: any) {
    const transactionToConvert = JSON.parse(JSON.stringify(original));
    const movements = transactionToConvert.transaction.movements || [];

    const targetMovement = movements.find((m: any) => !m.accountId);
    if (targetMovement) {
      targetMovement.accountId = accountId;
    }

    this.unprocessedTransactionService.convert(original.transaction.id!, transactionToConvert).subscribe({
      next: () => {
        this.snackBar.open('Transaction processed (account assigned)', 'Close', { duration: 3000 });
        if (dialogRef) {
          this.handleSuccess(original.transaction.id!, dialogRef);
        }
      },
      error: () => {
        this.snackBar.open('Failed to process transaction', 'Close', { duration: 3000 });
        if (dialogRef) dialogRef.componentInstance.setLoading(false);
      }
    });
  }

  deleteTransaction(transaction: UnprocessedTransaction, duplicateOfId?: string, dialogRef?: any): void {
    this.unprocessedTransactionService.delete(transaction.transaction.id!, duplicateOfId).subscribe({
      next: () => {
        this.snackBar.open('Transaction deleted', 'Close', { duration: 3000 });
        if (dialogRef) {
          this.handleSuccess(transaction.transaction.id!, dialogRef);
        } else {
          this.loadUnprocessedTransactions();
        }
      },
      error: () => {
        this.snackBar.open('Failed to delete transaction', 'Close', { duration: 3000 });
        if (dialogRef) dialogRef.componentInstance.setLoading(false);
      }
    });
  }

  private handleSuccess(processedId: string, dialogRef: any) {
    const currentList = this.unprocessedTransactionService.unprocessedTransactions();
    const nextTransaction = currentList.find(t => t.transaction.id !== processedId);

    if (nextTransaction) {
      dialogRef.componentInstance.updateTransaction(nextTransaction);
      this.loadUnprocessedTransactions();
    } else {
      dialogRef.close();
      this.loadUnprocessedTransactions();
    }
  }
}
