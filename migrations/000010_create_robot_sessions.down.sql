-- migrations/000010_create_robot_sessions.down.sql

DROP TRIGGER IF EXISTS set_robot_sessions_updated_at ON robot_sessions;
DROP TABLE IF EXISTS robot_sessions;
