import { Component, inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { DecimalPipe } from '@angular/common';

export interface ImportResultDialogData {
    title: string;
    message: string;
    status: 'success' | 'error';
    balances?: { amount: number; currency: string }[];
}

@Component({
    selector: 'app-import-result-dialog',
    standalone: true,
    imports: [MatDialogModule, MatButtonModule, MatIconModule, DecimalPipe],
    templateUrl: './import-result-dialog.component.html',
    styleUrl: './import-result-dialog.component.scss',
})
export class ImportResultDialogComponent {
    protected readonly data = inject<ImportResultDialogData>(MAT_DIALOG_DATA);
}
