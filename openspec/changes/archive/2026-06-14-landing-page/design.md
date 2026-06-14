## Context

The Angular frontend has two auth-related components:
- `LandingComponent` (`/landing`) — long scrollable page with hero, features, how-it-works, key-concepts, and a CTA section. The hero has "Get Started" (→ `/auth/login`) and "Learn More" (scroll to features) buttons.
- `LoginComponent` (`/auth/login`) — a standalone page with just a centered `mat-card` containing email/password form.

The home guard already handles smart routing: users with a `gb_logged_in_hint` in localStorage skip the landing entirely and go straight to `/dashboard`. The landing is therefore only shown on first visits or expired sessions.

## Goals / Non-Goals

**Goals:**
- Eliminate the `/landing` → `/auth/login` navigation step
- Embed the login form directly in the landing page hero (right column)
- Keep all existing informational sections (Features, How It Works, Key Concepts)
- Remove the bottom CTA section (redundant once login is at the top)

**Non-Goals:**
- Removing the `/auth/login` route (keep it for direct URL access)
- Changing the authentication flow, guards, or tokens
- Redesigning the Features / How It Works / Key Concepts sections
- Touching the Next.js (`app/`) frontend

## Decisions

### D1: Merge login logic into LandingComponent (not a sub-component)

The login form is small (2 fields + button + error) and tightly coupled to the hero layout. Pulling it into a separate `LandingLoginFormComponent` would add indirection with no reuse benefit. Instead, `LandingComponent` imports the same `ReactiveFormsModule` / `AuthService` that `LoginComponent` already uses and handles `onSubmit()` directly.

*Alternative considered:* Embed `<app-login>` inside the hero — rejected because `LoginComponent` wraps everything in a `mat-card` with its own padding and header, which would fight the hero layout.

### D2: Hero becomes a CSS Grid two-column layout; mobile uses a collapsible form

The hero section switches from a single centered `max-width: 900px` column to a CSS Grid:
```
grid-template-columns: 1fr 1fr
```
Left cell: badge + title + subtitle + "Learn More" button.
Right cell: login card (email, password, submit, error).

On mobile (`max-width: 768px`), the grid collapses to a single column. The login form is hidden by default — a "Sign In" button is shown instead. Tapping it reveals the form fields inline with a CSS/Angular animation. This keeps the hero clean on first glance while still allowing one-tap access to login.

*Alternative considered:* Show the login form below the tagline on mobile unconditionally — rejected because on typical phone viewports the form would be pushed off-screen, requiring a scroll to find it. The collapsible approach keeps the tagline visible and login accessible in one tap.

### D3: Keep `/auth/login` route and LoginComponent intact

`LoginComponent` remains as-is. The route stays in `app.routes.ts`. Only the landing page stops linking to it. This avoids breaking bookmarks or direct navigation.

### D4: Remove CTA section, keep footer

The bottom CTA card ("Ready to Master Your Finances?") is removed. The footer (`&copy; 2026 GeekBudget...`) remains. This keeps the page from ending abruptly without adding a redundant second login call-to-action.

## Risks / Trade-offs

- **Login form on a long page** → users who scroll past the hero might not find the form easily. Mitigation: the hero stays above the fold on all typical screen sizes; the form is the first thing visible.
- **Duplicate login logic** → `LandingComponent` and `LoginComponent` both handle form submission. Mitigation: the logic is trivial (2 fields + AuthService call); shared abstraction would be premature. If it grows, extract to a shared service method or a dedicated `LoginFormComponent` later.
- **Mobile layout** → stacking form below tagline on small screens is fine, but the form should come second (tagline first) so the app identity is clear before asking for credentials.
