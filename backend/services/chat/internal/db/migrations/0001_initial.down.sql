-- Drop triggers
DROP TRIGGER IF EXISTS trg_update_conversation_members_modtime ON conversation_members;
DROP TRIGGER IF EXISTS trg_update_messages_modtime ON messages;
DROP TRIGGER IF EXISTS trg_update_conversations_modtime ON conversations;
DROP TRIGGER IF EXISTS trg_update_conversation_last_message ON messages;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_timestamp();
DROP FUNCTION IF EXISTS update_conversation_last_message();

-- Drop tables
DROP TABLE IF EXISTS conversation_members;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
