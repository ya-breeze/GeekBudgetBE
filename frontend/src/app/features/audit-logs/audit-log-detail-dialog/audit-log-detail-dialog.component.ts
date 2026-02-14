import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { AuditLog } from '../../../core/api/models/audit-log';
import { AppDatePipe } from '../../../shared/pipes/app-date.pipe';

@Component({
    selector: 'app-audit-log-detail-dialog',
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        MatIconModule,
        MatListModule,
        MatDividerModule,
        AppDatePipe,
    ],
    templateUrl: './audit-log-detail-dialog.component.html',
    styleUrl: './audit-log-detail-dialog.component.scss',
})
export class AuditLogDetailDialogComponent {
    constructor(
        public dialogRef: MatDialogRef<AuditLogDetailDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: AuditLog,
    ) {}

    get formattedSnapshot(): string {
        if (!this.data.snapshot) return '';
        try {
            return JSON.stringify(JSON.parse(this.data.snapshot), null, 2);
        } catch {
            return this.data.snapshot;
        }
    }

    close(): void {
        this.dialogRef.close();
    }
}
