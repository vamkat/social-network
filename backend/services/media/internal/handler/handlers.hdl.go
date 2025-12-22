package handler

import (
	"context"
	"time"

	"social-network/services/media/internal/application"
	pb "social-network/shared/gen-go/media"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/mapping"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MediaHandler struct {
	pb.UnimplementedMediaServiceServer
	Application *application.MediaService
	// Configs     configs.Server
}

// Provides image id and an upload URL that can only be accessed through container DNS.
// All uploads are marked with a false validation tag and must be validated through ValidateUpload handler.
// Unvalidated uploads expire after the defined `lifecycle.Expiration` on file services configuration.
//
// Usage:
//
//	exp := time.Duration(10 * time.Minute).Seconds()
//	var MediaService media.MediaServiceClient
//
//	mediaRes, err := MediaService.UploadImage(r.Context(), &media.UploadImageRequest{
//		Filename:   httpReq.AvatarName,
//		MimeType:   httpReq.AvatarType,
//		SizeBytes:  httpReq.AvatarSize,
//		Visibility: media.FileVisibility_PUBLIC,
//		Variants: []media.ImgVariant{
//			media.FileVariant_THUMBNAIL,
//			media.FileVariant_LARGE,
//		},
//		ExpirationSeconds: int64(exp),
//	})
func (m *MediaHandler) UploadImage(ctx context.Context,
	req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request or file_meta is nil")
	}

	// Convert variants
	variants := make([]ct.FileVariant, len(req.Variants))
	for i, v := range req.Variants {
		variants[i] = mapping.PbToCtFileVariant(v)
	}
	appReq := application.UploadImageReq{
		Filename:   req.Filename,
		MimeType:   req.MimeType,
		SizeBytes:  req.SizeBytes,
		Visibility: mapping.PbToCtFileVisibility(req.Visibility),
	}
	// Call application
	fileId, upUrl, err := m.Application.UploadImage(
		ctx,
		appReq,
		time.Duration(req.ExpirationSeconds)*time.Second,
		variants,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	return &pb.UploadImageResponse{
		FileId:    int64(fileId),
		UploadUrl: upUrl,
	}, nil
}

// GetImage handles the gRPC request for retrieving an image download URL.
// Expiration time of link is set according to image visibility settings set on upload and
// is defined withing the methods of custom type 'FileVisibility'.
// Unvalidated uploads wont be fetched.
// If variant requested is not yet created the handler returns original
//
// Usage:
//
//	var MediaService media.MediaServiceClient
//	mediaRes, err := h.MediaService.GetImage(r.Context(), &media.GetImageRequest{
//		ImageId: 1,
//		Variant: media.FileVariant_ORIGINAL,
//	})
func (m *MediaHandler) GetImage(ctx context.Context,
	req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// Call application
	downUrl, err := m.Application.GetImage(ctx, ct.Id(req.ImageId), mapping.PbToCtFileVariant(req.Variant))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get image: %v", err)
	}

	return &pb.GetImageResponse{
		DownloadUrl: downUrl,
	}, nil
}

func (m *MediaHandler) GetImages(ctx context.Context,
	req *pb.GetImagesRequest) (*pb.GetImagesResponse, error) {
	if req == nil || req.ImgIds == nil {
		return nil, status.Error(codes.InvalidArgument, "request or img_ids is nil")
	}

	// Convert img_ids to ct.Ids
	ids := make(ct.Ids, len(req.ImgIds.ImgIds))
	for i, id := range req.ImgIds.ImgIds {
		ids[i] = ct.Id(id)
	}

	// Call application
	downUrls, failedIds, err := m.Application.GetImages(ctx, ids, mapping.PbToCtFileVariant(req.Variant))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get images: %v", err)
	}

	// Build response
	downloadUrls := make(map[int64]string, len(downUrls))
	for id, url := range downUrls {
		downloadUrls[int64(id)] = url
	}

	pbFailedIds := make([]*pb.FailedId, len(failedIds))
	for i, fid := range failedIds {
		pbFailedIds[i] = &pb.FailedId{
			FileId: int64(fid.Id),
			Status: mapping.CtToPbUploadStatus(fid.Status),
		}
	}

	return &pb.GetImagesResponse{
		DownloadUrls: downloadUrls,
		FailedIds:    pbFailedIds,
	}, nil
}

// Checks if the upload matches the pre defined file metadata and configs FileService file constraints.
// If validation fails file cannot be retrived and will be deleted from file service after 24 hours
func (m *MediaHandler) ValidateUpload(ctx context.Context,
	req *pb.ValidateUploadRequest) (*pb.ValidateUploadResponse, error) {
	if req == nil || req.FileId < 1 {
		return nil, status.Error(codes.InvalidArgument, "request or upload is nil")
	}

	// Call application
	url, err := m.Application.ValidateUpload(ctx, ct.Id(req.FileId), req.ReturnUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate upload: %v", err)
	}

	return &pb.ValidateUploadResponse{DownloadUrl: url}, nil
}
