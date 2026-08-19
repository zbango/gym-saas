# ZV2-022 — Gym and Branch identity model

**Status:** Awaiting human acceptance  
**Roadmap:** ZV2-022  
**Owner:** Product owner + Codex  
**Opened:** 2026-08-19

## Objective

Define the dependency-free Gym and Branch domain entities that establish
tenant ownership without introducing SaaS onboarding or persistence changes.

## Scope

### In scope

- Typed `GymID` and `BranchID` backed by canonical UUID strings.
- Gym and Branch constructors that validate IDs, names, IANA timezones, and
  created/updated timestamps.
- Offline-safe creation helpers and unit tests for every invariant.

### Out of scope

- SQLite tables, repositories, Wails methods, or React screens.
- Branch assignment on members, memberships, payments, or visits.
- SaaS onboarding, roles, and branch-specific operational rules.

## Plan

1. Add typed IDs and constructors in `go/core/domain` using the existing UUID
   and UTC conventions.
2. Validate names after trimming, validate timezones with Go's timezone
   database, and reject invalid timestamp order.
3. Add focused unit tests for valid creation and each rejected invariant.
4. Run core tests and vet, then record evidence and hand the task to the owner
   for manual acceptance.

## Acceptance criteria

- [x] A Gym has a valid UUID ID, nonblank display name, valid IANA timezone,
  and UTC created/updated timestamps.
- [x] A Branch has its own valid UUID, a valid parent Gym ID, nonblank name,
  valid IANA timezone, and UTC created/updated timestamps.
- [x] Constructor-created IDs are valid UUID v4 values.
- [x] Empty/whitespace names, malformed IDs, invalid timezones, zero
  timestamps, and `updated_at` before `created_at` are rejected by tests.
- [x] `go test ./... && go vet ./...` passes in `go/core`.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| The current first-slice SQLite schema has `gyms` but no `branches` table. | This task defines domain identity only. Adding persistence or assigning branches to operational rows remains a later, explicit schema task. |
| Gym names have no documented maximum length. | Reject only blank names; avoid inventing an unsynchronized storage rule. |
| Stored timestamps must be UTC. | Constructors normalize accepted timestamps to UTC but reject zero values and invalid ordering. |

## Implementation record

### Changed files

- `go/core/domain/gym.go`
- `go/core/domain/gym_test.go`
- `docs/journey/tasks/ZV2-022-gym-branch-identity.md`
- `docs/journey/README.md`

### Simplifications made

- The model is immutable and uses constructors/getters only; no generic base
  entity, repository, persistence abstraction, or new dependency was added.
- Branch is intentionally domain-only until an approved schema task adds
  durable branch storage and operational associations.

### Automated verification

```text
command: cd go/core && go test ./... && go vet ./...
result: passed
```

### Known limitations / follow-ups

- Branch persistence and operational association are intentionally deferred.

## Human acceptance

### Manual check

1. Run `cd go/core && go test ./domain -v`.
2. Confirm the Gym and Branch test groups pass and their names describe every
   accepted and rejected invariant.
3. Review `go/core/domain/gym.go` and confirm that the entities have no import
   from SQLite, Wails, HTTP, or device packages.

**Expected result:** The output demonstrates that tenant identity can be
created offline and invalid identity/timezone/timestamp input is refused before
persistence exists.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**  
**Notes:** 
