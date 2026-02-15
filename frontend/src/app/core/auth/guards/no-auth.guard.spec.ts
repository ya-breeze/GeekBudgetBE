import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { noAuthGuard } from './no-auth.guard';
import { AuthService } from '../services/auth.service';

import { of } from 'rxjs';

describe('noAuthGuard', () => {
    let authService: jasmine.SpyObj<AuthService>;
    let router: jasmine.SpyObj<Router>;

    beforeEach(() => {
        const authServiceSpy = jasmine.createSpyObj('AuthService', ['checkAuth', 'isLoggedIn']);
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

    it('should allow navigation when user is not authenticated', (done) => {
        authService.checkAuth.and.returnValue(of(false));

        const result = TestBed.runInInjectionContext(() => noAuthGuard(null as any, null as any));

        if (result instanceof of(false).constructor) {
            (result as any).subscribe((val: any) => {
                expect(val).toBe(true);
                expect(router.createUrlTree).not.toHaveBeenCalled();
                done();
            });
        }
    });

    it('should redirect to dashboard when user is authenticated', (done) => {
        localStorage.setItem('gb_logged_in_hint', 'true');
        authService.checkAuth.and.returnValue(of(true));
        const urlTree = {} as any;
        router.createUrlTree.and.returnValue(urlTree);

        const result = TestBed.runInInjectionContext(() => noAuthGuard(null as any, null as any));

        if (result instanceof of(true).constructor) {
            (result as any).subscribe((val: any) => {
                expect(val).toBe(urlTree);
                expect(router.createUrlTree).toHaveBeenCalledWith(['/dashboard']);
                done();
            });
        }
    });
});
