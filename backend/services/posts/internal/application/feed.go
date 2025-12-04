package application

import "context"

//FRONT: Are posts for groups the requester is a member of included here? I assumed no - NO
func (s *Application) GetPersonalizedFeed(ctx context.Context, req GetPersonalizedFeedReq) ([]Post, error) {
	//API Gateway needs to provide list of ids the requester follows
	//(and in case we include group posts, a list of groups they are a member of)
	return nil, nil
}

func (s *Application) GetPublicFeed(ctx context.Context, req GenericPaginatedReq) ([]Post, error) {
	return nil, nil
}
