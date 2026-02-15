import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '../services/auth.service';

import { map, of } from 'rxjs';

export const noAuthGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // Optimization: If no hint of being logged in, allow access immediately
    const hasLoginHint = localStorage.getItem('gb_logged_in_hint') === 'true';
    if (!hasLoginHint && !authService.isLoggedIn()) {
        return of(true);
    }

    return authService.checkAuth().pipe(
        map((isLoggedIn) => {
            if (!isLoggedIn) {
                return true;
            }
            return router.createUrlTree(['/dashboard']);
        }),
    );
};
