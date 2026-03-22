import { Component, inject, OnInit } from '@angular/core';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TemplateService } from '../services/template.service';
import { AccountService } from '../../accounts/services/account.service';
import { CurrencyService } from '../../currencies/services/currency.service';
import { TransactionTemplate } from '../../../core/api/models/transaction-template';
import { TransactionTemplateNoId } from '../../../core/api/models/transaction-template-no-id';
import { Movement } from '../../../core/api/models/movement';

export interface TemplateInitialValues {
    name?: string;
    description?: string;
    place?: string;
    partnerName?: string;
    extra?: string;
    movements?: Movement[];
}

export interface TemplateEditDialogData {
    /** Provide when editing an existing template. Sets isEditMode = true. */
    template?: TransactionTemplate;
    /** Provide when creating from a transaction. isEditMode stays false. */
    initialValues?: TemplateInitialValues;
}

@Component({
    selector: 'app-template-edit-dialog',
    standalone: true,
    imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatIconModule,
        MatSelectModule,
    ],
    templateUrl: './template-edit-dialog.component.html',
})
export class TemplateEditDialogComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly dialogRef = inject(MatDialogRef<TemplateEditDialogComponent>);
    protected readonly data = inject<TemplateEditDialogData>(MAT_DIALOG_DATA);
    private readonly templateService = inject(TemplateService);
    protected readonly accountService = inject(AccountService);
    protected readonly currencyService = inject(CurrencyService);

    protected readonly isEditMode = !!this.data?.template;

    protected readonly form = this.fb.group({
        name: ['', Validators.required],
        description: [''],
        place: [''],
        partnerName: [''],
        extra: [''],
        movements: this.fb.array([]),
    });

    get movements(): FormArray {
        return this.form.get('movements') as FormArray;
    }

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();
        this.currencyService.loadCurrencies().subscribe();

        if (this.data?.template) {
            const t = this.data.template;
            this.form.patchValue({
                name: t.name,
                description: t.description ?? '',
                place: t.place ?? '',
                partnerName: t.partnerName ?? '',
                extra: t.extra ?? '',
            });
            (t.movements ?? []).forEach((m) => this.addMovement(m));
        } else if (this.data?.initialValues) {
            const v = this.data.initialValues;
            this.form.patchValue({
                name: v.name ?? '',
                description: v.description ?? '',
                place: v.place ?? '',
                partnerName: v.partnerName ?? '',
                extra: v.extra ?? '',
            });
            (v.movements ?? []).forEach((m) => this.addMovement(m));
            if (!v.movements?.length) this.addMovement();
        } else {
            this.addMovement();
        }
    }

    protected addMovement(movement?: Partial<Movement>): void {
        this.movements.push(
            this.fb.group({
                amount: [
                    movement?.amount != null ? String(movement.amount) : '',
                    Validators.required,
                ],
                currencyId: [movement?.currencyId ?? '', Validators.required],
                accountId: [movement?.accountId ?? ''],
            }),
        );
    }

    protected removeMovement(index: number): void {
        if (this.movements.length > 1) {
            this.movements.removeAt(index);
        }
    }

    protected save(): void {
        if (this.form.invalid) return;

        const value = this.form.getRawValue();
        const payload: TransactionTemplateNoId = {
            name: value.name!,
            description: value.description || undefined,
            place: value.place || undefined,
            partnerName: value.partnerName || undefined,
            extra: value.extra || undefined,
            movements: (value.movements as any[]).map((m) => ({
                amount: Number(m.amount),
                currencyId: m.currencyId,
                accountId: m.accountId || undefined,
            })),
        };

        const obs = this.isEditMode
            ? this.templateService.update(this.data.template!.id, payload)
            : this.templateService.create(payload);

        obs.subscribe(() => this.dialogRef.close(true));
    }

    protected cancel(): void {
        this.dialogRef.close(false);
    }
}
