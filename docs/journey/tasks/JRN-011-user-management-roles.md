# JRN-011 — User management roles

**Status:** Awaiting human acceptance
**Roadmap:** Supporting user/auth work; user creation remains separate
**Owner:** Product owner + Codex
**Opened:** 2026-08-22

## Objective

Define canonical system roles in Go and recreate the supplied Gestión de
Usuarios list as a Tailwind visual mock.

## Role policy

| Role | Can manage employee users |
|---|---|
| `super_admin` | Yes |
| `gym_owner` | Yes |
| `gym_admin` | Yes |
| `trainer` | No |
| `receptionist` | No |
| `member` | No |

## Scope

### In scope

- Dependency-free Go `UserRole` parsing and `CanManageEmployeeUsers` policy.
- Gestión de Usuarios mock page with synthetic employee rows, role/status
  badges, metrics, controls, table, and inactive add-user affordance.

### Out of scope

- User accounts, authentication, sessions, invitations, passwords, persistence,
  real role assignment, and authorization middleware.
- A React authority check. Future user mutation use cases must enforce the Go
  domain policy.

## Acceptance criteria

- [x] All six roles parse in Go; invalid roles fail.
- [x] Only super admin, gym owner, and gym admin pass the employee-user policy.
- [x] The mock page matches the supplied hierarchy with synthetic data.
- [x] Go tests and frontend typecheck pass.

## Implementation record

### Changed files

- `go/core/domain/user_role.go`
- `go/core/domain/user_role_test.go`
- `apps/desktop/frontend/src/features/users/UserManagementPage.tsx`
- `apps/desktop/frontend/src/features/dashboard/DashboardRouteOutlet.tsx`
- `packages/ui/src/PageHeader.tsx`
- `packages/ui/src/MetricCard.tsx`
- `packages/ui/src/DataTable.tsx`
- This journey task page

### Simplifications made

- The page composes shared PageHeader, MetricCard, and DataTable components;
  screen code only supplies user-specific rows, columns, metrics, and actions.
- Role parsing and user-management policy are canonical Go domain code.
- The UI renders synthetic employees only and has no creation/authorization
  behavior.

### Automated verification

```text
cd go/core && go test ./...
result: passed

npm run typecheck
result: passed

Local browser visual check at 1511×920 logical pixels
result: Gestión de Usuarios displays four metrics, five user rows, role/status
badges, search, table, and pagination hierarchy matching the V1 reference.
```

### Known limitations / follow-ups

- User creation, authentication, persistence, invitations, and user-list
  integration remain separate tasks.
- The Add User button is visual-only until a Go application use case can
  enforce `CanManageEmployeeUsers`.

## Human acceptance

### Manual check

1. Open Gestión de Usuarios from the dashboard sidebar.
2. Compare metrics, search, role/status badges, and table hierarchy to the
   supplied V1 screen.
3. Confirm Add User is visual-only and creation remains unimplemented.

**Expected result:** The page is visually equivalent to V1 while role policy
lives in Go, ready for a later user-creation flow.

### Owner result

- [ ] Accepted
- [ ] Changes requested

**Date:**
**Notes:**
