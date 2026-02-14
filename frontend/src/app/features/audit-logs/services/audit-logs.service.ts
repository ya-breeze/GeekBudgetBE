import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { AuditLog } from '../../../core/api/models/audit-log';
import { getAuditLogs, GetAuditLogs$Params } from '../../../core/api/fn/audit-logs/get-audit-logs';

@Injectable({
    providedIn: 'root',
})
export class AuditLogsService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly auditLogs = signal<AuditLog[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadAuditLogs(params?: GetAuditLogs$Params): Observable<AuditLog[]> {
        this.loading.set(true);
        this.error.set(null);

        return getAuditLogs(this.http, this.apiConfig.rootUrl, params).pipe(
            map((response) => response.body),
            tap({
                next: (logs) => {
                    this.auditLogs.set(logs);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load audit logs');
                    this.loading.set(false);
                },
            }),
        );
    }
}
