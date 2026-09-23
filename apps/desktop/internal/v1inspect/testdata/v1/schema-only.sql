-- Sanitized structural fixture based on V1 migrations 001–028. It contains no
-- customer data or device credentials and is intentionally not a V2 schema.
CREATE TABLE migrations (version INTEGER PRIMARY KEY, name TEXT NOT NULL, applied_at TEXT NOT NULL);
INSERT INTO migrations (version, name, applied_at) VALUES
  (1, 'migration_001', '2026-01-01T00:00:00Z'), (2, 'migration_002', '2026-01-01T00:00:00Z'),
  (3, 'migration_003', '2026-01-01T00:00:00Z'), (4, 'migration_004', '2026-01-01T00:00:00Z'),
  (5, 'migration_005', '2026-01-01T00:00:00Z'), (6, 'migration_006', '2026-01-01T00:00:00Z'),
  (7, 'migration_007', '2026-01-01T00:00:00Z'), (8, 'migration_008', '2026-01-01T00:00:00Z'),
  (9, 'migration_009', '2026-01-01T00:00:00Z'), (10, 'migration_010', '2026-01-01T00:00:00Z'),
  (11, 'migration_011', '2026-01-01T00:00:00Z'), (12, 'migration_012', '2026-01-01T00:00:00Z'),
  (13, 'migration_013', '2026-01-01T00:00:00Z'), (14, 'migration_014', '2026-01-01T00:00:00Z'),
  (15, 'migration_015', '2026-01-01T00:00:00Z'), (16, 'migration_016', '2026-01-01T00:00:00Z'),
  (17, 'migration_017', '2026-01-01T00:00:00Z'), (18, 'migration_018', '2026-01-01T00:00:00Z'),
  (19, 'migration_019', '2026-01-01T00:00:00Z'), (20, 'migration_020', '2026-01-01T00:00:00Z'),
  (21, 'migration_021', '2026-01-01T00:00:00Z'), (22, 'migration_022', '2026-01-01T00:00:00Z'),
  (23, 'migration_023', '2026-01-01T00:00:00Z'), (24, 'migration_024', '2026-01-01T00:00:00Z'),
  (25, 'migration_025', '2026-01-01T00:00:00Z'), (26, 'migration_026', '2026-01-01T00:00:00Z'),
  (27, 'migration_027', '2026-01-01T00:00:00Z'), (28, 'migration_028', '2026-01-01T00:00:00Z');

CREATE TABLE system_users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, role TEXT, status TEXT, gym_id TEXT);
CREATE TABLE activation (id INTEGER PRIMARY KEY, gym_id TEXT, gym_name TEXT, region TEXT);
CREATE TABLE sessions (id INTEGER PRIMARY KEY, user_id INTEGER REFERENCES system_users(id), token TEXT, expires_at TEXT);
CREATE TABLE membership_plans (id INTEGER PRIMARY KEY, type TEXT, price REAL, duration INTEGER, visits INTEGER, gym_id TEXT);
CREATE TABLE survey_questions (id INTEGER PRIMARY KEY, question TEXT, rating_type TEXT);
CREATE TABLE customers (id INTEGER PRIMARY KEY, email TEXT, identification_number TEXT, external_id TEXT, has_access INTEGER, status TEXT, gym_id TEXT);
CREATE TABLE customer_progress_tracking (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id));
CREATE TABLE customer_memberships (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id), membership_plan_id INTEGER REFERENCES membership_plans(id), start_date TEXT, membership_price_usd REAL, visits_remaining INTEGER);
CREATE TABLE payments (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id), customer_membership_id INTEGER REFERENCES customer_memberships(id), amount_usd REAL, transaction_type TEXT);
CREATE TABLE products (id INTEGER PRIMARY KEY, price REAL, stock INTEGER, gym_id TEXT);
CREATE TABLE sales (id INTEGER PRIMARY KEY, total_amount REAL, customer_id INTEGER REFERENCES customers(id), gym_id TEXT);
CREATE TABLE sale_items (id INTEGER PRIMARY KEY, sale_id INTEGER REFERENCES sales(id), product_id INTEGER REFERENCES products(id), unit_price REAL, subtotal REAL);
CREATE TABLE device_config (id INTEGER PRIMARY KEY, gym_id TEXT, ip TEXT, username TEXT, password TEXT);
CREATE TABLE customer_visits (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id), membership_id INTEGER REFERENCES customer_memberships(id), visit_date TEXT, membership_status TEXT);
CREATE TABLE audit_log (id INTEGER PRIMARY KEY, action TEXT, entity_type TEXT);
CREATE TABLE alert_log (id INTEGER PRIMARY KEY, alert_type TEXT, gym_id TEXT);
CREATE TABLE marketing_fidelity_notifications (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id), gym_id TEXT);
CREATE TABLE marketing_survey_invitations (id INTEGER PRIMARY KEY, customer_id INTEGER REFERENCES customers(id), customer_membership_id INTEGER REFERENCES customer_memberships(id), gym_id TEXT, token TEXT);
CREATE TABLE survey_responses (id INTEGER PRIMARY KEY, invitation_id INTEGER REFERENCES marketing_survey_invitations(id), question_text TEXT, answer_value TEXT);
