-- migrations/000004_alter_users_for_sona.up.sql

-- Add enterprise_id foreign key for multi-tenant isolation.
ALTER TABLE users
    ADD COLUMN enterprise_id BIGINT REFERENCES enterprises(id) ON DELETE SET NULL;

-- Add Firebase UID for external auth provider linkage.
ALTER TABLE users
    ADD COLUMN firebase_uid VARCHAR(128) UNIQUE;

-- Migrate existing roles to SONA roles before changing the constraint.
UPDATE users SET role = 'mta' WHERE role = 'admin';
UPDATE users SET role = 'eta' WHERE role = 'operator';
UPDATE users SET role = 'caregiver' WHERE role = 'viewer';

-- Replace the role CHECK constraint with SONA roles.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('mta', 'eta', 'caregiver', 'family', 'robot'));

CREATE INDEX IF NOT EXISTS idx_users_enterprise_id ON users (enterprise_id);
CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users (firebase_uid);
