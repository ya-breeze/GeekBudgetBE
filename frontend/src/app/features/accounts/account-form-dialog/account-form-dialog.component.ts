import { Component, inject, ViewChild, ElementRef } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatIconModule } from '@angular/material/icon';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { Account } from '../../../core/api/models/account';
import { AccountNoId } from '../../../core/api/models/account-no-id';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';
import { OnInit } from '@angular/core';

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
        MatDatepickerModule,
        MatNativeDateModule,
    ],
    templateUrl: './account-form-dialog.component.html',
    styleUrl: './account-form-dialog.component.scss',
})
export class AccountFormDialogComponent implements OnInit {
    private readonly dialogRef = inject(MatDialogRef<AccountFormDialogComponent>);
    private readonly data = inject<AccountFormDialogData>(MAT_DIALOG_DATA);
    private readonly fb = inject(FormBuilder);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly currencyService = inject(CurrencyService);
    private readonly userService = inject(UserService);

    @ViewChild('fileInput') fileInput!: ElementRef<HTMLInputElement>;

    protected readonly form: FormGroup;
    protected readonly isEditMode = this.data.mode === 'edit';
    protected readonly currencies = this.currencyService.currencies;
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
            bankId: [this.data.account?.bankInfo?.bankId || ''],
            bankAccountId: [this.data.account?.bankInfo?.accountId || ''],
            ignoreUnprocessedBefore: [this.parseDate(this.data.account?.ignoreUnprocessedBefore)],
            openingDate: [this.parseDate(this.data.account?.openingDate)],
            closingDate: [this.parseDate(this.data.account?.closingDate)],
            balances: this.fb.array([]),
        });

        if (this.data.account?.image) {
            const root = this.apiConfig.rootUrl.endsWith('/')
                ? this.apiConfig.rootUrl.slice(0, -1)
                : this.apiConfig.rootUrl;
            this.imagePreview = `${root}/images/${this.data.account.image}`;
        }
    }

    ngOnInit(): void {
        this.currencyService.loadCurrencies().subscribe();
        this.userService.loadUser().subscribe((user) => {
            // Initialize balances
            if (
                this.data.account?.bankInfo?.balances &&
                this.data.account.bankInfo.balances.length > 0
            ) {
                // Editing existing account with balances
                this.data.account.bankInfo.balances.forEach((balance) => {
                    this.addBalance(balance);
                });
            } else if (
                this.data.account?.type === 'asset' ||
                this.form.get('type')?.value === 'asset'
            ) {
                // New asset account or creating asset - add one default balance
                this.addBalance({
                    currencyId: user.favoriteCurrencyId || '',
                    openingBalance: 0,
                });
            }
        });
    }

    get balances(): FormArray {
        return this.form.get('balances') as FormArray;
    }

    addBalance(initialData?: {
        currencyId?: string;
        openingBalance?: number;
        closingBalance?: number;
    }): void {
        const balanceGroup = this.fb.group({
            currencyId: [initialData?.currencyId || '', [Validators.required]],
            openingBalance: [initialData?.openingBalance || 0],
            closingBalance: [initialData?.closingBalance],
        });
        this.balances.push(balanceGroup);
    }

    removeBalance(index: number): void {
        this.balances.removeAt(index);
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

        if (this.fileInput && this.fileInput.nativeElement) {
            this.fileInput.nativeElement.value = '';
        }
    }

    onSubmit(): void {
        if (this.form.valid) {
            const formValue = this.form.value;
            const account: AccountNoId = {
                name: formValue.name,
                type: formValue.type,
                description: formValue.description,
                showInDashboardSummary: formValue.showInDashboardSummary,
                hideFromReports: formValue.hideFromReports,
                bankInfo:
                    formValue.type === 'asset'
                        ? {
                              ...this.data.account?.bankInfo,
                              bankId: formValue.bankId || undefined,
                              accountId: formValue.bankAccountId || undefined,
                              balances: formValue.balances?.map((b: any) => ({
                                  currencyId: b.currencyId || undefined,
                                  openingBalance: b.openingBalance || 0,
                                  closingBalance: b.closingBalance,
                              })),
                          }
                        : undefined,
                ignoreUnprocessedBefore:
                    formValue.type === 'asset' && formValue.ignoreUnprocessedBefore
                        ? formValue.ignoreUnprocessedBefore.toISOString()
                        : undefined,
                openingDate:
                    formValue.type === 'asset' && formValue.openingDate
                        ? formValue.openingDate.toISOString()
                        : undefined,
                closingDate:
                    formValue.type === 'asset' && formValue.closingDate
                        ? formValue.closingDate.toISOString()
                        : undefined,
            };

            this.dialogRef.close({
                account: account,
                image: this.selectedFile,
                deleteImage: this.deleteImage,
            });
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    private parseDate(dateStr?: string | null): Date | null {
        if (!dateStr || dateStr.startsWith('0001-01-01')) {
            return null;
        }
        return new Date(dateStr);
    }
}
