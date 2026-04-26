import { HttpInterceptorFn, HttpErrorResponse, HttpClient } from '@angular/common/http';
import { inject } from '@angular/core';
import { BehaviorSubject, catchError, filter, switchMap, take, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { ApiConfiguration } from '../../api/api-configuration';

let isRefreshing = false;
const refreshSubject = new BehaviorSubject<boolean | null>(null);

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
                        if (
                            !req.url.includes('/logout') &&
                            !req.url.includes('/auth/refresh')
                        ) {
                            if (!isRefreshing) {
                                isRefreshing = true;
                                refreshSubject.next(null);
                                return http
                                    .post(
                                        `${apiConfig.rootUrl}/auth/refresh`,
                                        {},
                                        { withCredentials: true },
                                    )
                                    .pipe(
                                        switchMap(() => {
                                            isRefreshing = false;
                                            refreshSubject.next(true);
                                            return next(req);
                                        }),
                                        catchError(() => {
                                            isRefreshing = false;
                                            refreshSubject.next(false);
                                            authService.logout();
                                            return throwError(
                                                () =>
                                                    new Error(
                                                        'Session expired. Please login again.',
                                                    ),
                                            );
                                        }),
                                    );
                            }
                            // Another request is already refreshing — wait for it
                            return refreshSubject.pipe(
                                filter((result) => result !== null),
                                take(1),
                                switchMap((success) => {
                                    if (success) {
                                        return next(req);
                                    }
                                    return throwError(
                                        () => new Error('Session expired. Please login again.'),
                                    );
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
