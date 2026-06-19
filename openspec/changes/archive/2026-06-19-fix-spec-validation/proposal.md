## Why

`openspec validate --specs` fails for 2 of 16 canonical specs, which blocks `openspec validate` from passing cleanly and risks `openspec archive` issues on future changes:

- `chart-palette` — uses the delta header `## ADDED Requirements` (an archive artifact) where a canonical spec must use `## Requirements`.
- `landing-page` — missing the `## Purpose` and `## Requirements` sections entirely; the file starts directly at `### Requirement:`.

The other 14 specs already conform. This brings the 2 stragglers in line so the whole suite validates.

## What Changes

- `chart-palette/spec.md` — rename the `## ADDED Requirements` header to `## Requirements`.
- `landing-page/spec.md` — add a title, a one-line `## Purpose`, and a `## Requirements` wrapper above the existing requirements.

Format-only: no requirement statement or scenario is added, removed, or semantically changed. Acceptance is `openspec validate --specs` passing for all 16 capabilities.

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
<!-- none — reformats existing canonical specs for tool compatibility; no requirement behavior changes, no spec deltas. Archive with `--skip-specs`. -->

## Impact

- `openspec/specs/chart-palette/spec.md`, `openspec/specs/landing-page/spec.md` (2 files).
- No application code, API, or behavior changes.
- Unblocks `openspec validate --specs` (16/16) for the project.
