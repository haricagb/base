-- migrations/000007_create_caregiver_residents.up.sql

CREATE TABLE IF NOT EXISTS caregiver_residents (
    caregiver_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resident_id     BIGINT      NOT NULL REFERENCES residents(id) ON DELETE CASCADE,
    assigned_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (caregiver_id, resident_id)
);

CREATE INDEX IF NOT EXISTS idx_caregiver_residents_caregiver ON caregiver_residents (caregiver_id);
CREATE INDEX IF NOT EXISTS idx_caregiver_residents_resident ON caregiver_residents (resident_id);

-- Enforce maximum 5 residents per caregiver.
CREATE OR REPLACE FUNCTION check_max_residents_per_caregiver()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM caregiver_residents WHERE caregiver_id = NEW.caregiver_id) >= 5 THEN
        RAISE EXCEPTION 'A caregiver cannot be assigned more than 5 residents';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_max_residents
    BEFORE INSERT ON caregiver_residents
    FOR EACH ROW
    EXECUTE FUNCTION check_max_residents_per_caregiver();
