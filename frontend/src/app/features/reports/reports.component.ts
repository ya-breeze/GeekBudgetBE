import { Component, signal, inject } from '@angular/core';
import { RouterOutlet, RouterLink, Router, NavigationEnd } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { filter } from 'rxjs';

@Component({
  selector: 'app-reports',
  imports: [RouterOutlet, RouterLink, MatCardModule, MatButtonModule, MatIconModule],
  template: `
    <div class="reports-container">
      @if (showCards()) {
      <h1>Reports & Analytics</h1>
      <p class="subtitle">Financial reports and insights</p>

      <div class="reports-grid">
        <mat-card>
          <mat-card-header>
            <mat-icon>bar_chart</mat-icon>
            <mat-card-title>Expense Report</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <p>View detailed expense breakdown by category and time period.</p>
          </mat-card-content>
          <mat-card-actions>
            <a mat-button color="primary" routerLink="expense">View Report</a>
          </mat-card-actions>
        </mat-card>

        <mat-card>
          <mat-card-header>
            <mat-icon>trending_up</mat-icon>
            <mat-card-title>Income Report</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <p>Analyze income sources and trends over time.</p>
          </mat-card-content>
          <mat-card-actions>
            <button mat-button color="primary" disabled>Coming Soon</button>
          </mat-card-actions>
        </mat-card>

        <mat-card>
          <mat-card-header>
            <mat-icon>account_balance</mat-icon>
            <mat-card-title>Balance Report</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <p>Track account balances and net worth over time.</p>
          </mat-card-content>
          <mat-card-actions>
            <button mat-button color="primary" disabled>Coming Soon</button>
          </mat-card-actions>
        </mat-card>

        <mat-card>
          <mat-card-header>
            <mat-icon>download</mat-icon>
            <mat-card-title>Export Data</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <p>Export your financial data in various formats.</p>
          </mat-card-content>
          <mat-card-actions>
            <button mat-button color="primary" disabled>Coming Soon</button>
          </mat-card-actions>
        </mat-card>
      </div>
      } @else {
      <router-outlet />
      }
    </div>
  `,
  styles: `
    .reports-container {
      padding: 24px;
    }
    .subtitle {
      margin: 0 0 24px 0;
      color: rgba(0, 0, 0, 0.6);
    }
    .reports-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 24px;
    }
    mat-card-header {
      display: flex;
      align-items: center;
      gap: 12px;
      mat-icon {
        font-size: 32px;
        width: 32px;
        height: 32px;
        color: #2196f3;
      }
    }
  `,
})
export class ReportsComponent {
  private readonly router = inject(Router);
  protected readonly showCards = signal(true);

  constructor() {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: NavigationEnd) => {
        this.showCards.set(event.url === '/reports' || event.url === '/reports/');
      });
  }
}
