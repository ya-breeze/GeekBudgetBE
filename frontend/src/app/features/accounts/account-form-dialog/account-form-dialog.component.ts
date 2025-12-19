import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { Account } from '../../../core/api/models/account';
import { AccountNoId } from '../../../core/api/models/account-no-id';

export interface AccountFormDialogData {
    mode: 'create' | 'edit';
    account?: Account;
}

@Component({
    selector: 'app-account-form-dialog',
    imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatSelectModule,
        MatSlideToggleModule,
    ],
    templateUrl: './account-form-dialog.component.html',
    styleUrl: './account-form-dialog.component.scss',
})
export class AccountFormDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<AccountFormDialogComponent>);
    private readonly data = inject<AccountFormDialogData>(MAT_DIALOG_DATA);
    private readonly fb = inject(FormBuilder);

    protected readonly form: FormGroup;
    protected readonly isEditMode = this.data.mode === 'edit';
    protected readonly accountTypes = [
        { value: 'expense', label: 'Expense' },
        { value: 'income', label: 'Income' },
        { value: 'asset', label: 'Asset' },
    ];

    constructor() {
        this.form = this.fb.group({
            name: [this.data.account?.name || '', [Validators.required, Validators.maxLength(100)]],
            type: [this.data.account?.type || 'expense', [Validators.required]],
            description: [this.data.account?.description || '', [Validators.maxLength(500)]],
            showInDashboardSummary: [this.data.account?.showInDashboardSummary ?? true],
        });
    }

    onSubmit(): void {
        if (this.form.valid) {
            const formValue: AccountNoId = this.form.value;
            this.dialogRef.close(formValue);
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
