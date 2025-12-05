------------------------------------------------------------
-- 1. Conversations
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS conversations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id BIGINT, -- In users service; NULL => DM
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- One conversation per group (group chat)
CREATE UNIQUE INDEX IF NOT EXISTS uq_conversations_group_id
    ON conversations(group_id)
    WHERE group_id IS NOT NULL;


------------------------------------------------------------
-- 2. Messages
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id BIGINT,
    message_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);


------------------------------------------------------------
-- 3. Conversation Members
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS conversation_members (
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    last_read_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_conversation_members_user ON conversation_members(user_id);
CREATE INDEX idx_conversation_members_last_read ON conversation_members(last_read_message_id);


------------------------------------------------------------
-- 5. Triggers: Auto-update updated_at
------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_conversations_modtime
BEFORE UPDATE ON conversations
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_messages_modtime
BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_update_conversation_members_modtime
BEFORE UPDATE ON conversation_members
FOR EACH ROW EXECUTE FUNCTION update_timestamp();


------------------------------------------------------------
-- 6. Soft Delete Logic
------------------------------------------------------------
CREATE OR REPLACE FUNCTION soft_delete_conversation()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at := CURRENT_TIMESTAMP;

    UPDATE messages 
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE conversation_id = OLD.id;

    UPDATE conversation_members
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE conversation_id = OLD.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_soft_delete_conversation
BEFORE UPDATE ON conversations
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION soft_delete_conversation();


-- Soft delete message
CREATE OR REPLACE FUNCTION soft_delete_message()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_soft_delete_message
BEFORE UPDATE ON messages
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION soft_delete_message();


-- Soft delete conversation member
CREATE OR REPLACE FUNCTION soft_delete_conversation_member()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_soft_delete_conversation_member
BEFORE UPDATE ON conversation_members
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION soft_delete_conversation_member();
