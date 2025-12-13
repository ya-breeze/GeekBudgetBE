import { Component, output, inject, signal } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { LayoutService } from '../services/layout.service';
import { Router } from '@angular/router';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
}

@Component({
  selector: 'app-header',
  imports: [MatToolbarModule, MatButtonModule, MatIconModule, MatMenuModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
})
export class HeaderComponent {
  private readonly layoutService = inject(LayoutService);
  private readonly router = inject(Router);

  readonly menuToggle = output<void>();

  protected readonly sidenavOpened = this.layoutService.sidenavOpened;

  protected readonly menuItems = signal<MenuItem[]>([
    { label: 'Dashboard', icon: 'dashboard', route: '/dashboard' },
    { label: 'Transactions', icon: 'receipt_long', route: '/transactions' },
    { label: 'Accounts', icon: 'account_balance', route: '/accounts' },
    { label: 'Currencies', icon: 'currency_exchange', route: '/currencies' },
    { label: 'Bank Importers', icon: 'cloud_upload', route: '/bank-importers' },
    { label: 'Unprocessed', icon: 'pending_actions', route: '/unprocessed' },
    { label: 'Matchers', icon: 'rule', route: '/matchers' },
    { label: 'Budget', icon: 'savings', route: '/budget' },
    { label: 'Reports', icon: 'assessment', route: '/reports' },
    { label: 'Notifications', icon: 'notifications', route: '/notifications' },
    { label: 'Settings', icon: 'settings', route: '/settings' },
  ]);

  onMenuToggle(): void {
    this.menuToggle.emit();
  }

  navigateTo(route: string): void {
    this.router.navigate([route]);
  }
}

