## 1. Create ChartPaletteService

- [x] 1.1 Create `frontend/src/app/shared/services/chart-palette.service.ts` with the 25-color palette constant and `getColor`, `getColorWithAlpha`, `colorsForN` methods
- [x] 1.2 Verify hex-to-rgba conversion in `getColorWithAlpha` works correctly for all 25 palette entries

## 2. Update dashboard.component.ts

- [x] 2.1 Remove the `HUES` constant and `oklch(0.62 0.14 …)` formula
- [x] 2.2 Inject `ChartPaletteService` and replace `accountColorMap` computed to use `getColor(i)`
- [x] 2.3 Verify stacked bar chart and rank dots still render with new colors

## 3. Update expense-report.component.ts

- [x] 3.1 Remove local `colors` array (7 rgba values)
- [x] 3.2 Inject `ChartPaletteService` and replace color references with `getColor` / `getColorWithAlpha`

## 4. Update tag-analytics.component.ts

- [x] 4.1 Remove local `colors` array (7 rgba values)
- [x] 4.2 Inject `ChartPaletteService` and replace color references with `getColor` / `getColorWithAlpha`

## 5. Update balance-report.component.ts

- [x] 5.1 Remove local `colors` array (5 rgba values)
- [x] 5.2 Inject `ChartPaletteService` and replace color references with `getColorWithAlpha(index, 0.5)` for fill and `getColor(index)` for border

## 6. Verify and lint

- [x] 6.1 Run `make lint` and fix any ESLint issues
- [ ] 6.2 Deploy to WIP stack and visually verify dashboard, expense report, tag analytics, and balance report all render with the new consistent palette
