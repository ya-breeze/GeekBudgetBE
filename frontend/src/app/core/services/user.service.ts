import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../api/api-configuration';
import { User } from '../api/models/user';
import { UserPatchBody } from '../api/models/user-patch-body';
import { getUser } from '../api/fn/user/get-user';
import { updateUserFavoriteCurrency } from '../api/fn/user/update-user-favorite-currency';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly user = signal<User | null>(null);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadUser(): Observable<User> {
    this.loading.set(true);
    this.error.set(null);

    return getUser(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (user) => {
          this.user.set(user);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load user');
          this.loading.set(false);
        },
      })
    );
  }

  updateFavoriteCurrency(favoriteCurrencyId: string | null): Observable<User> {
    this.loading.set(true);
    this.error.set(null);

    const body: UserPatchBody = {
      favoriteCurrencyId: favoriteCurrencyId ?? undefined,
    };

    return updateUserFavoriteCurrency(this.http, this.apiConfig.rootUrl, { body }).pipe(
      map((response) => response.body),
      tap({
        next: (user) => {
          this.user.set(user);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to update favorite currency');
          this.loading.set(false);
        },
      })
    );
  }
}

