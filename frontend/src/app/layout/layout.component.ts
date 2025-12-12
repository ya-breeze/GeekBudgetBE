import { Component, inject, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { HeaderComponent } from './header/header.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FooterComponent } from './footer/footer.component';
import { MatSidenavModule } from '@angular/material/sidenav';
import { LayoutService } from './services/layout.service';
import { UnprocessedTransactionService } from '../features/unprocessed-transactions/services/unprocessed-transaction.service';

@Component({
  selector: 'app-layout',
  imports: [
    RouterOutlet,
    HeaderComponent,
    SidebarComponent,
    FooterComponent,
    MatSidenavModule,
  ],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss',
})
export class LayoutComponent implements OnInit {
  private readonly layoutService = inject(LayoutService);
  private readonly unprocessedTransactionService = inject(UnprocessedTransactionService);
  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  ngOnInit(): void {
    this.unprocessedTransactionService.loadUnprocessedTransactions().subscribe();
  }

  toggleSidenav(): void {
    this.layoutService.toggleSidenav();
  }
}

