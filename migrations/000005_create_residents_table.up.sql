-- migrations/000005_create_residents_table.up.sql

CREATE TABLE IF NOT EXISTS residents (
    id                      BIGSERIAL       PRIMARY KEY,
    enterprise_id           BIGINT          NOT NULL REFERENCES enterprises(id) ON DELETE CASCADE,
    full_name               VARCHAR(200)    NOT NULL,
    date_of_birth           DATE,
    room_number             VARCHAR(20),
    care_level              VARCHAR(20)     NOT NULL DEFAULT 'standard'
                                            CHECK (care_level IN ('low', 'standard', 'high', 'critical')),
    medical_notes_encrypted TEXT,
    is_active               BOOLEAN         NOT NULL DEFAULT true,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_residents_enterprise_id ON residents (enterprise_id);
CREATE INDEX IF NOT EXISTS idx_residents_is_active ON residents (is_active);

CREATE TRIGGER set_residents_updated_at
    BEFORE UPDATE ON residents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
