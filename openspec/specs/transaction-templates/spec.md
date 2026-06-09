# transaction-templates Specification

## Purpose

Templates store reusable transaction patterns (recurring payments, common splits) so a user can
quickly create similar transactions. A template mirrors transaction fields plus a `Name`.

## Requirements

### Requirement: Template structure

A template SHALL have a `Name`, optional descriptive fields (`Description`, `Place`, `Tags`,
`PartnerName`, `Extra`), and a list of `Movements`, scoped to a family with a UUID id.

#### Scenario: Create a template
- **WHEN** a user creates a template with a name and movements
- **THEN** the template is created with a generated UUID scoped to the family

### Requirement: Query templates

Templates SHALL be listable, optionally filtered by an account referenced in their movements.

#### Scenario: List templates by account
- **WHEN** templates are requested filtered by an account id
- **THEN** only templates whose movements reference that account are returned

### Requirement: Update and delete

Templates SHALL be updatable and deletable by id.

#### Scenario: Delete a template
- **WHEN** a template is deleted by id
- **THEN** it is removed and no longer returned in listings
