import { Component, inject, OnInit, signal } from '@angular/core';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { DatePipe } from '@angular/common';
import { NotificationService } from './services/notification.service';
import { Notification } from '../../core/api/models/notification';

@Component({
  selector: 'app-notifications',
  imports: [
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatSnackBarModule,
    DatePipe,
  ],
  template: `
    <div class="notifications-container">
      <h1>Notifications</h1>
      <p class="subtitle">System notifications and alerts</p>

      @if (loading()) {
      <div class="loading-container">
        <mat-spinner></mat-spinner>
      </div>
      } @else {
      <mat-list>
        @for (notification of notifications(); track notification.id) {
        <mat-list-item>
          <mat-icon matListItemIcon [class]="getIconClass(notification.type)">
            {{ getIcon(notification.type) }}
          </mat-icon>
          <div matListItemTitle>{{ notification.title }}</div>
          <div matListItemLine>{{ notification.description }}</div>
          <div matListItemLine class="date">{{ notification.date | date : 'short' }}</div>
          <button mat-icon-button matListItemMeta (click)="deleteNotification(notification)">
            <mat-icon>delete</mat-icon>
          </button>
        </mat-list-item>
        <mat-divider></mat-divider>
        } @empty {
        <mat-list-item>
          <div matListItemTitle>No notifications</div>
        </mat-list-item>
        }
      </mat-list>
      }
    </div>
  `,
  styles: `
    .notifications-container {
      padding: 24px;
    }
    .subtitle {
      margin: 0 0 24px 0;
      color: rgba(0, 0, 0, 0.6);
    }
    .loading-container {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 300px;
    }
    .date {
      font-size: 12px;
      color: rgba(0, 0, 0, 0.6);
    }
    .icon-other {
      color: #2196f3;
    }
    .icon-match {
      color: #4caf50;
    }
    .icon-no-match {
      color: #f44336;
    }
  `,
})
export class NotificationsComponent implements OnInit {
  private readonly notificationService = inject(NotificationService);
  private readonly snackBar = inject(MatSnackBar);

  protected readonly notifications = this.notificationService.notifications;
  protected readonly loading = this.notificationService.loading;

  ngOnInit(): void {
    this.notificationService.loadNotifications().subscribe();
  }

  deleteNotification(notification: Notification): void {
    if (notification.id) {
      this.notificationService.delete(notification.id).subscribe({
        next: () => {
          this.snackBar.open('Notification deleted', 'Close', { duration: 3000 });
        },
        error: () => {
          this.snackBar.open('Failed to delete notification', 'Close', { duration: 3000 });
        },
      });
    }
  }

  getIcon(type: string): string {
    switch (type) {
      case 'balanceMatch':
        return 'check_circle';
      case 'balanceDoesntMatch':
        return 'error';
      default:
        return 'info';
    }
  }

  getIconClass(type: string): string {
    switch (type) {
      case 'balanceMatch':
        return 'icon-match';
      case 'balanceDoesntMatch':
        return 'icon-no-match';
      default:
        return 'icon-other';
    }
  }
}
