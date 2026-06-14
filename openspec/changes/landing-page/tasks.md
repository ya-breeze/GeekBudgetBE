## 1. LandingComponent — add login form logic

- [x] 1.1 Import `FormBuilder`, `FormGroup`, `Validators`, `ReactiveFormsModule`, `AuthService`, `Router`, `MatFormFieldModule`, `MatInputModule`, `MatProgressSpinnerModule` into `LandingComponent`
- [x] 1.2 Add `loginForm: FormGroup`, `isLoading = false`, `errorMessage = ''` fields and build the form in the constructor (email + required/email validators, password + required/minLength(4))
- [x] 1.3 Add `onSubmit()` method: validate form, call `authService.login()`, navigate to `/dashboard` on success, set `errorMessage` on error

## 2. LandingComponent — update hero template

- [x] 2.1 Replace the single-column hero content with a two-column grid: left cell keeps badge + title + subtitle + "Learn More" button; right cell contains the login form
- [x] 2.2 Add the login form markup in the right cell (email field, password field, submit button with spinner, error message display)
- [x] 2.3 Remove the "Get Started" button (`routerLink="/auth/login"`) from the hero

## 3. LandingComponent — update hero styles

- [x] 3.1 Change `.hero-content` from `max-width: 900px; margin: 0 auto` to a CSS Grid with `grid-template-columns: 1fr 1fr; gap: 3rem; align-items: center`
- [x] 3.2 Add mobile breakpoint (`max-width: 768px`) to collapse the grid to a single column
- [x] 3.3 Style the login card in the right column (use `mat-card` with appropriate padding; ensure it matches the app's Material theme)

## 4. LandingComponent — mobile collapsible login

- [x] 4.1 Add `showLoginForm = false` boolean field to `LandingComponent` and a `toggleLoginForm()` method
- [x] 4.2 In the mobile breakpoint template: show a "Sign In" button that calls `toggleLoginForm()`; wrap the form fields in `@if (showLoginForm)` so they only render when expanded
- [x] 4.3 Add Angular `animations` to `LandingComponent` (or use CSS `max-height` transition) to animate the form expansion/collapse smoothly

## 5. LandingComponent — remove CTA section

- [x] 5.1 Remove the entire `<section class="cta-section">` block from `landing.component.html`
- [x] 5.2 Remove `.cta-section` and `.cta-card` CSS rules from `landing.component.scss`

## 6. Verify

- [x] 6.1 Run `make lint` (frontend ESLint + Angular checks) and fix any issues
- [x] 6.2 Deploy to WIP stack and manually verify: login works from landing hero, error shows on bad credentials, desktop shows two-column hero, mobile shows tagline + "Sign In" toggle that expands the form, Key Concepts section is the last content section before the footer
