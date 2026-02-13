-- migrations/000007_create_caregiver_residents.down.sql

DROP TRIGGER IF EXISTS enforce_max_residents ON caregiver_residents;
DROP FUNCTION IF EXISTS check_max_residents_per_caregiver();
DROP TABLE IF EXISTS caregiver_residents;
