import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
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
        MatDialogModule,
        MatTooltipModule,
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
    protected readonly displayedColumns = signal([
        'outputAccount',
        'outputDescription',
        'confidence',
        'actions',
    ]);

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
            disableClose: true,
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result) {
                // List automatically updates via signals in service
            }
        });
    }

    protected deleteMatcher(matcher: Matcher): void {
        if (confirm(`Are you sure you want to delete matcher "${matcher.outputDescription}"?`)) {
            this.matcherService.delete(matcher.id).subscribe();
        }
    }

    protected getConfidenceBadge(matcher: Matcher): {
        text: string;
        class: string;
        tooltip: string;
    } {
        const count = matcher.confirmationsCount || 0;
        const total = matcher.confirmationsTotal || 0;

        if (total === 0) {
            return { text: 'New', class: 'badge-secondary', tooltip: 'No confirmation history' };
        }

        const percentage = (count / total) * 100;
        const isPerfect = percentage === 100;
        const isLargeSample = count >= 10;

        let badgeClass = 'badge-danger'; // <40%
        if (isPerfect && isLargeSample) {
            badgeClass = 'badge-perfect';
        } else if (percentage >= 70) {
            badgeClass = 'badge-success';
        } else if (percentage >= 40) {
            badgeClass = 'badge-warning';
        }

        return {
            text: `${count}/${total}`,
            class: badgeClass,
            tooltip: `${count} successful confirmations out of ${total} attempts (${percentage.toFixed(0)}%)`,
        };
    }

    private compareMatchers(
        a: Matcher & { outputAccountName?: string },
        b: Matcher & { outputAccountName?: string },
        active: string,
        direction: 'asc' | 'desc',
    ): number {
        const valueA = this.getMatcherSortValue(a, active);
        const valueB = this.getMatcherSortValue(b, active);
        return this.comparePrimitiveValues(valueA, valueB, direction);
    }

    private getMatcherSortValue(
        matcher: Matcher & { outputAccountName?: string },
        active: string,
    ): string | number | null {
        switch (active) {
            case 'outputAccount':
                return matcher.outputAccountName ?? matcher.outputAccountId ?? '';
            case 'outputDescription':
                return this.removeLeadingEmoji(matcher.outputDescription ?? '');
            case 'confidence':
                const total = matcher.confirmationsTotal || 0;
                if (total === 0) return -1;

                const count = matcher.confirmationsCount || 0;
                let score = count / total;

                // Boost score for proven perfect matches (>= 10 confirmations, 100%)
                // or penalize small perfect matches.
                // User wants: perfect (>=10) > perfect (<10)

                if (score === 1.0) {
                    if (count >= 10) {
                        return 2.0; // Top tier
                    }
                    // Small sample perfect: 1.0 - penalty
                    // Map 1..9 to 0.91..0.99 (approx)
                    // 9/9 -> 1.0 - 0.001 = 0.999
                    // 1/1 -> 1.0 - 0.009 = 0.991
                    return 1.0 - 0.001 * (10 - count);
                }

                return score;
            default:
                return null;
        }
    }

    private removeLeadingEmoji(text: string): string {
        // Remove leading emoji characters for sorting purposes
        return text
            .replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '')
            .trim();
    }

    private comparePrimitiveValues(
        a: string | number | Date | null | undefined,
        b: string | number | Date | null | undefined,
        direction: 'asc' | 'desc',
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
