# ADR-002: OpenAPI-First Development with Code Generation

## Status
Accepted

## Context and Problem Statement

The backend (Go) and frontends (Angular, Next.js) must agree on the same API contract. Hand-writing
request/response types on both sides drifts and causes integration bugs, especially around money
fields and strict update payloads.

## Decision Drivers

- Single source of truth for the API contract
- Type-safe clients and server stubs generated, not hand-maintained
- Consistent models across Go backend and TypeScript frontends

## Considered Options

- **OpenAPI spec-first + code generation** — write `api/openapi.yaml`, generate both sides
- **Code-first** — generate the spec from Go annotations
- **Hand-written types on each side** — no generation

## Decision Outcome

Chosen: **OpenAPI 3.0 spec-first** (`api/openapi.yaml`) with generated Go server interfaces
(`backend/pkg/generated/goserver`) and a generated Angular client. Workflow: edit the spec first,
then run `make generate`. Generated code under `backend/pkg/generated/` is never hand-edited.

### Pros

- The contract is explicit and reviewable in one file
- Backend and frontend models stay in sync
- Strict typing surfaces breaking changes at generation time
- Update payloads use generated no-id types, enforcing "no `id` in body" strictness

### Cons

- Changing the API is a two-step dance (edit spec, regenerate) rather than editing code directly
- Generator quirks occasionally need workarounds
- Contributors must learn the regeneration workflow and not touch generated files
