## Why

Each chart component in the Angular frontend defines its own color list independently — producing inconsistent palettes, repeated hues across the 22 expense categories, and visually cluttered graphs. A single canonical palette eliminates drift and gives every chart a coherent, professional appearance.

## What Changes

- Introduce a shared `ChartPaletteService` in `frontend/src/app/shared/` exposing a 25-color Tableau-inspired palette.
- Remove local color arrays from `dashboard.component.ts`, `expense-report.component.ts`, `tag-analytics.component.ts`, and `balance-report.component.ts`.
- Update all four components to inject `ChartPaletteService` and call its helpers (`getColor`, `getColorWithAlpha`, `colorsForN`) instead of inline arrays.

## Capabilities

### New Capabilities

- `chart-palette`: A shared Angular service that is the single source of truth for all chart/diagram colors in the frontend. Exposes a 25-color palette and helper methods for solid and semi-transparent variants.

### Modified Capabilities

*(none — no API or spec-level behavior changes)*

## Impact

- **Frontend only** — no backend, API, or OpenAPI changes.
- Affected files: `dashboard.component.ts`, `expense-report.component.ts`, `tag-analytics.component.ts`, `balance-report.component.ts`.
- New file: `frontend/src/app/shared/services/chart-palette.service.ts`.
- No new dependencies.
