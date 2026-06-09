# mcp-server Specification

## Purpose

GeekBudget exposes a Model Context Protocol (MCP) stdio server so AI assistants (Claude Desktop,
Claude Code) can query a family's financial data directly. It runs as a local process with exclusive
DB access; the stdio transport requires no network auth.

## Requirements

### Requirement: Stdio server with domain instructions

The MCP server SHALL run over stdio and provide an `instructions` block describing the data model
(accounts, currencies, transactions, movements, matchers, budgets, bank importers, reconciliation).

#### Scenario: Server starts and advertises instructions
- **WHEN** the MCP server is run
- **THEN** it serves over stdio and exposes the GeekBudget domain instructions

### Requirement: Family-scoped, read-only tools

The server SHALL be bound to a single family id and register read-oriented tool groups: accounts,
currencies, transactions, matchers, budget, reconciliation, notifications, and analysis.

#### Scenario: Tools operate within one family
- **GIVEN** the MCP server started for a family
- **WHEN** a tool queries data
- **THEN** only that family's data is accessible

#### Scenario: Registered tool groups
- **WHEN** the server initializes
- **THEN** account, currency, transaction, matcher, budget, reconciliation, notification, and analysis tools are registered

### Requirement: Configuration command

A CLI command SHALL generate/update the user's `.mcp.json` to register the GeekBudget MCP server.

#### Scenario: Generate MCP config
- **WHEN** the mcp-config command runs
- **THEN** an `.mcp.json` entry for the GeekBudget MCP server is produced
