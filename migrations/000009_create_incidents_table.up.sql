-- migrations/000009_create_incidents_table.up.sql

CREATE TABLE IF NOT EXISTS incidents (
    id              BIGSERIAL       PRIMARY KEY,
    resident_id     BIGINT          NOT NULL REFERENCES residents(id) ON DELETE CASCADE,
    reporter_id     BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    severity        VARCHAR(20)     NOT NULL DEFAULT 'low'
                                    CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    description     TEXT            NOT NULL,
    resolution      TEXT,
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_incidents_resident_id ON incidents (resident_id);
CREATE INDEX IF NOT EXISTS idx_incidents_reporter_id ON incidents (reporter_id);
CREATE INDEX IF NOT EXISTS idx_incidents_severity ON incidents (severity);

CREATE TRIGGER set_incidents_updated_at
    BEFORE UPDATE ON incidents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
