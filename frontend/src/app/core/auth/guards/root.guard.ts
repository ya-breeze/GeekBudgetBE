import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const rootGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    // Redirect based on authentication status
    if (authService.isLoggedIn()) {
        // Logged in users go straight to dashboard
        return router.createUrlTree(['/dashboard']);
    }

    // New/logged out users see the landing page
    return router.createUrlTree(['/landing']);
};
