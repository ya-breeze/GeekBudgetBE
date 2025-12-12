import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { NgClass } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { MatSortModule, Sort } from '@angular/material/sort';

import { AccountService } from './services/account.service';
import { Account } from '../../core/api/models/account';
import { AccountFormDialogComponent } from './account-form-dialog/account-form-dialog.component';
import { LayoutService } from '../../layout/services/layout.service';

@Component({
  selector: 'app-accounts',
  imports: [
    NgClass,
    MatTableModule,
    MatSortModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatSnackBarModule,
    MatChipsModule,
  ],
  templateUrl: './accounts.component.html',
  styleUrl: './accounts.component.scss',
})
export class AccountsComponent implements OnInit {
  private readonly accountService = inject(AccountService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);
  private readonly layoutService = inject(LayoutService);

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  protected readonly sortActive = signal<string | null>(null);
  protected readonly sortDirection = signal<'asc' | 'desc'>('asc');
  protected readonly sortedAccounts = computed(() => {
    const data = this.accounts();
    const columns = this.displayedColumns();

    if (!columns.length) {
      return data;
    }

    const active = this.sortActive() ?? columns[0];
    const direction = this.sortDirection();

    return [...data].sort((a, b) => this.compareAccounts(a, b, active, direction));
  });

  protected readonly accounts = this.accountService.accounts;
  protected readonly loading = this.accountService.loading;
  protected readonly displayedColumns = signal(['name', 'type', 'description', 'actions']);

  ngOnInit(): void {
    this.loadAccounts();
  }

  loadAccounts(): void {
    this.accountService.loadAccounts().subscribe();
  }

  openCreateDialog(): void {
    const dialogRef = this.dialog.open(AccountFormDialogComponent, {
      width: '600px',
      data: { mode: 'create' },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.accountService.create(result).subscribe({
          next: () => {
            this.snackBar.open('Account created successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to create account', 'Close', { duration: 3000 });
          },
        });
      }
    });
  }

  openEditDialog(account: Account): void {
    const dialogRef = this.dialog.open(AccountFormDialogComponent, {
      width: '600px',
      data: { mode: 'edit', account },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result && account.id) {
        this.accountService.update(account.id, result).subscribe({
          next: () => {
            this.snackBar.open('Account updated successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to update account', 'Close', { duration: 3000 });
          },
        });
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

  private compareAccounts(
    a: Account,
    b: Account,
    active: string,
    direction: 'asc' | 'desc'
  ): number {
    const valueA = this.getAccountSortValue(a, active);
    const valueB = this.getAccountSortValue(b, active);
    return this.comparePrimitiveValues(valueA, valueB, direction);
  }

  private getAccountSortValue(account: Account, active: string): string | null {
    switch (active) {
      case 'name':
        return this.removeLeadingEmoji(account.name ?? '');
      case 'type':
        return account.type ?? '';
      case 'description':
        return account.description ?? '';
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

  deleteAccount(account: Account): void {
    if (confirm(`Are you sure you want to delete "${account.name}"?`)) {
      if (account.id) {
        this.accountService.delete(account.id).subscribe({
          next: () => {
            this.snackBar.open('Account deleted successfully', 'Close', { duration: 3000 });
          },
          error: () => {
            this.snackBar.open('Failed to delete account', 'Close', { duration: 3000 });
          },
        });
      }
    }
  }
}
