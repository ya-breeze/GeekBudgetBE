# Transaction UI Patterns

## Implementation Guidelines

### 1. ID-to-Name Resolution
When displaying entities that have associated IDs (Accounts, Currencies, Matchers), always resolve these to human-readable names.
- **Signal-based Maps**: Use Angular `computed` signals to create lookup maps from services.
- **Safe Fallbacks**: If a name is missing or the ID is not found, fallback to the ID itself or 'N/A' rather than leaving an empty space.

### 2. Information Density and Icons
- **Icon Visibility**: Transactions should show status icons (Suspicious, Auto/Manual Matched, Merged) in both list and detail views. Use consistent icons and tooltips.
- **Avoid Redundant Badges**: Do not show an "Active" badge if most transactions share this state. Focus on highlighting non-standard states (e.g., Merged, Suspicious).
- **Column Merging**: In tables, prefer merging related information (e.g., Currency next to Amount) to reduce horizontally sprawling columns.
- **Account Icons**: Use the `AccountDisplayComponent` to show account names along with their icons for immediate visual recognition.

### 3. Detail View and Navigation
- **Accordion Layout**: Use `mat-accordion` for complex detail views to allow progressive disclosure of information (e.g., raw source data, metadata).
- **History Preservation**: Use `Location.back()` for "Back" buttons in detail views. This allows users who navigated between linked transactions to go back through their specific path rather than jumping straight back to the transaction list.
- **Drill-down Context**: When linking from overview components to reports, use query parameters to transmit filter state (Account ID, date range).

- **DRY Dialogs**: Centralize common complex dialog results (like account editing) in service methods (e.g., `AccountService.handleAccountDialogResult`) to avoid duplicating 50+ lines of callback logic.

### 4. Layout & Spacing
- **Flexbox Centering**: For `mat-icon-button` or similar small interactive elements, use `display: flex; align-items: center; justify-content: center;` instead of `line-height` for vertically centering icons. This is more reliable across different browsers and themes.
- **Card Icon Placement**: In card headers, place action icons (like settings) to the left of the title or use a toolbar-like arrangement. Use negative margins (e.g., `margin-left: -8px`) to pull icons closer to the card edge for a cleaner look.

### 5. Filtered Drill-down Pattern
- **Preserved Context**: When navigating from a dashboard summary (e.g., an account balance) to a detail report (e.g., Balance Report), always pass the context via query parameters:
  - `accountId`: Specific entity ID.
  - `from` and `to`: ISO timestamps for the relevant period (e.g., last 12 months).
- **Reactive Receipt**: Detail components should subscribe to `route.queryParams` to automatically initialize their filters, ensuring the user sees exactly what they clicked on without extra steps.
