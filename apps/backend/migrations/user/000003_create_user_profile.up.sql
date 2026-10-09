-- Module: user — user_profile (1:1 with users). Holds all Mifflin-St Jeor inputs.
-- Enum-like values use CHECK constraints (erd.md §1 enum strategy).
BEGIN;

CREATE TABLE IF NOT EXISTS user_profile (
    id             UUID PRIMARY KEY,
    user_id        UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    weight_kg      NUMERIC(5,2) NOT NULL,
    height_cm      NUMERIC(5,2) NOT NULL,
    age            INT          NOT NULL CHECK (age > 0 AND age < 150),
    sex            VARCHAR(10)  NOT NULL CHECK (sex IN ('male', 'female')),
    activity_level VARCHAR(20)  NOT NULL CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'active', 'very_active')),
    goal           VARCHAR(10)  NOT NULL CHECK (goal IN ('lose', 'maintain', 'gain')),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

COMMIT;
