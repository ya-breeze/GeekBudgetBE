## Context

14 of 16 specs already match the OpenSpec CLI schema (`## Purpose` + `## Requirements`, `### Requirement:` / `#### Scenario:`). Two do not:

- `chart-palette` has a correct `## Purpose` but its requirements sit under `## ADDED Requirements` — a *delta* header that should have become `## Requirements` when the originating change was archived.
- `landing-page` has correct heading levels but no `## Purpose` / `## Requirements` wrapper and no title.

## Goals / Non-Goals

**Goals:**
- `openspec validate --specs` passes for all 16 capabilities.
- Zero change to requirement/scenario meaning.

**Non-Goals:**
- Editing the other 14 specs (already valid).
- Touching application code or the unrelated uncommitted work in the repo.

## Decisions

- **`chart-palette`:** rename `## ADDED Requirements` → `## Requirements`. Single-line header change; requirement/scenario bodies untouched.
- **`landing-page`:** prepend `# Landing Page Specification`, a one-line `## Purpose`, and `## Requirements` before the first `### Requirement:`. No requirement/scenario edits.
- **Verify with the tool:** acceptance is `openspec validate --specs` → 16 passed, 0 failed.
- **No spec deltas:** archive with `openspec archive --skip-specs`.

## Risks / Trade-offs

- [Header rename in chart-palette accidentally drops a requirement] → Compare requirement/scenario counts before/after; rely on `openspec validate`.
