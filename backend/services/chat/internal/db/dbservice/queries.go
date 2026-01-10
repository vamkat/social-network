package dbservice

const (
	// ====================================
	// GROUP_CONVERSATIONS
	// ====================================

	createGroupConv = `
	INSERT INTO group_conversations (group_id)
	VALUES ($1)
	ON CONFLICT (group_id)
	DO UPDATE SET updated_at = group_conversations.updated_at
	RETURNING group_id;
	`

	createGroupMessage = `
	INSERT INTO group_messages (group_id, sender_id, message_text)
	SELECT c.group_id, $2, $3
	FROM group_conversations c
	WHERE c.group_id = $1
	AND c.deleted_at IS NULL
	RETURNING
		id,
		group_id,
		sender_id,
		message_text,
		created_at,
		updated_at,
		deleted_at;
	`

	getPrevGroupMsgs = `
    SELECT
        gm.id,
        gm.group_id,
        gm.sender_id,
        gm.message_text,
        gm.created_at,
        gm.updated_at,
        gm.deleted_at

    FROM group_conversations gc
    JOIN group_messages gm
        ON gm.group_id = gc.group_id
    WHERE gc.group_id = $1
      AND gm.deleted_at IS NULL
	  AND gc.deleted_at IS NULL
      AND gm.id < $3
    ORDER BY gm.id DESC
    LIMIT $4;
	`

	getNextGroupMsgs = `
    SELECT
        gm.id,
        gm.group_id,
        gm.sender_id,
        gm.message_text,
        gm.created_at,
        gm.updated_at,
        gm.deleted_at

    FROM group_conversations gc
    JOIN group_messages gm
        ON gm.group_id = gc.group_id
    WHERE gc.group_id = $1
      AND gm.deleted_at IS NULL
	  AND gc.deleted_at IS NULL
      AND gm.id > $3
    ORDER BY gm.id ASC
    LIMIT $4;
	`

	// ====================================
	// PRIVATE_CONVERSATIONS
	// ====================================
	getOrCreatePrivateConv = `
	WITH ins AS (
		INSERT INTO private_conversations (user_a, user_b)
		VALUES (
			LEAST($1::bigint, $2::bigint),
			GREATEST($1::bigint, $2::bigint)
		)
		ON CONFLICT (user_a, user_b) DO NOTHING
		RETURNING *
	)
	SELECT *
	FROM ins
	UNION ALL
	SELECT *
	FROM private_conversations c
	WHERE user_a = LEAST($1, $2)
		AND user_b = GREATEST($1, $2)
		AND c.deleted_at IS NULL 
		AND NOT EXISTS (SELECT 1 FROM ins);
	`

	newPrivateMessage = `
    WITH inserted_message AS (
        INSERT INTO private_messages (conversation_id, sender_id, message_text)
        SELECT
            c.id,
            $2 AS sender_id,
            $3 AS message_text
        FROM private_conversations c
        WHERE c.id = $1
          AND c.deleted_at IS NULL
          AND ($2 = c.user_a OR $2 = c.user_b)
        RETURNING
            id,
            conversation_id,
            sender_id,
            message_text,
            created_at,
            updated_at,
            deleted_at
    ),
    updated_conversation AS (
        UPDATE private_conversations pc
        SET
            last_read_message_id_a = CASE
                WHEN pc.user_a = im.sender_id THEN im.id
                ELSE pc.last_read_message_id_a
            END,
            last_read_message_id_b = CASE
                WHEN pc.user_b = im.sender_id THEN im.id
                ELSE pc.last_read_message_id_b
            END
        FROM inserted_message im
        WHERE pc.id = im.conversation_id
    )
    SELECT * FROM inserted_message;
    `

	getPrivateConvs = `
	WITH user_conversations AS (
		SELECT
			pc.id AS conversation_id,
			pc.updated_at,

			-- determine other user
			CASE
				WHEN pc.user_id_a = $1 THEN pc.user_id_b
				ELSE pc.user_id_a
			END AS other_user_id,

			-- determine last read message for this user
			CASE
				WHEN pc.user_id_a = $1 THEN pc.last_read_message_id_a
				ELSE pc.last_read_message_id_b
			END AS last_read_message_id

		FROM private_conversations pc
		WHERE $1 IN (pc.user_id_a, pc.user_id_b)
		AND pc.updated_at < $2
	)

	SELECT
		uc.conversation_id,
		uc.updated_at,
		uc.other_user_id,

		-- last message
		lm.id           AS last_message_id,
		lm.sender_id    AS last_message_sender_id,
		lm.message_text AS last_message_text,
		lm.created_at   AS last_message_created_at,

		-- unread count
		COUNT(pm.id) FILTER (
			WHERE pm.id > COALESCE(uc.last_read_message_id, 0)
		) AS unread_count

	FROM user_conversations uc

	-- last message per conversation
	LEFT JOIN LATERAL (
		SELECT pm.id, pm.sender_id, pm.message_text, pm.created_at
		FROM private_messages pm
		WHERE pm.conversation_id = uc.conversation_id
		AND pm.deleted_at IS NULL
		ORDER BY pm.id DESC
		LIMIT 1
	) lm ON true

	-- unread messages
	LEFT JOIN private_messages pm
		ON pm.conversation_id = uc.conversation_id
	AND pm.deleted_at IS NULL

	GROUP BY
		uc.conversation_id,
		uc.updated_at,
		uc.other_user_id,
		uc.last_read_message_id,
		lm.id,
		lm.sender_id,
		lm.message_text,
		lm.created_at

	ORDER BY uc.updated_at DESC
	LIMIT $3;
	`

	getPrevPrivateMsgs = `
	SELECT pm.*
	FROM private_conversations pc
	JOIN private_messages pm
	ON pm.conversation_id = pc.id
	WHERE pc.id = $1
	AND $2 IN (pc.user_a, pc.user_b)
	AND pm.deleted_at IS NULL
	AND pm.id < $3
	ORDER BY pm.id DESC
	LIMIT $4;
	`

	getNextPrivateMsgs = `
	SELECT pm.*
	FROM private_conversations pc
	JOIN private_messages pm
	ON pm.conversation_id = pc.id
	WHERE pc.id = $1
	AND $2 IN (pc.user_a, pc.user_b)
	AND pm.deleted_at IS NULL
	AND pm.id > $3
	ORDER BY pm.id ASC
	LIMIT $4;
	`

	updateLastReadMessage = `
	UPDATE private_conversations
	SET
		last_read_message_id_a = CASE
			WHEN user_a = $2 THEN $3
			ELSE last_read_message_id_a
		END,
		last_read_message_id_b = CASE
			WHEN user_b = $2 THEN $3
			ELSE last_read_message_id_b
		END,
	WHERE id = $1
	AND (user_a = $3 OR user_b = $3)
	AND deleted_at IS NULL;
	`
)
