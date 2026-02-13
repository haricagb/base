-- migrations/000003_create_enterprises_table.up.sql

CREATE TABLE IF NOT EXISTS enterprises (
    id              BIGSERIAL       PRIMARY KEY,
    name            VARCHAR(200)    NOT NULL UNIQUE,
    contact_email   VARCHAR(255)    NOT NULL,
    subscription_tier VARCHAR(20)   NOT NULL DEFAULT 'basic'
                                    CHECK (subscription_tier IN ('basic', 'professional', 'enterprise')),
    max_robots      INT             NOT NULL DEFAULT 10,
    is_active       BOOLEAN         NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_enterprises_is_active ON enterprises (is_active);

CREATE TRIGGER set_enterprises_updated_at
    BEFORE UPDATE ON enterprises
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
