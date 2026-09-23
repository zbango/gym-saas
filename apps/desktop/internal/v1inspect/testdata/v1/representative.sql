-- Sanitized V1 fixture data. It is applied after schema-only.sql.
INSERT INTO activation (id, gym_id, gym_name, region) VALUES (1, 'gym-zeus', 'Zeus Fixture Gym', 'America/Guayaquil');
INSERT INTO system_users (id, username, email, role, status, gym_id) VALUES (1, 'fixture.admin', 'fixture.admin@example.test', 'gym_admin', 'active', 'gym-zeus');
INSERT INTO sessions (id, user_id, token, expires_at) VALUES (1, 1, 'fixture-token-not-production', '2026-12-31T00:00:00Z');
INSERT INTO customers (id, email, identification_number, external_id, has_access, status, gym_id) VALUES (1, 'member@example.test', 'FIXTURE-001', 'device-1001', 1, 'active', 'gym-zeus');
INSERT INTO membership_plans (id, type, price, duration, visits, gym_id) VALUES (1, 'monthly', 45.00, 1, NULL, 'gym-zeus');
INSERT INTO customer_memberships (id, customer_id, membership_plan_id, start_date, membership_price_usd, visits_remaining) VALUES (1, 1, 1, '2026-08-01T00:00:00Z', 45.00, NULL);
INSERT INTO payments (id, customer_id, customer_membership_id, amount_usd, transaction_type) VALUES (1, 1, 1, 45.00, 'income'), (2, NULL, NULL, 12.50, 'expense');
INSERT INTO device_config (id, gym_id, ip, username, password) VALUES (1, 'gym-zeus', '192.0.2.10', 'fixture-device', 'REDACTED');
INSERT INTO customer_visits (id, customer_id, membership_id, visit_date, membership_status) VALUES (1, 1, 1, '2026-08-02T08:00:00Z', 'active');
INSERT INTO products (id, price, stock, gym_id) VALUES (1, 3.50, 8, 'gym-zeus');
INSERT INTO sales (id, total_amount, customer_id, gym_id) VALUES (1, 3.50, 1, 'gym-zeus');
INSERT INTO sale_items (id, sale_id, product_id, unit_price, subtotal) VALUES (1, 1, 1, 3.50, 3.50);
INSERT INTO survey_questions (id, question, rating_type) VALUES (1, 'Fixture survey question', '1-5');
INSERT INTO marketing_fidelity_notifications (id, customer_id, gym_id) VALUES (1, 1, 'gym-zeus');
INSERT INTO marketing_survey_invitations (id, customer_id, customer_membership_id, gym_id, token) VALUES (1, 1, 1, 'gym-zeus', 'fixture-survey-token');
INSERT INTO survey_responses (id, invitation_id, question_text, answer_value) VALUES (1, 1, 'Fixture survey question', '5');
INSERT INTO audit_log (id, action, entity_type) VALUES (1, 'fixture_create', 'customer');
INSERT INTO alert_log (id, alert_type, gym_id) VALUES (1, 'fixture_alert', 'gym-zeus');
