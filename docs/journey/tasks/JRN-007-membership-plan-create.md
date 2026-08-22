# JRN-007 — Membership-plan create

**Status:** Awaiting human acceptance
**Roadmap:** ZV2-024, ZV2-047, ZV2-085 (create mutation integration)
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Create a V1-style plan dialog that persists valid membership plans through the
existing Wails `CreateMembershipPlan` API and refreshes the real plan list.

## Schema alignment

| Dialog field | Current contract | Outcome |
|---|---|---|
| Nombre del plan | `name` | Persisted |
| Tipo de plan | `validityKind`, duration snapshot | Persisted through fixed type mapping |
| Precio USD | `priceCents`, `currency` | USD value converted to integer cents |
| Estado | `status` | Persisted |
| Plan por visitas | `visitLimit`, `durationValue`, `durationUnit` | Requires both visit limit and calendar window |
| Descripción | No field in domain, API, or SQLite | Omitted; must not be silently dropped |

## Scope

### In scope

- Reusable accessible dialog primitive.
- V1-style create form with fixed time-plan types and a schema-complete visit
  plan branch.
- Real Wails `CreateMembershipPlan` call followed by list refresh.
- Submission, loading, field error, and API-error states.

### Out of scope

- Description persistence, edit, delete, archive, custom plan types, and
  changes to the Go domain/schema/API.

## Plan

1. Add a shared dialog primitive with explicit close behavior.
2. Map visual plan types to the existing application input contract.
3. Call only the existing create API and refresh the existing read list after a
   successful response.
4. Verify Go adapter coverage, frontend typecheck, and dialog visual states.

## Acceptance criteria

- [x] The dialog contains only schema-backed fields and creates a plan through
  the Wails API.
- [x] Visit plans require visit limit and calendar validity.
- [x] Price is persisted in integer cents while shown in USD.
- [x] Success closes the dialog and refreshes the real list.
- [x] Go tests and frontend typecheck pass.

## Implementation record

### Changed files

- `packages/ui/src/Dialog.tsx`
- `packages/ui/src/index.tsx`
- `apps/desktop/frontend/src/features/plans/MembershipPlanCreateDialog.tsx`
- `apps/desktop/frontend/src/features/plans/MembershipPlanListPage.tsx`
- `apps/desktop/frontend/src/features/plans/membership-plan-list-page.css`
- This journey task and the visual inventory

### Simplifications made

- The existing create API and Go contracts were reused; no domain/schema/API
  change was needed.
- Fixed plan types map to approved validity tuples. Visit plans expose their
  required validity window instead of inventing an expiry.

### Automated verification

```text
cd apps/desktop && go test ./...
result: passed

cd go/core && go test ./...
result: passed

npm run typecheck
result: passed

Local browser dialog check at 1511×920 logical pixels
result: fixed plan has locked duration; visit plan shows quantity and a
four-unit validity selector; browser submit explicitly refuses persistence.
```

### Known limitations / follow-ups

- Final persistence/list-refresh confirmation requires creating a plan in the
  Wails desktop runtime, which uses the local SQLite database.
- Description cannot be included until it has an approved backend field.

## Human acceptance

### Manual check

1. Open Planes in the desktop app, then choose **Agregar Nuevo Plan**.
2. Create one fixed duration plan and one visit plan with a validity window.
3. Confirm both persist and appear in the refreshed list after closing the
   dialog.

**Expected result:** The dialog visually matches V1 while every submitted
field has a meaningful backend destination.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
