import { Component, Inject, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { AuditLog } from '../../../core/api/models/audit-log';
import { AppDatePipe } from '../../../shared/pipes/app-date.pipe';

interface DiffRow {
    key: string;
    before: string;
    after: string;
    status: 'unchanged' | 'added' | 'removed' | 'modified';
}

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
    encapsulation: ViewEncapsulation.None, // Force styles to apply globally (scoped by selectors)
})
export class AuditLogDetailDialogComponent {
    diffRows: DiffRow[] = [];

    constructor(
        public dialogRef: MatDialogRef<AuditLogDetailDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: AuditLog,
    ) {
        console.log('AuditLog Data:', data);
        this.diffRows = this.calculateDiff(data.before, data.after);
        console.log('Diff Rows:', this.diffRows);
    }

    private calculateDiff(beforeStr?: string | null, afterStr?: string | null): DiffRow[] {
        const parse = (str?: string | null) => {
            if (!str) return null;
            try {
                return JSON.parse(str);
            } catch {
                return null;
            }
        };

        const beforeObj = parse(beforeStr);
        const afterObj = parse(afterStr);

        // If both are null/empty, nothing to show
        if (!beforeObj && !afterObj) return [];

        // If one is missing (Create/Delete), we can just show everything as added/removed
        // But the flatten logic handles it if we pass {} for the missing one.

        const flatBefore = this.flattenObject(beforeObj || {});
        const flatAfter = this.flattenObject(afterObj || {});

        const allKeys = new Set([...Object.keys(flatBefore), ...Object.keys(flatAfter)]);
        const rows: DiffRow[] = [];

        allKeys.forEach((key) => {
            const valBefore = flatBefore[key];
            const valAfter = flatAfter[key];
            const strBefore = valBefore !== undefined ? JSON.stringify(valBefore) : '';
            const strAfter = valAfter !== undefined ? JSON.stringify(valAfter) : '';

            let status: DiffRow['status'] = 'unchanged';

            if (valBefore === undefined && valAfter !== undefined) {
                status = 'added';
            } else if (valBefore !== undefined && valAfter === undefined) {
                status = 'removed';
            } else if (strBefore !== strAfter) {
                status = 'modified';
            }

            // Only show rows that have changes
            if (status !== 'unchanged') {
                rows.push({
                    key,
                    before: valBefore !== undefined ? String(valBefore) : '',
                    after: valAfter !== undefined ? String(valAfter) : '',
                    status,
                });
            }
        });

        return rows.sort((a, b) => a.key.localeCompare(b.key));
    }

    private flattenObject(obj: any, prefix = ''): Record<string, any> {
        const result: Record<string, any> = {};

        // Handle null/undefined
        if (obj === null || obj === undefined) {
            return result;
        }

        // Handle arrays specifically to preserve indices
        if (Array.isArray(obj)) {
            obj.forEach((value, index) => {
                const newKey = prefix ? `${prefix}[${index}]` : `[${index}]`;
                if (value && typeof value === 'object') {
                    const nested = this.flattenObject(value, newKey);
                    Object.assign(result, nested);
                } else {
                    result[newKey] = value;
                }
            });
            return result;
        }

        // Handle objects
        for (const key in obj) {
            if (Object.prototype.hasOwnProperty.call(obj, key)) {
                const value = obj[key];
                const newKey = prefix ? `${prefix}.${key}` : key;

                if (value && typeof value === 'object' && !Array.isArray(value)) {
                    const nested = this.flattenObject(value, newKey);
                    Object.assign(result, nested);
                } else if (Array.isArray(value)) {
                    // Recurse into array property
                    const nested = this.flattenObject(value, newKey);
                    Object.assign(result, nested);
                } else {
                    result[newKey] = value;
                }
            }
        }
        return result;
    }

    close(): void {
        this.dialogRef.close();
    }
}
