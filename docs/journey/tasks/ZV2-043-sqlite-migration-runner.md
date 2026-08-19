# ZV2-043 — SQLite migration runner

**Status:** Historical—reviewable  
**Roadmap:** ZV2-043  
**Owner:** Product owner + Codex  
**Completed:** 2026-08-19

## Objective

Replace the desktop proof-table bootstrap with a safe, embedded V2 schema
migration runner.

## Delivered

- Embedded V2 migration SQL in the desktop binary.
- Transactional schema application and migration-record writes.
- Migration name/checksum protection and refusal of newer unknown schemas.
- SQLite foreign-key enforcement, WAL, busy timeout, and single-connection
  operation.
- Narrow compatibility for the earlier proof database containing only
  `hello_records`; V1 and other unmanaged databases remain refused.

## Acceptance criteria and evidence

- [x] A fresh desktop database receives all V2 first-slice tables exactly once.
- [x] A failing migration rolls back both schema changes and its migration
  record.
- [x] Foreign keys are enforced and WAL/busy-timeout settings are verified.
- [x] Applied migration drift blocks startup.

### Changed files

- `apps/desktop/internal/sqlite/store.go`
- `apps/desktop/internal/sqlite/store_test.go`
- `apps/desktop/internal/sqlite/migrations/001_v2_operational_core.sql`
- Desktop proof bindings/UI and related documentation

### Automated verification

```text
cd apps/desktop && go test ./... && go vet ./...
result: passed
```

### Known limitations / follow-ups

- V1-to-V2 business-data migration remains a separate R4 task.
- Repository and domain use cases are not implemented by this task.

## Human acceptance

### Manual check

1. Run `cd apps/desktop && go test ./internal/sqlite -v`.
2. Confirm all named migration tests pass, including fresh schema, rollback,
   foreign-key enforcement, unmanaged database rejection, legacy proof bridge,
   and migration drift.
3. Review [the embedded migration SQL](../../../apps/desktop/internal/sqlite/migrations/001_v2_operational_core.sql)
   and confirm it matches the intended first operational slice.

**Expected result:** All tests pass and the schema contains only the documented
V2 operational tables plus `schema_migrations`.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**  
**Notes:** 
