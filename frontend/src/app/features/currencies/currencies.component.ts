import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { CurrencyService } from './services/currency.service';
import { Currency } from '../../core/api/models/currency';
import { CurrencyFormDialogComponent } from './currency-form-dialog/currency-form-dialog.component';

@Component({
  selector: 'app-currencies',
  imports: [
    MatTableModule,
    MatSortModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatSnackBarModule,
  ],
  templateUrl: './currencies.component.html',
  styleUrl: './currencies.component.scss',
})
export class CurrenciesComponent implements OnInit {
  private readonly currencyService = inject(CurrencyService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);

  protected readonly sortActive = signal<string | null>(null);
  protected readonly sortDirection = signal<'asc' | 'desc'>('asc');
  protected readonly sortedCurrencies = computed(() => {
    const data = this.currencies();
    const columns = this.displayedColumns();

    if (!columns.length) {
      return data;
    }

    const active = this.sortActive() ?? columns[0];
    const direction = this.sortDirection();

    return [...data].sort((a, b) => this.compareCurrencies(a, b, active, direction));
  });

  protected readonly currencies = this.currencyService.currencies;
  protected readonly loading = this.currencyService.loading;
  protected readonly displayedColumns = signal(['name', 'description', 'actions']);

  ngOnInit(): void {
    this.loadCurrencies();
  }

  loadCurrencies(): void {
    this.currencyService.loadCurrencies().subscribe();
  }

  openCreateDialog(): void {
    const dialogRef = this.dialog.open(CurrencyFormDialogComponent, {
      width: '500px',
      data: { mode: 'create' },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.currencyService.create(result).subscribe({
          next: () => {
            this.snackBar.open('Currency created successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to create currency', 'Close', { duration: 3000 });
          },
        });
      }
    });
  }

  openEditDialog(currency: Currency): void {
    const dialogRef = this.dialog.open(CurrencyFormDialogComponent, {
      width: '500px',
      data: { mode: 'edit', currency },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result && currency.id) {
        this.currencyService.update(currency.id, result).subscribe({
          next: () => {
            this.snackBar.open('Currency updated successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to update currency', 'Close', { duration: 3000 });
          },
        });
      }
    });
  }

  deleteCurrency(currency: Currency): void {
    if (confirm(`Are you sure you want to delete "${currency.name}"?`)) {
      if (currency.id) {
        this.currencyService.delete(currency.id).subscribe({
          next: () => {
            this.snackBar.open('Currency deleted successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to delete currency', 'Close', { duration: 3000 });
          },
        });
      }
    }
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

  private compareCurrencies(
    a: Currency,
    b: Currency,
    active: string,
    direction: 'asc' | 'desc'
  ): number {
    const valueA = this.getCurrencySortValue(a, active);
    const valueB = this.getCurrencySortValue(b, active);
    return this.comparePrimitiveValues(valueA, valueB, direction);
  }

  private getCurrencySortValue(currency: Currency, active: string): string | null {
    switch (active) {
      case 'name':
        return this.removeLeadingEmoji(currency.name ?? '');
      case 'description':
        return currency.description ?? '';
      default:
        return null;
    }
  }

  private removeLeadingEmoji(text: string): string {
    // Remove leading emoji characters for sorting purposes
    return text.replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '').trim();
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
