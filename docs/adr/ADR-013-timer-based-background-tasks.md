# ADR-013: Timer-Based Background Tasks

## Status
Accepted

## Context and Problem Statement

Several jobs must run periodically and independently of HTTP requests: bank imports, CNB currency
rate fetching, duplicate detection, and database backups. They must shut down cleanly with the
server and, in some cases, respond to on-demand triggers.

## Decision Drivers

- Periodic execution without an external scheduler (cron, systemd timers)
- Clean shutdown tied to the server lifecycle
- Some tasks need an immediate, on-demand trigger (e.g. import on importer create/update)
- Avoid busy loops that burn CPU

## Considered Options

- **External scheduler** (cron / systemd) invoking CLI commands
- **In-process goroutines with `time.Timer`/`time.Ticker`** driven by the server context
- **A job-queue library**

## Decision Outcome

Chosen: **in-process goroutines** started from `Server()`, each owning a timer and selecting on the
server `context.Done()` plus task-specific triggers. The bank importer uses a `time.Timer` and also
selects on a forced-import channel so importer create/update triggers an immediate run. Each task
returns a finish channel the server waits on during shutdown. Background mutations run with a
context carrying `ChangeSourceSystem` so audit logs attribute them to the system. Individual tasks
can be disabled via config (`DisableImporters`, `DisableCurrenciesRatesFetch`). Backups use a direct
`database/sql` connection for `VACUUM INTO` rather than the GORM handle.

Intervals: bank import (timer-driven, ~hourly), CNB rates (24h, 1h retry on failure), duplicate
detection (24h after an initial delay), backup (24h, 30s after start).

### Pros

- No external scheduler to configure; self-contained binary
- Lifecycle is tied to the server — clean startup and shutdown via context + finish channels
- Forced-import channel gives immediate response to user actions
- Timer pattern (not `select { default: }`) avoids busy loops; system change-source keeps audit accurate

### Cons

- Tasks run in-process — they scale and fail with the single server instance
- No persistence of schedule state across restarts (e.g. detection waits its initial delay again)
- A panic in a task goroutine must be contained so it doesn't take down the server
- Multiple server instances would each run the tasks (not designed for horizontal scaling)
