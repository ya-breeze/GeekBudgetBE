import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Notification } from '../../../core/api/models/notification';
import { getNotifications } from '../../../core/api/fn/notifications/get-notifications';
import { deleteNotification } from '../../../core/api/fn/notifications/delete-notification';

@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly notifications = signal<Notification[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadNotifications(): Observable<Notification[]> {
    this.loading.set(true);
    this.error.set(null);

    return getNotifications(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (notifications) => {
          this.notifications.set(notifications);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load notifications');
          this.loading.set(false);
        },
      })
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
            notifications.filter((n) => n.id !== id)
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete notification');
          this.loading.set(false);
        },
      })
    );
  }
}

