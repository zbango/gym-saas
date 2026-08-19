# ZV2-023 — Member domain model

**Status:** Awaiting human acceptance
**Roadmap:** ZV2-023
**Owner:** Product owner + Codex
**Opened:** 2026-08-19

## Objective

Define a dependency-free Member entity that owns identity/profile input and
explicit member status without device or access shortcuts.

## Scope

### In scope

- Typed `MemberID`, valid Gym ownership, first/last name, phone, optional
  email, date of birth, address, identification number, and explicit member
  status.
- Validation/normalization that matches the first-slice schema without
  inventing phone or identification formats not yet specified by the product.
- UTC lifecycle timestamps, immutable construction, and unit tests.

### Out of scope

- SQLite uniqueness checks, repositories, Wails/React, and search queries.
- Hikvision/device external IDs and the legacy `has_access` shortcut.
- Membership validity or access decisions.
- Archive/deletion workflows, which are a later use-case task.

## Plan

1. Define Member ID/status types and constructors in `go/core/domain`.
2. Validate required identity fields, canonical UUID ownership, optional email
   syntax, canonical date-of-birth input, status, and timestamps while
   trimming presentation input.
3. Add unit tests for every accepted and rejected invariant.
4. Run core tests and vet, update this record, and hand it to the owner for
   manual acceptance.

## Acceptance criteria

- [x] A Member has valid Member/Gym UUIDs, nonblank first/last names and phone,
  explicit `active`, `inactive`, or `blocked` status, and UTC timestamps.
- [x] Optional email is either absent or syntactically valid; date of birth is
  either absent or canonical `YYYY-MM-DD`; address and identity data are
  trimmed but not given invented format rules.
- [x] Invalid IDs, blank required fields, unsupported status, invalid email,
  zero timestamps, and reversed timestamps are rejected by tests.
- [x] No Member API contains external device IDs or an access boolean.
- [x] `go test ./... && go vet ./...` passes in `go/core`.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| Phone and identification-number formats vary by gym and country. | Require nonblank phone and preserve trimmed values; introduce a stricter policy only after a product decision. |
| Database uniqueness is tenant-scoped. | Domain validates one Member only; SQLite repository work later enforces per-gym uniqueness. |
| Legacy access/device fields could leak into V2. | They are deliberately absent from the model and its tests. |

## Implementation record

### Changed files

- `go/core/domain/member.go`
- `go/core/domain/member_test.go`
- `docs/journey/tasks/ZV2-023-member-domain.md`
- `docs/journey/README.md`

### Simplifications made

- The entity keeps only the fields whose rules are already clear. It adds no
  device abstraction, access calculation, repository, or frontend DTO.

### Automated verification

```text
command: cd go/core && go test ./... && go vet ./...
result: passed after adding date of birth and address
```

### Known limitations / follow-ups

- Archive/deletion, notes, and persistence are separate tasks. Date of birth
  intentionally has no age/future-date policy until the product defines one.

## Human acceptance

### Manual check

1. Run `cd go/core && go test ./domain -v`.
2. Confirm the Member test groups pass and cover every acceptance criterion.
3. Review `go/core/domain/member.go` and confirm it imports no infrastructure
   package and contains neither a device external ID nor a `has_access` field.

**Expected result:** A Member can be created offline with clean identity/status,
date-of-birth, and address input, while malformed or legacy-shortcut input has
no path into the model.

### Change request resolution

The owner requested date of birth and address on 2026-08-19. Both fields are
now included: date of birth is optional canonical `YYYY-MM-DD`; address is
optional trimmed text. Their unit coverage is included in the Member test
group.

### Owner result

- [ ] Accepted
- [x] Changes requested

**Date:** 2026-08-19
**Notes:** Initial change request resolved; awaiting recheck.
