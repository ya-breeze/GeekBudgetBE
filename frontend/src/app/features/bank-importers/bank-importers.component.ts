import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';
import { DatePipe } from '@angular/common';
import { BankImporterService } from './services/bank-importer.service';
import { BankImporter } from '../../core/api/models/bank-importer';
import { BankImporterFormDialogComponent } from './bank-importer-form-dialog/bank-importer-form-dialog.component';
import {
    BankImporterUploadDialogComponent,
    BankImporterUploadDialogResult,
} from './bank-importer-upload-dialog/bank-importer-upload-dialog.component';
import {
    ImportResultDialogComponent,
    ImportResultDialogData,
} from './import-result-dialog/import-result-dialog.component';
import { AccountService } from '../accounts/services/account.service';
import { LayoutService } from '../../layout/services/layout.service';

import { CurrencyService } from '../currencies/services/currency.service';
import { AccountDisplayComponent } from '../../shared/components/account-display/account-display.component';

@Component({
    selector: 'app-bank-importers',
    imports: [
        MatTableModule,
        MatSortModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatDialogModule,
        MatSnackBarModule,
        MatChipsModule,
        MatTooltipModule,
        DatePipe,
        AccountDisplayComponent,
    ],
    templateUrl: './bank-importers.component.html',
    styleUrl: './bank-importers.component.scss',
})
export class BankImportersComponent implements OnInit {
    private readonly bankImporterService = inject(BankImporterService);
    private readonly accountService = inject(AccountService);
    private readonly dialog = inject(MatDialog);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);
    private readonly currenciesService = inject(CurrencyService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;

    protected readonly loading = this.bankImporterService.loading;
    protected readonly displayedColumns = signal([
        'name',
        'type',
        'account',
        'lastImport',
        'actions',
    ]);

    protected readonly sortActive = signal<string | null>(null);
    protected readonly sortDirection = signal<'asc' | 'desc'>('asc');

    // Computed signal that enriches bank importers with account names and sorts
    protected readonly bankImporters = computed(() => {
        const importers = this.bankImporterService.bankImporters();
        const accounts = this.accountService.accounts();
        const columns = this.displayedColumns();

        // Create a map of account IDs to names for quick lookup
        const accountMap = new Map(accounts.map((acc) => [acc.id, acc.name]));

        // Enrich importers with account names
        const enrichedImporters = importers.map((importer) => ({
            ...importer,
            accountName: accountMap.get(importer.accountId),
            accountImage: accounts.find((a) => a.id === importer.accountId)?.image,
            feeAccountName: importer.feeAccountId
                ? accountMap.get(importer.feeAccountId)
                : undefined,
        }));

        if (!columns.length) {
            return enrichedImporters;
        }

        const active = this.sortActive() ?? columns[0];
        const direction = this.sortDirection();

        // Sort by selected column
        return [...enrichedImporters].sort((a, b) =>
            this.compareBankImporters(a, b, active, direction),
        );
    });

    ngOnInit(): void {
        this.loadBankImporters();
    }

    loadBankImporters(): void {
        this.accountService.loadAccounts().subscribe();
        this.currenciesService.loadCurrencies().subscribe();
        this.bankImporterService.loadBankImporters().subscribe();
    }

    openCreateDialog(): void {
        const dialogRef = this.dialog.open(BankImporterFormDialogComponent, {
            width: '700px',
            data: { mode: 'create' },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result) {
                this.bankImporterService.create(result).subscribe({
                    next: () => {
                        this.snackBar.open('Bank importer created successfully', 'Close', {
                            duration: 3000,
                        });
                    },
                    error: () => {
                        this.snackBar.open('Failed to create bank importer', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        });
    }

    openEditDialog(bankImporter: BankImporter): void {
        const dialogRef = this.dialog.open(BankImporterFormDialogComponent, {
            width: '700px',
            data: { mode: 'edit', bankImporter },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result && bankImporter.id) {
                this.bankImporterService.update(bankImporter.id, result).subscribe({
                    next: () => {
                        this.snackBar.open('Bank importer updated successfully', 'Close', {
                            duration: 3000,
                        });
                    },
                    error: () => {
                        this.snackBar.open('Failed to update bank importer', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        });
    }

    openUploadDialog(bankImporter: BankImporter): void {
        const dialogRef = this.dialog.open<
            BankImporterUploadDialogComponent,
            { bankImporter: BankImporter },
            BankImporterUploadDialogResult
        >(BankImporterUploadDialogComponent, {
            width: '500px',
            data: { bankImporter },
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result && bankImporter.id) {
                this.bankImporterService
                    .upload(bankImporter.id, result.file, result.format)
                    .subscribe({
                        next: (response) => {
                            const currencies = this.currenciesService.currencies();
                            const currencyMap = new Map(currencies.map((c) => [c.id, c.name]));

                            const balances = response.balances?.map((b) => ({
                                amount: b.amount ?? 0,
                                currency: b.currencyId
                                    ? currencyMap.get(b.currencyId) || b.currencyId
                                    : '?',
                            }));

                            this.dialog.open<ImportResultDialogComponent, ImportResultDialogData>(
                                ImportResultDialogComponent,
                                {
                                    width: '500px',
                                    data: {
                                        title: 'Import Successful',
                                        message:
                                            response.description ||
                                            'Transactions uploaded successfully',
                                        status: 'success',
                                        balances,
                                    },
                                },
                            );
                        },
                        error: (err) => {
                            this.dialog.open<ImportResultDialogComponent, ImportResultDialogData>(
                                ImportResultDialogComponent,
                                {
                                    width: '500px',
                                    data: {
                                        title: 'Import Failed',
                                        message: err || 'Failed to upload transactions',
                                        status: 'error',
                                    },
                                },
                            );
                        },
                    });
            }
        });
    }

    deleteBankImporter(bankImporter: BankImporter): void {
        if (confirm(`Are you sure you want to delete "${bankImporter.name}"?`)) {
            if (bankImporter.id) {
                this.bankImporterService.delete(bankImporter.id).subscribe({
                    next: () => {
                        this.snackBar.open('Bank importer deleted successfully', 'Close', {
                            duration: 3000,
                        });
                    },
                    error: () => {
                        this.snackBar.open('Failed to delete bank importer', 'Close', {
                            duration: 3000,
                        });
                    },
                });
            }
        }
    }

    resumeBankImporter(bankImporter: BankImporter): void {
        if (!bankImporter.id) return;

        // Create full object for update (BankImporterNoId)
        const updateData: any = {
            name: bankImporter.name,
            type: bankImporter.type,
            accountId: bankImporter.accountId,
            feeAccountId: bankImporter.feeAccountId,
            description: bankImporter.description,
            extra: bankImporter.extra,
            fetchAll: bankImporter.fetchAll,
            isStopped: false, // Resume
        };

        this.bankImporterService.update(bankImporter.id, updateData).subscribe({
            next: () => {
                this.snackBar.open('Resuming bank importer...', 'Close', { duration: 3000 });
            },
            error: () => {
                this.snackBar.open('Failed to resume bank importer', 'Close', { duration: 3000 });
            },
        });
    }

    getBankTypeLabel(type?: string): string {
        switch (type) {
            case 'fio':
                return 'FIO Bank';
            case 'kb':
                return 'KB Bank';
            case 'revolut':
                return 'Revolut';
            default:
                return type || '-';
        }
    }

    protected onSortChange(sort: Sort): void {
        if (!sort.direction) {
            this.sortActive.set(null);
            this.sortDirection.set('asc');
            return;
        }

        this.sortActive.set(sort.active);
        this.sortDirection.set(sort.direction);
    }

    private compareBankImporters(
        a: BankImporter & {
            accountName?: string;
            accountImage?: string;
            feeAccountName?: string;
        },
        b: BankImporter & {
            accountName?: string;
            accountImage?: string;
            feeAccountName?: string;
        },
        active: string,
        direction: 'asc' | 'desc',
    ): number {
        const valueA = this.getBankImporterSortValue(a, active);
        const valueB = this.getBankImporterSortValue(b, active);
        return this.comparePrimitiveValues(valueA, valueB, direction);
    }

    private getBankImporterSortValue(
        importer: BankImporter & {
            accountName?: string;
            accountImage?: string;
            feeAccountName?: string;
        },
        active: string,
    ): string | Date | null {
        switch (active) {
            case 'name':
                return this.removeLeadingEmoji(importer.name ?? '');
            case 'type':
                return this.getBankTypeLabel(importer.type);
            case 'account':
                return importer.accountName ?? importer.accountId ?? '';
            case 'lastImport':
                return importer.lastSuccessfulImport
                    ? new Date(importer.lastSuccessfulImport)
                    : null;
            default:
                return null;
        }
    }

    private removeLeadingEmoji(text: string): string {
        // Remove leading emoji characters for sorting purposes
        return text
            .replace(/^[\p{Emoji}\p{Emoji_Presentation}\p{Emoji_Modifier_Base}\s]+/u, '')
            .trim();
    }

    protected getImportsSummary(importer: BankImporter): string {
        if (!importer.lastImports || importer.lastImports.length === 0) {
            return 'No imports yet';
        }
        const total = importer.lastImports.length;
        const success = importer.lastImports.filter((i) => i.status === 'success').length;
        const error = total - success;
        return `Last ${total}: ${success} success, ${error} error`;
    }

    protected getImportsStatusColor(importer: BankImporter): string {
        if (!importer.lastImports || importer.lastImports.length === 0) {
            return '';
        }
        const hasError = importer.lastImports.some((i) => i.status === 'error');
        if (hasError) return '#f44336'; // Red
        return '#4caf50'; // Green
    }

    private comparePrimitiveValues(
        a: string | number | Date | null | undefined,
        b: string | number | Date | null | undefined,
        direction: 'asc' | 'desc',
    ): number {
        const factor = direction === 'asc' ? 1 : -1;

        if (a == null && b == null) return 0;
        if (a == null) return 1 * factor;
        if (b == null) return -1 * factor;

        if (typeof a === 'string' && typeof b === 'string') {
            return a.localeCompare(b) * factor;
        }

        if (typeof a === 'number' && typeof b === 'number') {
            return (a - b) * factor;
        }

        if (a instanceof Date && b instanceof Date) {
            return (a.getTime() - b.getTime()) * factor;
        }

        return `${a}`.localeCompare(`${b}`) * factor;
    }
}
