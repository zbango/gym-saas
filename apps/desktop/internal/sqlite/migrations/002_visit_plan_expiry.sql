-- Visit plans now have a calendar expiry as well as a visit allowance. Refuse
-- to guess an expiry for pre-existing visit records; none existed before the
-- plan feature, and a future data migration must make that decision explicit.
CREATE TABLE membership_plan_expiry_guard (
  value INTEGER NOT NULL CHECK (value = 1)
);

INSERT INTO membership_plan_expiry_guard (value)
SELECT 0
WHERE EXISTS (SELECT 1 FROM membership_plans WHERE validity_kind = 'visits')
   OR EXISTS (
     SELECT 1
     FROM memberships
     WHERE validity_kind_snapshot = 'visits'
        OR (validity_kind_snapshot = 'time' AND ends_at IS NULL)
   );

DROP TABLE membership_plan_expiry_guard;

DROP INDEX idx_membership_plans_gym_status;
DROP INDEX idx_membership_plans_gym_display_order;
DROP INDEX idx_membership_plans_gym_name_active;

ALTER TABLE memberships RENAME TO memberships_legacy;
ALTER TABLE membership_plans RENAME TO membership_plans_legacy;

CREATE TABLE membership_plans (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  name TEXT NOT NULL,
  validity_kind TEXT NOT NULL CHECK (validity_kind IN ('time', 'visits')),
  duration_value INTEGER NOT NULL CHECK (duration_value > 0),
  duration_unit TEXT NOT NULL CHECK (duration_unit IN ('days', 'weeks', 'months', 'years')),
  visit_limit INTEGER,
  price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
  currency TEXT NOT NULL DEFAULT 'USD',
  status TEXT NOT NULL CHECK (status IN ('active', 'inactive', 'archived')),
  display_order INTEGER NOT NULL DEFAULT 0,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  deleted_at TEXT,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
  CHECK (
    (validity_kind = 'time' AND visit_limit IS NULL) OR
    (validity_kind = 'visits' AND visit_limit IS NOT NULL AND visit_limit > 0)
  )
);

CREATE TABLE memberships (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  member_id TEXT NOT NULL,
  membership_plan_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'active', 'expired', 'cancelled')),
  starts_at TEXT NOT NULL,
  ends_at TEXT NOT NULL,
  activated_at TEXT,
  cancelled_at TEXT,
  cancellation_reason TEXT,
  validity_kind_snapshot TEXT NOT NULL CHECK (validity_kind_snapshot IN ('time', 'visits')),
  duration_value_snapshot INTEGER NOT NULL CHECK (duration_value_snapshot > 0),
  duration_unit_snapshot TEXT NOT NULL CHECK (duration_unit_snapshot IN ('days', 'weeks', 'months', 'years')),
  visit_limit_snapshot INTEGER,
  visits_remaining INTEGER,
  price_cents_snapshot INTEGER NOT NULL CHECK (price_cents_snapshot >= 0),
  currency_snapshot TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
  FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE RESTRICT,
  FOREIGN KEY (membership_plan_id) REFERENCES membership_plans(id) ON DELETE RESTRICT,
  CHECK (
    (validity_kind_snapshot = 'time' AND visit_limit_snapshot IS NULL AND visits_remaining IS NULL) OR
    (validity_kind_snapshot = 'visits' AND visit_limit_snapshot IS NOT NULL AND visit_limit_snapshot > 0 AND visits_remaining IS NOT NULL AND visits_remaining >= 0)
  )
);

INSERT INTO membership_plans (
  id, gym_id, name, validity_kind, duration_value, duration_unit, visit_limit,
  price_cents, currency, status, display_order, version, created_at, updated_at, deleted_at
)
SELECT
  id, gym_id, name, validity_kind, duration_value, duration_unit, visit_limit,
  price_cents, currency, status, display_order, version, created_at, updated_at, deleted_at
FROM membership_plans_legacy;

INSERT INTO memberships (
  id, gym_id, member_id, membership_plan_id, status, starts_at, ends_at,
  activated_at, cancelled_at, cancellation_reason, validity_kind_snapshot,
  duration_value_snapshot, duration_unit_snapshot, visit_limit_snapshot,
  visits_remaining, price_cents_snapshot, currency_snapshot, version,
  created_at, updated_at
)
SELECT
  id, gym_id, member_id, membership_plan_id, status, starts_at, ends_at,
  activated_at, cancelled_at, cancellation_reason, validity_kind_snapshot,
  duration_value_snapshot, duration_unit_snapshot, visit_limit_snapshot,
  visits_remaining, price_cents_snapshot, currency_snapshot, version,
  created_at, updated_at
FROM memberships_legacy;

DROP TABLE memberships_legacy;
DROP TABLE membership_plans_legacy;

CREATE INDEX idx_membership_plans_gym_status
  ON membership_plans(gym_id, status);

CREATE INDEX idx_membership_plans_gym_display_order
  ON membership_plans(gym_id, display_order, created_at DESC);

CREATE UNIQUE INDEX idx_membership_plans_gym_name_active
  ON membership_plans(gym_id, name)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_memberships_member_status
  ON memberships(member_id, status, starts_at DESC);

CREATE INDEX idx_memberships_gym_status
  ON memberships(gym_id, status, starts_at DESC);

CREATE INDEX idx_memberships_gym_created_at
  ON memberships(gym_id, created_at DESC);
