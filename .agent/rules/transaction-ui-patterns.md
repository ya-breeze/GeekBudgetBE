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

### 3. Detail View Navigation
- **Accordion Layout**: Use `mat-accordion` for complex detail views to allow progressive disclosure of information (e.g., raw source data, metadata).
- **History Preservation**: Use `Location.back()` for "Back" buttons in detail views. This allows users who navigated between linked transactions to go back through their specific path rather than jumping straight back to the transaction list.

### 5. Manual Merge & Cross-Month Selection
- **Selection Persistence**: Use a dedicated `TransactionSelectionService` to manage selected transactions globally. This allows users to select a transaction in one month, navigate, and select another in a different month.
- **Floating Bar Feedback**: Always show a persistent floating bar (`selection-floating-bar`) when transactions are selected. It should provide:
    - Count of selected items.
    - Brief preview (Date/Description) of selected items.
    - Action buttons (Merge, Clear).
- **Selection Constraints**: 
    - Enforce a strict limit of 2 transactions for merging. 
    - Disable checkboxes and show informative tooltips when the limit is reached.
- **Responsive Comparison**: Manual merge dialogs must use responsive grids (`grid-template-columns: 1fr` on small screens) to prevent clipping when comparing dense transaction data side-by-side.
- **Layout Spacing**: Add appropriate bottom padding or spacers in lists where a floating bar may overlap the footer content.
