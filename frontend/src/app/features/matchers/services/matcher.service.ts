import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, map } from 'rxjs';
import { ApiConfiguration } from '../../../core/api/api-configuration';
import { Matcher } from '../../../core/api/models/matcher';
import { MatcherNoId } from '../../../core/api/models/matcher-no-id';
import { getMatchers } from '../../../core/api/fn/matchers/get-matchers';
import { createMatcher } from '../../../core/api/fn/matchers/create-matcher';
import { updateMatcher } from '../../../core/api/fn/matchers/update-matcher';
import { deleteMatcher } from '../../../core/api/fn/matchers/delete-matcher';

@Injectable({
  providedIn: 'root',
})
export class MatcherService {
  private readonly http = inject(HttpClient);
  private readonly apiConfig = inject(ApiConfiguration);

  readonly matchers = signal<Matcher[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  loadMatchers(): Observable<Matcher[]> {
    this.loading.set(true);
    this.error.set(null);

    return getMatchers(this.http, this.apiConfig.rootUrl).pipe(
      map((response) => response.body),
      tap({
        next: (matchers) => {
          this.matchers.set(matchers);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to load matchers');
          this.loading.set(false);
        },
      })
    );
  }

  create(matcher: MatcherNoId): Observable<Matcher> {
    this.loading.set(true);
    this.error.set(null);

    return createMatcher(this.http, this.apiConfig.rootUrl, { body: matcher }).pipe(
      map((response) => response.body),
      tap({
        next: (matcher) => {
          this.matchers.update((matchers) => [...matchers, matcher]);
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to create matcher');
          this.loading.set(false);
        },
      })
    );
  }

  update(id: string, matcher: MatcherNoId): Observable<Matcher> {
    this.loading.set(true);
    this.error.set(null);

    return updateMatcher(this.http, this.apiConfig.rootUrl, { id, body: matcher }).pipe(
      map((response) => response.body),
      tap({
        next: (updatedMatcher) => {
          this.matchers.update((matchers) =>
            matchers.map((m) => (m.id === id ? updatedMatcher : m))
          );
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to update matcher');
          this.loading.set(false);
        },
      })
    );
  }

  delete(id: string): Observable<void> {
    this.loading.set(true);
    this.error.set(null);

    return deleteMatcher(this.http, this.apiConfig.rootUrl, { id }).pipe(
      map(() => undefined),
      tap({
        next: () => {
          this.matchers.update((matchers) => matchers.filter((m) => m.id !== id));
          this.loading.set(false);
        },
        error: (err) => {
          this.error.set(err.message || 'Failed to delete matcher');
          this.loading.set(false);
        },
      })
    );
  }
}

