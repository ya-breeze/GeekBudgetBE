import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { noAuthGuard } from './no-auth.guard';
import { AuthService } from '../services/auth.service';

describe('noAuthGuard', () => {
    let authService: jasmine.SpyObj<AuthService>;
    let router: jasmine.SpyObj<Router>;

    beforeEach(() => {
        const authServiceSpy = jasmine.createSpyObj('AuthService', ['isLoggedIn']);
        const routerSpy = jasmine.createSpyObj('Router', ['navigate', 'createUrlTree']);

        TestBed.configureTestingModule({
            providers: [
                { provide: AuthService, useValue: authServiceSpy },
                { provide: Router, useValue: routerSpy },
            ],
        });

        authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
        router = TestBed.inject(Router) as jasmine.SpyObj<Router>;
    });

    it('should allow navigation when user is not authenticated', () => {
        authService.isLoggedIn.and.returnValue(false);

        const result = TestBed.runInInjectionContext(() => noAuthGuard(null as any, null as any));

        expect(result).toBe(true);
        expect(router.createUrlTree).not.toHaveBeenCalled();
    });

    it('should redirect to dashboard when user is authenticated', () => {
        authService.isLoggedIn.and.returnValue(true);
        const urlTree = {} as any;
        router.createUrlTree.and.returnValue(urlTree);

        const result = TestBed.runInInjectionContext(() => noAuthGuard(null as any, null as any));

        expect(result).toBe(urlTree);
        expect(router.createUrlTree).toHaveBeenCalledWith(['/dashboard']);
    });
});
