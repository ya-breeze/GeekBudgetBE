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

    private readonly TOKEN_KEY = 'auth_token';
    private readonly isAuthenticatedSubject = new BehaviorSubject<boolean>(this.hasToken());
    public readonly isAuthenticated$ = this.isAuthenticatedSubject.asObservable();

    constructor() {}

    login(email: string, password: string): Observable<string> {
        const authData: AuthData = { email, password };
        return authorize(this.http, this.apiConfig.rootUrl, { body: authData }).pipe(
            map((response) => response.body.token),
            tap((token) => {
                this.setToken(token);
                this.isAuthenticatedSubject.next(true);
            }),
        );
    }

    logout(): void {
        this.removeToken();
        this.isAuthenticatedSubject.next(false);
        this.router.navigate(['/auth/login']);
    }

    getToken(): string | null {
        return localStorage.getItem(this.TOKEN_KEY);
    }

    private setToken(token: string): void {
        localStorage.setItem(this.TOKEN_KEY, token);
    }

    private removeToken(): void {
        localStorage.removeItem(this.TOKEN_KEY);
    }

    private hasToken(): boolean {
        return !!this.getToken();
    }

    isLoggedIn(): boolean {
        return this.hasToken();
    }
}
