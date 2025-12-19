import {
    Component,
    computed,
    forwardRef,
    inject,
    Input,
    OnChanges,
    SimpleChanges,
    ViewChild,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
    ControlValueAccessor,
    FormControl,
    FormsModule,
    NG_VALUE_ACCESSOR,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { MatAutocompleteModule, MatAutocomplete } from '@angular/material/autocomplete';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatOptionModule } from '@angular/material/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { startWith, map } from 'rxjs';
import { Account } from '../../../core/api/models/account';
import { AccountService } from '../../../features/accounts/services/account.service';

export interface AccountGroup {
    name: string;
    accounts: Account[];
}

@Component({
    selector: 'app-account-select',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatAutocompleteModule,
        MatOptionModule,
    ],
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => AccountSelectComponent),
            multi: true,
        },
    ],
    templateUrl: './account-select.component.html',
    styleUrls: ['./account-select.component.scss'],
})
export class AccountSelectComponent implements ControlValueAccessor, OnChanges {
    private readonly accountService: AccountService = inject(AccountService);
    protected readonly accounts = this.accountService.accounts;

    @Input() label = 'Account';
    @Input() required = false;
    @Input() placeholder = 'Select an account';

    @ViewChild('auto') auto!: MatAutocomplete;
    readonly searchControl = new FormControl<string | Account>('');

    // Internal value tracking (Account ID)
    private _value: string | null = null;
    private onChange: (value: string | null) => void = () => {};
    private onTouched: () => void = () => {};

    protected readonly filterValue = toSignal(
        this.searchControl.valueChanges.pipe(
            startWith(''),
            map((value) => {
                if (typeof value === 'string') return value;
                return value?.name || '';
            }),
        ),
        { initialValue: '' },
    );

    protected readonly filteredAccountGroups = computed(() => {
        const accounts = this.accounts();
        const filterValue = (this.filterValue() || '').toLowerCase();

        // 1. Filter
        const filtered = accounts.filter((account: Account) =>
            account.name.toLowerCase().includes(filterValue),
        );

        // 2. Sort
        filtered.sort((a: Account, b: Account) => a.name.localeCompare(b.name));

        // 3. Group
        const groups: AccountGroup[] = [];
        const typeMap = new Map<string, Account[]>();

        filtered.forEach((account: Account) => {
            const type = account.type || 'Other';
            if (!typeMap.has(type)) {
                typeMap.set(type, []);
            }
            typeMap.get(type)!.push(account);
        });

        // Order of groups
        const order = ['expense', 'income', 'asset', 'liability', 'equity'];
        const pluralMap: Record<string, string> = {
            expense: 'Expenses',
            income: 'Incomes',
            asset: 'Assets',
            liability: 'Liabilities',
            equity: 'Equities',
        };

        const getGroupName = (type: string) => {
            return pluralMap[type] || type.charAt(0).toUpperCase() + type.slice(1) + 's';
        };

        // Add known types in order
        order.forEach((type) => {
            if (typeMap.has(type)) {
                groups.push({ name: getGroupName(type), accounts: typeMap.get(type)! });
                typeMap.delete(type);
            }
        });

        // Add remaining types
        typeMap.forEach((accounts, type) => {
            groups.push({ name: getGroupName(type), accounts: accounts });
        });

        return groups;
    });

    constructor() {
        // Ensure accounts are loaded (redundant if app initializes it, but safe)
        if (this.accounts().length === 0) {
            this.accountService.loadAccounts().subscribe();
        }

        // Listen for selection changes from the autocomplete
        this.searchControl.valueChanges.subscribe((value) => {
            if (typeof value === 'object' && value !== null) {
                this.emitValue(value.id);
            } else if (value === '') {
                // Clear selection if input is cleared
                this.emitValue(null);
            }
            // Note: If user types text that doesn't match, we don't necessarily clear immediately
            // until onBlur logic or valid option selected.
            // But for now, let's keep it simple: only Emit if it's a valid object.
        });
    }

    ngOnChanges(changes: SimpleChanges): void {
        if (changes['required']) {
            if (this.required) {
                this.searchControl.addValidators(Validators.required);
            } else {
                this.searchControl.removeValidators(Validators.required);
            }
            this.searchControl.updateValueAndValidity();
        }
    }

    // ControlValueAccessor implementation
    writeValue(obj: any): void {
        this._value = obj;
        if (obj) {
            const account = this.accounts().find((a: Account) => a.id === obj);
            if (account) {
                this.searchControl.setValue(account, { emitEvent: false });
            } else {
                // Handle case where account might not be loaded yet, or invalid
                // If accounts signal updates later, we might want to re-check.
                // For now, assume loaded.
            }
        } else {
            this.searchControl.reset(null);
        }
    }

    registerOnChange(fn: any): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: any): void {
        this.onTouched = fn;
    }

    setDisabledState?(isDisabled: boolean): void {
        if (isDisabled) {
            this.searchControl.disable();
        } else {
            this.searchControl.enable();
        }
    }

    displayFn(account: Account): string {
        return account && account.name ? account.name : '';
    }

    onBlur() {
        // If the panel is open, we don't want to clear the value yet.
        // We will validate when the panel closes.
        if (this.auto && this.auto.isOpen) {
            return;
        }
        this.validateSelection();
    }

    validateSelection() {
        this.onTouched();
        const val = this.searchControl.value;
        if (typeof val === 'string' && val !== '') {
            // Invalid selection (text only)
            this.searchControl.setValue(null);
            this.emitValue(null);
        }
    }

    private emitValue(value: string | null) {
        if (this._value !== value) {
            this._value = value;
            this.onChange(value);
        }
    }
}
