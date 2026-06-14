## 1. LandingComponent — add login form logic

- [ ] 1.1 Import `FormBuilder`, `FormGroup`, `Validators`, `ReactiveFormsModule`, `AuthService`, `Router`, `MatFormFieldModule`, `MatInputModule`, `MatProgressSpinnerModule` into `LandingComponent`
- [ ] 1.2 Add `loginForm: FormGroup`, `isLoading = false`, `errorMessage = ''` fields and build the form in the constructor (email + required/email validators, password + required/minLength(4))
- [ ] 1.3 Add `onSubmit()` method: validate form, call `authService.login()`, navigate to `/dashboard` on success, set `errorMessage` on error

## 2. LandingComponent — update hero template

- [ ] 2.1 Replace the single-column hero content with a two-column grid: left cell keeps badge + title + subtitle + "Learn More" button; right cell contains the login form
- [ ] 2.2 Add the login form markup in the right cell (email field, password field, submit button with spinner, error message display)
- [ ] 2.3 Remove the "Get Started" button (`routerLink="/auth/login"`) from the hero

## 3. LandingComponent — update hero styles

- [ ] 3.1 Change `.hero-content` from `max-width: 900px; margin: 0 auto` to a CSS Grid with `grid-template-columns: 1fr 1fr; gap: 3rem; align-items: center`
- [ ] 3.2 Add mobile breakpoint (`max-width: 768px`) to collapse the grid to a single column
- [ ] 3.3 Style the login card in the right column (use `mat-card` with appropriate padding; ensure it matches the app's Material theme)

## 4. LandingComponent — remove CTA section

- [ ] 4.1 Remove the entire `<section class="cta-section">` block from `landing.component.html`
- [ ] 4.2 Remove `.cta-section` and `.cta-card` CSS rules from `landing.component.scss`

## 5. Verify

- [ ] 5.1 Run `make lint` (frontend ESLint + Angular checks) and fix any issues
- [ ] 5.2 Deploy to WIP stack and manually verify: login works from landing hero, error shows on bad credentials, mobile layout stacks correctly, Key Concepts section is the last content section before the footer
