package application

// func TestGetAllGroupsPaginated_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(1),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.GetAllGroupsRow{
// 		{
// 			ID:               1,
// 			GroupTitle:       "Group 1",
// 			GroupDescription: "Description 1",
// 			MembersCount:     5,
// 		},
// 		{
// 			ID:               2,
// 			GroupTitle:       "Group 2",
// 			GroupDescription: "Description 2",
// 			MembersCount:     10,
// 		},
// 	}

// 	mockDB.On("GetAllGroups", ctx, sqlc.GetAllGroupsParams{
// 		Limit:  10,
// 		Offset: 0,
// 	}).Return(expectedRows, nil)

// 	// Mock the userInRelationToGroup calls for each group
// 	mockDB.On("IsUserGroupOwner", ctx, sqlc.IsUserGroupOwnerParams{
// 		ID:         1,
// 		GroupOwner: 1,
// 	}).Return(false, nil)
// 	mockDB.On("IsUserGroupMember", ctx, sqlc.IsUserGroupMemberParams{
// 		GroupID: 1,
// 		UserID:  1,
// 	}).Return(false, nil)
// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: 1,
// 		UserID:  1,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	mockDB.On("IsUserGroupOwner", ctx, sqlc.IsUserGroupOwnerParams{
// 		ID:         2,
// 		GroupOwner: 1,
// 	}).Return(false, nil)
// 	mockDB.On("IsUserGroupMember", ctx, sqlc.IsUserGroupMemberParams{
// 		GroupID: 2,
// 		UserID:  1,
// 	}).Return(false, nil)
// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: 2,
// 		UserID:  1,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	groups, err := service.GetAllGroupsPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, groups, 2)
// 	assert.Equal(t, "Group 1", groups[0].GroupTitle.String())
// 	assert.Equal(t, int32(5), groups[0].MembersCount)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetAllGroupsPaginated_Empty(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.Pagination{
// 		UserId: ct.Id(1),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	mockDB.On("GetAllGroups", ctx, sqlc.GetAllGroupsParams{
// 		Limit:  10,
// 		Offset: 0,
// 	}).Return([]sqlc.GetAllGroupsRow{}, nil)

// 	groups, err := service.GetAllGroupsPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, groups, 0)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetUserGroupsPaginated_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	userID := int64(1)
// 	req := models.Pagination{
// 		UserId: ct.Id(userID),
// 		Limit:  ct.Limit(10),
// 		Offset: ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.GetUserGroupsRow{
// 		{
// 			GroupID:          1,
// 			GroupTitle:       "Group 1",
// 			GroupDescription: "Description 1",
// 			MembersCount:     5,
// 			IsMember:         true,
// 			IsOwner:          true,
// 		},
// 		{
// 			GroupID:          2,
// 			GroupTitle:       "Group 2",
// 			GroupDescription: "Description 2",
// 			MembersCount:     10,
// 			IsMember:         true,
// 			IsOwner:          false,
// 		},
// 	}

// 	mockDB.On("GetUserGroups", ctx, sqlc.GetUserGroupsParams{
// 		GroupOwner: userID,
// 		Limit:      10,
// 		Offset:     0,
// 	}).Return(expectedRows, nil)

// 	// Mock the isGroupMembershipPending calls for each group
// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: 1,
// 		UserID:  userID,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: 2,
// 		UserID:  userID,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	groups, err := service.GetUserGroupsPaginated(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, groups, 2)
// 	assert.Equal(t, true, groups[0].IsOwner)
// 	assert.Equal(t, false, groups[1].IsOwner)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetGroupInfo_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	groupID := int64(1)
// 	userID := int64(1)
// 	req := models.GeneralGroupReq{
// 		GroupId: ct.Id(groupID),
// 		UserId:  ct.Id(userID),
// 	}

// 	expectedRow := sqlc.GetGroupInfoRow{
// 		ID:               groupID,
// 		GroupTitle:       "Test Group",
// 		GroupDescription: "Test Description",
// 		MembersCount:     15,
// 	}

// 	mockDB.On("GetGroupInfo", ctx, groupID).Return(expectedRow, nil)

// 	// Mock the userInRelationToGroup calls
// 	mockDB.On("IsUserGroupOwner", ctx, sqlc.IsUserGroupOwnerParams{
// 		ID:         groupID,
// 		GroupOwner: userID,
// 	}).Return(false, nil)
// 	mockDB.On("IsUserGroupMember", ctx, sqlc.IsUserGroupMemberParams{
// 		GroupID: groupID,
// 		UserID:  userID,
// 	}).Return(false, nil)
// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: groupID,
// 		UserID:  userID,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	group, err := service.GetGroupInfo(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "Test Group", group.GroupTitle.String())
// 	assert.Equal(t, int32(15), group.MembersCount)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetGroupInfo_NotFound(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	groupID := int64(999)
// 	userID := int64(1)
// 	req := models.GeneralGroupReq{
// 		GroupId: ct.Id(groupID),
// 		UserId:  ct.Id(userID),
// 	}

// 	mockDB.On("GetGroupInfo", ctx, groupID).Return(sqlc.GetGroupInfoRow{}, errors.New("group not found"))

// 	_, err := service.GetGroupInfo(ctx, req)

// 	assert.Error(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestGetGroupMembers_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	groupID := int64(1)
// 	userID := int64(1)
// 	req := models.GroupMembersReq{
// 		UserId:  ct.Id(userID),
// 		GroupId: ct.Id(groupID),
// 		Limit:   ct.Limit(10),
// 		Offset:  ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.GetGroupMembersRow{
// 		{
// 			ID:       1,
// 			Username: "user1",
// 			AvatarID: 2,
// 			Role: sqlc.NullGroupRole{
// 				GroupRole: "owner",
// 				Valid:     true,
// 			},
// 		},
// 		{
// 			ID:       2,
// 			Username: "user2",
// 			AvatarID: 3,
// 			Role: sqlc.NullGroupRole{
// 				GroupRole: "member",
// 				Valid:     true,
// 			},
// 		},
// 	}

// 	mockDB.On("GetGroupMembers", ctx, sqlc.GetGroupMembersParams{
// 		GroupID: groupID,
// 		Limit:   10,
// 		Offset:  0,
// 	}).Return(expectedRows, nil)

// 	// Mock isGroupMember call
// 	mockDB.On("IsUserGroupMember", ctx, sqlc.IsUserGroupMemberParams{
// 		GroupID: groupID,
// 		UserID:  userID,
// 	}).Return(true, nil)

// 	members, err := service.GetGroupMembers(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, members, 2)
// 	assert.Equal(t, "user1", members[0].Username.String())
// 	assert.Equal(t, "owner", members[0].GroupRole)
// 	mockDB.AssertExpectations(t)
// }

// func TestSearchGroups_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	searchTerm := "test"
// 	req := models.GroupSearchReq{
// 		SearchTerm: ct.SearchTerm(searchTerm),
// 		UserId:     ct.Id(1),
// 		Limit:      ct.Limit(10),
// 		Offset:     ct.Offset(0),
// 	}

// 	expectedRows := []sqlc.SearchGroupsFuzzyRow{
// 		{
// 			ID:               1,
// 			GroupTitle:       "Test Group",
// 			GroupDescription: "A test group",
// 			MembersCount:     5,
// 		},
// 	}

// 	mockDB.On("SearchGroupsFuzzy", ctx, sqlc.SearchGroupsFuzzyParams{
// 		Similarity: searchTerm,
// 		GroupOwner: 1,
// 		Limit:      10,
// 		Offset:     0,
// 	}).Return(expectedRows, nil)

// 	// Mock the isGroupMembershipPending call for the returned group
// 	mockDB.On("IsGroupMembershipPending", ctx, sqlc.IsGroupMembershipPendingParams{
// 		GroupID: 1,
// 		UserID:  1,
// 	}).Return(pgtype.Bool{Bool: false, Valid: true}, nil)

// 	groups, err := service.SearchGroups(ctx, req)

// 	assert.NoError(t, err)
// 	assert.Len(t, groups, 1)
// 	assert.Equal(t, "Test Group", groups[0].GroupTitle.String())
// 	mockDB.AssertExpectations(t)
// }

// func TestCreateGroup_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.CreateGroupRequest{
// 		OwnerId:          ct.Id(1),
// 		GroupTitle:       ct.Title("New Group"),
// 		GroupDescription: ct.About("New Description"),
// 	}

// 	mockDB.On("CreateGroup", ctx, sqlc.CreateGroupParams{
// 		GroupOwner:       1,
// 		GroupTitle:       "New Group",
// 		GroupDescription: "New Description",
// 	}).Return(int64(5), nil)

// 	// expect CreateGroupConversation to be called but return nil
// 	mockClients.On("CreateGroupConversation", ctx, int64(5), int64(1)).Return(nil)

// 	groupID, err := service.CreateGroup(ctx, &req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, models.GroupId(5), groupID)
// 	mockDB.AssertExpectations(t)
// }

// func TestCreateGroup_Error(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.CreateGroupRequest{
// 		OwnerId:          ct.Id(1),
// 		GroupTitle:       ct.Title("New Group"),
// 		GroupDescription: ct.About("New Description"),
// 	}

// 	mockDB.On("CreateGroup", ctx, sqlc.CreateGroupParams{
// 		GroupOwner:       1,
// 		GroupTitle:       "New Group",
// 		GroupDescription: "New Description",
// 	}).Return(int64(0), errors.New("database error"))

// 	// client shouldn't be called on DB error, but set expectation to be safe
// 	mockClients.On("CreateGroupConversation", ctx, int64(0), int64(1)).Return(nil)

// 	_, err := service.CreateGroup(ctx, &req)

// 	assert.Error(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestLeaveGroup_Success(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.GeneralGroupReq{
// 		GroupId: ct.Id(1),
// 		UserId:  ct.Id(2),
// 	}

// 	mockDB.On("LeaveGroup", ctx, sqlc.LeaveGroupParams{
// 		GroupID: 1,
// 		UserID:  2,
// 	}).Return(nil)

// 	err := service.LeaveGroup(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestRequestJoinGroupOrCancel_Request(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.GroupJoinRequest{
// 		GroupId:     ct.Id(1),
// 		RequesterId: ct.Id(2),
// 	}

// 	mockDB.On("SendGroupJoinRequest", ctx, sqlc.SendGroupJoinRequestParams{
// 		GroupID: 1,
// 		UserID:  2,
// 	}).Return(nil)

// 	err := service.RequestJoinGroup(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestRequestJoinGroupOrCancel_Cancel(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.GroupJoinRequest{
// 		GroupId:     ct.Id(1),
// 		RequesterId: ct.Id(2),
// 	}

// 	mockDB.On("CancelGroupJoinRequest", ctx, sqlc.CancelGroupJoinRequestParams{
// 		GroupID: 1,
// 		UserID:  2,
// 	}).Return(nil)

// 	err := service.CancelJoinGroupRequest(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestRespondToGroupInvite_Accept(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	mockClients := new(MockClients)
// 	service := NewApplicationWithMocks(mockDB, mockClients)

// 	ctx := context.Background()
// 	req := models.HandleGroupInviteRequest{
// 		GroupId:   ct.Id(1),
// 		InvitedId: ct.Id(2),
// 		Accepted:  true,
// 	}

// 	mockDB.On("AcceptGroupInvite", ctx, sqlc.AcceptGroupInviteParams{
// 		GroupID:    1,
// 		ReceiverID: 2,
// 	}).Return(nil)

// 	// expect AddMembersToGroupConversation to be called
// 	mockClients.On("AddMembersToGroupConversation", ctx, int64(1), []int64{2}).Return(nil)

// 	err := service.RespondToGroupInvite(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }

// func TestRespondToGroupInvite_Decline(t *testing.T) {
// 	mockDB := new(MockQuerier)
// 	service := NewApplication(mockDB, nil, nil)

// 	ctx := context.Background()
// 	req := models.HandleGroupInviteRequest{
// 		GroupId:   ct.Id(1),
// 		InvitedId: ct.Id(2),
// 		Accepted:  false,
// 	}

// 	mockDB.On("DeclineGroupInvite", ctx, sqlc.DeclineGroupInviteParams{
// 		GroupID:    1,
// 		ReceiverID: 2,
// 	}).Return(nil)

// 	err := service.RespondToGroupInvite(ctx, req)

// 	assert.NoError(t, err)
// 	mockDB.AssertExpectations(t)
// }
