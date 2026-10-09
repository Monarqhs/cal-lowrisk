-- Module: user — role (first-class entity per BRD §4; seeded master data).
-- PKs are native uuid (app-generated UUIDv7); see erd.md §2.
-- Each module owns a dedicated PostgreSQL schema (= module name). "user" is a
-- reserved word, so it MUST always be double-quoted. See product.md + add-module §1.
BEGIN;

CREATE SCHEMA IF NOT EXISTS "user";

CREATE TABLE IF NOT EXISTS "user".role (
    id          UUID PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

COMMIT;
