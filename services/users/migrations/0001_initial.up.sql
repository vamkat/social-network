-- Enable extensions
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Case-insensitive collation
CREATE COLLATION IF NOT EXISTS case_insensitive_ai (
  provider = icu,
  locale = 'und-u-ks-level1',
  deterministic = false
);

-----------------------------------------
-- Users table
-----------------------------------------
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    email CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar VARCHAR(255),
    about_me TEXT, 
    profile_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_id ON users(id);

-- Prepopulate system user
INSERT INTO users (
    username, email, first_name, last_name, date_of_birth
) VALUES (
    'system', 'system@example.com', 'System', 'User', '2000-01-01'
);

-----------------------------------------
-- Auth table (one-to-one with users)
-----------------------------------------
CREATE TABLE auth_user (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    last_login_at TIMESTAMPTZ
);

-----------------------------------------
-- Follows
-----------------------------------------
CREATE TABLE follows (
    follower_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);

-----------------------------------------
-- Follow requests
-----------------------------------------
CREATE TYPE follow_request_status AS ENUM ('pending','accepted','rejected');

CREATE TABLE follow_requests (
    requester_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    target_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status follow_request_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (requester_id, target_id)
);

-----------------------------------------
-- Groups & modular settings
-----------------------------------------
CREATE TYPE group_type AS ENUM ('public','custom');

CREATE TABLE group_type_settings (
    group_type group_type PRIMARY KEY,
    chat_enabled BOOLEAN NOT NULL,
    events_enabled BOOLEAN NOT NULL,
    privacy_enabled BOOLEAN NOT NULL,
    about_enabled BOOLEAN NOT NULL
);

-- Prepopulate group type settings
INSERT INTO group_type_settings VALUES
('public', FALSE, FALSE, TRUE, FALSE),
('custom', TRUE, TRUE, FALSE, TRUE);

CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    group_owner BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_title TEXT NOT NULL,
    group_description TEXT NOT NULL,
    group_type group_type NOT NULL DEFAULT 'custom',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_groups_owner ON groups(group_owner);

-- Create default public group owned by system user (id=1)
INSERT INTO groups (group_owner, group_title, group_description, group_type)
VALUES (1, 'General', 'The public group for all users', 'public');

-----------------------------------------
-- Trigger: auto-add new users to public groups
-----------------------------------------
CREATE OR REPLACE FUNCTION add_user_to_public_groups()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO group_members(group_id, user_id, role, joined_at)
    SELECT id, NEW.id, 'member', CURRENT_TIMESTAMP
    FROM groups
    WHERE group_type = 'public';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_user_to_public_groups
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION add_user_to_public_groups();

-----------------------------------------
-- Group members
-----------------------------------------
CREATE TYPE group_role AS ENUM ('member','admin','owner');

CREATE TABLE group_members (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role group_role DEFAULT 'member',
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_members_user ON group_members(user_id);

-----------------------------------------
-- Group join requests
-----------------------------------------
CREATE TYPE join_request_status AS ENUM ('pending','accepted','rejected');

CREATE TABLE group_join_requests (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status join_request_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_join_requests_status ON group_join_requests(status);

-----------------------------------------
-- Group invites
-----------------------------------------
CREATE TYPE group_invite_status AS ENUM ('pending','accepted','declined','expired');

CREATE TABLE group_invites (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    sender_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status group_invite_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, receiver_id)
);

CREATE INDEX idx_group_invites_status ON group_invites(status);
