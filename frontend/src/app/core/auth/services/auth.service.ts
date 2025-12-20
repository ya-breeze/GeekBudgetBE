import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, tap, map } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { authorize } from '../../api/fn/auth/authorize';
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

    constructor() {
        // Optimistically assume not logged in, or check an endpoint to confirm session
        // For strictness, one would call /v1/users/me (if it exists) here.
    }

    login(email: string, password: string): Observable<void> {
        const authData: AuthData = { email, password };
        return authorize(this.http, this.apiConfig.rootUrl, { body: authData }).pipe(
            map(() => void 0),
            tap(() => {
                this.isAuthenticatedSubject.next(true);
            }),
        );
    }

    logout(): void {
        this.isAuthenticatedSubject.next(false);
        // Call backend logout to clear cookie
        this.http.post(`${this.apiConfig.rootUrl}/v1/logout`, {}).subscribe();
        this.router.navigate(['/auth/login']);
    }

    isLoggedIn(): boolean {
        // Fallback or explicit check.
        // Without access to document.cookie (HttpOnly), we rely on state or 401s
        return this.isAuthenticatedSubject.value;
    }
}
