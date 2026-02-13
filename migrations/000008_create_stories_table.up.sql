-- migrations/000008_create_stories_table.up.sql

CREATE TABLE IF NOT EXISTS stories (
    id                  BIGSERIAL       PRIMARY KEY,
    resident_id         BIGINT          NOT NULL REFERENCES residents(id) ON DELETE CASCADE,
    author_id           BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_encrypted   TEXT            NOT NULL,
    media_url           TEXT,
    story_type          VARCHAR(30)     NOT NULL DEFAULT 'memory'
                                        CHECK (story_type IN ('memory', 'photo', 'video', 'audio', 'general')),
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stories_resident_id ON stories (resident_id);
CREATE INDEX IF NOT EXISTS idx_stories_author_id ON stories (author_id);

CREATE TRIGGER set_stories_updated_at
    BEFORE UPDATE ON stories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
