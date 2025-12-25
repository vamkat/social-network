package dbservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to cleanup conversation by id
func cleanupConversation(t *testing.T, ctx context.Context, convID int64) {
	_, _ = testPool.Exec(ctx, "DELETE FROM messages WHERE conversation_id = $1", convID)
	_, _ = testPool.Exec(ctx, "DELETE FROM conversation_members WHERE conversation_id = $1", convID)
	_, _ = testPool.Exec(ctx, "DELETE FROM conversations WHERE id = $1", convID)
}

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	// Setup test database connection
	ctx := context.Background()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/chat_test"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	testPool = pool

	code := m.Run()

	// Cleanup (optional: truncate tables)
	truncateTestTables(ctx, pool)

	os.Exit(code)
}

func truncateTestTables(ctx context.Context, pool *pgxpool.Pool) {
	tables := []string{
		"messages",
		"conversation_members",
		"conversations",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			fmt.Printf("Failed to truncate %s: %v\n", table, err)
		}
	}
}

func TestCreatePrivateConv_AddMembers_GetConversationMembers(t *testing.T) {
	ctx := context.Background()
	q := New(testPool)

	userA := ct.Id(1001)
	userB := ct.Id(1002)

	// Create private conversation
	convId, err := q.CreatePrivateConv(ctx, md.CreatePrivateConvParams{UserA: userA, UserB: userB})
	require.NoError(t, err)
	require.True(t, convId > 0)

	// Add members
	err = q.AddConversationMembers(ctx, md.AddConversationMembersParams{ConversationId: convId, UserIds: ct.Ids{userA, userB}})
	require.NoError(t, err)

	// Get conversation members from perspective of userA (should return userB)
	members, err := q.GetConversationMembers(ctx, md.GetConversationMembersParams{ConversationId: convId, UserID: userA})
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, ct.Id(userB), members[0])

	cleanupConversation(t, ctx, int64(convId))
}

func TestCreateGroupConv_AddMembersToGroupConversation_GetConversationMembers(t *testing.T) {
	ctx := context.Background()
	q := New(testPool)

	groupId := ct.Id(2001)
	// Create group conversation
	convId, err := q.CreateGroupConv(ctx, groupId)
	require.NoError(t, err)
	require.True(t, convId > 0)

	// Add members to group conversation
	users := ct.Ids{ct.Id(3001), ct.Id(3002), ct.Id(3003)}
	convFromAdd, err := q.AddMembersToGroupConversation(ctx, md.AddMembersToGroupConversationParams{GroupId: groupId, UserIds: users})
	require.NoError(t, err)
	require.True(t, convFromAdd > 0)

	// Use one member to fetch the other members
	members, err := q.GetConversationMembers(ctx, md.GetConversationMembersParams{ConversationId: convFromAdd, UserID: users[0]})
	require.NoError(t, err)
	// members should include the other two users
	// length may be 2 since GetConversationMembers excludes the caller
	assert.Len(t, members, 2)

	cleanupConversation(t, ctx, int64(convFromAdd))
}

func TestGetUserConversations_Basic(t *testing.T) {
	ctx := context.Background()
	q := New(testPool)

	// Ensure schema has expected columns (some migrations may omit these in older DBs)
	_, _ = testPool.Exec(ctx, `ALTER TABLE conversations ADD COLUMN IF NOT EXISTS last_message_id BIGINT REFERENCES messages(id)`)
	_, _ = testPool.Exec(ctx, `ALTER TABLE conversation_members ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP`)

	// Create a DM and add members
	userA := ct.Id(4001)
	userB := ct.Id(4002)
	convId, err := q.CreatePrivateConv(ctx, md.CreatePrivateConvParams{UserA: userA, UserB: userB})
	require.NoError(t, err)
	require.True(t, convId > 0)
	err = q.AddConversationMembers(ctx, md.AddConversationMembersParams{ConversationId: convId, UserIds: ct.Ids{userA, userB}})
	require.NoError(t, err)

	// Fetch user conversations for userA (groupId zero => DM)
	rows, err := q.GetUserConversations(ctx, md.GetUserConversationsParams{UserId: userA, GroupId: ct.Id(0), Limit: ct.Limit(10), Offset: ct.Offset(0)})
	require.NoError(t, err)
	// Expect at least one conversation
	require.GreaterOrEqual(t, len(rows), 1)

	cleanupConversation(t, ctx, int64(convId))
}

func TestUpdateLastReadMessage(t *testing.T) {
	ctx := context.Background()
	q := New(testPool)

	// Ensure conversation_members has updated_at so triggers won't fail
	_, _ = testPool.Exec(ctx, `ALTER TABLE conversation_members ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP`)

	// Create conversation and members
	userA := ct.Id(5001)
	userB := ct.Id(5002)
	convId, err := q.CreatePrivateConv(ctx, md.CreatePrivateConvParams{UserA: userA, UserB: userB})
	require.NoError(t, err)
	require.True(t, convId > 0)
	err = q.AddConversationMembers(ctx, md.AddConversationMembersParams{ConversationId: convId, UserIds: ct.Ids{userA, userB}})
	require.NoError(t, err)

	// Create a message as userA
	msg, err := q.CreateMessage(ctx, md.CreateMessageParams{ConversationId: convId, SenderId: userA, MessageText: "hello"})
	require.NoError(t, err)

	// Update last read for userB
	convMember, err := q.UpdateLastReadMessage(ctx, md.UpdateLastReadMessageParams{ConversationId: convId, UserID: userB, LastReadMessageId: msg.Id})
	require.NoError(t, err)
	assert.Equal(t, msg.Id, convMember.LastReadMessageId.Int64)

	cleanupConversation(t, ctx, int64(convId))
}
