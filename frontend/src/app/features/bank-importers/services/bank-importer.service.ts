import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { BankImporter } from '../../../core/api/models/bank-importer';
import { BankImporterNoId } from '../../../core/api/models/bank-importer-no-id';
import { ImportResult } from '../../../core/api/models/import-result';
import { getBankImporters } from '../../../core/api/fn/bank-importers/get-bank-importers';
import { createBankImporter } from '../../../core/api/fn/bank-importers/create-bank-importer';
import { updateBankImporter } from '../../../core/api/fn/bank-importers/update-bank-importer';
import { deleteBankImporter } from '../../../core/api/fn/bank-importers/delete-bank-importer';
import { uploadBankImporter } from '../../../core/api/fn/bank-importers/upload-bank-importer';

@Injectable({
    providedIn: 'root',
})
export class BankImporterService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);

    readonly bankImporters = signal<BankImporter[]>([]);
    readonly loading = signal(false);
    readonly error = signal<string | null>(null);

    loadBankImporters(): Observable<BankImporter[]> {
        this.loading.set(true);
        this.error.set(null);

        return getBankImporters(this.http, this.apiConfig.rootUrl).pipe(
            map((response) => response.body),
            tap({
                next: (bankImporters) => {
                    this.bankImporters.set(bankImporters);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to load bank importers');
                    this.loading.set(false);
                },
            }),
        );
    }

    create(bankImporter: BankImporterNoId): Observable<BankImporter> {
        this.loading.set(true);
        this.error.set(null);

        return createBankImporter(this.http, this.apiConfig.rootUrl, { body: bankImporter }).pipe(
            map((response) => response.body),
            tap({
                next: (bankImporter) => {
                    this.bankImporters.update((importers) => [...importers, bankImporter]);
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to create bank importer');
                    this.loading.set(false);
                },
            }),
        );
    }

    update(id: string, bankImporter: BankImporterNoId): Observable<BankImporter> {
        this.loading.set(true);
        this.error.set(null);

        return updateBankImporter(this.http, this.apiConfig.rootUrl, {
            id,
            body: bankImporter,
        }).pipe(
            map((response) => response.body),
            tap({
                next: (updatedImporter) => {
                    this.bankImporters.update((importers) =>
                        importers.map((i) => (i.id === id ? updatedImporter : i)),
                    );
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to update bank importer');
                    this.loading.set(false);
                },
            }),
        );
    }

    delete(id: string): Observable<void> {
        this.loading.set(true);
        this.error.set(null);

        return deleteBankImporter(this.http, this.apiConfig.rootUrl, { id }).pipe(
            map(() => undefined),
            tap({
                next: () => {
                    this.bankImporters.update((importers) => importers.filter((i) => i.id !== id));
                    this.loading.set(false);
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to delete bank importer');
                    this.loading.set(false);
                },
            }),
        );
    }

    upload(id: string, file: File, format: 'csv' | 'xlsx'): Observable<ImportResult> {
        this.loading.set(true);
        this.error.set(null);

        return uploadBankImporter(this.http, this.apiConfig.rootUrl, {
            id,
            format,
            body: { file },
        }).pipe(
            map((response) => response.body),
            tap({
                next: () => {
                    this.loading.set(false);
                    // Reload importers to update last import status
                    this.loadBankImporters().subscribe();
                },
                error: (err) => {
                    this.error.set(err.message || 'Failed to upload bank importer file');
                    this.loading.set(false);
                },
            }),
        );
    }
}
