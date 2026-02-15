import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, tap, map, catchError, of } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { authorize } from '../../api/fn/auth/authorize';
import { getUser } from '../../api/fn/user/get-user';
import { AuthData } from '../../api/models/auth-data';
import { ApiConfiguration } from '../../api/api-configuration';

@Injectable({
    providedIn: 'root',
})
export class AuthService {
    private readonly http = inject(HttpClient);
    private readonly apiConfig = inject(ApiConfiguration);
    private readonly router = inject(Router);

    // Track authentication state locally. ideally verification uses an endpoint /me
    private readonly isAuthenticatedSubject = new BehaviorSubject<boolean>(false);
    public readonly isAuthenticated$ = this.isAuthenticatedSubject.asObservable();
    private isInitialized = false;
    private readonly LOGGED_IN_HINT_KEY = 'gb_logged_in_hint';

    constructor() {
        // Optimistically assume not logged in
    }

    login(email: string, password: string): Observable<void> {
        const authData: AuthData = { email, password };
        return authorize(this.http, this.apiConfig.rootUrl, { body: authData }).pipe(
            map(() => void 0),
            tap(() => {
                this.isInitialized = true;
                this.isAuthenticatedSubject.next(true);
                localStorage.setItem(this.LOGGED_IN_HINT_KEY, 'true');
            }),
        );
    }

    logout(): void {
        this.isAuthenticatedSubject.next(false);
        this.isInitialized = true;
        localStorage.setItem(this.LOGGED_IN_HINT_KEY, 'false');
        // Call backend logout to clear cookie
        this.http.post(`${this.apiConfig.rootUrl}/v1/logout`, {}).subscribe();

        // Only redirect to login if we're not on a public page
        const currentUrl = this.router.url;
        const isPublicRoute =
            currentUrl === '/' || currentUrl === '/landing' || currentUrl.startsWith('/auth');

        if (!isPublicRoute) {
            this.router.navigate(['/auth/login']);
        }
    }

    isLoggedIn(): boolean {
        return this.isAuthenticatedSubject.value;
    }

    checkAuth(): Observable<boolean> {
        if (this.isInitialized) {
            return of(this.isAuthenticatedSubject.value);
        }

        return getUser(this.http, this.apiConfig.rootUrl).pipe(
            map(() => {
                this.isInitialized = true;
                this.isAuthenticatedSubject.next(true);
                localStorage.setItem(this.LOGGED_IN_HINT_KEY, 'true');
                return true;
            }),
            catchError(() => {
                this.isInitialized = true;
                this.isAuthenticatedSubject.next(false);
                localStorage.setItem(this.LOGGED_IN_HINT_KEY, 'false');
                return of(false);
            }),
        );
    }
}
