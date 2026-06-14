## Why

The landing page currently requires two navigations to sign in: `/landing` → `/auth/login` → `/dashboard`. For a personal tool used daily, this is unnecessary friction. The landing page is also structured as a SaaS marketing page, which feels out of place for a private family finance app.

## What Changes

- The landing page hero section becomes a two-column layout: tagline on the left, login form on the right
- The login form moves from `/auth/login` (separate page) into the landing page hero directly
- The bottom CTA section ("Ready to Master Your Finances?") is removed — it duplicates the login call-to-action now embedded at the top
- The "Get Started" button is removed (replaced by the inline login form)
- The "Learn More" button (scrolls to features) is kept on the left column
- The Features, How It Works, and Key Concepts sections remain unchanged

## Capabilities

### New Capabilities

- `landing-page`: Landing page UX including the hero layout and inline login form

### Modified Capabilities

<!-- No existing specs change requirements -->

## Impact

- `frontend/src/app/features/auth/landing/` — hero layout and login form logic added
- `frontend/src/app/features/auth/login/` — `LoginComponent` can be kept for direct URL access but is no longer linked from the landing page
- `frontend/src/app/app.routes.ts` — no route changes required; `/auth/login` route stays for direct access
