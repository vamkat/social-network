-- Enable citext extension for case-insensitive text
CREATE EXTENSION IF NOT EXISTS citext;

CREATE COLLATION IF NOT EXISTS case_insensitive_ai (
  provider = icu,
  locale = 'und-u-ks-level1',
  deterministic = false
);


-- Schema migrations table
CREATE TABLE IF NOT EXISTS schema_migrations (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    version TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--TODO check what actually needs to be in master index
-- Master index table
CREATE TABLE master_index (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content_type TEXT NOT NULL CHECK (content_type IN ('user', 'post', 'comment', 'group', 'message', 'event')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_master_type ON master_index(content_type);


-- Auth users table
CREATE TABLE auth_user (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    identifier CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    last_login_at TIMESTAMPTZ,
    CONSTRAINT ux_auth_user_user UNIQUE (user_id),
    CONSTRAINT ux_auth_user_identifier UNIQUE (identifier),
    CHECK (identifier = btrim(identifier))
);


-- Users table
CREATE TABLE users (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),  
    username CITEXT COLLATE case_insensitive_ai UNIQUE,
    email CITEXT COLLATE case_insensitive_ai UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    -- active INTEGER NOT NULL DEFAULT 1,
    date_of_birth DATE NOT NULL,
    avatar VARCHAR(255),
    about_me TEXT, 
    profile_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CONSTRAINT email_or_username_required CHECK (
        (username IS NOT NULL AND username <> '')
        OR (email IS NOT NULL AND email <> '')
    )
);

    CREATE UNIQUE INDEX idx_users_public_id ON users(public_id);

--follows table
CREATE TABLE follows (
    follower_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id)
);

--follow requests table
CREATE TABLE follow_requests (
    requester_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    target_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (requester_id, target_id)
);

-- Groups table -- TODO keep owner here or in members?
CREATE TABLE groups (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    group_owner BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_title TEXT NOT NULL,
    group_description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

-- Group members table
CREATE TABLE group_members (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role TEXT DEFAULT 'member' CHECK (role IN ('member', 'admin', 'owner')),
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

-- group join requests table
CREATE TABLE group_join_requests (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

-- group invites table
CREATE TABLE group_invites (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    sender_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'declined', 'expired')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, receiver_id)
);


CREATE TYPE post_visibility AS ENUM ('public', 'almost_private', 'private');

-- Posts table
CREATE TABLE posts (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    post_title TEXT NOT NULL,
    post_body TEXT NOT NULL,
    post_creator BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    visibility post_visibility NOT NULL DEFAULT 'public',
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE comments (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    comment_creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    parent_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    comment_body TEXT NOT NULL,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Events table
CREATE TABLE events (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    event_title TEXT NOT NULL,
    event_body TEXT NOT NULL,
    event_creator BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    event_date DATE NOT NULL,
    still_valid BOOLEAN NOT NULL DEFAULT TRUE,
    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Event response table
CREATE TABLE event_response (
     id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
     event_id BIGINT REFERENCES events(id) ON DELETE CASCADE,
     user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     going BOOLEAN,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     CONSTRAINT ux_event_user UNIQUE (event_id, user_id)
);

-- Images table
CREATE TABLE images (
     id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
     file_name TEXT,
     entity_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Conversations table 
--TODO check if better as two separate tables
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    dm BOOLEAN NOT NULL,
    group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT dm_group_constraint
        CHECK (
            (dm = TRUE AND group_id IS NULL) OR
            (dm = FALSE AND group_id IS NOT NULL)
        )
);

-- Conversation members table
CREATE TABLE conversation_members (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT conversation_member_unique UNIQUE (conversation_id, user_id)
);
CREATE INDEX idx_conversation_member_conversation ON conversation_member(conversation_id);
CREATE INDEX idx_conversation_member_user ON conversation_member(user_id);


-- Messages table
CREATE TABLE messages (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    message_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivered BOOLEAN NOT NULL DEFAULT TRUE,
    edited_at TIMESTAMPTZ
);


-- Reactions table
CREATE TABLE reactions (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    content_id BIGINT NOT NULL REFERENCES master_index(id) ON DELETE CASCADE,
    reaction_type TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_reaction_per_content UNIQUE (user_id, content_id, reaction_type)
);

-- Reaction details table
CREATE TABLE reaction_details (
    id BIGINT PRIMARY KEY REFERENCES master_index(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);



