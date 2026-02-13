-- migrations/000004_alter_users_for_sona.down.sql

DROP INDEX IF EXISTS idx_users_firebase_uid;
DROP INDEX IF EXISTS idx_users_enterprise_id;

-- Revert SONA roles back to original roles.
UPDATE users SET role = 'admin' WHERE role = 'mta';
UPDATE users SET role = 'operator' WHERE role = 'eta';
UPDATE users SET role = 'viewer' WHERE role IN ('caregiver', 'family', 'robot');

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
    CHECK (role IN ('admin', 'operator', 'viewer'));

ALTER TABLE users DROP COLUMN IF EXISTS firebase_uid;
ALTER TABLE users DROP COLUMN IF EXISTS enterprise_id;
