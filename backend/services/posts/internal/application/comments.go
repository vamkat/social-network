package application

import "context"

//FRONT: Do you prefer full Comment instead of just id?
func (s *Application) CreateComment(ctx context.Context, req CreateCommentReq) (err error) {
	// check requester can actually view parent entity? (probably not needed?)
	return nil
}

// FRONT: Do I return full comment or just error?
func (s *Application) EditComment(ctx context.Context, req EditCommentReq) error {
	//check requester is creator
	return nil
}

func (s *Application) DeleteComment(ctx context.Context, req GenericReq) error {
	//check requester is comment creator
	return nil
}

func (s *Application) GetCommentsByParentId(ctx context.Context, req GenericPaginatedReq) ([]Comment, error) {
	// check requester can actually view parent entity?
	return nil, nil
}

//FRONT: I assume this is a different endpoint that feed? Or do you prefer I include this in every Post in []Post I return?- IN POST
func (s *Application) GetLatestCommentForPostId(ctx context.Context, req GenericReq) (Comment, error) {
	// check requester can actually view parent entity?
	return Comment{}, nil
}
