package application

// func TestHasRightToView_Success(t *testing.T) {
// 	ctx := context.Background()

// 	dbMock := &dbmocks.MockQueries{}
// 	clientMock := &clientmocks.MockClients{}

// 	// entity id 42, creator id 100
// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(42)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 100, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(7), int64(100)).Return(true, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(7), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)

// 	app := NewApplicationWithMocks(dbMock, clientMock)

// 	ok, err := app.hasRightToView(ctx, accessContext{requesterId: 7, entityId: 42})
// 	assert.NoError(t, err)
// 	assert.True(t, ok)

// 	dbMock.AssertExpectations(t)
// 	clientMock.AssertExpectations(t)
// }

// func TestCreatePost_SelectedAudience_Success(t *testing.T) {
// 	ctx := context.Background()

// 	dbMock := &dbmocks.MockQueries{}
// 	txMock := &mocks.MockTxRunner{Queries: dbMock}
// 	clientMock := &clientmocks.MockClients{}

// 	// Expect CreatePost called inside transaction
// 	dbMock.On("CreatePost", mock.Anything, mock.Anything).Return(int64(123), nil)
// 	dbMock.On("InsertPostAudience", mock.Anything, mock.Anything).Return(int64(1), nil)
// 	txMock.On("RunTx", mock.Anything).Return(nil)

// 	app := NewApplicationWithMocksTx(dbMock, clientMock, txMock)

// 	req := models.CreatePostReq{
// 		CreatorId:   ct.Id(1),
// 		Body:        ct.PostBody("hello world"),
// 		Audience:    ct.Audience("selected"),
// 		AudienceIds: ct.FromInt64s([]int64{2}),
// 	}

// 	err := app.CreatePost(ctx, req)
// 	assert.NoError(t, err)

// 	dbMock.AssertExpectations(t)
// 	txMock.AssertExpectations(t)
// }

// fake retriever to avoid touching redis/clients during user retrieval calls
// type fakeRetriever struct{}

// func (f *fakeRetriever) GetUsers(ctx context.Context, userIDs []int64) (map[int64]models.User, error) {
// 	return nil, nil
// }

// func TestDeletePost_Success(t *testing.T) {
// 	ctx := context.Background()

// 	dbMock := &dbmocks.MockQueries{}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(10)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 5, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(2), int64(5)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(2), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("DeletePost", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	app := NewApplicationWithMocks(dbMock, clientMock)

// 	err := app.DeletePost(ctx, models.GenericReq{RequesterId: ct.Id(2), EntityId: ct.Id(10)})
// 	assert.NoError(t, err)

// 	dbMock.AssertExpectations(t)
// 	clientMock.AssertExpectations(t)
// }

// func TestEditPost_Success(t *testing.T) {
// 	ctx := context.Background()

// 	dbMock := &dbmocks.MockQueries{}
// 	txMock := &mocks.MockTxRunner{Querier: dbMock}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(20)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 3, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(3), int64(3)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(3), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)

// 	// inside tx
// 	dbMock.On("EditPostContent", mock.Anything, mock.Anything).Return(int64(1), nil)
// 	dbMock.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 	dbMock.On("UpdatePostAudience", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	txMock.On("RunTx", mock.Anything).Return(nil)

// 	app := NewApplicationWithMocksTx(dbMock, clientMock, txMock)

// 	req := models.EditPostReq{
// 		RequesterId: ct.Id(3),
// 		PostId:      ct.Id(20),
// 		NewBody:     ct.PostBody("new body"),
// 		Image:       ct.Id(1),
// 		Audience:    ct.Audience("everyone"),
// 	}

// 	err := app.EditPost(ctx, req)
// 	assert.NoError(t, err)

// 	dbMock.AssertExpectations(t)
// 	txMock.AssertExpectations(t)
// }

// func TestCreateComment_Success(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQueries{}
// 	txMock := &mocks.MockTxRunner{Querier: dbMock}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(30)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 4, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(6), int64(4)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(6), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)

// 	dbMock.On("CreateComment", mock.Anything, mock.Anything).Return(nil)
// 	dbMock.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 	txMock.On("RunTx", mock.Anything).Return(nil)

// 	app := NewApplicationWithMocksTx(dbMock, clientMock, txMock)

// 	req := models.CreateCommentReq{CreatorId: ct.Id(6), ParentId: ct.Id(30), Body: ct.CommentBody("hey there"), Image: ct.Id(1)}
// 	err := app.CreateComment(ctx, req)
// 	assert.NoError(t, err)

// 	dbMock.AssertExpectations(t)
// }

// func TestEditComment_Success(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQueries{}
// 	txMock := &mocks.MockTxRunner{Querier: dbMock}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(40)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 7, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(7), int64(7)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(7), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)

// 	dbMock.On("EditComment", mock.Anything, mock.Anything).Return(int64(1), nil)
// 	dbMock.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 	txMock.On("RunTx", mock.Anything).Return(nil)

// 	app := NewApplicationWithMocksTx(dbMock, clientMock, txMock)

// 	req := models.EditCommentReq{CreatorId: ct.Id(7), CommentId: ct.Id(40), Body: ct.CommentBody("updated comment"), Image: ct.Id(1)}
// 	err := app.EditComment(ctx, req)
// 	assert.NoError(t, err)
// }

// func TestDeleteComment_Success(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQueries{}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(50)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 9, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(8), int64(9)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(8), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("DeleteComment", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	app := NewApplicationWithMocks(dbMock, clientMock)

// 	err := app.DeleteComment(ctx, models.GenericReq{RequesterId: ct.Id(8), EntityId: ct.Id(50)})
// 	assert.NoError(t, err)
// }

// func TestCreateEvent_EditDeleteRespond_Success(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQueries{}
// 	txMock := &mocks.MockTxRunner{Querier: dbMock}
// 	clientMock := &clientmocks.MockClients{}

// 	// CreateEvent
// 	clientMock.On("IsGroupMember", mock.Anything, int64(11), int64(2)).Return(true, nil)
// 	dbMock.On("CreateEvent", mock.Anything, mock.Anything).Return(nil)

// 	app := NewApplicationWithMocksTx(dbMock, clientMock, txMock)

// 	createReq := models.CreateEventReq{CreatorId: ct.Id(11), GroupId: ct.Id(2), Title: ct.Title("t"), Body: ct.EventBody("event body"), EventDate: ct.EventDateTime(time.Now())}
// 	err := app.CreateEvent(ctx, createReq)
// 	assert.NoError(t, err)

// 	// EditEvent
// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(60)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 11, GroupID: 2}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(11), int64(11)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(11), int64(2)).Return(true, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)

// 	dbMock.On("EditEvent", mock.Anything, mock.Anything).Return(int64(1), nil)
// 	dbMock.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 	txMock.On("RunTx", mock.Anything).Return(nil)

// 	editReq := models.EditEventReq{EventId: ct.Id(60), RequesterId: ct.Id(11), Title: ct.Title("t2"), Body: ct.EventBody("event edit"), EventDate: ct.EventDateTime(time.Now()), Image: ct.Id(1)}
// 	err = app.EditEvent(ctx, editReq)
// 	assert.NoError(t, err)

// 	// DeleteEvent
// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(70)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 11, GroupID: 2}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(11), int64(11)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(11), int64(2)).Return(true, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("DeleteEvent", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	err = app.DeleteEvent(ctx, models.GenericReq{RequesterId: ct.Id(11), EntityId: ct.Id(70)})
// 	assert.NoError(t, err)

// 	// RespondToEvent
// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(80)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 5, GroupID: 2}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(12), int64(5)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(12), int64(2)).Return(true, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("UpsertEventResponse", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	err = app.RespondToEvent(ctx, models.RespondToEventReq{EventId: ct.Id(80), ResponderId: ct.Id(12), Going: true})
// 	assert.NoError(t, err)

// 	// RemoveEventResponse
// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(90)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 5, GroupID: 2}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(13), int64(5)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(13), int64(2)).Return(true, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("DeleteEventResponse", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	err = app.RemoveEventResponse(ctx, models.GenericReq{RequesterId: ct.Id(13), EntityId: ct.Id(90)})
// 	assert.NoError(t, err)
// }

// func TestToggleOrInsertReaction_Success(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQueries{}
// 	clientMock := &clientmocks.MockClients{}

// 	dbMock.On("GetEntityCreatorAndGroup", mock.Anything, int64(100)).Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 14, GroupID: 0}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(14), int64(14)).Return(false, nil)
// 	clientMock.On("IsGroupMember", mock.Anything, int64(14), int64(0)).Return(false, nil)
// 	dbMock.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 	dbMock.On("ToggleOrInsertReaction", mock.Anything, mock.Anything).Return(int64(1), nil)

// 	app := NewApplicationWithMocks(dbMock, clientMock)

// 	err := app.ToggleOrInsertReaction(ctx, models.GenericReq{RequesterId: ct.Id(14), EntityId: ct.Id(100)})
// 	assert.NoError(t, err)
// }

// func TestFeeds_Empty_NoHydrationError(t *testing.T) {
// 	ctx := context.Background()
// 	dbMock := &dbmocks.MockQuerier{}
// 	clientMock := &clientmocks.MockClients{}

// 	clientMock.On("GetFollowingIds", mock.Anything, int64(21)).Return([]int64{}, nil)
// 	clientMock.On("IsFollowing", mock.Anything, int64(21), int64(21)).Return(false, nil)
// 	dbMock.On("GetPersonalizedFeed", mock.Anything, mock.Anything).Return([]sqlc.GetPersonalizedFeedRow{}, nil)
// 	dbMock.On("GetPublicFeed", mock.Anything, mock.Anything).Return([]sqlc.GetPublicFeedRow{}, nil)
// 	dbMock.On("GetUserPostsPaginated", mock.Anything, mock.Anything).Return([]sqlc.GetUserPostsPaginatedRow{}, nil)
// 	dbMock.On("GetGroupPostsPaginated", mock.Anything, mock.Anything).Return([]sqlc.GetGroupPostsPaginatedRow{}, nil)

// 	app := NewApplicationWithMocks(dbMock, clientMock)
// 	app.userRetriever = &fakeRetriever{}

// 	_, err := app.GetPersonalizedFeed(ctx, models.GetPersonalizedFeedReq{RequesterId: ct.Id(21), Limit: ct.Limit(10), Offset: ct.Offset(0)})
// 	assert.NoError(t, err)

// 	_, err = app.GetPublicFeed(ctx, models.GenericPaginatedReq{RequesterId: ct.Id(21), Limit: ct.Limit(10), Offset: ct.Offset(0)})
// 	assert.NoError(t, err)

// 	_, err = app.GetUserPostsPaginated(ctx, models.GetUserPostsReq{CreatorId: ct.Id(21), RequesterId: ct.Id(21), Limit: ct.Limit(10), Offset: ct.Offset(0)})
// 	// GetUserPostsPaginated returns ErrNotFound on empty rows
// 	assert.Error(t, err)

// 	// GetGroupPostsPaginated requires a group id; passing 0 should return ErrNoGroupIdGiven
// 	_, err = app.GetGroupPostsPaginated(ctx, models.GetGroupPostsReq{RequesterId: ct.Id(21), GroupId: ct.Id(0)})
// 	assert.Error(t, err)
// }
