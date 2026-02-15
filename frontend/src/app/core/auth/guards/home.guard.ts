import { inject } from '@angular/core';
import { CanActivateFn, Router, UrlTree } from '@angular/router';
import { AuthService } from '../services/auth.service';

import { map, of } from 'rxjs';

export const homeGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // Optimization: If we have a hint that the user is definitely not logged in,
    // we can skip the backend check and solve the 'landing should not cause auth check' requirement.
    const hasLoginHint = localStorage.getItem('gb_logged_in_hint') === 'true';

    if (!hasLoginHint && !authService.isLoggedIn()) {
        return of(router.parseUrl('/landing'));
    }

    return authService.checkAuth().pipe(
        map((isLoggedIn) => {
            if (isLoggedIn) {
                return router.parseUrl('/dashboard');
            } else {
                return router.parseUrl('/landing');
            }
        }),
    );
};
