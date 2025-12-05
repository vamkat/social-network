-- Drop triggers
DROP TRIGGER IF EXISTS trg_soft_delete_conversation_member ON conversation_members;
DROP TRIGGER IF EXISTS trg_soft_delete_message ON messages;
DROP TRIGGER IF EXISTS trg_soft_delete_conversation ON conversations;
DROP TRIGGER IF EXISTS trg_update_conversation_members_modtime ON conversation_members;
DROP TRIGGER IF EXISTS trg_update_messages_modtime ON messages;
DROP TRIGGER IF EXISTS trg_update_conversations_modtime ON conversations;

-- Drop trigger functions
DROP FUNCTION IF EXISTS soft_delete_conversation_member();
DROP FUNCTION IF EXISTS soft_delete_message();
DROP FUNCTION IF EXISTS soft_delete_conversation();
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop tables
DROP TABLE IF EXISTS conversation_members;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
