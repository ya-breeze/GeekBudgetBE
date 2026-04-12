import { Component, effect, inject, OnInit } from '@angular/core';
import { RouterOutlet, Router, NavigationEnd } from '@angular/router';
import { filter } from 'rxjs';
import { HeaderComponent } from './header/header.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FooterComponent } from './footer/footer.component';
import { MobileNavComponent } from './mobile-nav/mobile-nav.component';
import { MatSidenavModule } from '@angular/material/sidenav';
import { LayoutService } from './services/layout.service';
import { UnprocessedTransactionService } from '../features/unprocessed-transactions/services/unprocessed-transaction.service';
import { NotificationService } from '../features/notifications/services/notification.service';

@Component({
    selector: 'app-layout',
    imports: [RouterOutlet, HeaderComponent, SidebarComponent, FooterComponent, MobileNavComponent, MatSidenavModule],
    templateUrl: './layout.component.html',
    styleUrl: './layout.component.scss',
})
export class LayoutComponent implements OnInit {
    private readonly layoutService = inject(LayoutService);
    private readonly router = inject(Router);
    private readonly unprocessedTransactionService = inject(UnprocessedTransactionService);
    private readonly notificationService = inject(NotificationService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;
    protected readonly isMobile = this.layoutService.isMobile;

    constructor() {
        // Close sidenav when viewport shrinks to mobile
        effect(() => {
            if (this.layoutService.isMobile()) {
                this.layoutService.closeSidenav();
            }
        });
    }

    ngOnInit(): void {
        this.unprocessedTransactionService.loadUnprocessedTransactions().subscribe();
        this.notificationService.startPolling();

        // Close sidenav after each navigation on mobile
        this.router.events.pipe(filter((e) => e instanceof NavigationEnd)).subscribe(() => {
            if (this.layoutService.isMobile()) {
                this.layoutService.closeSidenav();
            }
        });
    }

    toggleSidenav(): void {
        this.layoutService.toggleSidenav();
    }
}
