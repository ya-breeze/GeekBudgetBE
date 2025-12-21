import { Component, inject, signal } from '@angular/core';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { BankImporter } from '../../../core/api/models/bank-importer';

export interface BankImporterUploadDialogData {
    bankImporter: BankImporter;
}

export interface BankImporterUploadDialogResult {
    file: File;
    format: 'csv' | 'xlsx';
}

interface ImporterConfig {
    allowedExtensions: string[];
    accept: string;
    defaultFormat: 'csv' | 'xlsx';
    fixedFormat: boolean;
    description: string;
}

const IMPORTER_CONFIGS: Record<string, ImporterConfig> = {
    kb: {
        allowedExtensions: ['csv'],
        accept: '.csv',
        defaultFormat: 'csv',
        fixedFormat: true,
        description: 'Upload CSV file exported from KB internet banking.',
    },
    revolut: {
        allowedExtensions: ['xlsx', 'xls'],
        accept: '.xlsx, .xls',
        defaultFormat: 'xlsx',
        fixedFormat: true,
        description:
            'Upload <b>Excel</b> file exported from Revolut.<br>Make sure to export transactions in the Revolut <b>mobile app</b> not in the website. Because the website exports transactions in a different format and only for a single currency.',
    },
    default: {
        allowedExtensions: ['csv', 'xlsx', 'xls'],
        accept: '.csv, .xlsx, .xls',
        defaultFormat: 'csv',
        fixedFormat: false,
        description: 'Upload transaction file.',
    },
};

@Component({
    selector: 'app-bank-importer-upload-dialog',
    standalone: true,
    imports: [
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
        MatIconModule,
        ReactiveFormsModule,
        FormsModule,
    ],
    templateUrl: './bank-importer-upload-dialog.component.html',
    styleUrl: './bank-importer-upload-dialog.component.scss',
})
export class BankImporterUploadDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<BankImporterUploadDialogComponent>);
    private readonly data = inject<BankImporterUploadDialogData>(MAT_DIALOG_DATA);

    protected readonly bankImporter = this.data.bankImporter;
    protected readonly config =
        IMPORTER_CONFIGS[this.bankImporter.type || ''] || IMPORTER_CONFIGS['default'];

    protected readonly selectedFormat = signal<'csv' | 'xlsx'>(this.config.defaultFormat);
    protected readonly selectedFile = signal<File | null>(null);
    protected readonly isDragging = signal(false);

    protected onFileSelected(event: Event): void {
        const input = event.target as HTMLInputElement;
        if (input.files?.length) {
            this.handleFile(input.files[0]);
        }
    }

    protected onDragOver(event: DragEvent): void {
        event.preventDefault();
        event.stopPropagation();
        this.isDragging.set(true);
    }

    protected onDragLeave(event: DragEvent): void {
        event.preventDefault();
        event.stopPropagation();
        this.isDragging.set(false);
    }

    protected onFileDrop(event: DragEvent): void {
        event.preventDefault();
        event.stopPropagation();
        this.isDragging.set(false);

        if (event.dataTransfer?.files.length) {
            this.handleFile(event.dataTransfer.files[0]);
        }
    }

    private handleFile(file: File): void {
        const extension = file.name.split('.').pop()?.toLowerCase();

        if (extension && this.config.allowedExtensions.includes(extension)) {
            this.selectedFile.set(file);
        } else {
            // Ideally notify user about wrong format
            console.warn(
                `Invalid file format: ${extension}. Allowed: ${this.config.allowedExtensions.join(', ')}`,
            );
        }
    }

    protected onCancel(): void {
        this.dialogRef.close();
    }

    protected onUpload(): void {
        const file = this.selectedFile();
        if (file) {
            this.dialogRef.close({
                file: file,
                format: this.selectedFormat(),
            });
        }
    }
}
