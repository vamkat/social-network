package handler

import (
	"context"
	"fmt"
	"social-network/services/media/internal/application"
	"social-network/shared/gen-go/media"

	_ "github.com/lib/pq"
)

type MediaHandler struct {
	media.UnimplementedMediaServiceServer
	Application *application.MediaService
	Port        string
}

func (s *MediaHandler) UploadImage(ctx context.Context, req *media.UploadImageRequest) (*media.UploadImageResponse, error) {
	// Call your existing SaveImage method
	info, err := s.Application.SaveImage(ctx, req.FileContent, req.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	return &media.UploadImageResponse{
		Id:        info.Id,
		MimeType:  info.MimeType,
		SizeBytes: info.SizeBytes,
		Bucket:    info.Bucket,
		ObjectKey: info.ObjectKey,
	}, nil
}

func (h *MediaHandler) FetchImage() {}
