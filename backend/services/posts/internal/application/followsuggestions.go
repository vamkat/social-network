package application

import "context"

func (s *Application) SuggestUsersByPostActivity(ctx context.Context, userId int64) ([]int64, error) {
	//returns five random ids that fit one of the following criteria:
	//Users who liked one or more of *your public posts*
	// Users who commented on your public posts
	//  Users who liked the same posts as you
	// Users who commented on the same posts as you
	// Actual Basic User Info will be retrieved by Gateway from Users
	return nil, nil
}
