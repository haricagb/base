-- migrations/000009_create_incidents_table.down.sql

DROP TRIGGER IF EXISTS set_incidents_updated_at ON incidents;
DROP TABLE IF EXISTS incidents;
