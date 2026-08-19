CREATE TABLE schema_migrations (
  version INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  checksum TEXT NOT NULL,
  applied_at TEXT NOT NULL
);

CREATE TABLE gyms (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  timezone TEXT NOT NULL DEFAULT 'UTC',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE members (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT,
  phone TEXT NOT NULL,
  identification_number TEXT,
  date_of_birth TEXT,
  address TEXT,
  status TEXT NOT NULL CHECK (status IN ('active', 'inactive', 'blocked')),
  notes TEXT,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  deleted_at TEXT,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_members_email_unique
  ON members(gym_id, email)
  WHERE email IS NOT NULL AND deleted_at IS NULL;

CREATE UNIQUE INDEX idx_members_identification_unique
  ON members(gym_id, identification_number)
  WHERE identification_number IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_members_gym_status
  ON members(gym_id, status);

CREATE INDEX idx_members_gym_created_at
  ON members(gym_id, created_at DESC);

CREATE INDEX idx_members_gym_phone
  ON members(gym_id, phone);

CREATE TABLE membership_plans (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  name TEXT NOT NULL,
  validity_kind TEXT NOT NULL CHECK (validity_kind IN ('time', 'visits')),
  duration_value INTEGER,
  duration_unit TEXT CHECK (duration_unit IN ('days', 'weeks', 'months', 'years')),
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
    (validity_kind = 'time' AND duration_value IS NOT NULL AND duration_value > 0 AND duration_unit IS NOT NULL AND visit_limit IS NULL) OR
    (validity_kind = 'visits' AND visit_limit IS NOT NULL AND visit_limit > 0 AND duration_value IS NULL AND duration_unit IS NULL)
  )
);

CREATE INDEX idx_membership_plans_gym_status
  ON membership_plans(gym_id, status);

CREATE INDEX idx_membership_plans_gym_display_order
  ON membership_plans(gym_id, display_order, created_at DESC);

CREATE UNIQUE INDEX idx_membership_plans_gym_name_active
  ON membership_plans(gym_id, name)
  WHERE deleted_at IS NULL;

CREATE TABLE memberships (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  member_id TEXT NOT NULL,
  membership_plan_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'active', 'expired', 'cancelled')),
  starts_at TEXT NOT NULL,
  ends_at TEXT,
  activated_at TEXT,
  cancelled_at TEXT,
  cancellation_reason TEXT,
  validity_kind_snapshot TEXT NOT NULL CHECK (validity_kind_snapshot IN ('time', 'visits')),
  duration_value_snapshot INTEGER,
  duration_unit_snapshot TEXT CHECK (duration_unit_snapshot IN ('days', 'weeks', 'months', 'years')),
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
    (validity_kind_snapshot = 'time' AND duration_value_snapshot IS NOT NULL AND duration_value_snapshot > 0 AND duration_unit_snapshot IS NOT NULL AND visit_limit_snapshot IS NULL AND visits_remaining IS NULL) OR
    (validity_kind_snapshot = 'visits' AND visit_limit_snapshot IS NOT NULL AND visit_limit_snapshot > 0 AND visits_remaining IS NOT NULL AND visits_remaining >= 0 AND duration_value_snapshot IS NULL AND duration_unit_snapshot IS NULL AND ends_at IS NULL)
  )
);

CREATE INDEX idx_memberships_member_status
  ON memberships(member_id, status, starts_at DESC);

CREATE INDEX idx_memberships_gym_status
  ON memberships(gym_id, status, starts_at DESC);

CREATE INDEX idx_memberships_gym_created_at
  ON memberships(gym_id, created_at DESC);

CREATE TABLE payments (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  member_id TEXT NOT NULL,
  membership_id TEXT,
  status TEXT NOT NULL CHECK (status IN ('pending', 'posted', 'voided', 'refunded')),
  kind TEXT NOT NULL CHECK (kind IN ('initial', 'renewal', 'adjustment', 'other')),
  amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0),
  currency TEXT NOT NULL DEFAULT 'USD',
  payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'credit_card', 'debit_card', 'bank_transfer', 'check', 'other')),
  reference TEXT,
  notes TEXT,
  paid_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
  FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE RESTRICT,
  FOREIGN KEY (membership_id) REFERENCES memberships(id) ON DELETE SET NULL
);

CREATE INDEX idx_payments_member_created_at
  ON payments(member_id, created_at DESC);

CREATE INDEX idx_payments_gym_status_paid_at
  ON payments(gym_id, status, paid_at DESC);

CREATE INDEX idx_payments_membership
  ON payments(membership_id);

CREATE TABLE visits (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  member_id TEXT NOT NULL,
  membership_id TEXT,
  visited_at TEXT NOT NULL,
  source TEXT NOT NULL CHECK (source IN ('desktop', 'device', 'manual', 'import')),
  access_result TEXT NOT NULL CHECK (access_result IN ('granted', 'denied')),
  denial_reason TEXT,
  external_event_id TEXT,
  created_at TEXT NOT NULL,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
  FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE RESTRICT,
  FOREIGN KEY (membership_id) REFERENCES memberships(id) ON DELETE SET NULL
);

CREATE INDEX idx_visits_member_visited_at
  ON visits(member_id, visited_at DESC);

CREATE INDEX idx_visits_gym_visited_at
  ON visits(gym_id, visited_at DESC);

CREATE UNIQUE INDEX idx_visits_external_event_unique
  ON visits(gym_id, source, external_event_id)
  WHERE external_event_id IS NOT NULL;
