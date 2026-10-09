BEGIN;

DELETE FROM role WHERE name IN ('user', 'admin');

COMMIT;
