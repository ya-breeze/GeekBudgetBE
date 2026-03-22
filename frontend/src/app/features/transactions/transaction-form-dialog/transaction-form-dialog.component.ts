import { Component, inject, OnInit, signal } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule, MatDialog } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { Transaction } from '../../../core/api/models/transaction';
import { TransactionNoId } from '../../../core/api/models/transaction-no-id';
import { Movement } from '../../../core/api/models/movement';
import { AccountService } from '../../accounts/services/account.service';

import { CurrencyService } from '../../currencies/services/currency.service';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';
import { TemplatePickerDialogComponent } from '../../templates/template-picker/template-picker-dialog.component';
import { TemplateService } from '../../templates/services/template.service';

export interface TransactionFormDialogData {
    mode: 'create' | 'edit';
    transaction?: Transaction;
    initialValues?: TransactionNoId;
}

@Component({
    selector: 'app-transaction-form-dialog',
    imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatSelectModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatIconModule,
        MatChipsModule,
        AccountSelectComponent,
    ],
    templateUrl: './transaction-form-dialog.component.html',
    styleUrl: './transaction-form-dialog.component.scss',
})
export class TransactionFormDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<TransactionFormDialogComponent>);
    private readonly data = inject<TransactionFormDialogData>(MAT_DIALOG_DATA);
    private readonly fb = inject(FormBuilder);
    private readonly accountService = inject(AccountService);
    private readonly currencyService = inject(CurrencyService);
    private readonly dialog = inject(MatDialog);
    private readonly templateService = inject(TemplateService);

    protected readonly form: FormGroup;
    protected readonly isEditMode = this.data.mode === 'edit';
    protected readonly accounts = this.accountService.accounts;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly tags = signal<string[]>([]);

    constructor() {
        const src = this.data.transaction ?? this.data.initialValues;
        this.form = this.fb.group({
            date: [
                src?.date ? new Date(src.date) : new Date(),
                [Validators.required],
            ],
            description: [src?.description || '', [Validators.maxLength(500)]],
            movements: this.fb.array([], [Validators.required, Validators.minLength(1)]),
            partnerName: [src?.partnerName || ''],
            partnerAccount: [(this.data.transaction ?? this.data.initialValues as any)?.partnerAccount || ''],
            place: [src?.place || ''],
        });

        const tags = this.data.transaction?.tags ?? this.data.initialValues?.tags;
        if (tags) {
            this.tags.set([...tags]);
        }
    }

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();

        const movements = this.data.transaction?.movements ?? this.data.initialValues?.movements;
        if (movements?.length) {
            movements.forEach((movement) => {
                this.addMovement(movement);
            });
        } else {
            this.addMovement();
        }
    }

    get movements(): FormArray {
        return this.form.get('movements') as FormArray;
    }

    addMovement(movement?: Movement): void {
        const movementGroup = this.fb.group({
            accountId: [movement?.accountId || '', [Validators.required]],
            currencyId: [movement?.currencyId || '', [Validators.required]],
            amount: [movement?.amount || 0, [Validators.required]],
            description: [movement?.description || ''],
        });
        this.movements.push(movementGroup);
    }

    removeMovement(index: number): void {
        this.movements.removeAt(index);
    }

    addTag(event: Event): void {
        const input = event.target as HTMLInputElement;
        const value = input.value.trim();
        if (value) {
            this.tags.update((tags) => [...tags, value]);
            input.value = '';
        }
    }

    removeTag(tag: string): void {
        this.tags.update((tags) => tags.filter((t) => t !== tag));
    }

    protected useTemplate(): void {
        const dialogRef = this.dialog.open(TemplatePickerDialogComponent, { width: '400px', data: {} });
        dialogRef.afterClosed().subscribe((template) => {
            if (template) {
                const tx = this.templateService.templateToTransactionNoId(template);
                this.form.patchValue({
                    description: tx.description,
                    place: tx.place,
                    partnerName: tx.partnerName,
                });
                this.movements.clear();
                (tx.movements ?? []).forEach((m) => this.addMovement(m));
                this.tags.set(tx.tags ?? []);
            }
        });
    }

    onSubmit(): void {
        if (this.form.valid) {
            const formValue = this.form.value;

            // Start with base object if editing (excluding id)
            let baseTransaction = {};
            if (this.isEditMode && this.data.transaction) {
                // eslint-disable-next-line @typescript-eslint/no-unused-vars
                const { id, ...rest } = this.data.transaction;
                baseTransaction = rest;
            }

            const transaction: TransactionNoId = {
                ...baseTransaction,
                date: formValue.date.toISOString(),
                description: formValue.description,
                movements: formValue.movements,
                partnerName: formValue.partnerName || undefined,
                partnerAccount: formValue.partnerAccount || undefined,
                place: formValue.place || undefined,
                tags: this.tags().length > 0 ? this.tags() : undefined,
            };
            this.dialogRef.close(transaction);
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }
}
