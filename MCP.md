# MCP Stdio Server for GeekBudget

## Context

Add an MCP (Model Context Protocol) stdio server so AI assistants (Claude Desktop, Claude Code) can query GeekBudget data directly. Stdio transport — no auth needed, the process runs locally with exclusive DB access.

Read-only tools only. The MCP server's `instructions` field provides rich domain context about the data model, relationships, and app logic so the AI understands what it's working with.

A `geekbudget mcp-config` CLI command generates/updates the user's `.mcp.json` for easy setup.

## New Files

```
backend/
  cmd/commands/cmdmcp.go              # "mcp" Cobra command (stdio server)
  cmd/commands/cmdmcpconfig.go        # "mcp-config" Cobra command (generate .mcp.json)
  pkg/mcpserver/
    server.go                         # MCPServer struct, Run(), tool registration
    instructions.go                   # MCP instructions text (domain context)
    helpers.go                        # JSON serialization, date parsing, result builders
    tools_accounts.go                 # Account tools
    tools_currencies.go               # Currency tools
    tools_transactions.go             # Transaction tools
    tools_matchers.go                 # Matcher tools
    tools_budget.go                   # Budget tools
    tools_reconciliation.go           # Reconciliation status tools
    tools_notifications.go            # Notification tools
    tools_analysis.go                 # Composite financial_summary tool
```

## Modified Files

- `backend/cmd/main.go` — register `CmdMCP(logger)` and `CmdMCPConfig(logger)`
- `backend/go.mod` — add `github.com/modelcontextprotocol/go-sdk`

---

## Step 1: Add MCP SDK dependency

```bash
cd backend && go get github.com/modelcontextprotocol/go-sdk@latest
```

## Step 2: Cobra command — `backend/cmd/commands/cmdmcp.go`

Follow `cmdmatch.go` pattern (file: `backend/cmd/commands/cmdmatch.go`):
- `--username` flag (required) — resolved via `storage.GetUserID(username)`
- Extract config from Cobra context — **but discard the logger** from `createConfigAndLogger` (it writes to stdout, which is reserved for MCP JSON-RPC)
- Create stderr-only logger: `slog.New(slog.NewJSONHandler(os.Stderr, ...))`
- Pass stderr logger to `database.NewStorage()` so GORM logs go to stderr too
- Call `mcpserver.Run(cmd.Context(), logger, storage, userID)`

## Step 3: MCP config generator — `backend/cmd/commands/cmdmcpconfig.go`

New `geekbudget mcp-config` command that writes/updates `.mcp.json`:
- `--username` flag (required)
- `--output` flag (default: `.mcp.json` in current directory)
- Detects the path to the `geekbudget` binary (via `os.Executable()`)
- Reads `GB_DBPATH` from config
- Generates JSON:

```json
{
  "mcpServers": {
    "geekbudget": {
      "command": "/absolute/path/to/geekbudget",
      "args": ["mcp", "--username", "test@test.com"],
      "env": {
        "GB_DBPATH": "/absolute/path/to/geekbudget.db"
      }
    }
  }
}
```

- If file exists, reads it, merges the `geekbudget` key (preserving other servers), writes back
- Prints the path written to stderr

## Step 4: Register commands in `backend/cmd/main.go`

Add to `rootCmd.AddCommand(...)` block:

```go
commands.CmdMCP(logger),
commands.CmdMCPConfig(logger),
```

## Step 5: MCP server core — `backend/pkg/mcpserver/server.go`

```go
type MCPServer struct {
    logger  *slog.Logger
    storage database.Storage
    userID  string
}
```

`Run()` function:
1. Creates `mcp.NewServer(...)` with `mcp.WithInstructions(instructions)`
2. Calls `s.registerXxxTools(server)` for each tool group
3. Runs `server.Run(ctx, &mcp.StdioTransport{})`

## Step 6: MCP instructions — `backend/pkg/mcpserver/instructions.go`

A `const instructions` string passed via `mcp.WithInstructions()`. Provides the AI with domain context:

```
GeekBudget is a personal finance management application.

## Data Model

- **Accounts** have a type: "asset" (bank accounts, cash), "expense" (categories like groceries, rent), or "income" (salary, etc.)
- **Currencies** are user-defined (e.g. CZK, EUR, USD). All amounts reference a currency ID.
- **Transactions** represent financial events. Each transaction has a date, description, optional place/tags/partner info, and a list of **Movements**.
- **Movements** are the core of double-entry bookkeeping: each movement transfers an amount in a specific currency to/from an account. A transaction typically has 2+ movements that balance out (e.g. -100 CZK from "Cash" account, +100 CZK to "Groceries" account).
- **Matchers** are regex-based rules that auto-categorize imported bank transactions. They match on description, partner name, partner account number, currency, place, or keywords. Matchers have a confirmation history tracking their accuracy.
- **Budget Items** define expected spending per account/category.
- **Bank Importers** connect to banks (FIO, Revolut, KB) to fetch transactions automatically.
- **Reconciliation** compares the app's computed balance against the bank's reported balance for asset accounts.

## Key Relationships

- Transactions → Movements → Accounts + Currencies (each movement references one account and one currency)
- Matchers → Accounts (output account for auto-categorization)
- Bank Importers → Accounts (target account for imported transactions)
- Reconciliation → Accounts + Currencies

## Conventions

- All IDs are UUIDs
- Money amounts are decimal numbers (never floating point)
- Dates are ISO 8601 (YYYY-MM-DD)
- Tags are string arrays on transactions for custom labeling
- "Suspicious" transactions have issues flagged by the system (e.g. unbalanced movements)
- "Unprocessed" transactions are imported but not yet categorized
- Duplicate detection flags transactions with similar dates (±2 days) and amounts from different import sources

## Common Queries

- To understand spending: list transactions filtered by date, look at movements going to expense accounts
- To check balances: use get_account_balance for a specific account+currency, or financial_summary for an overview
- To find categorization issues: list_transactions with onlySuspicious=true, or check matchers
- To verify bank sync: get_reconciliation_status shows delta between app and bank balances
```

## Step 7: Helpers — `backend/pkg/mcpserver/helpers.go`

Shared utilities:
- `toJSON(v any) string` — indented JSON for readable responses
- `textResult(text string)` — wrap in `*mcp.CallToolResult`
- `jsonResult(data any)` — serialize + wrap
- `errorResult(err error)` — error result via `result.SetError()`
- `parseOptionalDate(s string) (time.Time, error)` — parse `YYYY-MM-DD` or zero
- `defaultDateFrom(s string, daysBack int)` — parse or default to N days ago

## Step 8: Tool definitions (all read-only)

All tools set `ToolAnnotations{ReadOnlyHint: true}`.

| Tool | File | Args | Storage call |
|------|------|------|--------------|
| `list_accounts` | tools_accounts.go | — | `GetAccounts(userID)` |
| `get_account` | tools_accounts.go | `id` | `GetAccount(userID, id)` |
| `get_account_balance` | tools_accounts.go | `accountId`, `currencyId` | `GetAccountBalance(userID, accID, curID)` |
| `get_account_history` | tools_accounts.go | `accountId` | `GetAccountHistory(userID, accID)` |
| `list_currencies` | tools_currencies.go | — | `GetCurrencies(userID)` |
| `list_transactions` | tools_transactions.go | `dateFrom?`, `dateTo?`, `onlySuspicious?` | `GetTransactions(userID, from, to, flag)` |
| `get_transaction` | tools_transactions.go | `id` | `GetTransaction(userID, id)` |
| `search_transactions` | tools_transactions.go | `query`, `dateFrom?`, `dateTo?` | `GetTransactions` + in-memory filter on Description/PartnerName/Place/Tags |
| `get_duplicate_transactions` | tools_transactions.go | `transactionId` | `GetDuplicateTransactionIDs(userID, id)` |
| `list_matchers` | tools_matchers.go | — | `GetMatchers(userID)` |
| `get_matcher` | tools_matchers.go | `id` | `GetMatcher(userID, id)` |
| `list_budget_items` | tools_budget.go | — | `GetBudgetItems(userID)` |
| `get_reconciliation_status` | tools_reconciliation.go | — | Composite: replicates `api_reconciliation.go:24-123` logic |
| `get_reconciliation_history` | tools_reconciliation.go | `accountId` | `GetReconciliationsForAccount(userID, accID)` |
| `list_notifications` | tools_notifications.go | — | `GetNotifications(userID)` |
| `financial_summary` | tools_analysis.go | — | Composite: accounts + currencies + bulk data + notifications |

**16 read-only tools total.**

### Tool input pattern

Input args are Go structs with `json` + `jsonschema` tags. SDK auto-infers JSON Schema.

```go
type ListTransactionsArgs struct {
    DateFrom       string `json:"dateFrom,omitempty" jsonschema:"start date YYYY-MM-DD, defaults to 30 days ago"`
    DateTo         string `json:"dateTo,omitempty" jsonschema:"end date YYYY-MM-DD, defaults to today"`
    OnlySuspicious bool   `json:"onlySuspicious,omitempty" jsonschema:"if true only return suspicious transactions"`
}
```

### `list_transactions` defaults

If `dateFrom` omitted → 30 days ago. If `dateTo` omitted → zero time (storage treats as no upper bound).

### `get_reconciliation_status` logic

Replicates the reconciliation assembly from `backend/pkg/server/api/api_reconciliation.go` lines 24-123. ~60 lines of data assembly from 4 storage calls (`GetAccounts`, `GetBankImporters`, `GetCurrencies`, `GetBulkReconciliationData`). Duplicated rather than imported to avoid coupling to HTTP layer's `context.Value`/`ImplResponse` patterns.

### `financial_summary` — the "what should I know" tool

Returns a single response combining:
- All accounts (name, type, id)
- All currencies (name, id)
- Per-account/currency balances from `GetBulkReconciliationData`
- Unprocessed transaction counts
- Active notifications
- Reconciliation flags

## userID flow

```
CLI: --username "test@test.com"
  → storage.GetUserID(username) → userID string
    → MCPServer{storage, userID, logger}
      → each tool: s.storage.GetAccounts(s.userID)
```

## Logging

- `createConfigAndLogger` (in `cmdserver.go`) writes to stdout — **cannot use for MCP**
- Create dedicated: `slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))`
- Pass to both `database.NewStorage()` and `MCPServer`

## Critical files to reference

| File | Purpose |
|------|---------|
| `backend/cmd/commands/cmdmatch.go` | Pattern for new Cobra command (username flag, storage init) |
| `backend/cmd/commands/cmdserver.go` | `createConfigAndLogger` helper to reuse for config extraction |
| `backend/cmd/main.go` | Command registration |
| `backend/pkg/database/storage.go:40-131` | Storage interface — full API surface for tools |
| `backend/pkg/server/api/api_reconciliation.go:24-123` | Reconciliation status logic to replicate |
| `backend/pkg/database/models/transaction.go` | Transaction model, Movement structure |
| `backend/pkg/database/models/account.go` | Account model, BankAccountInfo |
| `backend/pkg/database/models/matcher.go` | Matcher model, confirmation logic |

## Implementation order

1. `go.mod` — add dependency
2. `cmdmcp.go` + `cmdmcpconfig.go` + `main.go` registration
3. `server.go` + `instructions.go` + `helpers.go` — core structure
4. `tools_analysis.go` — `financial_summary` (most useful, test first)
5. `tools_accounts.go` + `tools_currencies.go`
6. `tools_transactions.go` (list, get, search, duplicates)
7. `tools_matchers.go` + `tools_budget.go` + `tools_notifications.go` + `tools_reconciliation.go`
8. `make all` — verify build + tests + lint pass

## Verification

1. Build: `cd backend && go build -o bin/geekbudget ./cmd`
2. Generate config: `./bin/geekbudget mcp-config --username test@test.com`
3. Manual protocol test:
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | \
     GB_DBPATH=./geekbudget.db ./bin/geekbudget mcp --username test@test.com
   ```
4. Place generated `.mcp.json` in project root, restart Claude Code, verify tools appear
5. Test in Claude: "give me a financial summary", "list my accounts", "show transactions from last week"
6. `make all` must pass
