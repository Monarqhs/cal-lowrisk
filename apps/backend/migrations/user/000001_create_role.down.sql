BEGIN;

DROP TABLE IF EXISTS "user".role;

-- Drop the module schema only if empty (safe: later migrations' tables are gone by now).
DROP SCHEMA IF EXISTS "user" RESTRICT;

COMMIT;
