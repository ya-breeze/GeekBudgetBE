import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatcherService } from './services/matcher.service';
import { Matcher } from '../../core/api/models/matcher';
import { AccountService } from '../accounts/services/account.service';
import { LayoutService } from '../../layout/services/layout.service';
import { MatcherEditDialogComponent } from './matcher-edit-dialog/matcher-edit-dialog.component';

@Component({
  selector: 'app-matchers',
  standalone: true,
  imports: [
    MatTableModule,
    MatSortModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatDialogModule
  ],
  templateUrl: './matchers.component.html',
  styleUrl: './matchers.component.css',
})
export class MatchersComponent implements OnInit {
  private readonly matcherService = inject(MatcherService);
  private readonly accountService = inject(AccountService);
  private readonly layoutService = inject(LayoutService);
  private readonly dialog = inject(MatDialog);

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  protected readonly loading = this.matcherService.loading;
  protected readonly displayedColumns = signal(['name', 'outputAccount', 'outputDescription', 'actions']);

  protected readonly sortActive = signal<string | null>(null);
  protected readonly sortDirection = signal<'asc' | 'desc'>('asc');

  // Computed signal that enriches matchers with account names and sorts
  protected readonly matchers = computed(() => {
    const matchers = this.matcherService.matchers();
    const accounts = this.accountService.accounts();
    const columns = this.displayedColumns();

    // Create a map of account IDs to names for quick lookup
    const accountMap = new Map(accounts.map((acc) => [acc.id, acc.name]));

    // Enrich matchers with account names
    const enrichedMatchers = matchers.map((matcher) => ({
      ...matcher,
      outputAccountName: accountMap.get(matcher.outputAccountId),
    }));

    if (!columns.length) {
      return enrichedMatchers;
    }

    const active = this.sortActive() ?? columns[0];
    const direction = this.sortDirection();

    // Sort by selected column
    return [...enrichedMatchers].sort((a, b) => this.compareMatchers(a, b, active, direction));
  });

  ngOnInit(): void {
    this.accountService.loadAccounts().subscribe();
    this.matcherService.loadMatchers().subscribe();
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

  protected openMatcherDialog(matcher?: Matcher): void {
    const dialogRef = this.dialog.open(MatcherEditDialogComponent, {
      data: matcher ? { matcher } : undefined,
      width: '600px',
      disableClose: true
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // List automatically updates via signals in service
      }
    });
  }

  protected deleteMatcher(matcher: Matcher): void {
    if (confirm(`Are you sure you want to delete matcher "${matcher.name}"?`)) {
      this.matcherService.delete(matcher.id).subscribe();
    }
  }

  private compareMatchers(
    a: Matcher & { outputAccountName?: string },
    b: Matcher & { outputAccountName?: string },
    active: string,
    direction: 'asc' | 'desc'
  ): number {
    const valueA = this.getMatcherSortValue(a, active);
    const valueB = this.getMatcherSortValue(b, active);
    return this.comparePrimitiveValues(valueA, valueB, direction);
  }

  private getMatcherSortValue(
    matcher: Matcher & { outputAccountName?: string },
    active: string
  ): string | null {
    switch (active) {
      case 'name':
        return this.removeLeadingEmoji(matcher.name ?? '');
      case 'outputAccount':
        return matcher.outputAccountName ?? matcher.outputAccountId ?? '';
      case 'outputDescription':
        return matcher.outputDescription ?? '';
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
