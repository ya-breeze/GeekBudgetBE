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
    protected readonly selectedFormat = signal<'csv' | 'xlsx'>('csv');
    protected readonly selectedFile = signal<File | null>(null);
    protected readonly isDragging = signal(false);

    protected onFileSelected(event: Event): void {
        const input = event.target as HTMLInputElement;
        if (input.files?.length) {
            this.selectedFile.set(input.files[0]);
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
            const file = event.dataTransfer.files[0];
            const extension = file.name.split('.').pop()?.toLowerCase();

            if (this.selectedFormat() === 'csv' && extension === 'csv') {
                this.selectedFile.set(file);
            } else if (
                this.selectedFormat() === 'xlsx' &&
                (extension === 'xlsx' || extension === 'xls')
            ) {
                this.selectedFile.set(file);
            } else {
                // Ideally show an error, but for now just ignore or let the user know via logic
                // If the format doesn't match, maybe just set it anyway or check extension strictly?
                // For simplicity, let's just accept it and let the backend/user handle mismatch,
                // or align with the input accept attribute logic.
                // The input has accept=".csv" or ".xlsx, .xls".
                if (this.selectedFormat() === 'csv') {
                    if (extension === 'csv') this.selectedFile.set(file);
                } else {
                    if (extension === 'xlsx' || extension === 'xls') this.selectedFile.set(file);
                }
            }
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
