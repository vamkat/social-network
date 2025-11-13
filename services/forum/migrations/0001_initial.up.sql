CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

------------------------------------------
-- Stub tables for external references
------------------------------------------
CREATE TABLE ext_users (
    id BIGINT PRIMARY KEY,
    username TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ext_groups (
    id BIGINT PRIMARY KEY,
    title TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

------------------------------------------
-- Master Index
------------------------------------------
CREATE TYPE content_type AS ENUM ('post', 'comment', 'event');

CREATE TABLE master_index (
    id BIGSERIAL PRIMARY KEY,
    content_type content_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_master_type ON master_index(content_type);

------------------------------------------
-- Posts
------------------------------------------
CREATE TYPE post_visibility AS ENUM ('public', 'almost_private', 'private');


CREATE TABLE posts (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    post_title TEXT NOT NULL,
    post_body TEXT NOT NULL,
    post_creator BIGINT NOT NULL REFERENCES ext_users(id),
    group_id BIGINT REFERENCES ext_groups(id),
    visibility post_visibility DEFAULT 'public',
    comments_count INT DEFAULT 0; --auto-updated via triggers
    --reactions_count INT DEFAULT 0; --auto-updated via triggers
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_creator ON posts(post_creator);
CREATE INDEX idx_posts_group ON posts(group_id);

------------------------------------------
-- Comments
------------------------------------------
CREATE TABLE comments (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    comment_creator_id BIGINT NOT NULL REFERENCES ext_users(id), 
    parent_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment_body TEXT NOT NULL,
    --group_id BIGINT NOT NULL REFERENCES ext_groups(id), 
    --reactions_count INT DEFAULT 0; --auto-updated via triggers
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_parent ON comments(parent_id);
CREATE INDEX idx_comments_creator ON comments(comment_creator_id);

------------------------------------------
-- Events
------------------------------------------
CREATE TABLE events (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    event_title TEXT NOT NULL,
    event_body TEXT NOT NULL,
    event_creator BIGINT NOT NULL REFERENCES ext_users(id), 
    group_id BIGINT NOT NULL REFERENCES ext_groups(id), 
    event_date DATE NOT NULL,
    still_valid BOOLEAN  DEFAULT TRUE,
    --reactions_count INT DEFAULT 0; --auto-updated via triggers
    --responses_count INT DEFAULT 0; --auto-updated via triggers
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_creator ON events(event_creator);
CREATE INDEX idx_events_date ON events(event_date);

------------------------------------------
-- Event responses
------------------------------------------
CREATE TABLE event_response (
     id BIGSERIAL PRIMARY KEY,
     event_id BIGINT REFERENCES events(id) ON DELETE CASCADE,
     user_id BIGINT NOT NULL REFERENCES ext_users(id),  
     going BOOLEAN,
     created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
     CONSTRAINT ux_event_user UNIQUE (event_id, user_id)
);

CREATE INDEX idx_event_response_event ON event_response(event_id);

------------------------------------------
-- Images
------------------------------------------
CREATE TABLE images (
     id BIGSERIAL PRIMARY KEY,
     file_name TEXT,
     entity_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_images_entity ON images(entity_id);

------------------------------------------
-- Reactions
------------------------------------------
CREATE TABLE reactions (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
    reaction_type TEXT NOT NULL,
    user_id BIGINT NOT NULL, -- in users service
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_reaction_per_content UNIQUE (user_id, content_id, reaction_type)
);

CREATE INDEX idx_reactions_content ON reactions(content_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);


------------------------------------------
-- Trigger to auto instert to master_index
------------------------------------------
CREATE OR REPLACE FUNCTION add_to_master_index()
RETURNS TRIGGER AS $$
DECLARE
    new_id BIGINT;
BEGIN
    INSERT INTO master_index (content_type)
    VALUES (TG_ARGV[0])
    RETURNING id INTO new_id;
    NEW.id := new_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach triggers to each content table
CREATE TRIGGER before_insert_post
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION add_to_master_index('post');

CREATE TRIGGER before_insert_comment
BEFORE INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION add_to_master_index('comment');

CREATE TRIGGER before_insert_event
BEFORE INSERT ON events
FOR EACH ROW
EXECUTE FUNCTION add_to_master_index('event');

------------------------------------------
-- Trigger to auto-update updated_at timestamps
------------------------------------------

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to main tables
CREATE TRIGGER update_post_modtime
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_comment_modtime
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_event_modtime
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

------------------------------------------
-- Trigger to auto-update comments_count in posts
------------------------------------------

CREATE OR REPLACE FUNCTION update_post_comment_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts
        SET comments_count = comments_count + 1
        WHERE id = NEW.parent_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts
        SET comments_count = GREATEST(comments_count - 1, 0)
        WHERE id = OLD.parent_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_post_comment_count
AFTER INSERT OR DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_post_comment_count();

-- think about adding triggers for reactions and event responses