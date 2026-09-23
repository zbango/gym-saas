-- Sanitized malformed V1 data. It is applied after schema-only.sql to prove
-- that inspection reports, rather than converts, unsafe source rows.
INSERT INTO customers (id, email, identification_number, external_id, has_access, status, gym_id) VALUES
  (1, 'one@example.test', 'DUPLICATE-ID', 'DUPLICATE-DEVICE', 2, 'active', 'gym-zeus'),
  (2, 'two@example.test', 'DUPLICATE-ID', 'DUPLICATE-DEVICE', 1, 'active', 'gym-zeus');
INSERT INTO membership_plans (id, type, price, duration, visits, gym_id) VALUES (1, 'count-based', 40.00, 1, NULL, 'gym-zeus');
INSERT INTO customer_memberships (id, customer_id, membership_plan_id, start_date, membership_price_usd, visits_remaining) VALUES (1, 999, 999, '2026-08-01T00:00:00Z', 40.00, 4);
INSERT INTO payments (id, customer_id, customer_membership_id, amount_usd, transaction_type) VALUES
  (1, 1, NULL, -10.00, 'income'),
  (2, 1, NULL, 10.005, 'income'),
  (3, 1, NULL, 5.00, 'mystery');
CREATE TABLE payments_backup AS SELECT * FROM payments;
