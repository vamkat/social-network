package application

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"social-network/services/notifications/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation of the database layer for testing
type MockDB struct {
	notifications []sqlc.GetNotificationsByUserIdAllRow
}

func (m *MockDB) GetNotificationsByUserIdAll(ctx context.Context, arg sqlc.GetNotificationsByUserIdAllParams) ([]sqlc.GetNotificationsByUserIdAllRow, error) {
	// Mock implementation - return all notifications for the user
	var result []sqlc.GetNotificationsByUserIdAllRow
	userId := int64(0)
	if arg.Column1.Valid {
		userId = arg.Column1.Int64
	}

	for _, n := range m.notifications {
		if n.UserID == userId {
			result = append(result, n)
		}
	}

	// Sort by CreatedAt DESC (most recent first) as the SQL query does
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].CreatedAt.Time.Before(result[j].CreatedAt.Time) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// Apply limit and offset
	limit := int32(0)
	offset := int32(0)
	if arg.Column2.Valid {
		limit = int32(arg.Column2.Int64)
	}
	if arg.Column3.Valid {
		offset = int32(arg.Column3.Int64)
	}

	start := int(offset)
	end := start + int(limit)
	if start > len(result) {
		start = len(result)
	}
	if end > len(result) {
		end = len(result)
	}

	if start > end {
		start = end
	}

	return result[start:end], nil
}

func (m *MockDB) GetNotificationsByUserIdFilteredByStatus(ctx context.Context, arg sqlc.GetNotificationsByUserIdFilteredByStatusParams) ([]sqlc.GetNotificationsByUserIdFilteredByStatusRow, error) {
	// Mock implementation - return notifications for the user with specified seen status
	var result []sqlc.GetNotificationsByUserIdFilteredByStatusRow
	userId := int64(0)
	if arg.Column1.Valid {
		userId = arg.Column1.Int64
	}

	for _, n := range m.notifications {
		if n.UserID == userId && n.Seen.Bool == arg.Column2.Bool {
			// Convert to the filtered type (they have the same structure for our purposes)
			result = append(result, sqlc.GetNotificationsByUserIdFilteredByStatusRow{
				ID:             n.ID,
				UserID:         n.UserID,
				NotifType:      n.NotifType,
				SourceService:  n.SourceService,
				SourceEntityID: n.SourceEntityID,
				Seen:           n.Seen,
				NeedsAction:    n.NeedsAction,
				Acted:          n.Acted,
				Payload:        n.Payload,
				CreatedAt:      n.CreatedAt,
				ExpiresAt:      n.ExpiresAt,
			})
		}
	}

	// Apply limit and offset
	limit := int32(0)
	offset := int32(0)
	if arg.Column3.Valid {
		limit = int32(arg.Column3.Int64)
	}
	if arg.Column4.Valid {
		offset = int32(arg.Column4.Int64)
	}

	start := int(offset)
	end := start + int(limit)
	if start > len(result) {
		start = len(result)
	}
	if end > len(result) {
		end = len(result)
	}

	if start > end {
		start = end
	}

	return result[start:end], nil
}

func (m *MockDB) GetNotificationCountAll(ctx context.Context, userID pgtype.Int8) (pgtype.Int8, error) {
	userId := int64(0)
	if userID.Valid {
		userId = userID.Int64
	}

	count := int64(0)
	for _, n := range m.notifications {
		if n.UserID == userId {
			count++
		}
	}
	return pgtype.Int8{Int64: count, Valid: true}, nil
}

func (m *MockDB) GetNotificationCountFilteredByStatus(ctx context.Context, arg sqlc.GetNotificationCountFilteredByStatusParams) (pgtype.Int8, error) {
	userId := int64(0)
	if arg.Column1.Valid {
		userId = arg.Column1.Int64
	}

	count := int64(0)
	for _, n := range m.notifications {
		if n.UserID == userId && n.Seen.Bool == arg.Column2.Bool {
			count++
		}
	}
	return pgtype.Int8{Int64: count, Valid: true}, nil
}

func (m *MockDB) MarkNotificationsAsRead(ctx context.Context, arg sqlc.MarkNotificationsAsReadParams) error {
	userId := int64(0)
	if arg.Column1.Valid {
		userId = arg.Column1.Int64
	}

	for i, n := range m.notifications {
		for _, id := range arg.Column2 { // arg.Column2 is the slice of notification IDs
			if n.ID == id && n.UserID == userId {
				m.notifications[i].Seen = pgtype.Bool{Bool: true, Valid: true}
				break
			}
		}
	}
	return nil
}

func (m *MockDB) CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.CreateNotificationRow, error) {
	userId := int64(0)
	if arg.Column1.Valid {
		userId = arg.Column1.Int64
	}

	notifType := ""
	if arg.Column2.Valid {
		notifType = arg.Column2.String
	}

	sourceService := ""
	if arg.Column3.Valid {
		sourceService = arg.Column3.String
	}

	var sourceEntityID pgtype.Int8
	if arg.Column4.Valid {
		sourceEntityID = arg.Column4
	} else {
		sourceEntityID = pgtype.Int8{Valid: false}
	}

	dbNotification := sqlc.CreateNotificationRow{
		ID:             int64(len(m.notifications) + 1),
		UserID:         userId,
		NotifType:      notifType,
		SourceService:  sourceService,
		SourceEntityID: sourceEntityID,
		Seen:           pgtype.Bool{Bool: false, Valid: true},
		NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
		Acted:          pgtype.Bool{Bool: false, Valid: true},
		Payload:        arg.Column5,
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ExpiresAt:      pgtype.Timestamptz{Time: time.Now().AddDate(0, 0, 30), Valid: true}, // 30 days from now
	}

	// Also add to our internal notifications array for testing purposes
	m.notifications = append(m.notifications, sqlc.GetNotificationsByUserIdAllRow{
		ID:             dbNotification.ID,
		UserID:         dbNotification.UserID,
		NotifType:      dbNotification.NotifType,
		SourceService:  dbNotification.SourceService,
		SourceEntityID: dbNotification.SourceEntityID,
		Seen:           dbNotification.Seen,
		NeedsAction:    dbNotification.NeedsAction,
		Acted:          dbNotification.Acted,
		Payload:        dbNotification.Payload,
		CreatedAt:      dbNotification.CreatedAt,
		ExpiresAt:      dbNotification.ExpiresAt,
	})
	return dbNotification, nil
}

func (m *MockDB) MarkNotificationAsActed(ctx context.Context, id pgtype.Int8) error {
	notificationId := int64(0)
	if id.Valid {
		notificationId = id.Int64
	}

	for i, n := range m.notifications {
		if n.ID == notificationId {
			m.notifications[i].Acted = pgtype.Bool{Bool: true, Valid: true}
			return nil
		}
	}
	return nil // Not found, but we'll return nil to match interface
}


func TestGetNotifications(t *testing.T) {
	payload1, _ := json.Marshal(map[string]string{"follower_id": "200"})
	payload2, _ := json.Marshal(map[string]string{"post_id": "300"})

	mockDB := &MockDB{
		notifications: []sqlc.GetNotificationsByUserIdAllRow{
			{
				ID:             1,
				UserID:         100,
				NotifType:      string(NewFollower),
				SourceService:  string(UsersService),
				SourceEntityID: pgtype.Int8{Int64: 200, Valid: true},
				Seen:           pgtype.Bool{Bool: false, Valid: true},
				NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
				Acted:          pgtype.Bool{Bool: false, Valid: true},
				Payload:        payload1,
				CreatedAt:      pgtype.Timestamptz{Time: time.Now().Add(-1 * time.Hour), Valid: true}, // Earlier
				ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(29 * 24 * time.Hour), Valid: true},
			},
			{
				ID:             2,
				UserID:         100,
				NotifType:      string(Like),
				SourceService:  string(PostsService),
				SourceEntityID: pgtype.Int8{Int64: 300, Valid: true},
				Seen:           pgtype.Bool{Bool: true, Valid: true},
				NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
				Acted:          pgtype.Bool{Bool: false, Valid: true},
				Payload:        payload2,
				CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true}, // Most recent
				ExpiresAt:      pgtype.Timestamptz{Time: time.Now().Add(29 * 24 * time.Hour), Valid: true},
			},
		},
	}

	service := &Service{
		db: mockDB,
	}

	ctx := context.Background()

	// Test getting all notifications for a user
	req := &GetNotificationsRequest{
		UserID: 100,
		Limit:  10,
		Offset: 0,
	}

	resp, err := service.GetNotifications(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	assert.Equal(t, 2, len(resp.Notifications))

	// Check that the first notification (most recent) is the one we expect
	// Since the query orders by created_at DESC, the most recent should be first
	assert.Equal(t, int64(2), resp.Notifications[0].ID) // Most recent notification first due to DESC order
	assert.Equal(t, Like, resp.Notifications[0].Type)
	assert.Equal(t, true, resp.Notifications[0].Seen)
}

func TestMarkAsRead(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{"follower_id": "200"})

	mockDB := &MockDB{
		notifications: []sqlc.GetNotificationsByUserIdAllRow{
			{
				ID:             1,
				UserID:         100,
				NotifType:      string(NewFollower),
				SourceService:  string(UsersService),
				SourceEntityID: pgtype.Int8{Int64: 200, Valid: true},
				Seen:           pgtype.Bool{Bool: false, Valid: true},
				NeedsAction:    pgtype.Bool{Bool: false, Valid: true},
				Acted:          pgtype.Bool{Bool: false, Valid: true},
				Payload:        payload,
				CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
				ExpiresAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
		},
	}

	service := &Service{
		db: mockDB,
	}

	ctx := context.Background()

	req := &MarkAsReadRequest{
		UserID:          100,
		NotificationIDs: []int64{1},
	}

	resp, err := service.MarkAsRead(ctx, req)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, int32(1), resp.UpdatedCount)
}

func TestCreateNotification(t *testing.T) {
	mockDB := &MockDB{
		notifications: []sqlc.GetNotificationsByUserIdAllRow{},
	}

	service := &Service{
		db: mockDB,
	}

	ctx := context.Background()

	req := &CreateNotificationRequest{
		UserID:         100,
		Type:           NewFollower,
		SourceService:  UsersService,
		SourceEntityID: 200,
		Payload:        map[string]string{"follower_id": "200"},
		NeedsAction:    false,
	}

	resp, err := service.CreateNotification(ctx, req)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, int64(1), resp.ID) // First notification created
}