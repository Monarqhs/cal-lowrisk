-- Module: user — users. Auth identity + role FK + soft delete (admin deactivate, ADM-5).
-- Lives in the "user" schema (reserved word → always double-quoted).
BEGIN;

CREATE TABLE IF NOT EXISTS "user".users (
    id            UUID PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id       UUID NOT NULL REFERENCES "user".role (id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ  -- NULL = active; set = deactivated (soft delete)
);

CREATE INDEX IF NOT EXISTS idx_users_role_id ON "user".users (role_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON "user".users (deleted_at);

COMMIT;
