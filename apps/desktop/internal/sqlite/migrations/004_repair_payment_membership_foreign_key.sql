-- Migration 002 rebuilt memberships, and SQLite retargeted this pre-existing
-- foreign key to memberships_legacy during the table rename. Rebuild payments
-- so new and upgraded databases both reference the final memberships table.
DROP INDEX idx_payments_member_created_at;
DROP INDEX idx_payments_gym_status_paid_at;
DROP INDEX idx_payments_membership;

ALTER TABLE payments RENAME TO payments_legacy;

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

INSERT INTO payments (
  id, gym_id, member_id, membership_id, status, kind, amount_cents, currency,
  payment_method, reference, notes, paid_at, created_at, updated_at
)
SELECT
  id, gym_id, member_id, membership_id, status, kind, amount_cents, currency,
  payment_method, reference, notes, paid_at, created_at, updated_at
FROM payments_legacy;

DROP TABLE payments_legacy;

CREATE INDEX idx_payments_member_created_at
  ON payments(member_id, created_at DESC);

CREATE INDEX idx_payments_gym_status_paid_at
  ON payments(gym_id, status, paid_at DESC);

CREATE INDEX idx_payments_membership
  ON payments(membership_id);
