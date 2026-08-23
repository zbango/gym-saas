# JRN-009 — Tailwind style migration

**Status:** Awaiting human acceptance
**Roadmap:** Supporting frontend delivery
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Make Tailwind the styling implementation for the current desktop/web screens
and shared UI primitives.

## Scope

### In scope

- Current desktop feature styles: login, dashboard, member list, and
  membership-plan list/create dialog.
- Shared UI primitives and current web-host inline styles.
- Deleting superseded feature stylesheet files/imports.
- The Tailwind entry stylesheets as the only remaining stylesheet per host.

### Out of scope

- Changing UI behavior, Go/Wails contracts, app data, or route structure.

## Cleanup plan

1. Lock desktop/web builds and visual baseline.
2. Migrate one screen group at a time, delete its stylesheet, and typecheck.
3. Migrate shared primitives/host inline styles to remove the second styling
   API.
4. Rebuild, lint, and visually compare Login, Dashboard, Members, Plans, and
   the plan-create dialog.

## Acceptance criteria

- [x] Current screen and shared UI presentation use Tailwind class names, not
  feature CSS files or inline style objects.
- [x] Static Zeus styling remains available to utility classes.
- [x] Existing desktop/web builds, typecheck, and visual behavior pass.
- [x] Only Tailwind entry stylesheets remain for the two hosts.

## Implementation record

### Changed files

- Desktop feature components, dashboard, shared UI primitives, and web host
  entry styling.
- Deleted four feature stylesheet files and removed their imports.
- `packages/ui/src/DataTable.tsx` utility-slot API, removing its final inline
  presentation style object.

### Simplifications made

- Tailwind utility classes own current Zeus presentation values. JRN-010
  subsequently removed runtime token/provider usage entirely.

### Automated verification

```text
npm run build -w @gym-saas/desktop-frontend
result: passed

npm run build -w @gym-saas/web
result: passed

npm run typecheck
result: passed

npm run lint:ts
result: passed with existing warnings outside the migration scope

Source scan
result: only host `tailwind.css` entry files remain; no feature CSS import or
inline presentation style object remains.

Local desktop visual check at 1511×920 logical pixels
result: Login, Members, Plans, and plan-create dialog retained their Zeus
presentation after utility migration.
```

### Known limitations / follow-ups

- JRN-010 supersedes runtime token usage with a static Tailwind Zeus palette.

## Human acceptance

### Manual check

1. Open Login, Dashboard, Clientes, Planes, and the create-plan dialog.
2. Confirm their Zeus layout and interactions remain unchanged.
3. Inspect the frontend source and confirm feature CSS files are absent.

**Expected result:** Current presentation is Tailwind-based, theme-aware, and
visually equivalent to the pre-migration screens.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
