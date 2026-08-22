# Desktop vertical slices

This document records the boundary rules used by the desktop app. A feature is
complete only when its user-facing action reaches the durable local database
and the UI is refreshed from the result.

## Member flow

```text
React feature component
  -> feature API client
  -> feature-specific Wails API (delivery adapter)
  -> application service (use case)
  -> repository port
  -> SQLite repository adapter
  -> SQLite database

                       application service
                             <-> domain
```

`domain` is not the last step in a linear pipeline. It sits at the center:
the application service uses domain constructors and value types to validate
commands before persistence, while the SQLite adapter reconstructs valid
domain objects when it reads data.

For Members, the concrete route is:

```text
MemberPanel -> features/members/api.ts -> MemberAPI.CreateMember
-> application.MemberService.Create -> ports.MemberRepository.Create
-> sqlite.MemberRepository.Create -> members table
```

Membership Plans use the same route through `MembershipPlanPanel`, the
feature-local Plans API, `MembershipPlanAPI.CreateMembershipPlan`, and
`application.MembershipPlanService`. The domain requires every visit plan to
include both a visit limit and a calendar duration; a future purchased
membership will snapshot those values and expire when either bound is reached.

After a mutation, `useMembers` reloads the list through the same route in
reverse (`ListMembers`), so the screen reflects database state rather than a
locally guessed update. Archive is a soft deletion: the repository sets
`deleted_at`, and normal reads intentionally exclude archived records.

## Layer responsibilities

| Layer | Owns | Must not own |
| --- | --- | --- |
| React feature | Form state, presentation, explicit UI confirmation | Business validation or SQL |
| Feature API client | Typed calls for one Wails feature | Other feature bindings or UI state |
| Wails feature API | Request/response DTO conversion and context handoff for one feature | Business rules or other features' use cases |
| Application service | Use-case sequencing, tenant scope, domain validation, error context | SQLite or Wails imports |
| Domain | Invariants and valid `Member` construction | HTTP, SQLite, Wails, or device imports |
| Repository port | Storage contract needed by an application service | SQL implementation details |
| SQLite adapter | SQL, row mapping, soft-delete predicates, transactions | UI decisions |

## Wails binding organization

Wails provides one global namespace, but the desktop binds separate API
objects: `window.go.main.MemberAPI`, `window.go.main.MembershipPlanAPI`, and
`window.go.main.DesktopAPI`. This prevents a growing all-purpose `App` object.

- `src/platform/wails.ts` is the small, framework-specific global boundary.
- `src/features/members/api.ts` owns only Member calls and Member transport
  types.
- `src/features/plans/api.ts` owns only Membership Plan calls and transport
  types.
- `src/features/updates/api.ts` owns only updater calls and types.
- A new feature gets its own `features/<feature>/api.ts`; it does not extend a
  global list of bindings or an existing API object.

Wails generates the low-level modules under `frontend/wailsjs/` during desktop
development and production builds. Feature API clients wrap and type those
generated modules; components do not import them directly. Do not edit the
generated files, and do not recreate a hand-written global binding list.

## Mutation/query state: no RTK Query yet

RTK Query is designed around Redux and remote data caches. This desktop app
currently uses direct Wails RPC calls to a local SQLite database and has no
Redux store, so adding RTK Query now would introduce a second state framework
without solving a present problem.

`features/members/useMembers.ts` is the current light-weight pattern:

1. load the canonical list on feature mount;
2. submit a typed Wails mutation;
3. reload after a successful mutation;
4. keep pending/error/form state local to the feature.

Reassess a query library (more likely TanStack Query than RTK Query) when
several screens share the same read models, need invalidation across features,
or require background/refetch behavior. Introduce it only with a concrete
cache and invalidation policy.

## SQL and ORM decision

Use explicit SQL in the SQLite adapter for now. The queries are small,
reviewable, use bound parameters, and keep transaction and soft-delete logic
visible. This matches the project rule to prefer simple Go and not add
dependencies without a demonstrated need.

An ORM is not a replacement for the repository boundary; it would live only
inside an infrastructure adapter. Consider one only if repeated, measurable
mapping/query boilerplate outweighs its migration, transaction, and query
visibility costs. Do not let an ORM model become the domain model.

## Rules for the next feature

1. Add or extend domain rules first, with domain tests.
2. Define only the repository operations required by the application use case.
3. Implement the application service using domain constructors, not UI checks.
4. Keep SQL in the SQLite adapter and test persistence/restart behavior.
5. Add a feature-local Wails API client and feature-local state hook.
6. Refresh data from the authoritative use case after mutations.
7. If a mutation becomes cloud-replicated, write its outbox record in the same
   SQLite transaction; do not simulate sync success in the UI.
