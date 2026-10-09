-- Module: user — role (first-class entity per BRD §4; seeded master data).
-- PKs are native uuid (app-generated UUIDv7); see erd.md §2.
BEGIN;

CREATE TABLE IF NOT EXISTS role (
    id          UUID PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

COMMIT;
