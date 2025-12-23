import { Component, inject, OnInit, signal, HostListener } from '@angular/core';
import { FormBuilder, Validators, ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatCheckboxModule } from '@angular/material/checkbox';
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
import { Transaction } from '../../../core/api/models/transaction';
import { AccountSelectComponent } from '../../../shared/components/account-select/account-select.component';
import { ApiConfiguration } from '../../../core/api/api-configuration';

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
        MatProgressSpinnerModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        AccountSelectComponent,
    ],
    templateUrl: './matcher-edit-dialog.component.html',
    styleUrls: ['./matcher-edit-dialog.component.scss'],
})
export class MatcherEditDialogComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    private readonly dialogRef = inject(MatDialogRef<MatcherEditDialogComponent>);
    private readonly matcherService = inject(MatcherService);
    private readonly accountService = inject(AccountService);
    private readonly apiConfig = inject(ApiConfiguration);
    readonly data = inject<{ matcher?: Matcher; transaction?: Transaction } | undefined>(
        MAT_DIALOG_DATA,
    );

    protected readonly accounts = this.accountService.accounts;
    protected readonly loading = signal(false);

    // Regex/Matcher Testing Signals
    protected readonly testString = signal('');
    protected readonly matchResult = signal<{ result?: boolean; reason?: string } | null>(null);
    protected readonly regexResult = signal<{
        isValid: boolean;
        isMatch: boolean;
        error?: string;
    } | null>(null);
    protected readonly testingMatch = signal(false);
    protected readonly testingRegex = signal(false);

    protected selectedFile: File | null = null;
    protected imagePreview: string | null = null;
    protected deleteImage = false;

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

    checkMatch(): void {
        if (!this.data?.transaction) return;

        this.testingMatch.set(true);
        const formValue = this.form.value;

        // Construct matcher object from form (excluding output fields which don't affect matching)
        // We need to be careful to constructing it correctly for the backend check
        const matcherToCheck: MatcherNoId = {
            descriptionRegExp: formValue.descriptionRegExp || undefined,
            partnerNameRegExp: formValue.partnerNameRegExp || undefined,
            partnerAccountNumberRegExp: formValue.partnerAccountNumberRegExp || undefined,
            currencyRegExp: formValue.currencyRegExp || undefined,
            extraRegExp: formValue.extraRegExp || undefined,
            placeRegExp: formValue.placeRegExp || undefined,
            outputAccountId: 'temp', // Not relevant for matching check
            outputDescription: 'temp',
            outputTags: [],
        };

        const t = this.data.transaction;
        const transactionToCheck: any = {
            date: t.date,
            description: t.description,
            externalIds: t.externalIds,
            extra: t.extra,
            movements: t.movements,
            partnerAccount: t.partnerAccount,
            partnerInternalId: t.partnerInternalId,
            partnerName: t.partnerName,
            place: t.place,
            tags: t.tags,
            unprocessedSources: t.unprocessedSources,
        };

        this.matcherService.checkMatcher(matcherToCheck, transactionToCheck).subscribe({
            next: (result) => {
                this.matchResult.set(result);
                this.testingMatch.set(false);
            },
            error: () => {
                this.matchResult.set({ result: false, reason: 'Network error checking match' });
                this.testingMatch.set(false);
            },
        });
    }

    testRegex(): void {
        const pattern = this.form.controls.descriptionRegExp.value;
        const test = this.testString();

        if (!pattern || !test) {
            this.regexResult.set(null);
            return;
        }

        this.testingMatch.set(true); // Re-use same loading signal
        this.matcherService.checkRegex(pattern, test).subscribe({
            next: (result) => {
                this.regexResult.set(result);
                this.testingMatch.set(false);
            },
            error: () => {
                this.regexResult.set({ isValid: false, isMatch: false, error: 'Network error' });
                this.testingMatch.set(false);
            },
        });
    }

    protected readonly form = this.fb.group({
        descriptionRegExp: [''],
        partnerNameRegExp: [''],
        partnerAccountNumberRegExp: [''],
        currencyRegExp: [''],
        extraRegExp: [''],
        placeRegExp: [''],
        outputAccountId: ['', Validators.required],
        outputDescription: ['', Validators.required],
        outputTags: [''], // Comma separated for simplicity initially
    });

    ngOnInit(): void {
        this.accountService.loadAccounts().subscribe();

        if (this.data?.matcher) {
            this.form.patchValue({
                descriptionRegExp: this.data.matcher.descriptionRegExp,
                partnerNameRegExp: this.data.matcher.partnerNameRegExp,
                partnerAccountNumberRegExp: this.data.matcher.partnerAccountNumberRegExp,
                currencyRegExp: this.data.matcher.currencyRegExp,
                extraRegExp: this.data.matcher.extraRegExp,
                placeRegExp: (this.data.matcher as any).placeRegExp, // Cast to any to bypass strict type checking until models are updated in Frontend
                outputAccountId: this.data.matcher.outputAccountId,
                outputDescription: this.data.matcher.outputDescription,
                outputTags: this.data.matcher.outputTags?.join(', '),
            });

            if (this.data.matcher.image) {
                const root = this.apiConfig.rootUrl.endsWith('/')
                    ? this.apiConfig.rootUrl.slice(0, -1)
                    : this.apiConfig.rootUrl;
                this.imagePreview = `${root}/images/${this.data.matcher.image}`;
            }
        } else if (this.data?.transaction) {
            // Pre-fill from transaction
            const t = this.data.transaction;
            this.form.patchValue({
                // Suggest name from description/partner

                // Pre-fill regexps with exact matches or simplistic logic (user will edit)
                descriptionRegExp: this.escapeRegExp(t.description || ''),
                partnerNameRegExp: t.partnerName ? this.escapeRegExp(t.partnerName) : '',
                partnerAccountNumberRegExp: t.partnerAccount
                    ? this.escapeRegExp(t.partnerAccount)
                    : '',
                placeRegExp: t.place ? this.escapeRegExp(t.place) : '',
                outputDescription: t.description, // Suggest keeping description or edit
            });
        }
    }

    get isCaseInsensitive(): boolean {
        const val = this.form.controls.descriptionRegExp.value || '';
        return val.startsWith('(?i)');
    }

    toggleCaseInsensitive(checked: boolean): void {
        const val = this.form.controls.descriptionRegExp.value || '';
        if (checked) {
            if (!val.startsWith('(?i)')) {
                this.form.controls.descriptionRegExp.setValue('(?i)' + val);
            }
        } else {
            if (val.startsWith('(?i)')) {
                this.form.controls.descriptionRegExp.setValue(val.substring(4));
            }
        }
    }

    get isWholeWord(): boolean {
        let val = this.form.controls.descriptionRegExp.value || '';
        if (val.startsWith('(?i)')) {
            val = val.substring(4);
        }
        return val.startsWith('\\b') && val.endsWith('\\b');
    }

    toggleWholeWord(checked: boolean): void {
        let val = this.form.controls.descriptionRegExp.value || '';
        let prefix = '';
        if (val.startsWith('(?i)')) {
            prefix = '(?i)';
            val = val.substring(4);
        }

        if (checked) {
            if (!val.startsWith('\\b')) val = '\\b' + val;
            if (!val.endsWith('\\b')) val = val + '\\b';
        } else {
            if (val.startsWith('\\b')) val = val.substring(2);
            if (val.endsWith('\\b')) val = val.substring(0, val.length - 2);
        }

        this.form.controls.descriptionRegExp.setValue(prefix + val);
    }

    private escapeRegExp(string: string): string {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }

    // Place RegExp Helpers
    get isPlaceCaseInsensitive(): boolean {
        const val = this.form.controls.placeRegExp.value || '';
        return val.startsWith('(?i)');
    }

    togglePlaceCaseInsensitive(checked: boolean): void {
        const val = this.form.controls.placeRegExp.value || '';
        if (checked) {
            if (!val.startsWith('(?i)')) {
                this.form.controls.placeRegExp.setValue('(?i)' + val);
            }
        } else {
            if (val.startsWith('(?i)')) {
                this.form.controls.placeRegExp.setValue(val.substring(4));
            }
        }
    }

    get isPlaceWholeWord(): boolean {
        let val = this.form.controls.placeRegExp.value || '';
        if (val.startsWith('(?i)')) {
            val = val.substring(4);
        }
        return val.startsWith('\\b') && val.endsWith('\\b');
    }

    togglePlaceWholeWord(checked: boolean): void {
        let val = this.form.controls.placeRegExp.value || '';
        let prefix = '';
        if (val.startsWith('(?i)')) {
            prefix = '(?i)';
            val = val.substring(4);
        }

        if (checked) {
            if (!val.startsWith('\\b')) val = '\\b' + val;
            if (!val.endsWith('\\b')) val = val + '\\b';
        } else {
            if (val.startsWith('\\b')) val = val.substring(2);
            if (val.endsWith('\\b')) val = val.substring(0, val.length - 2);
        }

        this.form.controls.placeRegExp.setValue(prefix + val);
    }

    getAccount(id: string): any {
        return this.accounts().find((a) => a.id === id);
    }

    save(): void {
        if (this.form.invalid) return;

        this.loading.set(true);
        const formValue = this.form.value;

        const matcherData: MatcherNoId = {
            descriptionRegExp: formValue.descriptionRegExp || undefined,
            partnerNameRegExp: formValue.partnerNameRegExp || undefined,
            partnerAccountNumberRegExp: formValue.partnerAccountNumberRegExp || undefined,
            currencyRegExp: formValue.currencyRegExp || undefined,
            extraRegExp: formValue.extraRegExp || undefined,
            placeRegExp: formValue.placeRegExp || undefined,
            outputAccountId: formValue.outputAccountId!,
            outputDescription: formValue.outputDescription!,
            outputTags: formValue.outputTags
                ? formValue.outputTags
                      .split(',')
                      .map((t) => t.trim())
                      .filter((t) => !!t)
                : [],
        };

        const request = this.data?.matcher
            ? this.matcherService.update(this.data.matcher.id, matcherData)
            : this.matcherService.create(matcherData);

        request.subscribe({
            next: (savedMatcher) => {
                // Handle Image Upload/Delete
                if (savedMatcher.id) {
                    if (this.deleteImage && !this.selectedFile) {
                        this.matcherService.deleteImage(savedMatcher.id).subscribe({
                            next: () => {
                                this.loading.set(false);
                                this.dialogRef.close(true);
                            },
                            error: () => {
                                this.loading.set(false);
                                // Show warning or just close?
                                this.dialogRef.close(true);
                            },
                        });
                        return;
                    } else if (this.selectedFile) {
                        this.matcherService
                            .uploadImage(savedMatcher.id, this.selectedFile)
                            .subscribe({
                                next: () => {
                                    this.loading.set(false);
                                    this.dialogRef.close(true);
                                },
                                error: () => {
                                    this.loading.set(false);
                                    this.dialogRef.close(true);
                                },
                            });
                        return;
                    }
                }

                this.loading.set(false);
                this.dialogRef.close(true);
            },
            error: () => {
                this.loading.set(false);
                // Error handling usually done in service via snackbar or global handler?
                // If not, we should show message here. Service has error signal.
            },
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
