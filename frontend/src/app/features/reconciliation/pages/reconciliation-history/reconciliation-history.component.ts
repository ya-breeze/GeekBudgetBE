import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ReconciliationService } from '../../services/reconciliation.service';
import { Reconciliation } from '../../../../core/api/models/reconciliation';

@Component({
    selector: 'app-reconciliation-history',
    standalone: true,
    imports: [
        CommonModule,
        RouterModule,
        MatTableModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './reconciliation-history.component.html',
    styleUrls: ['./reconciliation-history.component.scss'],
})
export class ReconciliationHistoryComponent implements OnInit {
    private route = inject(ActivatedRoute);
    private reconciliationService = inject(ReconciliationService);

    accountId = '';
    currencyId = '';
    history = signal<Reconciliation[]>([]);
    loading = signal(true);
    error = signal<string | null>(null);

    accountName = signal<string>('');
    currencySymbol = signal<string>('');

    displayedColumns: string[] = ['date', 'reconciledBalance', 'expectedBalance', 'delta', 'type'];

    ngOnInit(): void {
        this.accountId = this.route.snapshot.paramMap.get('accountId') || '';
        this.currencyId = this.route.snapshot.paramMap.get('currencyId') || '';
        this.loadLookups();
        this.loadHistory();
    }

    loadLookups(): void {
        this.reconciliationService.getAccounts().subscribe((accounts) => {
            const acc = accounts.find((a) => a.id === this.accountId);
            if (acc) this.accountName.set(acc.name);
        });
        this.reconciliationService.getCurrencies().subscribe((currencies) => {
            const curr = currencies.find((c) => c.id === this.currencyId);
            if (curr) this.currencySymbol.set(curr.name);
        });
    }

    loadHistory(): void {
        this.loading.set(true);
        this.error.set(null);
        this.reconciliationService.getHistory(this.accountId, this.currencyId).subscribe({
            next: (data) => {
                this.history.set(data);
                this.loading.set(false);
            },
            error: (err) => {
                this.error.set(err.message || 'Failed to load history');
                this.loading.set(false);
            },
        });
    }

    getDelta(item: Reconciliation): number {
        return (item.reconciledBalance || 0) - (item.expectedBalance || 0);
    }
}
