# JRN-010 — Static Zeus Tailwind theme

**Status:** Awaiting human acceptance
**Roadmap:** Supporting frontend delivery
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Make Tailwind the only active presentation system by replacing runtime theme
provider/token usage with a single static Zeus theme.

## Scope

### In scope

- Shared static Zeus palette/font values in one Tailwind v4 brand profile.
- Removing active `ThemeProvider`, `ThemeSwitcher`, and `--gs-*` usage.
- Removing unused shared theme-provider/theme-definition source files.

### Out of scope

- Reintroducing multiple runtime-selectable themes.
- Changing screen behavior, domain logic, or the selected Zeus visual palette.

## Acceptance criteria

- [x] Active frontend source has no `ThemeProvider`, `ThemeSwitcher`, or
  `--gs-*` reference.
- [x] Desktop and web build with static Tailwind Zeus styling.
- [x] Login and Plans visually retain their Zeus presentation.

## Implementation record

### Changed files

- `packages/ui/src/brand.css`
- `apps/desktop/frontend/src/tailwind.css`
- `apps/web/src/tailwind.css`
- Both host entry points and all current Zeus utility consumers
- `packages/ui/src/index.tsx`
- `packages/shared/src/index.ts`
- Deleted `packages/ui/src/theme.tsx` and `packages/shared/src/theme.ts`

### Simplifications made

- Static colors replace runtime theme selection and CSS-token propagation.
- Both hosts import the same `brand.css`; a future gym's palette is changed in
  that one semantic profile rather than in component classes.

### Automated verification

```text
npm run build -w @gym-saas/desktop-frontend
result: passed

npm run build -w @gym-saas/web
result: passed

npm run typecheck
result: passed

Source scan
result: no ThemeProvider, ThemeSwitcher, useTheme, --gs-* or inline
presentation style object in active frontend/shared UI source.

Local visual check at 1511×920 logical pixels
result: Login flow and Zeus membership-plan screen render normally.
```

### Known limitations / follow-ups

- The application no longer supports runtime theme switching. Add a new,
  separately scoped Tailwind theme strategy only if multiple themes return as a
  product requirement.

## Human acceptance

### Manual check

1. Open Login and Planes.
2. Confirm the static Zeus appearance matches the approved current visual.
3. Confirm no runtime theme selector is present.

**Expected result:** Tailwind is the sole active styling system with a static
Zeus palette.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
