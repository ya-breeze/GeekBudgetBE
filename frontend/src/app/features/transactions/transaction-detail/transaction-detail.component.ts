import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTableModule } from '@angular/material/table';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDividerModule } from '@angular/material/divider';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ActivatedRoute, Router } from '@angular/router';
import { TransactionService } from '../services/transaction.service';
import { MergedTransactionService } from '../../merged-transactions/services/merged-transaction.service';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { MatcherService } from '../../matchers/services/matcher.service';
import { Transaction } from '../../../core/api/models/transaction';
import { MergedTransaction } from '../../../core/api/models/merged-transaction';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ImageUrlPipe } from '../../../shared/pipes/image-url.pipe';
import { AccountDisplayComponent } from '../../../shared/components/account-display/account-display.component';

@Component({
    selector: 'app-transaction-detail',
    standalone: true,
    imports: [
        CommonModule,
        MatCardModule,
        MatButtonModule,
        MatIconModule,
        MatChipsModule,
        MatExpansionModule,
        MatTableModule,
        MatProgressSpinnerModule,
        MatDividerModule,
        MatTooltipModule,
        ImageUrlPipe,
        AccountDisplayComponent,
    ],
    templateUrl: './transaction-detail.component.html',
    styleUrl: './transaction-detail.component.scss',
})
export class TransactionDetailComponent implements OnInit {
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);
    private readonly location = inject(Location);
    private readonly transactionService = inject(TransactionService);
    private readonly mergedTransactionService = inject(MergedTransactionService);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly matcherService = inject(MatcherService);
    private readonly snackBar = inject(MatSnackBar);

    readonly transactionId = signal<string | null>(null);
    readonly transaction = signal<Transaction | null>(null);
    readonly mergedTransaction = signal<MergedTransaction | null>(null);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    // Derived state
    readonly isMerged = computed(() => !!this.mergedTransaction());
    readonly displayedTransaction = computed(() => {
        return this.mergedTransaction()?.transaction || this.transaction();
    });

    readonly movements = computed(() => this.displayedTransaction()?.movements || []);
    readonly displayedColumns: string[] = ['amount', 'account', 'comment'];

    // Maps for ID resolution
    readonly accountMap = computed(() => {
        const map = new Map<string, import('../../../core/api/models/account').Account>();
        this.accountService.accounts().forEach((a) => map.set(a.id!, a));
        return map;
    });

    readonly currencyMap = computed(() => {
        const map = new Map<string, string>();
        this.currencyService.currencies().forEach((c) => map.set(c.id!, c.name || c.id!));
        return map;
    });

    readonly matcherMap = computed(() => {
        const map = new Map<string, import('../../../core/api/models/matcher').Matcher>();
        this.matcherService.matchers().forEach((m) => map.set(m.id!, m));
        return map;
    });

    getMatcherName(id: string): string {
        const matcher = this.matcherMap().get(id);
        return matcher?.outputDescription || matcher?.descriptionRegExp || id;
    }

    ngOnInit(): void {
        this.route.paramMap.subscribe((params) => {
            const id = params.get('id');
            if (id) {
                this.transactionId.set(id);
                this.loadTransactionData(id);
                // Ensure auxiliary data is loaded
                if (this.accountService.accounts().length === 0)
                    this.accountService.loadAccounts().subscribe();
                if (this.currencyService.currencies().length === 0)
                    this.currencyService.loadCurrencies().subscribe();
                if (this.matcherService.matchers().length === 0)
                    this.matcherService.loadMatchers().subscribe();
            }
        });
    }

    loadTransactionData(id: string): void {
        this.loading.set(true);
        this.error.set(null);
        this.transaction.set(null);
        this.mergedTransaction.set(null);

        // Try fetching as active transaction first
        this.transactionService.getTransaction(id).subscribe({
            next: (t) => {
                this.transaction.set(t);
                this.loading.set(false);
            },
            error: () => {
                // If failed, try fetching as merged transaction
                this.mergedTransactionService.getMergedTransaction(id).subscribe({
                    next: (mt) => {
                        this.mergedTransaction.set(mt);
                        this.loading.set(false);
                    },
                    error: (err) => {
                        this.error.set('Transaction not found');
                        this.loading.set(false);
                        console.error('Failed to load transaction details', err);
                    },
                });
            },
        });
    }

    goBack(): void {
        this.location.back();
    }

    editTransaction(): void {
        // Navigate to list with open edit dialog or implement edit logic here?
        // For now, maybe just log or show not implemented, or navigate to list and open dialog
        // Since existing edit is a dialog on list, reusing it might be tricky without moving logic.
        // We'll leave it as a TODO or just navigate back for now.
        this.snackBar.open('Editing from detail view is not yet implemented', 'Close', {
            duration: 3000,
        });
    }

    viewMergedTransaction(id: string): void {
        this.router.navigate(['/transactions', id]);
    }

    getFormattedRawData(): string {
        const raw = this.displayedTransaction()?.unprocessedSources;
        if (!raw) return '';
        try {
            return JSON.stringify(JSON.parse(raw), null, 2);
        } catch {
            return raw;
        }
    }
}
