import { Component, inject, signal } from '@angular/core';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { Currency } from '../../../core/api/models/currency';

export interface CurrencyDeleteDialogData {
    currencyId: string;
    currencyName: string;
    message: string;
    isUsed: boolean;
    currencies: Currency[];
    selectedReplacementId?: string;
}

@Component({
    selector: 'app-currency-delete-dialog',
    standalone: true,
    imports: [
        MatDialogModule,
        MatFormFieldModule,
        MatSelectModule,
        MatButtonModule,
        ReactiveFormsModule,
    ],
    templateUrl: './currency-delete-dialog.component.html',
    styles: [
        `
            .full-width {
                width: 100%;
                margin-top: 16px;
            }
            .warning {
                color: #f44336;
                font-weight: 500;
                margin-bottom: 16px;
            }
        `,
    ],
})
export class CurrencyDeleteDialogComponent {
    private readonly dialogRef = inject(MatDialogRef<CurrencyDeleteDialogComponent>);
    protected readonly data = inject<CurrencyDeleteDialogData>(MAT_DIALOG_DATA);

    protected readonly replacementControl = new FormControl<string | null>(
        this.data.selectedReplacementId || null,
    );

    protected readonly usableCurrencies = signal(
        this.data.currencies.filter((c) => c.id !== this.data.currencyId),
    );

    constructor() {
        if (this.data.isUsed) {
            this.replacementControl.setValidators([Validators.required]);
        }
    }

    onConfirm(): void {
        if (this.data.isUsed && this.replacementControl.invalid) {
            this.replacementControl.markAsTouched();
            return;
        }

        this.dialogRef.close({
            replaceWithCurrencyId: this.replacementControl.value,
        });
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
