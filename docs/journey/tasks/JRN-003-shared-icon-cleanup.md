# JRN-003 — Shared icon cleanup

**Status:** Awaiting human acceptance
**Roadmap:** Supporting visual-parity work
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Make the supplied V1 SVG icons a shared, Tailwind-independent UI asset that
desktop and web can use consistently.

## Scope

### In scope

- `packages/ui/src/icons/**` and the public UI export surface.
- Removing Tailwind-only default class names and the Tailwind-styled file icon
  wrapper.
- A natural SVG fallback size and `currentColor` inheritance for every icon.
- Replacing duplicate login SVGs with the shared user, lock, and eye icons.

### Out of scope

- Adding Tailwind, a new icon dependency, or an icon framework.
- Changing the artwork/path data of the supplied V1 icons.
- Reworking unrelated dashboard or login layout styles.

## Cleanup plan

1. Lock type and visual behavior with the existing frontend verification.
2. Remove Tailwind-only defaults and normalize icon rendering one smell at a
   time: dead wrapper styles, duplicate defaults, then API consistency.
3. Replace duplicate consumers in login with shared icons.
4. Re-run typecheck, lint, and local visual checks.

## Acceptance criteria

- [x] Every shared icon has a usable non-Tailwind default rendering path.
- [x] No icon source contains Tailwind utility class names.
- [x] Login and dashboard consume shared icons where a matching one exists.
- [x] `npm run typecheck` passes and the icon rendering is visually verified.

## Implementation record

### Changed files

- `packages/ui/src/icons/**`
- `packages/ui/src/index.tsx`
- `apps/desktop/frontend/src/features/auth/LoginPage.tsx`
- `apps/desktop/frontend/src/features/dashboard/DashboardPage.tsx`
- This journey task page

### Simplifications made

- Removed 50 Tailwind utility-class references and made every SVG naturally
  size to `1em` without an external CSS framework.
- Replaced the Tailwind-composed `FileIcon` wrapper with a pure SVG, and moved
  spinner animation into SVG rather than a Tailwind animation class.
- Deleted the duplicate login user/lock/eye SVG paths in favor of shared icons.

### Automated verification

```text
npm run typecheck
result: passed

npm run lint:ts
result: passed with pre-existing warnings outside this cleanup scope

Tailwind utility scan of `packages/ui/src/icons`
result: no matches

Local browser visual checks at 1511×920 logical pixels
result: login shared icons render at 18×18; all 16 dashboard navigation icons
and four toolbar-action icons render at nonzero dimensions.
```

### Known limitations / follow-ups

- The library intentionally remains a minimal SVG set; individual screens own
  their sizing and layout classes.

## Human acceptance

### Manual check

1. Run the desktop frontend mock.
2. Verify that the login fields and dashboard navigation/toolbars render their
   icons at the expected size and color without Tailwind loaded.

**Expected result:** Shared V1 icons work unchanged in the Zeus mock and are
ready for the web host to import from `@gym-saas/ui`.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
