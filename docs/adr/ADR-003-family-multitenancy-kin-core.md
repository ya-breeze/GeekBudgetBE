# ADR-003: Family-Based Multi-Tenancy via kin-core

## Status
Accepted

## Context and Problem Statement

Household finances are shared: spouses and family members need to see and edit the same accounts,
transactions, and budgets. The tenancy boundary is therefore the *family*, not the individual
user. Auth and user management should be reusable across sibling apps in the same ecosystem.

## Decision Drivers

- Multiple users must share one financial dataset
- Isolation must be enforced at the data layer, not just the UI
- Auth/JWT/user management should be shared with sibling apps (Diary, KinCart)

## Considered Options

- **Per-user scoping** — every row keyed by `UserID`
- **Family scoping via `kin-core`** — every row keyed by `FamilyID`, shared auth library
- **Third-party auth service** — offload auth entirely

## Decision Outcome

Chosen: **family scoping via `github.com/ya-breeze/kin-core`**. Every domain model carries a
`FamilyID uuid.UUID`. Users belong to a family (`models.Family` embeds `kin-core` Family with a
`Users` slice). Middleware extracts the family from the JWT and injects it into the request
context; handlers retrieve it via `constants.GetFamilyID(ctx)`. All storage methods take
`familyID` as their first scoping argument.

### Pros

- Household sharing is the default and requires no per-user ACLs
- Isolation is enforced uniformly — every query filters by `FamilyID`
- Auth, password hashing, JWT, and refresh/blacklist tokens are reused via kin-core

### Cons

- No per-user data isolation within a family — all members see everything
- Adding per-user privacy later would be a breaking schema and API change
- External dependency on an internal library; kin-core changes can break the app
