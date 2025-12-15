package handler

import (
	"social-network/services/media/internal/application"
	pb "social-network/shared/gen-go/media"

	_ "github.com/lib/pq"
)

type MediaHandler struct {
	pb.UnimplementedMediaServiceServer
	Application *application.MediaService
	Port        string
}

func PresignedUpload()   {}
func PresignedDownload() {}
func VerifyUpload()      {}

// func (s *MediaHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest,
// ) (*pb.UploadImageResponse, error) {
// 	// Call your existing SaveImage method
// 	info, err := s.Application.SaveImage(ctx, req.FileContent, req.Filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to save image: %w", err)
// 	}

// 	return &pb.UploadImageResponse{
// 		Id:        info.Id,
// 		MimeType:  info.MimeType,
// 		SizeBytes: info.SizeBytes,
// 		Bucket:    info.Bucket,
// 		ObjectKey: info.ObjectKey,
// 	}, nil
// }

// func (s *MediaHandler) RetrieveImage(
// 	req *pb.RetrieveImageRequest, stream pb.MediaService_RetrieveImageServer,
// ) error {

// 	ctx := stream.Context()

// 	// --- 1. Call your domain service ---
// 	reader, meta, err := s.Application.RetriveImageById(ctx, customtypes.Id(req.ImageId))
// 	if err != nil {
// 		return status.Errorf(codes.NotFound, "cannot retrieve image: %v", err)
// 	}
// 	defer reader.Close()

// 	// --- 2. Send metadata first ---
// 	metaMsg := &pb.ImageMeta{
// 		Id:        meta.Id,
// 		Filename:  meta.Filename,
// 		MimeType:  meta.MimeType,
// 		SizeBytes: meta.SizeBytes,
// 		Bucket:    meta.Bucket,
// 		ObjectKey: meta.ObjectKey,
// 	}

// 	if err := stream.Send(&pb.RetrieveImageResponse{
// 		Payload: &pb.RetrieveImageResponse_Meta{
// 			Meta: metaMsg,
// 		},
// 	}); err != nil {
// 		return status.Errorf(codes.Internal, "failed to send meta: %v", err)
// 	}

// 	// --- 3. Now stream chunks ---
// 	buf := make([]byte, 32*1024) // 32KB recommended

// 	for {
// 		n, readErr := reader.Read(buf)
// 		if readErr == io.EOF {
// 			break
// 		}
// 		if readErr != nil {
// 			return status.Errorf(codes.Internal, "failed to read image: %v", readErr)
// 		}

// 		chunk := &pb.ImageChunk{Data: buf[:n]}

// 		if err := stream.Send(&pb.RetrieveImageResponse{
// 			Payload: &pb.RetrieveImageResponse_Chunk{
// 				Chunk: chunk,
// 			},
// 		}); err != nil {
// 			return status.Errorf(codes.Internal, "failed to send chunk: %v", err)
// 		}
// 	}

// 	return nil
// }
