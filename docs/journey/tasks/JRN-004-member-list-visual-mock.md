# JRN-004 — Member list (V1: Clientes) visual mock

**Status:** Awaiting human acceptance
**Roadmap:** Supporting visual-parity work; member query/use-case integration remains separate
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Recreate the supplied Member list (V1: Clientes) and extract the page-header, metric-card,
and generic data-table presentation primitives needed by future V1 screens.

## Scope

### In scope

- Shared `PageHeader`, `MetricCard`, and generic `DataTable` components in
  `@gym-saas/ui`.
- Members summary cards, search and status controls, table, row action icon,
  membership/expiry display, and pagination shell.
- Static synthetic member rows and presentation-only state for selected page,
  sort indicator, filter selection, and action buttons.
- The Clientes item in the dashboard mock opens this screen.

### Out of scope

- Real member data, filtering, sorting, pagination, row menus, create/edit
  flows, persistence, Wails bindings, or domain/application rules.
- Reusing the real personal data visible in the supplied screenshot.

## Plan

1. Add shared primitives with typed, controlled presentation props.
2. Compose the Member screen from static synthetic data and shared V1 icons.
3. Connect the existing mock navigation entry to the screen.
4. Render at desktop reference scale and verify table geometry and states.

## Acceptance criteria

- [x] The Member list header, metrics, controls, table, and pagination visually
  match the supplied references.
- [x] Shared page-header, metric-card, and generic table components have no
  Member-specific business logic.
- [x] All row data is synthetic and actions remain mock-only.
- [x] Typecheck and visual verification pass.

## Implementation record

### Changed files

- `packages/ui/src/PageHeader.tsx`
- `packages/ui/src/MetricCard.tsx`
- `packages/ui/src/DataTable.tsx`
- `packages/ui/src/icons/MoreVerticalIcon.tsx`
- `packages/ui/src/icons/index.ts`
- `packages/ui/src/index.tsx`
- `apps/desktop/frontend/src/features/members/MemberListPage.tsx`
- `apps/desktop/frontend/src/features/members/member-list-page.css`
- `apps/desktop/frontend/src/features/dashboard/DashboardPage.tsx`
- This journey task and the visual inventory

### Simplifications made

- `DataTable` owns only typed rows and column rendering; it does not know a
  member, query, sort rule, or pagination policy.
- Search/filter/page controls keep visual state only; rows remain synthetic
  until the backend supplies authoritative query results.

### Automated verification

```text
npm run typecheck
result: passed

Local browser visual check at 1511×920 logical pixels
result: five metric cards, ten table rows, shared page header, table, and
pagination visible at the reference layout density.

Local mock control check
result: search state, selected status, and page indicator changed as expected.
```

### Known limitations / follow-ups

- Query controls deliberately do not filter, sort, or paginate rows yet.
- Row action, add-client, and refresh controls are visual only pending their
  V1 flow references and backend contracts.

## Human acceptance

### Manual check

1. Open the V1 “Clientes” screen from the dashboard sidebar.
2. Compare the header, metrics, controls, table rows, and pagination shell to
   the three supplied V1 references.
3. Confirm that member data is clearly mock data and actions do not claim a
   real mutation.

**Expected result:** Clientes visually matches V1 and establishes reusable
data-display primitives without implementing member-management behavior.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
