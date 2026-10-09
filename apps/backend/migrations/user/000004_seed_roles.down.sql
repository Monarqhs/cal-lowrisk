BEGIN;

DELETE FROM "user".role WHERE name IN ('user', 'admin');

COMMIT;
