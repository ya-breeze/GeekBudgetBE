import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { LoginComponent } from './login.component';
import { AuthService } from '../../../core/auth/services/auth.service';
import { of, throwError } from 'rxjs';

describe('LoginComponent', () => {
    let component: LoginComponent;
    let fixture: ComponentFixture<LoginComponent>;
    let authService: jasmine.SpyObj<AuthService>;
    let router: jasmine.SpyObj<Router>;

    beforeEach(async () => {
        const authServiceSpy = jasmine.createSpyObj('AuthService', ['login']);
        const routerSpy = jasmine.createSpyObj('Router', ['navigate']);

        await TestBed.configureTestingModule({
            imports: [LoginComponent, ReactiveFormsModule],
            providers: [
                { provide: AuthService, useValue: authServiceSpy },
                { provide: Router, useValue: routerSpy },
            ],
        }).compileComponents();

        authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
        router = TestBed.inject(Router) as jasmine.SpyObj<Router>;

        fixture = TestBed.createComponent(LoginComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should display login form', () => {
        const compiled = fixture.nativeElement;
        const emailInput = compiled.querySelector('input[type="email"]');
        const passwordInput = compiled.querySelector('input[type="password"]');
        const submitButton = compiled.querySelector('button[type="submit"]');

        expect(emailInput).toBeTruthy();
        expect(passwordInput).toBeTruthy();
        expect(submitButton).toBeTruthy();
    });

    it('should validate email format', () => {
        const emailControl = component.loginForm.get('email');

        emailControl?.setValue('invalid-email');
        expect(emailControl?.hasError('email')).toBe(true);

        emailControl?.setValue('valid@email.com');
        expect(emailControl?.hasError('email')).toBe(false);
    });

    it('should validate required fields', () => {
        const emailControl = component.loginForm.get('email');
        const passwordControl = component.loginForm.get('password');

        emailControl?.setValue('');
        passwordControl?.setValue('');

        expect(emailControl?.hasError('required')).toBe(true);
        expect(passwordControl?.hasError('required')).toBe(true);
    });

    it('should call auth service on submit', () => {
        authService.login.and.returnValue(of(void 0));

        component.loginForm.setValue({
            email: 'test@example.com',
            password: 'password123',
        });

        component.onSubmit();

        expect(authService.login).toHaveBeenCalledWith('test@example.com', 'password123');
    });

    it('should navigate to dashboard on successful login', () => {
        authService.login.and.returnValue(of(void 0));

        component.loginForm.setValue({
            email: 'test@example.com',
            password: 'password123',
        });

        component.onSubmit();

        expect(router.navigate).toHaveBeenCalledWith(['/dashboard']);
    });

    it('should display error message on failed login', () => {
        const mockError = { status: 401, message: 'Invalid credentials' };
        authService.login.and.returnValue(throwError(() => mockError));

        component.loginForm.setValue({
            email: 'test@example.com',
            password: 'wrong-password',
        });

        component.onSubmit();

        expect(component.errorMessage).toBeTruthy();
    });

    it('should disable submit button when form is invalid', () => {
        component.loginForm.setValue({
            email: '',
            password: '',
        });

        fixture.detectChanges();

        const compiled = fixture.nativeElement;
        const submitButton = compiled.querySelector('button[type="submit"]');

        expect(submitButton.disabled).toBe(true);
    });

    it('should show loading state during login', (done) => {
        let loadingDuringCall = false;

        authService.login.and.callFake(() => {
            loadingDuringCall = component.isLoading;
            return of(void 0);
        });

        component.loginForm.setValue({
            email: 'test@example.com',
            password: 'password123',
        });

        component.onSubmit();

        // Loading should be true during the call
        expect(loadingDuringCall).toBe(true);
        // Loading should be false after completion
        expect(component.isLoading).toBe(false);
        done();
    });

    it('should enable submit button when form is valid', () => {
        component.loginForm.setValue({
            email: 'test@example.com',
            password: 'password123',
        });

        fixture.detectChanges();

        const compiled = fixture.nativeElement;
        const submitButton = compiled.querySelector('button[type="submit"]');

        expect(submitButton.disabled).toBe(false);
    });
});
