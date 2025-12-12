import { Component, signal, output, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { AuthService } from '../../core/auth/services/auth.service';

interface MenuItem {
  label: string;
  icon: string;
  route: string;
}

@Component({
  selector: 'app-sidebar',
  imports: [RouterLink, RouterLinkActive, MatListModule, MatIconModule, MatButtonModule, MatMenuModule],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
})
export class SidebarComponent {
  private readonly authService = inject(AuthService);

  menuToggle = output<void>();

  protected readonly currentYear = signal(new Date().getFullYear());
  protected readonly version = signal('0.0.1');

  onMenuToggle(): void {
    this.menuToggle.emit();
  }

  protected readonly menuItems = signal<MenuItem[]>([
    { label: 'Dashboard', icon: 'dashboard', route: '/dashboard' },
    { label: 'Unprocessed', icon: 'pending_actions', route: '/unprocessed' },
    { label: 'Budget', icon: 'savings', route: '/budget' },
    { label: 'Transactions', icon: 'receipt_long', route: '/transactions' },

    { label: 'Settings', icon: 'settings', route: '/settings' },
    { label: 'Accounts', icon: 'account_balance', route: '/accounts' },
    { label: 'Currencies', icon: 'currency_exchange', route: '/currencies' },
    { label: 'Bank Importers', icon: 'cloud_upload', route: '/bank-importers' },

    { label: 'Matchers', icon: 'rule', route: '/matchers' },
    { label: 'Reports', icon: 'assessment', route: '/reports' },
    
    { label: 'Notifications', icon: 'notifications', route: '/notifications' },
  ]);

  logout(): void {
    this.authService.logout();
  }
}
