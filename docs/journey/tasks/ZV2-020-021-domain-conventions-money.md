# ZV2-020/021 — Domain conventions and Money

**Status:** Historical—reviewable
**Roadmap:** ZV2-020, ZV2-021
**Owner:** Product owner + Codex
**Completed:** 2026-08-19

## Objective

Create a dependency-free domain foundation for offline-safe UUIDs, canonical
UTC timestamps, and integer-cent money.

## Delivered

- UUID version 4 generation and canonical UUID validation.
- Canonical UTC RFC 3339 timestamp formatting/parsing.
- Non-negative `Money` values with explicit three-letter uppercase currency
  codes, safe addition, comparison, mismatch checks, and overflow protection.
- Unit tests for each invariant.

## Acceptance criteria and evidence

- [x] Domain code imports no SQLite, Wails, HTTP, or device SDK package.
- [x] New identifiers are UUID v4 values and malformed IDs are rejected.
- [x] Persisted timestamps normalize to UTC and noncanonical inputs are
  rejected.
- [x] Money uses integer minor units only; negative values, currency mismatch,
  and addition overflow are rejected.

### Changed files

- `go/core/domain/id.go` and `id_test.go`
- `go/core/domain/time.go` and `time_test.go`
- `go/core/domain/money.go` and `money_test.go`
- `go/core/domain/doc.go`

### Automated verification

```text
cd go/core && go test ./... && go vet ./...
result: passed
```

### Known limitations / follow-ups

- The domain currently validates currency-code shape, not a business-approved
  currency allowlist.
- Entity-specific ID types and all membership rules belong to later domain
  tasks.

## Human acceptance

### Manual check

1. Run `cd go/core && go test ./domain -v`.
2. Confirm the UUID, timestamp, and Money test groups pass.
3. Review the public API in [money.go](../../../go/core/domain/money.go),
   [id.go](../../../go/core/domain/id.go), and [time.go](../../../go/core/domain/time.go).

**Expected result:** You can see that cents—not floats—are the only money
representation, UTC is the persistence boundary, and the domain package has no
infrastructure dependency.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
