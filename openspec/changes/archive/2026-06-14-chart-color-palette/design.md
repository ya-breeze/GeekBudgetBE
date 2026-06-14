## Context

Four Angular components each contain a local color array used to color Chart.js datasets:

| Component | Colors defined | Chart types |
|---|---|---|
| `dashboard.component.ts` | 15 oklch hues | Stacked bar, rank dots |
| `expense-report.component.ts` | 7 rgba values | Line, pie |
| `tag-analytics.component.ts` | 7 rgba values (copy of above) | Line, pie |
| `balance-report.component.ts` | 5 rgba values | Stacked area |

No shared abstraction exists. Each component independently decides what colors to use.

## Goals / Non-Goals

**Goals:**
- Single source of truth: one file defines the 25-color palette
- All four chart components use `ChartPaletteService` instead of local arrays
- Helpers for solid hex and rgba-with-alpha cover all existing usage patterns

**Non-Goals:**
- Theming or dark-mode support
- Dynamic or user-configurable palettes
- Next.js (`app/`) frontend — out of scope for now (Angular only)
- Backend changes of any kind

## Decisions

### Decision: Injectable Angular service, not a plain constant

**Chosen:** `@Injectable({ providedIn: 'root' }) ChartPaletteService`

**Alternatives considered:**
- Plain `export const CHART_PALETTE = [...]` in a shared file — works, but not consistent with how other shared concerns (auth, currency, layout) are structured in this codebase. Services are the Angular convention here.
- CSS custom properties — attractive for theming, but Chart.js doesn't read CSS variables for dataset colors without custom plugins.

**Rationale:** A service fits the existing pattern (`AccountService`, `CurrencyService`, etc.), is easily injected into components, and is straightforward to mock in tests.

---

### Decision: Hex strings as the canonical format, rgba computed on demand

**Chosen:** Store palette as `string[]` of hex values. `getColorWithAlpha` parses hex → rgba at call time.

**Alternatives considered:**
- Store as `[r, g, b]` tuples — more flexible but verbose to define and harder to read/verify visually.
- Store rgba strings directly — can't derive solid colors from them cleanly; the existing workaround (`replace('0.6', '1')`) is exactly the bug this change removes.

**Rationale:** Hex is the most readable format for a color palette, trivial to convert to rgba, and visually verifiable against design references (Tableau's published palette uses hex).

---

### Decision: Palette order follows Tableau 20 pairs, extras appended at the end

**Chosen:** Indices 0–19 = Tableau 20 (10 dark/light pairs), indices 20–24 = 5 extras (cyan, light cyan, chartreuse, dark indigo, magenta).

**Rationale:** This gives the first 10 categories (the most-spent, highest-priority ones) the most visually prominent, well-tested Tableau colors. The extras are for long-tail categories.

---

### Decision: `colorsForN` returns hex (not rgba)

Components that need transparency call `getColorWithAlpha` explicitly. `colorsForN` is for cases like pie charts where Chart.js accepts an array of colors and transparency isn't needed (or is applied uniformly).

## Risks / Trade-offs

- **Color order changes for existing users** — the dashboard stacked bar will reorder colors when the oklch palette is replaced. This is expected and desired (the fix is the point), but users will notice the change. → No mitigation needed; this is the intent.
- **Hex-to-rgba parsing** — `getColorWithAlpha` must correctly parse 6-digit hex. Edge cases (shorthand 3-digit hex) are not in the palette, so a simple regex is sufficient. → Keep palette hex values always 6 digits.

## Migration Plan

1. Create `ChartPaletteService` in `frontend/src/app/shared/services/`.
2. Update `dashboard.component.ts`: remove `HUES` array and oklch formula; inject service.
3. Update `expense-report.component.ts`: remove local `colors` array; inject service.
4. Update `tag-analytics.component.ts`: same as expense-report.
5. Update `balance-report.component.ts`: same pattern.
6. Run `make lint` to verify no ESLint issues.
7. Deploy to WIP, visually verify all four chart views render with the new palette.

Rollback: revert the four component files to their previous local color arrays. No data migration needed.

## Open Questions

*(none)*
