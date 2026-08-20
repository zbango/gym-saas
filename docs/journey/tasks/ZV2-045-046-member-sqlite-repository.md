# ZV2-045/046 — Member SQLite repository

**Status:** Awaiting human acceptance
**Roadmap:** ZV2-045, ZV2-046
**Owner:** Product owner + Codex
**Opened:** 2026-08-19

## Objective

Persist and reload valid Members through a narrow domain port and SQLite
adapter, proving that the V2 Member model survives database restart.

## Scope

### In scope

- A minimal Member repository interface owned by the domain/application
  boundary.
- SQLite `Create`, `Get`, and tenant-scoped `List` operations.
- Mapping every currently modeled Member field, including DOB and address.
- Temporary-database integration tests and a restart-persistence test.

### Out of scope

- Wails/React screens and application commands.
- Update/archive/delete, pagination, search, sync outbox, and transactions
  spanning multiple aggregates.
- Database uniqueness policy beyond constraints already in the V2 schema.

## Plan

1. Define only the Member repository methods required for create/get/list.
2. Implement the SQLite adapter without importing SQLite into the domain.
3. Map Member fields losslessly to the existing `members` table, including
   UTC timestamps and nullable values.
4. Test create/reload/list behavior, missing records, tenant isolation, and
   reopening the database.
5. Record automated evidence and require human acceptance before the next task.

## Acceptance criteria

- [x] A valid Member can be persisted and reloaded with every modeled field
  unchanged.
- [x] A Member belongs only to its Gym in list queries.
- [x] The repository returns a clear not-found result for an unknown Member.
- [x] A reopened SQLite database returns previously saved Members.
- [x] Domain code has no SQLite import; only the adapter owns SQL.
- [x] `go test ./... && go vet ./...` passes in the affected Go modules.

## Implementation record

### Changed files

- `go/core/ports/member_repository.go`
- `apps/desktop/internal/sqlite/member_repository.go`
- `apps/desktop/internal/sqlite/member_repository_test.go`
- This journey task page

### Simplifications made

- The port exposes only `Create`, `Get`, and tenant-scoped `List`.
- No application command, UI binding, update/archive behavior, or sync event
  was introduced.

### Automated verification

```text
cd go/core && go test ./... && go vet ./...
result: passed

cd apps/desktop && go test ./... && go vet ./...
result: passed
```

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| The `members` foreign key requires a Gym row. | Tests seed a Gym directly through the existing database setup until a Gym repository exists. |
| Domain Member currently omits schema fields such as notes and soft deletion. | This adapter persists only the approved domain model; later tasks extend the model and repository deliberately. |
| SQLite uniqueness failures need useful errors. | Preserve the database error now; classify business conflict errors only when application use cases need that distinction. |

## Human acceptance

### Manual check

1. Run the documented repository test command after implementation.
2. Confirm a test creates a Member, closes/reopens SQLite, and reloads the same
   DOB, address, status, and tenant ownership.
3. Review the adapter and verify SQL is outside `go/core/domain`.

**Expected result:** Member data is durable across restart while domain rules
remain independent from SQLite.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
