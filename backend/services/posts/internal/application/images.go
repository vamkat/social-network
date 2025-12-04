package application

import "context"

// FRONT: If user edits post to change the image, what do you send me exactly?
func (s *Application) InsertImages(ctx context.Context, req ImageReq) error {
	//runs in transaction from createPost (edit signature to include tx?)
	// also if user edits post to change image
	//accepts multiple images for scalability
	return nil
}

func (s *Application) DeleteImage(ctx context.Context, imageId int64) error {
	return nil
}

//FRONT: this is supposed to edit filename and/or sort order for image
// If we're doing only one image we're probably better off deleting and reinserting
//so probably not needed?
func (s *Application) UpdateImage(ctx context.Context) error {
	return nil
}

// FRONT: If only one image per post, this is probably never needed
func (s *Application) getImages(ctx context.Context, req GenericReq) ([]string, error) {
	return nil, nil
}

//ONE IMAGE!!!
