import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { Currency } from '../../../core/api/models/currency';
import { CurrencyNoId } from '../../../core/api/models/currency-no-id';

export interface CurrencyFormDialogData {
    mode: 'create' | 'edit';
    currency?: Currency;
}

@Component({
    selector: 'app-currency-form-dialog',
    imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
    ],
    templateUrl: './currency-form-dialog.component.html',
    styleUrl: './currency-form-dialog.component.scss',
})
export class CurrencyFormDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<CurrencyFormDialogComponent>);
    private readonly data = inject<CurrencyFormDialogData>(MAT_DIALOG_DATA);
    private readonly fb = inject(FormBuilder);

    protected readonly form: FormGroup;
    protected readonly isEditMode = this.data.mode === 'edit';

    constructor() {
        this.form = this.fb.group({
            name: [
                this.data.currency?.name || '',
                [Validators.required, Validators.maxLength(100)],
            ],
            description: [this.data.currency?.description || '', [Validators.maxLength(500)]],
        });
    }

    onSubmit(): void {
        if (this.form.valid) {
            const formValue: CurrencyNoId = this.form.value;
            this.dialogRef.close(formValue);
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
