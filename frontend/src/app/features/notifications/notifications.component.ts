import { Component, inject, OnInit, signal, WritableSignal } from '@angular/core';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { AppDatePipe } from '../../shared/pipes/app-date.pipe';
import { NotificationService } from './services/notification.service';
import { Notification } from '../../core/api/models/notification';
import { LayoutService } from '../../layout/services/layout.service';

@Component({
    selector: 'app-notifications',
    imports: [
        MatListModule,
        MatIconModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MatChipsModule,
        MatSnackBarModule,
        AppDatePipe,
    ],
    template: `
        <div class="notifications-container">
            @if (!sidenavOpened()) {
                <h1 class="page-title">Notifications</h1>
            }

            @if (loading()) {
                <div class="loading-container">
                    <mat-spinner></mat-spinner>
                </div>
            } @else {
                <mat-list>
                    @for (notification of notifications(); track notification.id) {
                        <mat-list-item
                            (click)="toggleExpand(notification)"
                            [class.expanded-item]="expandedIds().has(notification.id)"
                            class="notification-item"
                        >
                            <mat-icon matListItemIcon [class]="getIconClass(notification.type)">
                                {{ getIcon(notification.type) }}
                            </mat-icon>
                            <div matListItemTitle>{{ notification.title }}</div>
                            <div
                                class="description"
                                [class.expanded]="expandedIds().has(notification.id)"
                            >
                                {{ notification.description }}
                            </div>
                            <div matListItemLine class="date">
                                {{ notification.date | appDate: 'short' }}
                            </div>
                            <button
                                mat-icon-button
                                matListItemMeta
                                (click)="deleteNotification(notification, $event)"
                            >
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
            padding: 0;
        }
        .page-title {
            margin: 0 0 16px 0;
            font-size: 24px;
            font-weight: 500;
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
        .notification-item {
            cursor: pointer;
        }
        .description {
            white-space: normal;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
            margin-bottom: 4px;
            word-break: break-word;
            line-height: 1.4;
        }
        .description.expanded {
            -webkit-line-clamp: unset;
            display: block;
        }
        /* Crucial: Override MDC list item fixed height to allow expansion */
        .notification-item {
            height: auto !important;
            padding-top: 8px;
            padding-bottom: 8px;
        }
        /* Ensure the content container handles variable height */
        ::ng-deep .mdc-list-item__content {
            height: auto !important;
            align-self: flex-start !important; /* Fix alignment for variable height */
        }
        /* Allow lines to wrap and expand */
        ::ng-deep .mat-mdc-list-item-unscoped-content {
            display: flex;
            flex-direction: column;
            white-space: normal !important;
        }
    `,
})
export class NotificationsComponent implements OnInit {
    private readonly notificationService = inject(NotificationService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;

    protected readonly notifications = this.notificationService.notifications;
    protected readonly loading = this.notificationService.loading;
    protected readonly expandedIds: WritableSignal<Set<string>> = signal(new Set<string>());

    ngOnInit(): void {
        this.notificationService.loadNotifications().subscribe();
    }

    deleteNotification(notification: Notification, event: Event): void {
        event.stopPropagation();
        if (notification.id) {
            this.notificationService.delete(notification.id).subscribe({
                next: () => {
                    this.snackBar.open('Notification deleted', 'Close', { duration: 3000 });
                },
                error: () => {
                    this.snackBar.open('Failed to delete notification', 'Close', {
                        duration: 3000,
                    });
                },
            });
        }
    }

    toggleExpand(notification: Notification): void {
        const id = notification.id;
        if (!id) return;
        this.expandedIds.update((set) => {
            const newSet = new Set(set);
            if (newSet.has(id)) {
                newSet.delete(id);
            } else {
                newSet.add(id);
            }
            return newSet;
        });
    }

    getIcon(type: string): string {
        switch (type) {
            case 'balanceMatch':
                return 'check_circle';
            case 'balanceDoesntMatch':
                return 'error';
            case 'error':
                return 'report';
            case 'info':
                return 'info';
            default:
                return 'notifications';
        }
    }

    getIconClass(type: string): string {
        switch (type) {
            case 'balanceMatch':
                return 'icon-match';
            case 'balanceDoesntMatch':
            case 'error':
                return 'icon-no-match';
            case 'info':
                return 'icon-other';
            default:
                return 'icon-other';
        }
    }
}
