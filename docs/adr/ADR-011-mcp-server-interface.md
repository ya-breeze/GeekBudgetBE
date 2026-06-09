# ADR-011: MCP Server as a First-Class Interface

## Status
Accepted

## Context and Problem Statement

Beyond the REST API and web UI, there is value in letting AI assistants (Claude Desktop, Claude
Code) query a family's finances directly — "how much did I spend on groceries last month?" — without
building a bespoke integration each time.

## Decision Drivers

- Enable AI assistants to read GeekBudget data natively
- Avoid exposing a network surface or new auth path for this
- Give the assistant enough domain context to use the data correctly

## Considered Options

- **No AI interface** — assistants would scrape the REST API with custom glue
- **MCP over a network transport** — needs auth, TLS, hosting
- **MCP over stdio, local process** — no network, no auth, exclusive local DB access

## Decision Outcome

Chosen: an **MCP stdio server** (`pkg/mcpserver`) run as a local process bound to a single
`familyID`, registering read-oriented tool groups (accounts, currencies, transactions, matchers,
budget, reconciliation, notifications, analysis). A rich `instructions` block documents the data
model and relationships so the assistant understands accounts/movements/double-entry. A
`mcp-config` CLI command generates the user's `.mcp.json`.

### Pros

- Assistants get structured, domain-aware access without bespoke integrations
- Stdio transport means no network exposure and no separate auth to secure
- The instructions block centralizes domain knowledge for the AI

### Cons

- Local-only: the process needs direct DB access on the same host
- Family id is fixed at launch; multi-family use means multiple server instances
- Tool surface must be maintained in parallel with the REST API as features evolve
- Read-oriented by design — not a path for mutations
