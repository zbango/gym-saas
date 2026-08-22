# JRN-006 — Dashboard route outlet

**Status:** Awaiting human acceptance
**Roadmap:** Supporting visual-parity work
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Replace the dashboard's growing conditional render chain with a centralized,
extensible route-view registry while preserving all current visual routes.

## Scope

### In scope

- Dashboard home, Members, and Plans route rendering.
- Fallback mock route behavior for routes not implemented yet.
- Moving dashboard-home rendering out of the shell component where that makes
  the route boundary explicit.

### Out of scope

- A third-party router dependency, URL history, deep linking, or authorization
  behavior.
- Changing sidebar labels, route IDs, or any page's business/data behavior.

## Cleanup plan

1. Lock current route components with frontend typecheck and route-level visual
   checks.
2. Remove the nested conditional routing smell.
3. Extract route composition into a route-view registry and outlet.
4. Verify dashboard, Members, Plans, and fallback presentation still render.

## Acceptance criteria

- [x] DashboardPage contains no route-specific conditional render chain.
- [x] Adding a route requires a registry entry rather than editing a ternary.
- [x] Existing dashboard, Members, Plans, and fallback screens remain reachable.
- [x] Typecheck passes.

## Implementation record

### Changed files

- `apps/desktop/frontend/src/features/dashboard/DashboardPage.tsx`
- `apps/desktop/frontend/src/features/dashboard/DashboardHomePage.tsx`
- `apps/desktop/frontend/src/features/dashboard/DashboardRouteOutlet.tsx`
- This journey task page

### Simplifications made

- DashboardPage now owns only shell, navigation state, and chrome.
- The route outlet owns a single `routeViews` registry plus the fallback view.
- Dashboard-home content is a standalone page, preventing a circular import
  between the route registry and shell.

### Automated verification

```text
npm run typecheck
result: passed

npm run lint:ts
result: passed with pre-existing warnings outside this refactor scope

Local browser route check at 1511×920 logical pixels
result: dashboard, Members, Plans, and unimplemented Attendance fallback routes
all rendered through the route outlet.
```

### Known limitations / follow-ups

- This is an in-memory visual router only. URL history, deep links, and
  authorization-aware route guards are deferred until real application routing
  is scoped.

## Human acceptance

### Manual check

1. From the dashboard shell, open Panel de Control, Clientes, Planes, and an
   unimplemented sidebar item.
2. Confirm each route renders its current expected content.

**Expected result:** The shell remains unchanged while page selection is driven
by an extensible route registry.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
