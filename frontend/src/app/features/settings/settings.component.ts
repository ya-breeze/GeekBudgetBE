import { Component, inject, OnInit, signal } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { UserService } from '../../core/services/user.service';
import { CurrencyService } from '../currencies/services/currency.service';
import { LayoutService } from '../../layout/services/layout.service';

@Component({
    selector: 'app-settings',
    imports: [
        ReactiveFormsModule,
        MatCardModule,
        MatFormFieldModule,
        MatSelectModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
    ],
    templateUrl: './settings.component.html',
    styleUrl: './settings.component.scss',
})
export class SettingsComponent implements OnInit {
    private readonly fb = inject(FormBuilder);
    protected readonly userService = inject(UserService);
    private readonly currencyService = inject(CurrencyService);
    private readonly snackBar = inject(MatSnackBar);
    private readonly layoutService = inject(LayoutService);

    protected readonly sidenavOpened = this.layoutService.sidenavOpened;

    protected readonly settingsForm: FormGroup;
    protected readonly user = this.userService.user;
    protected readonly currencies = this.currencyService.currencies;
    protected readonly loading = signal(false);

    constructor() {
        this.settingsForm = this.fb.group({
            favoriteCurrencyId: [''],
        });
    }

    ngOnInit(): void {
        this.loading.set(true);

        // Load currencies and user data
        this.currencyService.loadCurrencies().subscribe();
        this.userService.loadUser().subscribe({
            next: (user) => {
                this.settingsForm.patchValue({
                    favoriteCurrencyId: user.favoriteCurrencyId || '',
                });
                this.loading.set(false);
            },
            error: () => {
                this.loading.set(false);
                this.snackBar.open('Failed to load user settings', 'Close', { duration: 3000 });
            },
        });
    }

    onSave(): void {
        if (this.settingsForm.valid) {
            const favoriteCurrencyId = this.settingsForm.value.favoriteCurrencyId || null;

            this.userService.updateFavoriteCurrency(favoriteCurrencyId).subscribe({
                next: () => {
                    this.snackBar.open('Settings saved successfully', 'Close', { duration: 3000 });
                },
                error: () => {
                    this.snackBar.open('Failed to save settings', 'Close', { duration: 3000 });
                },
            });
        }
    }
}
