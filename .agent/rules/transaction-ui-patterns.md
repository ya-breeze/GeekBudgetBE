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

### 4. Archived (Merged) Transactions
- **Separate Retrieval**: Merged transactions are archived in a separate table. Use the specific `mergedTransactions` endpoint (e.g., `GET /v1/mergedTransactions/{id}`) when a standard transaction lookup fails or when navigating from a "Merged" link.
- **Visual Distinction**: Use a specific badge or border color (e.g., `archive` icon, Pink/Accent border) to distinguish archived transactions from active ones.
