import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '../services/auth.service';

import { map } from 'rxjs';

export const authGuard: CanActivateFn = () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    return authService.checkAuth().pipe(
        map((isLoggedIn) => {
            if (isLoggedIn) {
                return true;
            }
            return router.createUrlTree(['/landing']);
        }),
    );
};
