-- migrations/000008_create_stories_table.down.sql

DROP TRIGGER IF EXISTS set_stories_updated_at ON stories;
DROP TABLE IF EXISTS stories;
