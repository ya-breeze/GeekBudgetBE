import { Component, Inject, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { FormsModule } from '@angular/forms';
import { Account } from '../../../core/api/models/account';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';
import { AccountDisplayComponent } from '../../../shared/components/account-display/account-display.component';

export interface AccountDeleteConfirmDialogData {
    accountToDelete: Account;
    availableAccounts: Account[];
}

export interface AccountDeleteConfirmDialogResult {
    confirmed: boolean;
    replaceWithAccountId?: string;
}

@Component({
    selector: 'app-account-delete-confirm-dialog',
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        FormsModule,
        AccountSelectComponent,
        AccountDisplayComponent,
    ],
    templateUrl: './account-delete-confirm-dialog.component.html',
    styles: [
        `
            mat-form-field {
                width: 100%;
                margin-top: 16px;
            }
        `,
    ],
})
export class AccountDeleteConfirmDialogComponent {
    private dialogRef = inject(MatDialogRef<AccountDeleteConfirmDialogComponent>);

    protected replaceWithAccountId = signal<string | null>(null);
    protected accountToDelete: Account;
    protected availableAccounts: Account[];

    constructor(@Inject(MAT_DIALOG_DATA) data: AccountDeleteConfirmDialogData) {
        this.accountToDelete = data.accountToDelete;
        this.availableAccounts = data.availableAccounts;
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    onConfirm(): void {
        this.dialogRef.close({
            confirmed: true,
            replaceWithAccountId: this.replaceWithAccountId(),
        } as AccountDeleteConfirmDialogResult);
    }
}
