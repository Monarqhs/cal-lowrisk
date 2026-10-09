-- Module: user — seed roles (master data). Fixed, hardcoded UUIDs so FKs
-- (users.role_id) can reference known values. Idempotent (erd.md §2, add-module §4).
BEGIN;

INSERT INTO role (id, name, description) VALUES
    ('00000000-0000-7000-8000-000000000001', 'user',  'Default role for end users (mobile app)'),
    ('00000000-0000-7000-8000-000000000002', 'admin', 'Administrator (web admin): catalog + account management')
ON CONFLICT (name) DO NOTHING;

COMMIT;
