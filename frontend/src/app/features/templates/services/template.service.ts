import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map, tap } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';
import { TransactionTemplateNoId } from '../../../core/api/models/transaction-template-no-id';
import { TransactionNoId } from '../../../core/api/models/transaction-no-id';
import { getTemplates } from '../../../core/api/fn/templates/get-templates';
import { createTemplate } from '../../../core/api/fn/templates/create-template';
import { updateTemplate } from '../../../core/api/fn/templates/update-template';
import { deleteTemplate } from '../../../core/api/fn/templates/delete-template';

@Injectable({
    providedIn: 'root',
})
export class TemplateService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly templates = signal<TransactionTemplate[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadTemplates(accountId?: string): Observable<TransactionTemplate[]> {
        this.loading.set(true);
        this.error.set(null);

        return getTemplates(this.http, this.apiConfig.rootUrl, accountId ? { accountId } : undefined).pipe(
            map((response) => response.body ?? []),
            tap({
                next: (templates) => {
                    this.templates.set(templates);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load templates');
                    this.loading.set(false);
                },
            }),
        );
    }

    create(template: TransactionTemplateNoId): Observable<TransactionTemplate> {
        this.loading.set(true);
        this.error.set(null);

        return createTemplate(this.http, this.apiConfig.rootUrl, { body: template }).pipe(
            map((response) => response.body),
            tap({
                next: (t) => {
                    this.templates.update((templates) => [...templates, t]);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to create template');
                    this.loading.set(false);
                },
            }),
        );
    }

    update(id: string, template: TransactionTemplateNoId): Observable<TransactionTemplate> {
        this.loading.set(true);
        this.error.set(null);

        return updateTemplate(this.http, this.apiConfig.rootUrl, { id, body: template }).pipe(
            map((response) => response.body),
            tap({
                next: (updated) => {
                    this.templates.update((templates) =>
                        templates.map((t) => (t.id === id ? updated : t)),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to update template');
                    this.loading.set(false);
                },
            }),
        );
    }

    delete(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return deleteTemplate(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.templates.update((templates) => templates.filter((t) => t.id !== id));
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to delete template');
                    this.loading.set(false);
                },
            }),
        );
    }

    /**
     * Converts a template to a TransactionNoId payload for creating a transaction.
     * Sets date to today. Zeroes out all import-only fields.
     */
    templateToTransactionNoId(template: TransactionTemplate): TransactionNoId {
        return {
            date: new Date().toISOString(),
            description: template.description ?? '',
            place: template.place ?? '',
            tags: template.tags ?? [],
            partnerName: template.partnerName ?? '',
            extra: template.extra ?? '',
            movements: template.movements ?? [],
            // Zero out import-only fields
            externalIds: [],
            isAuto: false,
        };
    }
}
