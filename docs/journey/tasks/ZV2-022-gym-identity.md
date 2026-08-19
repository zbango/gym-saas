# ZV2-022 — Gym identity model

**Status:** Accepted
**Roadmap:** ZV2-022 (Branch deferred)
**Owner:** Product owner + Codex
**Opened:** 2026-08-19

## Objective

Define the dependency-free Gym entity that establishes tenant ownership for
the current single-location product.

## Scope

### In scope

- A typed `GymID` backed by canonical UUID strings.
- Gym construction that validates the ID, name, IANA timezone, and lifecycle
  timestamps.
- Offline-safe Gym creation and unit tests for every invariant.

### Out of scope

- Branches or multi-location behavior.
- SQLite repositories, Wails methods, and React screens.
- SaaS onboarding, roles, and branch-specific operational rules.

## Plan

1. Add a typed Gym ID and constructor in `go/core/domain` using the existing
   UUID and UTC conventions.
2. Validate names after trimming, validate timezones with Go's timezone
   database, and reject invalid timestamp order.
3. Add focused unit tests for valid creation and every rejected invariant.
4. Run core tests and vet, then record evidence and hand the task to the owner
   for manual acceptance.

## Acceptance criteria

- [x] A Gym has a valid UUID ID, nonblank display name, valid IANA timezone,
  and UTC created/updated timestamps.
- [x] Constructor-created IDs are valid UUID v4 values.
- [x] Empty/whitespace names, malformed IDs, invalid timezones, zero
  timestamps, and `updated_at` before `created_at` are rejected by tests.
- [x] `go test ./... && go vet ./...` passes in `go/core`.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| The roadmap includes a Branch model to support multi-location gyms. | Explicitly deferred: the current product is single-location, so no speculative Branch entity or persistence is retained. Revisit only when a real second location requires it. |
| Gym names have no documented maximum length. | Reject only blank names; avoid inventing an unsynchronized storage rule. |
| Stored timestamps must be UTC. | Constructors normalize accepted timestamps to UTC but reject zero values and invalid ordering. |

## Implementation record

### Changed files

- `go/core/domain/gym.go`
- `go/core/domain/gym_test.go`
- `docs/journey/tasks/ZV2-022-gym-identity.md`
- `docs/journey/README.md`

### Simplifications made

- The model is immutable and uses constructors/getters only; no generic base
  entity, repository, persistence abstraction, multi-location model, or new
  dependency was added.

### Automated verification

```text
command: cd go/core && go test ./... && go vet ./...
result: passed
```

### Known limitations / follow-ups

- Add a Branch model only after a real multi-location requirement is accepted.

## Human acceptance

### Manual check

1. Run `cd go/core && go test ./domain -v`.
2. Confirm the Gym test groups pass and their names describe every accepted and
   rejected invariant.
3. Review `go/core/domain/gym.go` and confirm it has no import from SQLite,
   Wails, HTTP, device packages, or a Branch model.

**Expected result:** The output demonstrates that a valid single-location Gym
can be created offline and invalid identity/timezone/timestamp input is
refused before persistence exists.

### Owner result

- [x] Accepted
- [ ] Changes requested

**Date:** 2026-08-19
**Notes:** Owner accepted the Gym-only scope in chat.
