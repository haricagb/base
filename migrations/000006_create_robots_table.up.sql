-- migrations/000006_create_robots_table.up.sql

CREATE TABLE IF NOT EXISTS robots (
    id                  BIGSERIAL       PRIMARY KEY,
    serial_number       VARCHAR(100)    NOT NULL UNIQUE,
    enterprise_id       BIGINT          NOT NULL REFERENCES enterprises(id) ON DELETE CASCADE,
    assigned_resident_id BIGINT         REFERENCES residents(id) ON DELETE SET NULL,
    status              VARCHAR(20)     NOT NULL DEFAULT 'provisioned'
                                        CHECK (status IN ('provisioned', 'active', 'idle', 'maintenance', 'decommissioned')),
    firmware_version    VARCHAR(50),
    last_heartbeat      TIMESTAMPTZ,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_robots_enterprise_id ON robots (enterprise_id);
CREATE INDEX IF NOT EXISTS idx_robots_assigned_resident_id ON robots (assigned_resident_id);
CREATE INDEX IF NOT EXISTS idx_robots_status ON robots (status);
CREATE INDEX IF NOT EXISTS idx_robots_serial_number ON robots (serial_number);

CREATE TRIGGER set_robots_updated_at
    BEFORE UPDATE ON robots
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
