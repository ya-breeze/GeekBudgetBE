import { Component, inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { TemplatePickerComponent } from './template-picker.component';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';

@Component({
    selector: 'app-template-picker-dialog',
    standalone: true,
    imports: [MatDialogModule, MatButtonModule, TemplatePickerComponent],
    template: `
        <h2 mat-dialog-title>Choose a template</h2>
        <mat-dialog-content>
            <app-template-picker
                [accountId]="data.accountId"
                (templateSelected)="onSelect($event)"
            />
        </mat-dialog-content>
        <mat-dialog-actions align="end">
            <button mat-button mat-dialog-close>Cancel</button>
        </mat-dialog-actions>
    `,
})
export class TemplatePickerDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<TemplatePickerDialogComponent>);
    protected readonly data = inject<{ accountId?: string }>(MAT_DIALOG_DATA);

    protected onSelect(template: TransactionTemplate): void {
        this.dialogRef.close(template);
    }
}
