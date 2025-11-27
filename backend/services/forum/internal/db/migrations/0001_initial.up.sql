-- Enable extensions
CREATE EXTENSION IF NOT EXISTS citext;

------------------------------------------
-- Master Index
------------------------------------------
CREATE TYPE content_type AS ENUM ('post','comment','event');

CREATE TABLE IF NOT EXISTS master_index (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content_type content_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_master_type ON master_index(content_type);

------------------------------------------
-- Posts
------------------------------------------
CREATE TYPE intended_audience AS ENUM ('everyone','followers','selected','group');

CREATE TABLE IF NOT EXISTS posts (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    post_body TEXT NOT NULL,
    creator_id BIGINT NOT NULL, -- in user service
    group_id BIGINT, -- in user service, null for user posts
    audience intended_audience NOT NULL DEFAULT 'everyone',
    comments_count INT DEFAULT 0,
    reactions_count INT DEFAULT 0,
    last_commented_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
); --image here or always join?

CREATE INDEX idx_posts_creator ON posts(creator_id);
CREATE INDEX idx_posts_group ON posts(group_id);
CREATE INDEX idx_posts_audience_created ON posts(audience, created_at DESC);

------------------------------------------
-- Post_audience (for 'selected' audience)
------------------------------------------
CREATE TABLE IF NOT EXISTS post_audience (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    allowed_user_id BIGINT, -- in user service
    PRIMARY KEY (post_id, allowed_user_id)
);

------------------------------------------
-- Feed_entries
------------------------------------------
CREATE TABLE IF NOT EXISTS feed_entries ( --check how to update lazily
    user_id BIGINT NOT NULL, -- in user service
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    seen BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY(user_id, post_id)
);

CREATE INDEX idx_feed_user_created ON feed_entries(user_id, created_at DESC);

------------------------------------------
-- Comments
------------------------------------------
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    comment_creator_id BIGINT NOT NULL, -- in users service
    parent_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment_body TEXT NOT NULL,
    reactions_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_comments_parent_created ON comments(parent_post_id, created_at DESC);
CREATE INDEX idx_comments_creator ON comments(comment_creator_id);

------------------------------------------
-- Events
------------------------------------------
CREATE TABLE IF NOT EXISTS events (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    event_title TEXT NOT NULL,
    event_body TEXT NOT NULL,
    event_creator_id BIGINT NOT NULL, -- in users service
    group_id BIGINT NOT NULL, -- in user service
    event_date DATE NOT NULL,
    still_valid BOOLEAN DEFAULT TRUE, --do we need this?
    going_count INT DEFAULT 0,
    not_going_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_creator ON events(event_creator_id);
CREATE INDEX idx_events_date ON events(event_date);

------------------------------------------
-- Event responses
------------------------------------------
CREATE TABLE IF NOT EXISTS event_responses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL, -- in users service
    going BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT ux_event_user UNIQUE (event_id, user_id)
);

CREATE INDEX idx_event_responses_event ON event_responses(event_id);

------------------------------------------
-- Images
------------------------------------------
CREATE TABLE IF NOT EXISTS images (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    file_name TEXT NOT NULL,
    entity_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
    sort_order INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_image_sort_order UNIQUE(entity_id, sort_order)
);

CREATE INDEX idx_images_entity ON images(entity_id);

------------------------------------------
-- Reactions (likes only)
------------------------------------------
CREATE TABLE IF NOT EXISTS reactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL, -- in users service
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_user_reaction_per_content UNIQUE (user_id, content_id)
);

CREATE INDEX idx_reactions_content ON reactions(content_id);
CREATE INDEX idx_reactions_user ON reactions(user_id);


------------------------------------------
-- Trigger to auto insert to master_index
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

CREATE TRIGGER trg_before_insert_post
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION add_to_master_index('post');

CREATE TRIGGER trg_before_insert_comment
BEFORE INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION add_to_master_index('comment');

CREATE TRIGGER trg_before_insert_event
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

CREATE TRIGGER trg_update_post_modtime
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_comment_modtime
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_event_modtime
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_event_responses_modtime
BEFORE UPDATE ON event_responses
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_reactions_modtime
BEFORE UPDATE ON reactions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_images_modtime
BEFORE UPDATE ON images
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

------------------------------------------
-- Trigger to maintain comments_count and last_commented_at
------------------------------------------
CREATE OR REPLACE FUNCTION update_post_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE posts
        SET comments_count = comments_count + 1,
            last_commented_at = NEW.created_at
        WHERE id = NEW.parent_post_id;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE posts
        SET comments_count = comments_count - 1,
            last_commented_at = (SELECT MAX(created_at) 
                                 FROM comments 
                                 WHERE parent_post_id = OLD.parent_post_id 
                                   AND deleted_at IS NULL)
        WHERE id = OLD.parent_post_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_comments_insert
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION update_post_comments_count();

CREATE TRIGGER trg_comments_delete
AFTER DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_post_comments_count();


------------------------------------------
-- Trigger to maintain reactions_count
------------------------------------------
CREATE OR REPLACE FUNCTION update_reactions_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Update posts reactions_count
    IF EXISTS (SELECT 1 FROM posts WHERE id = NEW.content_id) THEN
        IF TG_OP = 'INSERT' THEN
            UPDATE posts
            SET reactions_count = reactions_count + 1
            WHERE id = NEW.content_id;
        ELSIF TG_OP = 'DELETE' THEN
            UPDATE posts
            SET reactions_count = reactions_count - 1
            WHERE id = OLD.content_id;
        END IF;

    -- Update comments reactions_count
    ELSIF EXISTS (SELECT 1 FROM comments WHERE id = NEW.content_id) THEN
        IF TG_OP = 'INSERT' THEN
            UPDATE comments
            SET reactions_count = reactions_count + 1
            WHERE id = NEW.content_id;
        ELSIF TG_OP = 'DELETE' THEN
            UPDATE comments
            SET reactions_count = reactions_count - 1
            WHERE id = OLD.content_id;
        END IF;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Attach the triggers
CREATE TRIGGER trg_reactions_insert
AFTER INSERT ON reactions
FOR EACH ROW
EXECUTE FUNCTION update_reactions_count();

CREATE TRIGGER trg_reactions_delete
AFTER DELETE ON reactions
FOR EACH ROW
EXECUTE FUNCTION update_reactions_count();


------------------------------------------
-- Soft delete cascade for posts
------------------------------------------
CREATE OR REPLACE FUNCTION soft_delete_post_cascade()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at := CURRENT_TIMESTAMP;

    UPDATE comments SET deleted_at = CURRENT_TIMESTAMP WHERE parent_post_id = OLD.id;
    UPDATE reactions SET deleted_at = CURRENT_TIMESTAMP WHERE content_id = OLD.id;
    UPDATE feed_entries SET deleted_at = CURRENT_TIMESTAMP WHERE post_id = OLD.id;
    UPDATE images SET deleted_at = CURRENT_TIMESTAMP WHERE entity_id = OLD.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_soft_delete_post
BEFORE UPDATE ON posts
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION soft_delete_post_cascade();

------------------------------------------
-- Event responses count trigger
------------------------------------------
CREATE OR REPLACE FUNCTION update_event_response_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        IF NEW.going THEN
            UPDATE events SET going_count = going_count + 1 WHERE id = NEW.event_id;
        ELSE
            UPDATE events SET not_going_count = not_going_count + 1 WHERE id = NEW.event_id;
        END IF;

    ELSIF (TG_OP = 'UPDATE') THEN
        IF OLD.going <> NEW.going THEN
            IF NEW.going THEN
                UPDATE events
                SET going_count = going_count + 1,
                    not_going_count = not_going_count - 1
                WHERE id = NEW.event_id;
            ELSE
                UPDATE events
                SET going_count = going_count - 1,
                    not_going_count = not_going_count + 1
                WHERE id = NEW.event_id;
            END IF;
        END IF;

        -- Handle soft-delete / restore
        IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
            IF OLD.going THEN
                UPDATE events SET going_count = going_count - 1 WHERE id = NEW.event_id;
            ELSE
                UPDATE events SET not_going_count = not_going_count - 1 WHERE id = NEW.event_id;
            END IF;
        ELSIF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
            IF OLD.going THEN
                UPDATE events SET going_count = going_count + 1 WHERE id = NEW.event_id;
            ELSE
                UPDATE events SET not_going_count = not_going_count + 1 WHERE id = NEW.event_id;
            END IF;
        END IF;

    ELSIF (TG_OP = 'DELETE') THEN
        IF OLD.deleted_at IS NULL THEN
            IF OLD.going THEN
                UPDATE events SET going_count = going_count - 1 WHERE id = OLD.event_id;
            ELSE
                UPDATE events SET not_going_count = not_going_count - 1 WHERE id = OLD.event_id;
            END IF;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_event_responses_counts_insert
AFTER INSERT ON event_responses
FOR EACH ROW
EXECUTE FUNCTION update_event_response_counts();

CREATE TRIGGER trg_event_responses_counts_update
AFTER UPDATE ON event_responses
FOR EACH ROW
EXECUTE FUNCTION update_event_response_counts();

CREATE TRIGGER trg_event_responses_counts_delete
AFTER DELETE ON event_responses
FOR EACH ROW
EXECUTE FUNCTION update_event_response_counts();

------------------------------------------
-- Images sort_order trigger
------------------------------------------
CREATE OR REPLACE FUNCTION set_next_sort_order()
RETURNS TRIGGER AS $$
DECLARE
    max_order INT;
BEGIN
    IF NEW.sort_order IS NULL THEN
        SELECT COALESCE(MAX(sort_order),0) 
        INTO max_order
        FROM images
        WHERE entity_id = NEW.entity_id
        FOR UPDATE;

        NEW.sort_order := max_order + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_sort_order
BEFORE INSERT ON images
FOR EACH ROW
EXECUTE FUNCTION set_next_sort_order();