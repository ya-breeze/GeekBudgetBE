import { Component, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
}

@Component({
  selector: 'app-sidebar',
  imports: [RouterLink, RouterLinkActive, MatListModule, MatIconModule],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
})
export class SidebarComponent {
  protected readonly currentYear = signal(new Date().getFullYear());
  protected readonly version = signal('0.0.1');

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
}
