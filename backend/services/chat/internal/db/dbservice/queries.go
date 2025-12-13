package dbservice

const (
	// CONVERSATIONS
	createGroupConv = `
	INSERT INTO conversations (group_id)
	VALUES ($1)
	RETURNING id
	`

	addConversationMembers = `
	INSERT INTO conversation_members (conversation_id, user_id, last_read_message_id)
	SELECT $1, UNNEST($2::bigint[]), NULL
	`

	createPrivateConv = `
	WITH existing AS (
		SELECT c.id
		FROM conversations c
		JOIN conversation_members cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
		JOIN conversation_members cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
		WHERE c.group_id IS NULL
	)
	INSERT INTO conversations (group_id)
	SELECT NULL
	WHERE NOT EXISTS (SELECT 1 FROM existing)
	RETURNING id
	`

	deleteConversationByExactMembers = `
	WITH target_members AS (
		SELECT unnest($1::bigint[]) AS user_id
	),
	matched_conversation AS (
		SELECT cm.conversation_id
		FROM conversation_members cm
		JOIN target_members tm ON tm.user_id = cm.user_id
		GROUP BY cm.conversation_id
		HAVING 
			-- same count of overlapping members
			COUNT(*) = (SELECT COUNT(*) FROM target_members)
			-- and the conversation has no extra members
			AND COUNT(*) = (
				SELECT COUNT(*) 
				FROM conversation_members cm2 
				WHERE cm2.conversation_id = cm.conversation_id
				AND cm2.deleted_at IS NULL
			)
	)
	UPDATE conversations c
	SET deleted_at = NOW(),
		updated_at = NOW()
	WHERE c.id = (SELECT conversation_id FROM matched_conversation)
	RETURNING id, group_id, created_at, updated_at, deleted_at
	`

	getUserConversations = `
	WITH user_conversations AS (
		SELECT c.id AS conversation_id,
			c.created_at,
			c.updated_at,
			c.last_message_id,
			cm.last_read_message_id
		FROM conversations c
		JOIN conversation_members cm
			ON cm.conversation_id = c.id
		WHERE cm.user_id = $1
			AND cm.deleted_at IS NULL
			AND c.group_id IS NOT DISTINCT FROM $4
		ORDER BY c.last_message_id DESC
		LIMIT $2 OFFSET $3
	),

	member_list AS (
		SELECT uc.conversation_id,
			json_agg(cm.user_id) FILTER (WHERE cm.user_id != $1) AS member_ids
		FROM user_conversations uc
		JOIN conversation_members cm
			ON cm.conversation_id = uc.conversation_id
		GROUP BY uc.conversation_id
	),

	unread AS (
		SELECT uc.conversation_id,
			COUNT(m.id) AS unread_count,
			MIN(m.id) AS first_unread_message_id
		FROM user_conversations uc
		LEFT JOIN messages m
		ON m.conversation_id = uc.conversation_id
		AND m.id > COALESCE(uc.last_read_message_id, 0)
		AND m.deleted_at IS NULL
		GROUP BY uc.conversation_id
	)

	SELECT
		uc.conversation_id,
		uc.created_at,
		uc.updated_at,
		ml.member_ids,
		u.unread_count,
		u.first_unread_message_id
	FROM user_conversations uc
	JOIN member_list ml ON ml.conversation_id = uc.conversation_id
	LEFT JOIN unread u ON u.conversation_id = uc.conversation_id
	ORDER BY uc.last_message_id DESC;
	`

	// GROUP CONVERSATIONS
	addMembersToGroupConversation = `
	WITH convo AS (
		SELECT id
		FROM conversations
		WHERE group_id = $1
		AND deleted_at IS NULL
	),
	insert_members AS (
		INSERT INTO conversation_members (conversation_id, user_id)
		SELECT (SELECT id FROM convo), unnest($2::bigint[])
		ON CONFLICT (conversation_id, user_id) DO NOTHING
		RETURNING conversation_id
	)
	SELECT id FROM convo
`

	// MEMBERS
	getConversationMembers = `
	SELECT cm2.user_id
	FROM conversation_members cm1
	JOIN conversation_members cm2
	ON cm2.conversation_id = cm1.conversation_id
	WHERE cm1.user_id = $2
		AND cm2.conversation_id = $1
		AND cm2.user_id <> $2
		AND cm2.deleted_at IS NULL
	`

	deleteConversationMember = `
	UPDATE conversation_members to_delete
	SET deleted_at = NOW()
	FROM conversation_members owner
	WHERE to_delete.conversation_id = $1
		AND to_delete.user_id = $2
		AND to_delete.deleted_at IS NULL
		AND owner.conversation_id = $1
		AND owner.user_id = $3
		AND owner.deleted_at IS NULL
	RETURNING to_delete.conversation_id, to_delete.user_id, to_delete.last_read_message_id, cm_target.joined_at, cm_target.deleted_at
	`

	// MESSAGES
	createMessage = `
	INSERT INTO messages (conversation_id, sender_id, message_text)
	SELECT $1, $2, $3
	FROM conversation_members
	WHERE conversation_id = $1
	AND user_id = $2
	AND deleted_at IS NULL
	RETURNING id, conversation_id, sender_id, message_text, created_at, updated_at, deleted_at
`
	getPrevMessages = `
	SELECT 
		m.id,
		m.conversation_id,
		m.sender_id,
		m.message_text,
		m.created_at,
		m.updated_at,
		m.deleted_at,
		c.first_message_id
	FROM messages m
	JOIN conversations c
		ON c.id = m.conversation_id
	JOIN conversation_members cm
		ON cm.conversation_id = m.conversation_id
	WHERE m.conversation_id = $2
	AND cm.user_id = $3
	AND m.deleted_at IS NULL
  	AND (
        ($1 IS NULL AND m.id <= c.last_message_id)  -- inclusive when $1 is null
        OR
        ($1 IS NOT NULL AND m.id < $1)              -- strict when $1 is supplied
    	)
	ORDER BY m.id DESC
	LIMIT $4;
`
	getNextMessages = `
	SELECT 
		m.id,
		m.conversation_id,
		m.sender_id,
		m.message_text,
		m.created_at,
		m.updated_at,
		m.deleted_at,
		c.last_message_id
	FROM messages m
	JOIN conversations c
		ON c.id = m.conversation_id
	JOIN conversation_members cm
		ON cm.conversation_id = m.conversation_id
	WHERE m.conversation_id = $2
	AND cm.user_id = $3
	AND m.deleted_at IS NULL
	AND m.id > $1
	ORDER BY m.id ASC
	LIMIT $4;
`
	updateLastReadMessage = `
	UPDATE conversation_members cm
	SET last_read_message_id = $3
	WHERE cm.conversation_id = $1
	AND cm.user_id = $2
	AND cm.deleted_at IS NULL
	RETURNING conversation_id, user_id, last_read_message_id, joined_at, deleted_at
`
)
