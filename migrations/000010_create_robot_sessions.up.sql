-- migrations/000010_create_robot_sessions.up.sql

CREATE TABLE IF NOT EXISTS robot_sessions (
    id                              BIGSERIAL       PRIMARY KEY,
    robot_id                        BIGINT          NOT NULL REFERENCES robots(id) ON DELETE CASCADE,
    resident_id                     BIGINT          NOT NULL REFERENCES residents(id) ON DELETE CASCADE,
    session_type                    VARCHAR(30)     NOT NULL DEFAULT 'conversation'
                                                    CHECK (session_type IN ('conversation', 'activity', 'reminder', 'emergency')),
    started_at                      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    ended_at                        TIMESTAMPTZ,
    conversation_summary_encrypted  TEXT,
    created_at                      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_robot_sessions_robot_id ON robot_sessions (robot_id);
CREATE INDEX IF NOT EXISTS idx_robot_sessions_resident_id ON robot_sessions (resident_id);
CREATE INDEX IF NOT EXISTS idx_robot_sessions_started_at ON robot_sessions (started_at);

CREATE TRIGGER set_robot_sessions_updated_at
    BEFORE UPDATE ON robot_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
