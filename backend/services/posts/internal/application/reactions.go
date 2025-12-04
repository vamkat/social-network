package application

import "context"

// FRONT: reaction count changes. What do you need returned? If only error, you'll change count on your side optimistically - OPTMISTICALLY
func (s *Application) InsertReaction(ctx context.Context, req GenericReq) error {
	//runs in transaction
	//tries toggle reaction if exists first, then inserts reaction if it didn't exist
	return nil
}

// FRONT: Do you want to show this info or skip?
func (s *Application) GetWhoLikedEntityId(ctx context.Context, req GenericReq) ([]int64, error) {
	return nil, nil
}
