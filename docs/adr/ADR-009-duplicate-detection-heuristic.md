# ADR-009: Duplicate Detection Heuristic

## Status
Accepted

## Context and Problem Statement

The same financial event often appears twice: a transfer between two of the user's own accounts is
imported from both banks, producing two transactions. These should be surfaced as duplicates — but
without false-flagging genuinely distinct transactions (e.g. two coffees of the same price on
consecutive days).

## Decision Drivers

- Catch cross-source duplicates (especially inter-account transfers)
- Avoid false positives on independent same-amount transactions
- Keep detection cheap enough to run periodically over a recent window

## Considered Options

- **Exact-match only** (same external id) — misses cross-source duplicates entirely
- **Amount + date proximity** — too loose, flags unrelated same-price purchases
- **Amount + date proximity + different source + opposite direction** — targeted heuristic

## Decision Outcome

Chosen: a background scan (every 24h, 30-day window per family) flags a pair as potential
duplicates only when **all** hold:

1. dates within **±2 days**
2. matching movement amounts
3. **different** external ids (came from different import sources)
4. **opposite** net directions (one net-incoming, one net-outgoing)

Condition 3 avoids re-flagging same-source transactions already deduped at import time. Condition 4
is the key false-positive filter: a real transfer always shows opposite signs across the two
accounts, whereas two independent purchases are both negative. Matches are linked bidirectionally
in `TransactionDuplicate` and marked in `SuspiciousReasons`.

### Pros

- Targets the real-world case (cross-account transfers) precisely
- The opposite-direction filter eliminates the most common false positive
- Bounded cost: recent window, per-family

### Cons

- ±2 days and 30-day window are fixed heuristics, not configurable per user
- Pairwise comparison is O(n²) within the window — fine now, could grow costly
- Cannot detect duplicates older than the window without a manual run
- A genuine duplicate with the *same* direction (rare) would be missed
