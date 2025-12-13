import { Component, inject, OnInit, signal, computed, HostListener } from '@angular/core';
import { FormBuilder, Validators, ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { CommonModule } from '@angular/common';
import { MatcherService } from '../services/matcher.service';
import { AccountService } from '../../accounts/services/account.service';
import { Matcher } from '../../../core/api/models/matcher';
import { MatcherNoId } from '../../../core/api/models/matcher-no-id';

@Component({
    selector: 'app-matcher-edit-dialog',
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        FormsModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatIconModule,
        MatChipsModule,
        MatProgressSpinnerModule
    ],
    templateUrl: './matcher-edit-dialog.component.html',
    styleUrls: ['./matcher-edit-dialog.component.scss']
})
export class MatcherEditDialogComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly dialogRef = inject(MatDialogRef<MatcherEditDialogComponent>);
    private readonly matcherService = inject(MatcherService);
    private readonly accountService = inject(AccountService);
    readonly data = inject<Matcher | undefined>(MAT_DIALOG_DATA);

    protected readonly accounts = this.accountService.accounts;
    protected readonly loading = signal(false);

    // Regex Testing Signals
    protected readonly testString = signal('');
    protected readonly regexResult = signal<{ isValid: boolean, isMatch: boolean, error?: string } | null>(null);
    protected readonly testingRegex = signal(false);

    testRegex(): void {
        const pattern = this.form.controls.descriptionRegExp.value;
        const test = this.testString();

        if (!pattern || !test) {
            this.regexResult.set(null);
            return;
        }

        this.testingRegex.set(true);
        this.matcherService.checkRegex(pattern, test).subscribe({
            next: (result) => {
                this.regexResult.set(result);
                this.testingRegex.set(false);
            },
            error: () => {
                this.regexResult.set({ isValid: false, isMatch: false, error: 'Network error' });
                this.testingRegex.set(false);
            }
        });
    }

    protected readonly form = this.fb.group({
        name: ['', Validators.required],
        descriptionRegExp: [''],
        partnerNameRegExp: [''],
        partnerAccountNumberRegExp: [''],
        currencyRegExp: [''],
        extraRegExp: [''],
        outputAccountId: ['', Validators.required],
        outputDescription: ['', Validators.required],
        outputTags: [''] // Comma separated for simplicity initially
    });

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();

        if (this.data) {
            this.form.patchValue({
                name: this.data.name,
                descriptionRegExp: this.data.descriptionRegExp,
                partnerNameRegExp: this.data.partnerNameRegExp,
                partnerAccountNumberRegExp: this.data.partnerAccountNumberRegExp,
                currencyRegExp: this.data.currencyRegExp,
                extraRegExp: this.data.extraRegExp,
                outputAccountId: this.data.outputAccountId,
                outputDescription: this.data.outputDescription,
                outputTags: this.data.outputTags?.join(', ')
            });
        }
    }

    save(): void {
        if (this.form.invalid) return;

        this.loading.set(true);
        const formValue = this.form.value;

        const matcherData: MatcherNoId = {
            name: formValue.name!,
            descriptionRegExp: formValue.descriptionRegExp || undefined,
            partnerNameRegExp: formValue.partnerNameRegExp || undefined,
            partnerAccountNumberRegExp: formValue.partnerAccountNumberRegExp || undefined,
            currencyRegExp: formValue.currencyRegExp || undefined,
            extraRegExp: formValue.extraRegExp || undefined,
            outputAccountId: formValue.outputAccountId!,
            outputDescription: formValue.outputDescription!,
            outputTags: formValue.outputTags ? formValue.outputTags.split(',').map(t => t.trim()).filter(t => !!t) : []
        };

        const request = this.data
            ? this.matcherService.update(this.data.id, matcherData)
            : this.matcherService.create(matcherData);

        request.subscribe({
            next: () => {
                this.loading.set(false);
                this.dialogRef.close(true);
            },
            error: () => {
                this.loading.set(false);
                // Error handling usually done in service via snackbar or global handler?
                // If not, we should show message here. Service has error signal.
            }
        });
    }

    close(): void {
        this.dialogRef.close();
    }
    @HostListener('window:keydown.esc')
    onEsc(): void {
        if (!this.loading()) {
            this.close();
        }
    }
}
