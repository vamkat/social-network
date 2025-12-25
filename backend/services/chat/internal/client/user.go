package client

import (
	"context"
	cm "social-network/shared/gen-go/common"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

// Calls user client to convert a slice of ct.Ids representing users to a
// map[ct.Id]models.User.
func (c *Clients) UserIdsToMap(ctx context.Context,
	ids ct.Ids) (map[ct.Id]md.User, error) {
	// Deduplicate IDs
	uniq := make(map[ct.Id]struct{}, len(ids))
	cleaned := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := uniq[id]; !ok {
			uniq[id] = struct{}{}
			cleaned = append(cleaned, id.Int64())
		}
	}

	if len(cleaned) == 0 {
		return map[ct.Id]md.User{}, nil
	}
	// Call redis first

	// gRPC request
	req := &cm.UserIds{Values: cleaned}
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert response → map
	out := make(map[ct.Id]md.User, len(resp.Users))
	for _, u := range resp.Users {
		uid := ct.Id(u.UserId)
		out[uid] = md.User{
			UserId:   uid,
			Username: ct.Username(u.Username),
			AvatarId: ct.Id(u.Avatar),
		}
	}

	return out, nil
}

// Converts a slice of ct.Ids representing users to models.User slice.
func (c *Clients) UserIdsToUsers(ctx context.Context,
	ids ct.Ids) (userInfo []md.User, err error) {
	req := &cm.UserIds{Values: ids.Int64()}
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, u := range resp.Users {
		userInfo = append(userInfo, md.User{
			UserId:   ct.Id(u.UserId),
			Username: ct.Username(u.Username),
			AvatarId: ct.Id(u.Avatar),
		})
	}
	return userInfo, nil
}

// Function to implemeted by Hydrator
func (c *Clients) GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*cm.ListUsers, error) {
	req := &cm.UserIds{
		Values: userIds,
	}
	resp, err := c.UserClient.GetBatchBasicUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
