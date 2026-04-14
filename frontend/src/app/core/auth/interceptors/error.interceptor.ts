import { HttpInterceptorFn, HttpErrorResponse, HttpClient } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { ApiConfiguration } from '../../api/api-configuration';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
    const authService = inject(AuthService);
    const http = inject(HttpClient);
    const apiConfig = inject(ApiConfiguration);

    return next(req).pipe(
        catchError((error: HttpErrorResponse) => {
            let errorMessage = 'An error occurred';

            if (error.error instanceof ErrorEvent) {
                // Client-side error
                errorMessage = error.error.message;
            } else {
                // Server-side error
                switch (error.status) {
                    case 401:
                        if (!req.url.includes('/logout') && !req.url.includes('/v1/user') && !req.url.includes('/auth/refresh')) {
                            return http.post(`${apiConfig.rootUrl}/auth/refresh`, {}, { withCredentials: true }).pipe(
                                switchMap(() => next(req)),
                                catchError(() => {
                                    authService.logout();
                                    return throwError(() => new Error('Session expired. Please login again.'));
                                }),
                            );
                        }
                        errorMessage = 'Session expired. Please login again.';
                        break;
                    case 403:
                        errorMessage = 'Access denied';
                        break;
                    case 404:
                        errorMessage = 'Resource not found';
                        break;
                    case 429:
                        errorMessage = 'Too many requests. Please try again later.';
                        break;
                    case 500:
                        errorMessage = 'Server error. Please try again later.';
                        break;
                    default:
                        errorMessage = error.error?.message || error.message;
                }
            }

            console.error('HTTP Error:', errorMessage, error);
            return throwError(() => new Error(errorMessage));
        }),
    );
};
