import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { BankImporter } from '../../../core/api/models/bank-importer';
import { BankImporterNoId } from '../../../core/api/models/bank-importer-no-id';
import { AccountService } from '../../accounts/services/account.service';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';

export interface BankImporterFormDialogData {
  mode: 'create' | 'edit';
  bankImporter?: BankImporter;
}

@Component({
  selector: 'app-bank-importer-form-dialog',
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatCheckboxModule,
    AccountSelectComponent
  ],
  templateUrl: './bank-importer-form-dialog.component.html',
  styleUrl: './bank-importer-form-dialog.component.scss',
})
export class BankImporterFormDialogComponent implements OnInit {
  private readonly dialogRef = inject(MatDialogRef<BankImporterFormDialogComponent>);
  private readonly data = inject<BankImporterFormDialogData>(MAT_DIALOG_DATA);
  private readonly fb = inject(FormBuilder);
  private readonly accountService = inject(AccountService);

  protected readonly form: FormGroup;
  protected readonly isEditMode = this.data.mode === 'edit';
  protected readonly accounts = this.accountService.accounts;
  protected readonly bankTypes = [
    { value: 'fio', label: 'FIO Bank' },
    { value: 'kb', label: 'KB Bank' },
    { value: 'revolut', label: 'Revolut' },
  ];

  constructor() {
    this.form = this.fb.group({
      name: [this.data.bankImporter?.name || '', [Validators.required, Validators.maxLength(100)]],
      type: [this.data.bankImporter?.type || 'fio', [Validators.required]],
      accountId: [this.data.bankImporter?.accountId || '', [Validators.required]],
      feeAccountId: [this.data.bankImporter?.feeAccountId || ''],
      description: [this.data.bankImporter?.description || '', [Validators.maxLength(500)]],
      extra: [this.data.bankImporter?.extra || ''],
      fetchAll: [this.data.bankImporter?.fetchAll || false],
    });
  }

  ngOnInit(): void {
    this.accountService.loadAccounts().subscribe();
  }

  onSubmit(): void {
    if (this.form.valid) {
      const formValue = this.form.value;
      const bankImporter: BankImporterNoId = {
        name: formValue.name,
        type: formValue.type,
        accountId: formValue.accountId,
        feeAccountId: formValue.feeAccountId || undefined,
        description: formValue.description || undefined,
        extra: formValue.extra || undefined,
        fetchAll: formValue.fetchAll,
      };
      this.dialogRef.close(bankImporter);
    }
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}

