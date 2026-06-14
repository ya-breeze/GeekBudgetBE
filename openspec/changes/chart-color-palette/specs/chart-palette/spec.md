# chart-palette Specification

## Purpose

`ChartPaletteService` is the single source of truth for all chart and diagram colors in the
Angular frontend. It exposes a fixed 25-color palette and helper methods so every chart
component produces consistent, visually distinct colors without defining its own palette.

## ADDED Requirements

### Requirement: Canonical 25-color palette

The frontend SHALL define a single canonical palette of exactly 25 hex colors derived from the
Tableau 20 set (10 dark/light pairs) extended with 5 additional hues to cover the full set of
expense categories. No component SHALL define its own color array for charts or diagrams.

The palette in index order:
`#4E79A7`, `#A0CBE8`, `#F28E2B`, `#FFBE7D`, `#59A14F`, `#8CD17D`, `#B6992D`, `#F1CE63`,
`#499894`, `#86BCB6`, `#E15759`, `#FF9D9A`, `#79706E`, `#BAB0AC`, `#D37295`, `#FABFD2`,
`#B07AA1`, `#D4A6C8`, `#9D7660`, `#D7B5A6`, `#17BECF`, `#9EDAE5`, `#BCBD22`, `#393B79`,
`#CE6DBD`.

#### Scenario: Palette wraps for index beyond 25

- **WHEN** a chart requests color at index ≥ 25
- **THEN** the service returns the color at `index % 25`, cycling through the palette

### Requirement: Solid color lookup

The service SHALL expose `getColor(index: number): string` returning the hex color at
`index % 25`.

#### Scenario: First category gets first palette color

- **WHEN** a component calls `getColor(0)`
- **THEN** the service returns `#4E79A7`

#### Scenario: Arbitrary index returns correct color

- **WHEN** a component calls `getColor(4)`
- **THEN** the service returns `#59A14F`

### Requirement: Semi-transparent color variant

The service SHALL expose `getColorWithAlpha(index: number, alpha: number): string` returning
an `rgba(r, g, b, alpha)` string derived from the palette hex at `index % 25`. Alpha is a
0–1 float.

#### Scenario: Alpha variant matches base color

- **WHEN** a component calls `getColorWithAlpha(0, 0.6)`
- **THEN** the service returns `rgba(78, 121, 167, 0.6)`

### Requirement: Bulk color array

The service SHALL expose `colorsForN(n: number): string[]` returning an array of `n` hex
colors cycling through the palette.

#### Scenario: Fewer than 25 colors requested

- **WHEN** a component calls `colorsForN(7)`
- **THEN** the service returns the first 7 palette colors as an array

#### Scenario: More than 25 colors requested

- **WHEN** a component calls `colorsForN(30)`
- **THEN** the service returns 30 colors, cycling through the palette from index 0 again after index 24

### Requirement: All chart components use ChartPaletteService

Every Angular component that renders a Chart.js chart SHALL inject `ChartPaletteService` and
use its methods to assign `backgroundColor` and `borderColor`. No component SHALL contain a
local `colors`, `HUES`, or `PALETTE` array for chart coloring purposes.

#### Scenario: Dashboard stacked bar uses palette

- **WHEN** the dashboard renders the month-over-month stacked bar chart
- **THEN** each expense category bar segment uses the color at the category's index in the palette

#### Scenario: Expense report charts use palette

- **WHEN** the expense report renders its line and pie charts
- **THEN** each dataset/slice uses `getColor(index)` or `getColorWithAlpha(index, alpha)`

#### Scenario: Tag analytics charts use palette

- **WHEN** the tag analytics view renders its monthly line chart and pie chart
- **THEN** each tag series uses colors drawn from `ChartPaletteService`

#### Scenario: Balance report area chart uses palette

- **WHEN** the balance report renders its stacked area chart
- **THEN** each account series uses `getColorWithAlpha(index, 0.5)` for fill and `getColor(index)` for border
