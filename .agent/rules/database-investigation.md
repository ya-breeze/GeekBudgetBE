---
trigger: always_on
description: Guidelines for investigating data issues using the local development database
---

# Local Database Investigation

When investigating complex bugs, data inconsistencies, or verifying the state of the application, you are encouraged to check the local SQLite database.

## Database Location
The local database is located at the root of the project: `geekbudget.db`.

## Querying the Database
Use the `sqlite3` command-line tool to run queries. Always use `.header on` and `.mode column` (or `.mode markdown`) for readable output.

## Guidelines
1. **Read-Only First**: Prioritize SELECT queries for investigation. 
2. **NO Mutations**: NEVER perform INSERT/UPDATE/DELETE operations.
3. **Data Privacy**: the local DB contains sensitive personal data.
4. **UserID Isolation**: Remember that most tables have a `user_id`. When querying for a specific user, include `WHERE user_id = ...`.
5. **JSON Fields**: Some fields like `suspicious_reasons` are stored as JSON strings. Use `json_extract` or similar SQLite functions if needed, or just pipe to `jq` if you are comfortable.