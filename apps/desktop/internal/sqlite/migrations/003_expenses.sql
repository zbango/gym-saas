CREATE TABLE expenses (
  id TEXT PRIMARY KEY,
  gym_id TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'posted', 'voided')),
  amount_cents INTEGER NOT NULL CHECK (amount_cents > 0),
  currency TEXT NOT NULL DEFAULT 'USD',
  payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'credit_card', 'debit_card', 'bank_transfer', 'check', 'other')),
  reference TEXT,
  notes TEXT,
  paid_at TEXT,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
  CHECK (
    (status = 'posted' AND paid_at IS NOT NULL) OR
    (status IN ('pending', 'voided') AND paid_at IS NULL)
  )
);

CREATE INDEX idx_expenses_gym_status_paid_at
  ON expenses(gym_id, status, paid_at DESC);

CREATE INDEX idx_expenses_gym_created_at
  ON expenses(gym_id, created_at DESC);
