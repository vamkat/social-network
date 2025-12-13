------------------------------------------------------------
-- 1. Conversations
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS conversations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id BIGINT, -- In users service; NULL => DM
    last_message_id BIGINT,
    first_message_id BIGINT,
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
CREATE INDEX idx_messages_conversation_id_id
    ON messages(conversation_id, id);


--------------------------------------------------------------
-- 2.1 ALTER CONVERSATIONS
--------------------------------------------------------------

ALTER TABLE conversations
ADD CONSTRAINT conversations_last_message_id_fkey
    FOREIGN KEY (last_message_id)
    REFERENCES messages(id)
    ON DELETE SET NULL;

ALTER TABLE conversations
    ADD CONSTRAINT conversations_first_message_id_fkey
        FOREIGN KEY (first_message_id)
        REFERENCES messages(id)
        ON DELETE SET NULL;




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
-- UPDATE TIMESTAMP
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


-- UPDATE CONVERSATION FIRST AND LAST MESSAGE

CREATE OR REPLACE FUNCTION update_conversation_first_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE conversations
    SET first_message_id = NEW.id
    WHERE id = NEW.conversation_id
      AND first_message_id IS NULL;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_first_message
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION update_conversation_first_message();

CREATE OR REPLACE FUNCTION update_conversation_last_message()
RETURNS trigger AS $$
BEGIN
    UPDATE conversations
       SET last_message_id = NEW.id,
           updated_at = NEW.created_at
     WHERE id = NEW.conversation_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_update_conversation_last_message
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION update_conversation_last_message();