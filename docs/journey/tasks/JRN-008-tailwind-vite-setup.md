# JRN-008 — Tailwind Vite setup

**Status:** Awaiting human acceptance
**Roadmap:** Supporting frontend delivery
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Add Tailwind CSS v4 through its official Vite integration in the desktop and
web hosts before new screens introduce utility classes.

## Scope

### In scope

- `tailwindcss` and `@tailwindcss/vite` development dependencies in both Vite
  hosts.
- Shared Tailwind imports and Vite plugins for desktop and web.
- Explicit source discovery for `packages/ui/src` so future shared components
  can use utilities.

### Out of scope

- Rewriting existing token CSS or components to Tailwind utilities.
- Replacing the Zeus theme provider or its CSS variables.
- Adding a JavaScript Tailwind configuration file unless a later need proves it
  necessary.

## Acceptance criteria

- [x] Desktop and web builds compile Tailwind through Vite.
- [x] Existing desktop and web typechecks/builds remain green.
- [x] Existing Zeus CSS remains visually intact.

## Implementation record

### Changed files

- `package-lock.json`
- `apps/desktop/frontend/package.json`
- `apps/web/package.json`
- `apps/desktop/frontend/vite.config.ts`
- `apps/web/vite.config.ts`
- `apps/desktop/frontend/src/tailwind.css`
- `apps/web/src/tailwind.css`
- Both Vite entry points and this journey task

### Simplifications made

- Tailwind uses v4's CSS-first/Vite-plugin setup; no JavaScript config or
  additional styling layer was introduced.
- Existing Zeus variables and feature CSS remain the active styling system
  during gradual utility adoption.

### Automated verification

```text
npm run build -w @gym-saas/desktop-frontend
result: passed; Tailwind compiled through Vite

npm run build -w @gym-saas/web
result: passed; Tailwind compiled through Vite

npm run typecheck
result: passed

Local desktop visual check at 1511×920 logical pixels
result: Zeus login card remained 450px wide with its original gradient and
visual hierarchy intact.
```

### Known limitations / follow-ups

- Existing screens have not been converted to utility classes yet; adopt
  Tailwind feature by feature rather than mixing large rewrites into new work.

## Human acceptance

### Manual check

1. Start desktop and web development servers.
2. Confirm existing screens render normally.
3. Add a Tailwind utility class to a future component and confirm Vite compiles
   it in both hosts.

**Expected result:** Tailwind utilities are ready to use alongside existing
Zeus theme variables and feature CSS.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
