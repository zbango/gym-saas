# JRN-002 — Super-admin dashboard visual mock

**Status:** Awaiting human acceptance
**Roadmap:** Supporting visual-parity work; feature use cases remain separate
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Recreate the supplied V1 super-admin dashboard shell with its full navigation,
device-connection status, breadcrumbs, contextual actions, version display,
and mocked dashboard content.

## Scope

### In scope

- Zeus-themed desktop sidebar and the complete super-admin visual navigation.
- Expanded Memberships and Store child menus shown in the supplied reference.
- Device connection banner; breadcrumb/action bar; displayed application
  version; dashboard cards, segmented financial period control, and alert UI.
- Presentation-only view navigation from the login mock to the dashboard mock.

### Out of scope

- Authentication, authorization, permission evaluation, device connections,
  retry behavior, monitoring, financial summaries, navigation routes, and data
  persistence.
- Implementing any business rule or role decision in React. The mock accepts a
  displayed user profile; backend integration will supply permitted navigation
  and authoritative content later.

## Plan

1. Record the screenshots in the visual inventory and model the dashboard
   shell as presentation components.
2. Use shared Zeus tokens for all visual states, including status colors.
3. Add only mock UI state needed to demonstrate menu expansion and local
   dashboard controls.
4. Render at the reference desktop scale, compare visual hierarchy, and record
   verification.

## Acceptance criteria

- [x] The visual super-admin sidebar includes every supplied top-level and
  child menu entry, plus profile, version, and logout affordances.
- [x] The device error banner, breadcrumb/action bar, and dashboard content
  match the supplied hierarchy and Zeus treatment.
- [x] The screen exposes no real role, device, financial, or routing logic.
- [x] Typechecking passes and visual comparison evidence is recorded.
- [ ] The owner accepts the visual result against the supplied references.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| Displaying all super-admin options could be mistaken for authorization. | Explicitly retain this as a static visual profile; backend will provide permitted navigation. |
| Dashboard controls imply data changes. | Buttons and selectors demonstrate UI only; no persistence or device calls are introduced. |
| Status colors could bypass theming. | Add semantic status tokens to the shared theme provider instead of page-local colors. |

## Implementation record

### Changed files

- `packages/shared/src/theme.ts`
- `packages/ui/src/theme.tsx`
- `packages/ui/src/icons/`
- `packages/ui/src/index.tsx`
- `apps/desktop/frontend/src/features/dashboard/DashboardPage.tsx`
- `apps/desktop/frontend/src/features/dashboard/dashboard-page.css`
- `apps/desktop/frontend/src/main.tsx`
- `docs/journey/v1-visual-inventory.md`
- This journey task page

### Simplifications made

- The dashboard shell uses static super-admin display data and local state for
  active navigation, expanded menu groups, alert visibility, and period
  selection; it makes no backend or device call.
- The login transition is a mock screen transition only, not authentication.
- Dashboard icons are shared `@gym-saas/ui` SVG components styled with local
  CSS classes, not Tailwind utilities.

### Automated verification

```text
npm run typecheck
result: passed

Local browser visual check at 1511×920 logical pixels
result: sidebar 253px wide and 920px high; profile/version/logout footer,
three toolbar actions, and both expanded child-menu groups visible.

Local mock interaction check
result: Memberships collapses/re-expands; Planes opens its mock route; displayed
desktop version is `0.1.2`.

Shared icon render check
result: 16 dashboard navigation icons and four toolbar-action icons rendered
with nonzero dimensions and inherited their CSS color.
```

### Known limitations / follow-ups

- The dashboard uses mock text and a static visual super-admin profile.
- Backend integration must supply authenticated identity, permitted navigation,
  device state, real version metadata, dashboard read models, and route data.
- Several un-used legacy icon components still default to Tailwind utility
  class names. Normalize their defaults before consuming them without an
  explicit project CSS class.

## Human acceptance

### Manual check

1. Run the desktop frontend mock and submit the login form.
2. Compare the dashboard and expanded sidebar groups to both supplied V1
   screenshots.
3. Confirm that actions are visibly present but do not claim to perform a real
   operation.

**Expected result:** A Zeus-themed super-admin dashboard shell matches the
provided V1 presentation and remains unmistakably a visual mock.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
