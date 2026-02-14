# Frontend Development Patterns

## API Client Generation

The frontend uses `ng-openapi-gen` to generate the Angular API client from the OpenAPI specification.

### Workflow
1. **Modify API**: Update `api/openapi.yaml`.
2. **Generate Backend**: Run `make generate` in the root directory.
3. **Generate Frontend**:
   - Navigate to the frontend directory: `cd frontend`.
   - Run the generation script: `npm run generate-api`.
   - This script uses `ng-openapi-gen.json` to configure the input (`../api/openapi.yaml`) and output (`src/app/core/api`).

### Important Notes
- **Do NOT** rely on the `make generate` command to update the frontend code in `src/app/core/api`. The `Makefile` generates into `backend/pkg/generated/angular`, which is **NOT** the source of truth for the active frontend application.
- **Verification**: Always verify that the models in `src/app/core/api/models` reflect your changes in `openapi.yaml`.

## State Changes & Audit Logs
When displaying audit logs or history, use the `before` and `after` fields to show a diff.
- **Null Checks**: Both `before` and `after` can be null (e.g. for creation or deletion).
- **JSON Parsing**: These fields are stored as JSON strings. You must parse them safely (e.g., using `JSON.parse` inside a try-catch block) before displaying.
