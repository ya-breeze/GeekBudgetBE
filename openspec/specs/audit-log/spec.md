# audit-log Specification

## Purpose

Every create, update, and delete in the storage layer records an audit entry capturing what changed,
by whom (user vs. system), and the before/after state, scoped to a family.

## Requirements

### Requirement: Audit entry on mutations

Create, update, and delete operations SHALL record an `AuditLog` with `EntityType`, `EntityID`,
`Action`, `ChangeSource`, and JSON `Before`/`After` snapshots.

#### Scenario: Create records audit
- **WHEN** an entity is created
- **THEN** an audit entry is recorded with `before` empty and `after` set to the new entity

#### Scenario: Update records before and after
- **WHEN** an entity is updated
- **THEN** an audit entry is recorded with both the prior and new state

#### Scenario: Delete records before
- **WHEN** an entity is deleted
- **THEN** an audit entry is recorded with `before` set and `after` empty

### Requirement: Change source attribution

Each audit entry SHALL record whether the change came from a user action or a system/background
process, derived from the request context and defaulting to `system`.

#### Scenario: User-initiated change
- **GIVEN** a request carrying the user change source in context
- **WHEN** a mutation occurs
- **THEN** the audit entry's change source is `user`

#### Scenario: Background change
- **GIVEN** a background task without a user change source
- **WHEN** a mutation occurs
- **THEN** the audit entry's change source is `system`

### Requirement: Query audit logs

Audit logs SHALL be queryable with filters for entity type, entity id, date range, and pagination
(limit/offset).

#### Scenario: Filter by entity
- **WHEN** audit logs are requested filtered by entity type and id
- **THEN** only matching entries are returned, paginated
