package dbservice

import (
	"context"
	"social-network/shared/go/customtypes"
)

type Querier interface {
	CreateFile(ctx context.Context, fm File) (fileId customtypes.Id, err error)
	GetFileById(ctx context.Context, fileId customtypes.Id) (fm File, err error)
	CreateVariant(ctx context.Context, fm Variant) (fileId customtypes.Id, err error)
	GetVariant(ctx context.Context, fileId customtypes.Id,
		variant customtypes.ImgVariant) (fm File, err error)
	UpdateStatus(ctx context.Context, fileId customtypes.Id, status customtypes.UploadStatus) error
}

var _ Querier = (*Queries)(nil)
