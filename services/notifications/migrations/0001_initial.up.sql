------------------------------------------
-- Notification Types
------------------------------------------

CREATE TABLE IF NOT EXISTS notification_types (
    notif_type TEXT PRIMARY KEY,     -- e.g. "new_follow request"
    category TEXT,             -- e.g. "social", "chat", "forum"
    default_enabled BOOLEAN    -- for user notification settings
);

INSERT INTO notification_types (notif_type, category, default_enabled)
VALUES
  ('new_follower', 'social',  TRUE),
  ('follow_request', 'social',  TRUE),
  ('group_invite', 'group',  TRUE),
  ('group_join_request', 'group',  TRUE),
  ('new_event', 'group',  TRUE),
  ('new_message', 'chat',  TRUE),
  ('chat_mention', 'chat',  TRUE),
  ('forum_reply', 'forum', TRUE),
  ('forum_mention', 'forum',  TRUE),
  ('forum_like', 'forum', TRUE);

------------------------------------------
-- Notifications
------------------------------------------

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL, --recipient user id
    notif_type TEXT NOT NULL REFERENCES notification_types(notif_type),
    source_service TEXT NOT NULL CHECK (source_service IN ('users', 'chat', 'forum')),
    source_entity_id BIGINT,
    seen BOOLEAN DEFAULT FALSE,
    read BOOLEAN DEFAULT FALSE,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '30 days'),
    deleted_at TIMESTAMPTZ   
);

CREATE INDEX idx_notifications_user_unread 
ON notifications (user_id, seen, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_notifications_user_created 
ON notifications (user_id, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_notifications_user_type 
ON notifications (user_id, notif_type, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_notifications_source 
ON notifications (source_service, source_entity_id);