-- migrations/000006_create_robots_table.down.sql

DROP TRIGGER IF EXISTS set_robots_updated_at ON robots;
DROP TABLE IF EXISTS robots;
