import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AccountFormDialogComponent, AccountFormDialogData } from './account-form-dialog.component';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { of } from 'rxjs';
import { CurrencyService } from '../../currencies/services/currency.service';
import { UserService } from '../../../core/services/user.service';

describe('AccountFormDialogComponent', () => {
    let component: AccountFormDialogComponent;
    let fixture: ComponentFixture<AccountFormDialogComponent>;
    let mockDialogRef: jasmine.SpyObj<MatDialogRef<AccountFormDialogComponent>>;
    let mockCurrencyService: jasmine.SpyObj<CurrencyService>;
    let mockUserService: jasmine.SpyObj<UserService>;

    const mockDialogData: AccountFormDialogData = {
        mode: 'create',
    };

    beforeEach(async () => {
        mockDialogRef = jasmine.createSpyObj('MatDialogRef', ['close']);
        mockCurrencyService = jasmine.createSpyObj('CurrencyService', ['loadCurrencies']);
        mockCurrencyService.loadCurrencies.and.returnValue(of([]));
        Object.defineProperty(mockCurrencyService, 'currencies', { get: () => () => [] });

        mockUserService = jasmine.createSpyObj('UserService', ['loadUser']);
        mockUserService.loadUser.and.returnValue(
            of({ id: '1', email: 'test@test.com', startDate: '2025-01-01' }),
        );

        await TestBed.configureTestingModule({
            imports: [
                AccountFormDialogComponent,
                HttpClientTestingModule,
                NoopAnimationsModule,
                ReactiveFormsModule,
                MatSelectModule,
                MatInputModule,
                MatFormFieldModule,
                MatDatepickerModule,
                MatNativeDateModule,
                MatSlideToggleModule,
            ],
            providers: [
                { provide: MatDialogRef, useValue: mockDialogRef },
                { provide: MAT_DIALOG_DATA, useValue: mockDialogData },
                { provide: CurrencyService, useValue: mockCurrencyService },
                { provide: UserService, useValue: mockUserService },
            ],
        }).compileComponents();

        fixture = TestBed.createComponent(AccountFormDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should allow uploading the same file again after removing it', () => {
        // Simulate file input element
        const inputElement = document.createElement('input');
        inputElement.type = 'file';
        // Access the private/protected property if needed for testing, or set the ElementRef manually if ViewChild logic requires it.
        // However, since we used ViewChild, Angular should have populated it if the template was rendered.
        // Let's check if component.fileInput is populated.

        // In our case, the file input is inside an @if block (if (imagePreview) is NOT true initially).
        // The input is always there in the template: <input #fileInput ... />

        expect(component.fileInput).toBeDefined();

        // Simulate setting a value
        component.fileInput.nativeElement = inputElement; // Mocking nativeElement for test if needed, or rely on actual DOM
        // Actually, let's use the real DOM element if possible.
        const nativeInput = component.fileInput.nativeElement;

        // Simulate user selecting a file (browser sets the value)
        // We can't easily set 'files' property due to security, but we can set 'value' logic that we are testing clearing of.
        // Note: setting value to a non-empty string on input type=file throws security error in browsers,
        // but in jsdom/testing environments it might be restricted too.
        // However, the component relies on resetting it to "".

        // Let's spy on the native element property setter or just assume we can set it for the mock.
        // Since we can't set value to a path, we can try to verify the logic via spy or by checking if it handles the reset.

        // A better approach for this test environment:
        // 1. Trigger method removeImage()
        // 2. Verify value is empty string.

        // Let's forcefully set a value if permitted, or mock the element.
        Object.defineProperty(nativeInput, 'value', {
            value: 'C:\\fakepath\\test.png',
            writable: true,
        });

        expect(nativeInput.value).toBe('C:\\fakepath\\test.png');

        component.removeImage();

        expect(nativeInput.value).toBe('');
    });
});
