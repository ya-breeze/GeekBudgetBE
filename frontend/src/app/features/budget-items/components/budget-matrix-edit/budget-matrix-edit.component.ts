import { Component, Inject, inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { CurrencyPipe, DatePipe } from '@angular/common';
import { Router } from '@angular/router';

export interface BudgetMatrixEditData {
  accountId: string;
  accountName: string;
  month: string; // ISO string from caller
  currentAmount: number;
}

@Component({
  selector: 'app-budget-matrix-edit',
  standalone: true,
  imports: [
    MatDialogModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    CurrencyPipe,
    DatePipe
  ],
  template: `
    <h2 mat-dialog-title>Budget for {{ data.accountName }}</h2>
    <div mat-dialog-content>
      <p class="month-label">{{ data.month | date:'MMMM yyyy' }}</p>
      
      <form [formGroup]="form" (ngSubmit)="save()">
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Planned Amount</mat-label>
          <input matInput type="number" formControlName="amount" min="0" cdkFocusInitial>
        </mat-form-field>
      </form>
    </div>
    <div mat-dialog-actions class="actions-container">
      <button mat-button color="accent" (click)="viewTransactions()">
        <mat-icon>list_alt</mat-icon> View Transactions
      </button>
      <span class="spacer"></span>
      <button mat-button (click)="cancel()">Cancel</button>
      <button mat-flat-button color="primary" (click)="save()" [disabled]="form.invalid">Save</button>
    </div>
  `,
  styles: [`
    .month-label {
      margin-top: 0;
      margin-bottom: 16px;
      color: rgba(0, 0, 0, 0.6);
    }
    .full-width {
      width: 100%;
    }
    .actions-container {
      display: flex;
      width: 100%;
      padding: 0 24px 24px 24px; // Adjust padding to match dialog convention if needed
    }
    .spacer {
      flex: 1;
    }
  `]
})
export class BudgetMatrixEditComponent {
  private fb = inject(FormBuilder);
  private dialogRef = inject(MatDialogRef<BudgetMatrixEditComponent>);
  private router = inject(Router);
  public data: BudgetMatrixEditData = inject(MAT_DIALOG_DATA);

  protected form = this.fb.group({
    amount: [this.data.currentAmount, [Validators.required, Validators.min(0)]]
  });

  constructor() { }

  save() {
    if (this.form.valid) {
      this.dialogRef.close(this.form.value.amount);
    }
  }

  cancel() {
    this.dialogRef.close();
  }

  viewTransactions() {
    const date = new Date(this.data.month);
    this.router.navigate(['/transactions'], {
      queryParams: {
        accountId: this.data.accountId,
        month: date.getMonth(),
        year: date.getFullYear()
      }
    });
    this.dialogRef.close();
  }
}
