# JRN-001 — Zeus login visual mock

**Status:** Awaiting human acceptance
**Roadmap:** Supporting visual-parity work; authentication implementation is ZV2-124
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Recreate the supplied V1 Zeus Gym desktop login screen as a responsive,
presentation-only React view using the shared theme provider.

## Scope

### In scope

- Full-screen Zeus gradient canvas, centered dark login card, supplied logo,
  Spanish copy, field icons, password visibility control, and primary action.
- A reusable auth-specific field composition with accessible labels and native
  credential autofill semantics.
- A `zeus` theme definition and a login route that explicitly selects it.
- Local UI state only: username, password, and password visibility.

### Out of scope

- Credential verification, sessions, user roles, lockouts, password recovery,
  persistence, Go/Wails bindings, and navigation after sign-in.
- A generic cross-application text-field abstraction. The current field owns
  auth-only affordances and will be generalized only after a shared use case
  exists.
- Any business rule in React.

## Plan

1. Inspect the V1 reference and the supplied logo asset.
2. Define Zeus tokens through the existing shared theme provider.
3. Implement the responsive React screen and presentation-only interactions.
4. Compare a local render at the reference's 1511×920 logical desktop scale.
5. Record automated evidence and require owner visual acceptance.

## Acceptance criteria

- [x] The login screen reproduces the supplied V1 hierarchy, card geometry,
  amber canvas, controls, Spanish copy, and primary button.
- [x] The supplied `logo.png` appears in the desktop asset bundle rather than
  a recreated CSS mark.
- [x] The screen uses Zeus theme variables; no screen-specific color system is
  introduced.
- [x] Username and password have unambiguous accessible labels, and the
  password visibility control changes the input type.
- [x] `npm run typecheck` passes.
- [ ] The owner confirms the visual result against the supplied V1 screenshot.

## Risks and decisions

| Risk or decision | Mitigation or outcome |
|---|---|
| A prior persisted theme could override the reference look. | `ThemeProvider` supports `fixedTheme`; login uses `zeus` explicitly. |
| Mock UI could become accidental product logic. | The screen exposes no authentication, persistence, authorization, or post-login navigation. |
| Input needs may differ outside authentication. | Keep `AuthField` local until a real shared primitive is demonstrated. |

## Implementation record

### Changed files

- `packages/shared/src/theme.ts`
- `packages/ui/src/theme.tsx`
- `apps/desktop/frontend/src/assets/logo.png`
- `apps/desktop/frontend/src/features/auth/LoginPage.tsx`
- `apps/desktop/frontend/src/features/auth/login-page.css`
- `apps/desktop/frontend/src/main.tsx`
- `apps/desktop/frontend/index.html`
- This journey record and the visual inventory

### Simplifications made

- The root currently renders the login mock directly rather than adding a
  routing or mock-session system.
- Password recovery remains a visible, inert action until its backend contract
  is designed.

### Automated verification

```text
npm run typecheck
result: passed

Local browser check at 1511×920 logical pixels
result: login card rendered at 450×655; supplied logo rendered from
/src/assets/logo.png at 160×160; username/password fields and password
visibility control verified.
```

### Known limitations / follow-ups

- Authentication and recovery behavior are deliberately not implemented.
- Future screen tasks must add themselves to the visual inventory before their
  implementation begins.

## Human acceptance

### Manual check

1. Run `npm run dev -w @gym-saas/desktop-frontend`.
2. Compare the login screen to the supplied Zeus-theme V1 screenshot at a
   desktop-sized window.
3. Confirm the background, card proportions, supplied logo, Spanish copy,
   input affordances, and button hierarchy match the intended reference.

**Expected result:** The desktop starts on a Zeus-themed visual mock that
closely matches the supplied V1 login page without claiming to authenticate a
user.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
