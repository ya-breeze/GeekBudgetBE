import { Injectable, inject, signal, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {
    Observable,
    Subscription,
    interval,
    switchMap,
    map,
    tap,
    catchError,
    of,
    startWith,
} from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Notification } from '../../../core/api/models/notification';
import { getNotifications } from '../../../core/api/fn/notifications/get-notifications';
import { deleteNotification } from '../../../core/api/fn/notifications/delete-notification';

@Injectable({
    providedIn: 'root',
})
export class NotificationService implements OnDestroy {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly snackBar = inject(MatSnackBar);

    readonly notifications = signal<Notification[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    private pollingSubscription?: Subscription;
    private initialLoadDone = false;

    constructor() {
        // Polling will be started explicitly to avoid multiple subscriptions if the service is injected multiple times
        // although it's a singleton.
    }

    ngOnDestroy(): void {
        this.stopPolling();
    }

    loadNotifications(): Observable<Notification[]> {
        this.loading.set(true);
        this.error.set(null);

        return getNotifications(this.http, this.apiConfig.rootUrl).pipe(
            map((response) => response.body),
            tap({
                next: (notifications) => {
                    if (this.initialLoadDone) {
                        const currentIds = new Set(this.notifications().map((n) => n.id));
                        const newNotifications = notifications.filter((n) => !currentIds.has(n.id));

                        if (newNotifications.length > 0) {
                            const latest = newNotifications[0];
                            this.snackBar
                                .open(`${latest.title}: ${latest.description}`, 'View', {
                                    duration: 5000,
                                })
                                .onAction()
                                .subscribe(() => {
                                    // Potentially navigate or handle action
                                });
                        }
                    }

                    this.notifications.set(notifications);
                    this.loading.set(false);
                    this.initialLoadDone = true;
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load notifications');
                    this.loading.set(false);
                },
            }),
            catchError((err) => {
                this.error.set(err.message || 'Failed to load notifications');
                this.loading.set(false);
                return of([]);
            }),
        );
    }

    delete(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return deleteNotification(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.notifications.update((notifications) =>
                        notifications.filter((n) => n.id !== id),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to delete notification');
                    this.loading.set(false);
                },
            }),
            catchError((err) => {
                this.error.set(err.message || 'Failed to delete notification');
                this.loading.set(false);
                return of(void 0);
            }),
        );
    }

    startPolling(): void {
        if (this.pollingSubscription) {
            return;
        }
        // Poll every 1 minute (60,000 ms)
        this.pollingSubscription = interval(60000)
            .pipe(
                startWith(0),
                switchMap(() => this.loadNotifications()),
            )
            .subscribe();
    }

    stopPolling(): void {
        if (this.pollingSubscription) {
            this.pollingSubscription.unsubscribe();
            this.pollingSubscription = undefined;
        }
    }
}
