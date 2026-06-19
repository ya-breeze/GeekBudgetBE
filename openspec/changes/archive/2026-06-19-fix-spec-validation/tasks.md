## 1. Fix the two specs

- [x] 1.1 `chart-palette/spec.md` — rename `## ADDED Requirements` → `## Requirements`
- [x] 1.2 `landing-page/spec.md` — add title + `## Purpose` + `## Requirements` above the requirements

## 2. Verify

- [x] 2.1 `openspec validate --specs` passes for all 16 capabilities (0 failed)
- [x] 2.2 Confirm requirement/scenario counts in the 2 files are unchanged vs `main`

## 3. Finalize

- [ ] 3.1 Get user approval, then archive with `openspec archive fix-spec-validation --skip-specs`
- [ ] 3.2 Squash-merge to `main`
