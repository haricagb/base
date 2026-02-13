-- migrations/000005_create_residents_table.down.sql

DROP TRIGGER IF EXISTS set_residents_updated_at ON residents;
DROP TABLE IF EXISTS residents;
