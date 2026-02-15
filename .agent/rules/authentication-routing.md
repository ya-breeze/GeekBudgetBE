# Authentication and Routing Strategy

The application uses a lazy authentication initialization strategy to prioritize showing a landing page without hitting the backend, while still ensuring authenticated users are redirected to the dashboard.

## Key Mechanisms

### 1. LocalStorage Login Hint
- **Key**: `gb_logged_in_hint`
- **Values**: `'true'` or `'false'`.
- **Purpose**: Allows guards to skip the mandatory backend `checkAuth()` call on the first page load for new or explicitly logged-out users. This ensures the landing page loads instantly without backend interference.
- **Maintenance**: Updated by `AuthService` during `login()`, `logout()`, and successful/failed `checkAuth()` calls.

### 2. Guard Flow
- **`homeGuard`**: Used on the root path `/` and wildcard `**`.
    - If `gb_logged_in_hint` is missing or `'false'`, it immediately redirects to `/landing` (asynchronously returns `of(urlTree)`).
    - If `gb_logged_in_hint` is `'true'`, it calls `authService.checkAuth()` to verify the session.
- **`authGuard`**: Protects dashboard routes. Always calls `authService.checkAuth()`.
- **`noAuthGuard`**: Protects login and landing pages. Redirects authenticated users to `/dashboard`. Uses the same `localStorage` optimization as `homeGuard`.

### 3. Startup Auth (Lazy Initialization)
- **Do NOT** use `APP_INITIALIZER` to check authentication. This causes a mandatory delay and backend request on every page load, including the landing page.
- Authentication is initialized "just in time" by the first guard that requires it.

### 4. Interceptor Handling
- The `errorInterceptor` must **ignore** `401 Unauthorized` errors from the `/v1/user` endpoint. These are handled gracefully by `AuthService.checkAuth()` via `catchError`. Calling `logout()` (which triggers a redirect) in response to a 401 during the initial auth check will break the redirection flow.

## Testing Async Guards
When testing guards that now return `Observable<boolean | UrlTree>`:
1. Use `(done)` in the `it` block.
2. Check if the result is an `Observable` and `.subscribe()` to it.
3. Call `done()` inside the subscription.
4. Ensure `localStorage` hints are mocked or set appropriately in tests if the guard depends on them.
