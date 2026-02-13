-- migrations/000003_create_enterprises_table.down.sql

DROP TRIGGER IF EXISTS set_enterprises_updated_at ON enterprises;
DROP TABLE IF EXISTS enterprises;
