# JRN-005 — Membership-plan list

**Status:** Awaiting human acceptance
**Roadmap:** ZV2-024, ZV2-047, ZV2-085 (read-only UI integration)
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Render the supplied membership-plan card screen using the real desktop
read-only path: Wails `MembershipPlanAPI` → application service → SQLite
repository.

## Scope

### In scope

- Membership-plan page header and responsive plan-card grid matching V1.
- Desktop listing through `ListMembershipPlans` only, including loading, empty,
  and error states.
- Explicit browser-preview samples only when Wails is unavailable, isolated
  from the desktop runtime.

### Out of scope

- Create, update, archive/delete, status changes, or data seeding.
- Changing membership-plan domain, service, Wails, or SQLite contracts.
- Treating browser preview samples as backend data.

## Plan

1. Extract a read-only plan-list hook that calls the existing Wails API in the
   desktop runtime.
2. Compose V1 cards from the API view model and shared UI/icons.
3. Route the Membership Plans navigation item to the new list page.
4. Verify API/service/repository tests, frontend typecheck, and visual states.

## Acceptance criteria

- [x] Desktop listing calls only `ListMembershipPlans` and renders its result.
- [x] Cards display validity type, active/inactive status, price, duration, and
  read-only edit/delete affordances matching V1.
- [x] Loading, empty, and error states are visible and do not fabricate data.
- [x] Existing Go plan-list coverage and frontend typecheck pass.

## Implementation record

### Changed files

- `apps/desktop/frontend/src/features/plans/useMembershipPlanList.ts`
- `apps/desktop/frontend/src/features/plans/MembershipPlanListPage.tsx`
- `apps/desktop/frontend/src/features/plans/membership-plan-list-page.css`
- `apps/desktop/frontend/src/features/dashboard/DashboardPage.tsx`
- `apps/desktop/app_membership_plan_test.go`
- This journey task and the visual inventory

### Simplifications made

- The existing Go domain/service/repository/Wails API was reused unchanged.
- Browser preview samples are selected only outside Wails; a desktop runtime
  renders the API response or a true empty/error state.

### Automated verification

```text
cd apps/desktop && go test ./...
result: passed; adapter regression confirms ListMembershipPlans returns the
persisted plan's ID, name, price, and duration unit before archive.

cd go/core && go test ./...
result: passed

npm run typecheck
result: passed

Local browser visual check at 1511×920 logical pixels
result: six preview cards, V1 card hierarchy, and no loading/error state after
the browser-only preview resolves.
```

### Known limitations / follow-ups

- A fresh desktop database correctly shows the empty state until plans are
  actually persisted through a separately scoped mutation flow.
- Add/edit/delete controls are visual-only in this screen.

## Human acceptance

### Manual check

1. Open Planes from the Membresías sidebar group in the desktop app.
2. Confirm cards equal the actual plans returned by the local SQLite-backed
   desktop API.
3. Compare the card hierarchy to the supplied V1 reference.

**Expected result:** Wails-backed plans appear as V1-style cards; add, edit,
and delete controls remain visual-only.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
