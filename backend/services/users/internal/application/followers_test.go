package application

// func TestGetFollowersPaginated_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(1),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.GetFollowersRow{
// 		{
// 			ID:            2,
// 			Username:      "follower1",
// 			AvatarID:      3,
// 			ProfilePublic: true,
// 		},
// 		{
// 			ID:            3,
// 			Username:      "follower2",
// 			AvatarID:      4,
// 			ProfilePublic: true,
// 		},
// 	}

// 	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
// 		FollowingID: 1,
// 		Limit:       10,
// 		Offset:      0,
// 	}).Return(expectedRows, nil)

// 	users, err := service.GetFollowersPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, users, 2)
// 	assert.Equal(t, "follower1", users[0].Username.String())
// 	assert.Equal(t, "follower2", users[1].Username.String())
// 	mockDB.AssertExpectations(t)
// }

// func TestGetFollowersPaginated_Empty(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(999),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	mockDB.On("GetFollowers", ctx, sqlc.GetFollowersParams{
// 		FollowingID: 999,
// 		Limit:       10,
// 		Offset:      0,
// 	}).Return([]sqlc.GetFollowersRow{}, nil)

// 	users, err := service.GetFollowersPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, users, 0)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetFollowingPaginated_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(1),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.GetFollowingRow{
// 		{
// 			ID:            2,
// 			Username:      "following1",
// 			AvatarID:      3,
// 			ProfilePublic: true,
// 		},
// 		{
// 			ID:            3,
// 			Username:      "following2",
// 			AvatarID:      4,
// 			ProfilePublic: false,
// 		},
// 	}

// 	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
// 		FollowerID: 1,
// 		Limit:      10,
// 		Offset:     0,
// 	}).Return(expectedRows, nil)

// 	users, err := service.GetFollowingPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, users, 2)
// 	assert.Equal(t, "following1", users[0].Username.String())
// 	assert.Equal(t, "following2", users[1].Username.String())
// 	mockDB.AssertExpectations(t)
// }

// func TestGetFollowingPaginated_Empty(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(999),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	mockDB.On("GetFollowing", ctx, sqlc.GetFollowingParams{
// 		FollowerID: 999,
// 		Limit:      10,
// 		Offset:     0,
// 	}).Return([]sqlc.GetFollowingRow{}, nil)

// 	users, err := service.GetFollowingPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, users, 0)
// 	mockDB.AssertExpectations(t)
// }

// func TestFollowUser_Immediate(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
// 		PFollower: 1,
// 		PTarget:   2,
// 	}).Return("accepted", nil)

// 	// ensure AreFollowingEachOther is mocked so createPrivateConversation doesn't trigger an unexpected call
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: false, User2FollowsUser1: false}, nil)

// 	resp, err := service.FollowUser(ctx, req)

// 	assert.NoError(t, err)
// 	assert.False(t, resp.IsPending)
// 	assert.True(t, resp.ViewerIsFollowing)
// 	mockDB.AssertExpectations(t)
// }

// func TestFollowUser_Pending(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
// 		PFollower: 1,
// 		PTarget:   2,
// 	}).Return("requested", nil)

// 	resp, err := service.FollowUser(ctx, req)

// 	assert.NoError(t, err)
// 	assert.True(t, resp.IsPending)
// 	assert.False(t, resp.ViewerIsFollowing)
// 	mockDB.AssertExpectations(t)
// }

// func TestFollowUser_Error(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	mockDB.On("FollowUser", ctx, sqlc.FollowUserParams{
// 		PFollower: 1,
// 		PTarget:   2,
// 	}).Return("", errors.New("database error"))

// 	_, err := service.FollowUser(ctx, req)

// 	assert.Error(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestUnFollowUser_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   ct.Id(1),
// 		TargetUserId: ct.Id(2),
// 	}

// 	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(nil)

// 	// ensure AreFollowingEachOther is mocked so deletePrivateConversation doesn't trigger an unexpected call
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  1,
// 		FollowingID: 2,
// 	}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: true, User2FollowsUser1: false}, nil)

// 	isFollowing, err := service.UnFollowUser(ctx, req)

// 	assert.NoError(t, err)
// 	assert.True(t, isFollowing)
// 	mockDB.AssertExpectations(t)
// }

// func TestUnFollowUser_Error(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.FollowUserReq{
// 		FollowerId:   1,
// 		TargetUserId: 999,
// 	}

// 	mockDB.On("UnfollowUser", ctx, sqlc.UnfollowUserParams{
// 		FollowerID:  1,
// 		FollowingID: 999,
// 	}).Return(errors.New("not following"))

// 	isFollowing, err := service.UnFollowUser(ctx, req)

// 	assert.Error(t, err)
// 	assert.False(t, isFollowing)
// 	mockDB.AssertExpectations(t)
// }

// func TestHandleFollowRequest_Accept(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.HandleFollowRequestReq{
// 		UserId:      ct.Id(1),
// 		RequesterId: ct.Id(2),
// 		Accept:      true,
// 	}

// 	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
// 		RequesterID: 2,
// 		TargetID:    1,
// 	}).Return(nil)

// 	// when accepting, the service may try to create a private conversation; mock AreFollowingEachOther
// 	mockDB.On("AreFollowingEachOther", ctx, sqlc.AreFollowingEachOtherParams{
// 		FollowerID:  2,
// 		FollowingID: 1,
// 	}).Return(sqlc.AreFollowingEachOtherRow{User1FollowsUser2: false, User2FollowsUser1: false}, nil)

// 	err := service.HandleFollowRequest(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestHandleFollowRequest_Reject(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.HandleFollowRequestReq{
// 		UserId:      ct.Id(1),
// 		RequesterId: ct.Id(2),
// 		Accept:      false,
// 	}

// 	mockDB.On("RejectFollowRequest", ctx, sqlc.RejectFollowRequestParams{
// 		RequesterID: 2,
// 		TargetID:    1,
// 	}).Return(nil)

// 	err := service.HandleFollowRequest(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestHandleFollowRequest_AcceptError(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.HandleFollowRequestReq{
// 		UserId:      ct.Id(1),
// 		RequesterId: ct.Id(2),
// 		Accept:      true,
// 	}

// 	mockDB.On("AcceptFollowRequest", ctx, sqlc.AcceptFollowRequestParams{
// 		RequesterID: 2,
// 		TargetID:    1,
// 	}).Return(errors.New("request not found"))

// 	err := service.HandleFollowRequest(ctx, req)

// 	assert.Error(t, err)
// 	mockDB.AssertExpectations(t)
// }
