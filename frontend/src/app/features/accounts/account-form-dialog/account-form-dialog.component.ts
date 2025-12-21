import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatIconModule } from '@angular/material/icon';
import { Account } from '../../../core/api/models/account';
import { AccountNoId } from '../../../core/api/models/account-no-id';
import { ApiConfiguration } from '../../../core/api/api-configuration';

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
        MatIconModule,
    ],
    templateUrl: './account-form-dialog.component.html',
    styleUrl: './account-form-dialog.component.scss',
})
export class AccountFormDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<AccountFormDialogComponent>);
    private readonly data = inject<AccountFormDialogData>(MAT_DIALOG_DATA);
    private readonly fb = inject(FormBuilder);
    private readonly apiConfig = inject(ApiConfiguration);

    protected readonly form: FormGroup;
    protected readonly isEditMode = this.data.mode === 'edit';
    protected readonly accountTypes = [
        { value: 'expense', label: 'Expense' },
        { value: 'income', label: 'Income' },
        { value: 'asset', label: 'Asset' },
    ];

    protected selectedFile: File | null = null;
    protected imagePreview: string | null = null;
    protected deleteImage = false;

    constructor() {
        this.form = this.fb.group({
            name: [this.data.account?.name || '', [Validators.required, Validators.maxLength(100)]],
            type: [this.data.account?.type || 'expense', [Validators.required]],
            description: [this.data.account?.description || '', [Validators.maxLength(500)]],
            showInDashboardSummary: [this.data.account?.showInDashboardSummary ?? true],
            hideFromReports: [this.data.account?.hideFromReports ?? false],
        });

        if (this.data.account?.image) {
            const root = this.apiConfig.rootUrl.endsWith('/')
                ? this.apiConfig.rootUrl.slice(0, -1)
                : this.apiConfig.rootUrl;
            this.imagePreview = `${root}/images/${this.data.account.image}`;
        }
    }

    onFileSelected(event: Event): void {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            this.selectedFile = input.files[0];
            this.deleteImage = false;

            // Create preview
            const reader = new FileReader();
            reader.onload = () => {
                this.imagePreview = reader.result as string;
            };
            reader.readAsDataURL(this.selectedFile);
        }
    }

    removeImage(): void {
        this.selectedFile = null;
        this.imagePreview = null;
        this.deleteImage = true;
    }

    onSubmit(): void {
        if (this.form.valid) {
            const formValue: AccountNoId = this.form.value;
            this.dialogRef.close({
                account: formValue,
                image: this.selectedFile,
                deleteImage: this.deleteImage,
            });
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
