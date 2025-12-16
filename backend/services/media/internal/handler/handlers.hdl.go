package handler

import (
	"context"
	"time"

	"social-network/services/media/internal/application"
	pb "social-network/shared/gen-go/media"
	ct "social-network/shared/go/customtypes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MediaHandler struct {
	pb.UnimplementedMediaServiceServer
	Application *application.MediaService
	Port        string
}

// pbToCtImgVariant converts protobuf ImgVariant to customtypes ImgVariant
func pbToCtImgVariant(v pb.ImgVariant) ct.FileVariant {
	switch v {
	case pb.ImgVariant_THUMBNAIL:
		return ct.ImgThumbnail
	case pb.ImgVariant_SMALL:
		return ct.ImgSmall
	case pb.ImgVariant_MEDIUM:
		return ct.ImgMedium
	case pb.ImgVariant_LARGE:
		return ct.ImgLarge
	case pb.ImgVariant_ORIGINAL:
		return ct.Original
	default:
		return ct.FileVariant("") // invalid, but handle gracefully
	}
}

// ctToPbImgVariant converts customtypes ImgVariant to protobuf ImgVariant
// func ctToPbImgVariant(v ct.ImgVariant) pb.ImgVariant {
// 	switch v {
// 	case ct.Thumbnail:
// 		return pb.ImgVariant_THUMBNAIL
// 	case ct.Small:
// 		return pb.ImgVariant_SMALL
// 	case ct.Medium:
// 		return pb.ImgVariant_MEDIUM
// 	case ct.Large:
// 		return pb.ImgVariant_LARGE
// 	case ct.Original:
// 		return pb.ImgVariant_ORIGINAL
// 	default:
// 		return pb.ImgVariant_IMG_VARIANT_UNSPECIFIED
// 	}
// }

// pbToCtFileVisibility converts protobuf FileVisibility to customtypes FileVisibility
func pbToCtFileVisibility(v pb.FileVisibility) ct.FileVisibility {
	switch v {
	case pb.FileVisibility_PRIVATE:
		return ct.Private
	case pb.FileVisibility_PUBLIC:
		return ct.Public
	default:
		return ct.FileVisibility("") // invalid
	}
}

// ctToPbFileVisibility converts customtypes FileVisibility to protobuf FileVisibility
// func ctToPbFileVisibility(v ct.FileVisibility) pb.FileVisibility {
// 	switch v {
// 	case ct.Private:
// 		return pb.FileVisibility_PRIVATE
// 	case ct.Public:
// 		return pb.FileVisibility_PUBLIC
// 	default:
// 		return pb.FileVisibility_FILE_VISIBILITY_UNSPECIFIED
// 	}
// }

// pbToMdFileMeta converts protobuf FileMeta to models FileMeta
// func pbToMdFileMeta(fm *pb.FileMeta) md.FileMeta {
// 	return md.FileMeta{
// 		Id:         ct.Id(fm.Id),
// 		Filename:   fm.Filename,
// 		MimeType:   fm.MimeType,
// 		SizeBytes:  fm.SizeBytes,
// 		Bucket:     fm.Bucket,
// 		ObjectKey:  fm.ObjectKey,
// 		Visibility: pbToCtFileVisibility(fm.Visibility),
// 		Variant:    pbToCtImgVariant(fm.Variant),
// 	}
// }

// mdToPbFileMeta converts models FileMeta to protobuf FileMeta
// func mdToPbFileMeta(fm md.FileMeta) *pb.FileMeta {
// 	return &pb.FileMeta{
// 		Id:         int64(fm.Id),
// 		Filename:   fm.Filename,
// 		MimeType:   fm.MimeType,
// 		SizeBytes:  fm.SizeBytes,
// 		Bucket:     fm.Bucket,
// 		ObjectKey:  fm.ObjectKey,
// 		Visibility: ctToPbFileVisibility(fm.Visibility),
// 		Variant:    ctToPbImgVariant(fm.Variant),
// 	}
// }

// UploadImage handles the gRPC request for uploading an image
func (m *MediaHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request or file_meta is nil")
	}

	// Convert variants
	variants := make([]ct.FileVariant, len(req.Variants))
	for i, v := range req.Variants {
		variants[i] = pbToCtImgVariant(v)
	}
	appReq := application.UploadImageReq{
		Filename:   req.Filename,
		MimeType:   req.MimeType,
		SizeBytes:  req.SizeBytes,
		Visibility: pbToCtFileVisibility(req.Visibility),
	}
	// Call application
	fileId, upUrl, err := m.Application.UploadImage(ctx, appReq, time.Duration(req.ExpirationSeconds)*time.Second, variants)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload image: %v", err)
	}

	return &pb.UploadImageResponse{
		FileId:    int64(fileId),
		UploadUrl: upUrl,
	}, nil
}

// GetImage handles the gRPC request for retrieving an image download URL
func (m *MediaHandler) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// Call application
	downUrl, err := m.Application.GetImage(ctx, ct.Id(req.ImageId), pbToCtImgVariant(req.Variant))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get image: %v", err)
	}

	return &pb.GetImageResponse{
		DownloadUrl: downUrl,
	}, nil
}

// ValidateUpload handles the gRPC request for validating upload metadata
func (m *MediaHandler) ValidateUpload(ctx context.Context, req *pb.ValidateUploadRequest) (*emptypb.Empty, error) {
	if req == nil || req.FileId < 1 {
		return nil, status.Error(codes.InvalidArgument, "request or upload is nil")
	}

	// Call application
	err := m.Application.ValidateUpload(ctx, ct.Id(req.FileId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate upload: %v", err)
	}

	return &emptypb.Empty{}, nil
}
