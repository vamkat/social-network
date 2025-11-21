-- Enable extensions
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Case-insensitive collation
CREATE COLLATION IF NOT EXISTS case_insensitive_ai (
  provider = icu,
  locale = 'und-u-ks-level1',
  deterministic = false
);

-----------------------------------------
-- Users table
-----------------------------------------
CREATE TYPE user_status AS ENUM ('active', 'banned', 'deleted');

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar VARCHAR(255) NOT NULL,
    about_me TEXT NOT NULL, 
    profile_public BOOLEAN NOT NULL DEFAULT TRUE,
    current_status user_status NOT NULL DEFAULT 'active',
    ban_ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_id ON users(id);
CREATE INDEX idx_users_status ON users(current_status);


-----------------------------------------
-- Auth table (one-to-one with users)
-----------------------------------------
CREATE TABLE IF NOT EXISTS auth_user (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    email CITEXT COLLATE case_insensitive_ai UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    last_login_at TIMESTAMPTZ
);


-----------------------------------------
-- Follows
-----------------------------------------
CREATE TABLE IF NOT EXISTS follows (
    follower_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT no_self_follow CHECK (follower_id <> following_id)
);


CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);

-----------------------------------------
-- Follow requests
-----------------------------------------
CREATE TYPE follow_request_status AS ENUM ('pending','accepted','rejected');

CREATE TABLE IF NOT EXISTS follow_requests (
    requester_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    target_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status follow_request_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (requester_id, target_id)
);

CREATE INDEX idx_follow_requests_target_status ON follow_requests(target_id, status);


-----------------------------------------
-- Groups
-----------------------------------------
CREATE TABLE IF NOT EXISTS groups (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_owner BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    group_title TEXT NOT NULL,
    group_description TEXT NOT NULL,
    members_count INT NOT NULL DEFAULT 0, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_groups_owner ON groups(group_owner);

-----------------------------------------
-- Group members
-----------------------------------------
CREATE TYPE group_role AS ENUM ('member','owner');

CREATE TABLE IF NOT EXISTS group_members (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role group_role DEFAULT 'member',
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (group_id, user_id)
);

-- Ensure exactly one owner per group
CREATE UNIQUE INDEX IF NOT EXISTS idx_group_one_owner
ON group_members(group_id)
WHERE role='owner' AND deleted_at IS NULL;

CREATE INDEX idx_group_members_user ON group_members(user_id);


-----------------------------------------
-- Group join requests
-----------------------------------------
CREATE TYPE join_request_status AS ENUM ('pending','accepted','rejected');

CREATE TABLE IF NOT EXISTS group_join_requests (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    status join_request_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_join_requests_status ON group_join_requests(status);


-----------------------------------------
-- Group invites
-----------------------------------------
CREATE TYPE group_invite_status AS ENUM ('pending','accepted','declined','expired');

CREATE TABLE IF NOT EXISTS group_invites (
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status group_invite_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (group_id, receiver_id)
);

CREATE INDEX idx_group_invites_status ON group_invites(status);

-----------------------------------------
-- Trigger to auto-update updated_at timestamps
-----------------------------------------

-- Single trigger function for updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach to multiple tables
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_auth_user_updated_at
BEFORE UPDATE ON auth_user
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_groups_updated_at
BEFORE UPDATE ON groups
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_group_members_updated_at
BEFORE UPDATE ON group_members
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_group_join_requests_updated_at
BEFORE UPDATE ON group_join_requests
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_group_invites_updated_at
BEFORE UPDATE ON group_invites
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_follow_requests_updated_at
BEFORE UPDATE ON follow_requests
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-----------------------------------------
-- Soft delete cascade for users
-----------------------------------------
CREATE OR REPLACE FUNCTION soft_delete_user_cascade()
RETURNS TRIGGER AS $$
BEGIN
    -- Soft-delete follows
    UPDATE follows
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE follower_id = OLD.id OR following_id = OLD.id;

    -- Soft-delete follow requests
    UPDATE follow_requests
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE requester_id = OLD.id OR target_id = OLD.id;

    -- Soft-delete group memberships
    UPDATE group_members
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE user_id = OLD.id;

    -- Soft-delete group join requests
    UPDATE group_join_requests
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE user_id = OLD.id;

    -- Soft-delete group invites (sent or received)
    UPDATE group_invites
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE sender_id = OLD.id OR receiver_id = OLD.id;

    -- Optional: soft-delete owned groups
    UPDATE groups
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE group_owner = OLD.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_soft_delete_user
BEFORE UPDATE ON users
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION soft_delete_user_cascade();

-----------------------------------------
-- Function to follow a user regardless of privacy setting
-----------------------------------------

CREATE OR REPLACE FUNCTION follow_user(p_follower BIGINT, p_target BIGINT)
RETURNS TEXT AS $$
DECLARE
    is_public BOOLEAN;
BEGIN
    SELECT profile_public INTO is_public FROM users WHERE id = p_target;

    IF is_public THEN
        INSERT INTO follows (follower_id, following_id)
        VALUES (p_follower, p_target)
        ON CONFLICT DO NOTHING;

        RETURN 'followed';
    ELSE
        INSERT INTO follow_requests (requester_id, target_id, status)
        VALUES (p_follower, p_target, 'pending')
        ON CONFLICT DO NOTHING;

        RETURN 'requested';
    END IF;
END;
$$ LANGUAGE plpgsql;


-----------------------------------------
-- Trigger to update members count for group
-----------------------------------------
CREATE OR REPLACE FUNCTION update_group_members_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Member added or restored
    IF (TG_OP = 'INSERT')
       OR (TG_OP = 'UPDATE' AND OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL)
    THEN
        UPDATE groups
        SET members_count = (
            SELECT COUNT(*) 
            FROM group_members 
            WHERE group_id = NEW.group_id AND deleted_at IS NULL
        )
        WHERE id = NEW.group_id;
        RETURN NEW;
    END IF;

    -- Member soft-deleted
    IF (TG_OP = 'UPDATE' AND OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
    THEN
        UPDATE groups
        SET members_count = (
            SELECT COUNT(*) 
            FROM group_members 
            WHERE group_id = NEW.group_id AND deleted_at IS NULL
        )
        WHERE id = NEW.group_id;
        RETURN NEW;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- When a member is added
CREATE TRIGGER trg_group_members_count_insert
AFTER INSERT ON group_members
FOR EACH ROW
EXECUTE FUNCTION update_group_members_count();

-- When a member is soft-deleted or restored
CREATE TRIGGER trg_group_members_count_update
AFTER UPDATE ON group_members
FOR EACH ROW
EXECUTE FUNCTION update_group_members_count();

-----------------------------------------
-- Trigger to add follower when follow request is accepted
-----------------------------------------
CREATE OR REPLACE FUNCTION add_follower_on_accept()
RETURNS TRIGGER AS $$
BEGIN
    -- Only act if the request changed to 'accepted'
    IF NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted' THEN
        -- Insert follower, ignore conflicts
        INSERT INTO follows (follower_id, following_id, created_at)
        VALUES (NEW.requester_id, NEW.target_id, CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_follower_on_accept
AFTER UPDATE ON follow_requests
FOR EACH ROW
WHEN (NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted')
EXECUTE FUNCTION add_follower_on_accept();

-----------------------------------------
-- Trigger to accept pending follow requests when a profile changes to public
-----------------------------------------
CREATE OR REPLACE FUNCTION accept_pending_requests_on_public()
RETURNS TRIGGER AS $$
BEGIN
    -- Only act if profile switches from private to public
    IF OLD.profile_public = FALSE AND NEW.profile_public = TRUE THEN
        -- Bulk update: mark all pending requests as accepted
        UPDATE follow_requests
        SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
        WHERE target_id = NEW.id
          AND status = 'pending'
          AND deleted_at IS NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_accept_pending_requests_on_public
AFTER UPDATE ON users
FOR EACH ROW
WHEN (OLD.profile_public = FALSE AND NEW.profile_public = TRUE)
EXECUTE FUNCTION accept_pending_requests_on_public();



-----------------------------------------
-- Trigger to add user as group member when join request accepted
-----------------------------------------
CREATE OR REPLACE FUNCTION add_group_member_on_join_accept()
RETURNS TRIGGER AS $$
BEGIN
    -- Only act when a join request is accepted
    IF NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted' THEN
        INSERT INTO group_members (group_id, user_id, role, joined_at)
        VALUES (NEW.group_id, NEW.user_id, 'member', CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;  -- Already a member? No problem.
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_group_member_on_join_accept
AFTER UPDATE ON group_join_requests
FOR EACH ROW
WHEN (NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted')
EXECUTE FUNCTION add_group_member_on_join_accept();



-----------------------------------------
-- Trigger to add user as group member when group invite accepted
-----------------------------------------
CREATE OR REPLACE FUNCTION add_group_member_on_invite_accept()
RETURNS TRIGGER AS $$
BEGIN
    -- Only act when an invite is accepted
    IF NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted' THEN
        INSERT INTO group_members (group_id, user_id, role, joined_at)
        VALUES (NEW.group_id, NEW.receiver_id, 'member', CURRENT_TIMESTAMP)
        ON CONFLICT DO NOTHING;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_group_member_on_invite_accept
AFTER UPDATE ON group_invites
FOR EACH ROW
WHEN (NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted')
EXECUTE FUNCTION add_group_member_on_invite_accept();


-----------------------------------------
-- Trigger to add group owner as member on group creation
-----------------------------------------
CREATE OR REPLACE FUNCTION add_group_owner_as_member()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO group_members (group_id, user_id, role, joined_at)
    VALUES (NEW.id, NEW.group_owner, 'owner', CURRENT_TIMESTAMP)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_group_owner_as_member
AFTER INSERT ON groups
FOR EACH ROW
EXECUTE FUNCTION add_group_owner_as_member();
