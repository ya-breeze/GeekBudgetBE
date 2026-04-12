import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatBadgeModule } from '@angular/material/badge';
import { UnprocessedTransactionService } from '../../features/unprocessed-transactions/services/unprocessed-transaction.service';

@Component({
    selector: 'app-mobile-nav',
    imports: [RouterLink, RouterLinkActive, MatIconModule, MatBadgeModule],
    templateUrl: './mobile-nav.component.html',
    styleUrl: './mobile-nav.component.scss',
})
export class MobileNavComponent {
    private readonly unprocessedTransactionService = inject(UnprocessedTransactionService);
    protected readonly unprocessedCount = this.unprocessedTransactionService.unprocessedTransactions;
}
