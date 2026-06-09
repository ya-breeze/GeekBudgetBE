# ADR-008: Matcher Trust via Rolling Confirmation History

## Status
Accepted

## Context and Problem Statement

Matchers auto-categorize imported transactions. Letting a brand-new or unreliable matcher
auto-convert transactions risks miscategorizing data at scale. The system needs a way to decide
*when* a matcher is trustworthy enough to act without the user.

## Decision Drivers

- Avoid auto-applying immature or error-prone matchers
- Let a matcher earn (and lose) trust based on real outcomes
- Keep the trust signal simple and local to the matcher

## Considered Options

- **Manual "auto" toggle** — user explicitly marks a matcher as trusted
- **Rolling confirmation history with a perfect-match threshold** — trust derived from recent
  accept/reject outcomes
- **ML confidence scoring** — learned probability per matcher

## Decision Outcome

Chosen: each matcher keeps a rolling `ConfirmationHistory []bool`, capped at a configurable maximum
(`MatcherConfirmationHistoryMax`, default 10). A matcher is a **"perfect match"** — eligible for
auto-conversion — only when it has **at least 10 confirmations and all are `true`**. Accepting a
suggestion appends `true`; the cap drops the oldest entry, so a matcher can both earn trust and
have it eroded by recent failures.

### Pros

- Trust is grounded in actual recent user confirmations, not a guess
- Self-correcting: a recent bad outcome removes auto-match eligibility
- Simple to compute and reason about; no model to train

### Cons

- The "10 and all true" threshold is a hard-coded heuristic, not tuned per user
- The window is short — a matcher with a long clean record is treated the same as one with exactly
  10 confirmations
- Requires enough manual confirmations before any automation kicks in
- See ADR-009: auto-matching still guards against duplicates before converting
