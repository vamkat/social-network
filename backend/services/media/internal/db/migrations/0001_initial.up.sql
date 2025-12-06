-- 001_create_images_table.sql

CREATE TABLE images (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    original_name TEXT NOT NULL,
    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL UNIQUE,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    width INT,
    height INT,
    checksum TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Auto-update timestamp
CREATE TRIGGER update_images_updated_at
BEFORE UPDATE ON images
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

