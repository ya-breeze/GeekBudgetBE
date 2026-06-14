## ADDED Requirements

### Requirement: Inline login form in hero
The landing page SHALL display the login form (email, password, submit button) in the right column of the hero section, visible without scrolling, so users can authenticate without navigating to a separate page.

#### Scenario: Successful login from landing page
- **WHEN** user enters valid email and password in the hero login form and submits
- **THEN** system authenticates the user and navigates to `/dashboard`

#### Scenario: Failed login shows inline error
- **WHEN** user enters invalid credentials in the hero login form and submits
- **THEN** system displays an error message below the form fields without navigating away

#### Scenario: Form validation prevents empty submission
- **WHEN** user submits the hero login form with empty or invalid fields
- **THEN** system shows field-level validation errors and does not call the auth API

#### Scenario: Loading state during authentication
- **WHEN** user submits valid credentials and the auth request is in flight
- **THEN** the submit button shows a spinner and is disabled to prevent duplicate submissions

### Requirement: Two-column hero layout
The hero section SHALL use a two-column layout with the app tagline on the left and the login form on the right on desktop screens (>768px). On mobile screens (≤768px) the form SHALL be hidden by default behind a "Sign In" toggle button, expanding inline when tapped.

#### Scenario: Desktop hero layout
- **WHEN** the landing page is viewed on a screen wider than 768px
- **THEN** the hero shows tagline/subtitle on the left and the login form on the right, side by side, with no toggle needed

#### Scenario: Mobile hero — form collapsed by default
- **WHEN** the landing page is viewed on a screen 768px wide or narrower
- **THEN** the hero shows the tagline/subtitle and a "Sign In" button; the email and password fields are not visible

#### Scenario: Mobile hero — form expanded on tap
- **WHEN** user taps the "Sign In" button on mobile
- **THEN** the email and password fields expand inline below the tagline (animated), and the "Sign In" toggle button is replaced by the submit button

### Requirement: CTA section removed
The bottom "Ready to Master Your Finances?" CTA section SHALL NOT appear on the landing page, as the inline hero login form replaces its function.

#### Scenario: Page ends after Key Concepts section
- **WHEN** user scrolls to the bottom of the landing page
- **THEN** the page ends with the Key Concepts section followed directly by the footer
