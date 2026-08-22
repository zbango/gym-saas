# V2 First Slice Schema

This document defines the first real V2 operational schema for:

- `members`
- `membership_plans`
- `memberships`
- `payments`
- `visits`

The SQL source of truth is the ordered migration set:

- [001_v2_operational_core.sql](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/sqlite/migrations/001_v2_operational_core.sql)
- [002_visit_plan_expiry.sql](/Users/zbango/Documents/ChatGPT/gym/apps/desktop/internal/sqlite/migrations/002_visit_plan_expiry.sql)

## Runtime migration behavior

The desktop embeds this SQL into the Go binary and applies migrations in version order.
Each migration and its `schema_migrations` record commit in the same SQLite
transaction. At startup, the runner refuses to continue when an applied
migration's name or checksum differs from the embedded source, when the
database requires a newer unknown migration, or when an unmanaged database
already has application tables.

The only compatibility exception is the earlier desktop proof database with
its sole `hello_records` table. It can receive the V2 baseline unchanged;
the obsolete table is retained but no longer used. V1 databases remain
unmanaged and are intentionally refused until the dedicated read-only V1-to-V2
migration is implemented.

Migration 002 refuses to invent an expiry for any legacy visit plan or visit
membership. That failure is deliberate: an operator must choose the missing
calendar window rather than silently receiving an incorrect expiry date.

The SQLite connection configures foreign-key enforcement, WAL journaling, a
five-second busy timeout, and `synchronous=NORMAL`. The single desktop store
uses one connection, preventing in-process write contention while those
connection-level settings remain consistent.

## Reuse vs Rewrite

V1 schema is a **reference**, not a base layer to reuse directly.

Why we are rewriting instead of porting V1 migrations:

- V1 relies on `INTEGER PRIMARY KEY AUTOINCREMENT`, which is a poor fit for offline-first sync.
- V1 stores money as `REAL` in multiple tables.
- V1 mixes operational core with marketing, surveys, products, sales, activation, and other secondary concerns.
- V1 carries business shortcuts like `customers.has_access` and `customers.external_id` that should not survive the V2 model.
- V1 plan validity (`type + duration + visits`) is ambiguous for future domain logic.

## Core Design Decisions

### IDs

Every primary key is `TEXT`.

Intended runtime format:

- UUID or ULID

Reason:

- offline-safe creation
- easier future sync/cloud merge behavior

### Time

Every timestamp column is `TEXT` and should hold UTC ISO-8601 values.

Reason:

- simple SQLite storage
- explicit conversion in presentation layers

### Money

All money is stored as integer cents:

- `price_cents`
- `price_cents_snapshot`
- `amount_cents`

Reason:

- avoid float drift
- preserve deterministic accounting behavior

### Tenant ownership

Each operational row is owned by a `gym_id`.

Even in a single-gym local install, this keeps the model future-safe and prevents a single-tenant trap.

### Snapshot semantics

`memberships` store:

- `validity_kind_snapshot`
- `duration_value_snapshot`
- `duration_unit_snapshot`
- `visit_limit_snapshot`
- `visits_remaining`
- `price_cents_snapshot`
- `currency_snapshot`

Reason:

- editing a plan later must not rewrite historical memberships

## Tables

### `gyms`

Minimal tenant anchor.

Columns:

- `id`
- `name`
- `timezone`
- `created_at`
- `updated_at`

### `members`

The V2 replacement for V1 `customers`.

Notable changes from V1:

- no `external_id`
- no `has_access`
- sync-safe soft deletion via `deleted_at`
- `status` is explicit and limited to `active | inactive | blocked`

Indexes:

- by `(gym_id, status)`
- by `(gym_id, created_at desc)`
- by `(gym_id, phone)`
- unique partial indexes for email and identification number

### `membership_plans`

Plan validity is explicit instead of inferred.

Two modes:

- `time`
  requires `duration_value + duration_unit`
- `visits`
  requires both `visit_limit` and `duration_value + duration_unit`

Visit plans expire when their visits are exhausted or their calendar window
ends, whichever happens first. This avoids the V1 ambiguity where
monthly-like plans and duration counts drift apart.

Indexes:

- `(gym_id, status)`
- `(gym_id, display_order, created_at desc)`
- unique active name per gym

### `memberships`

Represents a purchased membership instance.

Tracks:

- status
- start/end dates
- cancellation
- plan snapshots
- remaining visits for visit-based memberships

Important rules:

- every membership has `ends_at` from its duration snapshot
- visit-based memberships also keep `visits_remaining`
- a visit-based membership expires at the earlier of `ends_at` and no
  remaining visits

Indexes:

- `(member_id, status, starts_at desc)`
- `(gym_id, status, starts_at desc)`
- `(gym_id, created_at desc)`

### `payments`

Represents member-side incoming payments only for this first slice.

It does **not** model general expenses yet.

Tracks:

- status
- kind (`initial`, `renewal`, `adjustment`, `other`)
- amount/currency
- payment method
- optional membership link

Indexes:

- `(member_id, created_at desc)`
- `(gym_id, status, paid_at desc)`
- `(membership_id)`

### `visits`

Auditable visit/access facts.

Tracks:

- member
- membership used, if any
- visit timestamp
- source
- granted/denied result
- denial reason
- optional external event identity

Important rule:

- `external_event_id` is unique per `(gym_id, source)` when present
- this gives a dedupe path for future device/webhook integrations

Indexes:

- `(member_id, visited_at desc)`
- `(gym_id, visited_at desc)`
- unique `(gym_id, source, external_event_id)` partial index

## What is intentionally not in this first slice

- products
- sales
- expenses
- surveys
- marketing tables
- activation/licensing
- sessions/auth
- device jobs
- sync outbox

Those should be added in later slices, not imported from V1 by default.

## Query patterns this schema is optimized for

Primary expected queries:

- latest members for a gym
- active/blocked members for a gym
- active/latest memberships for a member
- latest payments for a member
- latest visits for a member
- recent visits for a gym
- plans ordered for presentation

The first slice indexes are chosen for those patterns, not for every hypothetical report.

## Recommendation

Use this schema as the V2 baseline for the first operational slice and treat V1 migrations only as behavioral evidence.
