package application

import (
	"context"
	"net/url"
	"social-network/services/media/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
	"time"

	"github.com/google/uuid"
)

func (m *MediaService) UploadImage(ctx context.Context,
	fm md.FileMeta, exp time.Duration) (fileId ct.Id, upUrl string, err error) {
	if err := ct.ValidateStruct(fm); err != nil {
		return fileId, "", err
	}

	var url *url.URL
	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			fm.ObjectKey = uuid.NewString()
			fileId, err = m.Queries.CreateFile(ctx, dbservice.File{
				Filename:   fm.Filename,
				MimeType:   fm.MimeType,
				SizeBytes:  fm.SizeBytes,
				Bucket:     m.Cfgs.FileService.Buckets.Originals,
				ObjectKey:  fm.ObjectKey,
				Visibility: fm.Visibility,
				Status:     ct.Pending,
				Variant:    ct.Large,
			})
			if err != nil {
				return err
			}

			url, err = m.Clients.GenerateUploadURL(ctx, fm.Bucket, fm.ObjectKey, exp)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if errTx != nil {
		return 0, "", errTx
	}
	return fileId, url.String(), nil
}

func (m *MediaService) GetImage(ctx context.Context,
	imgId ct.Id, variant ct.ImgVariant,
) (downUrl string, err error) {
	if !imgId.IsValid() || !variant.IsValid() {
		return "", ct.ErrValidation
	}

	var fm dbservice.File
	var url *url.URL

	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			switch variant {
			case ct.Large:
				fm, err = m.Queries.GetFileById(ctx, imgId)
			default:
				fm, err = m.Queries.GetVariant(ctx, imgId, variant)
			}
			if err != nil {
				return err
			}

			url, err = m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
			if err != nil {
				return err
			}
			return nil
		},
	)

	if errTx != nil {
		return "", err
	}
	return url.String(), err
}

func (m *MediaService) ValidateUpload(ctx context.Context,
	upload md.FileMeta) error {
	if err := ct.ValidateStruct(upload); err != nil {
		return err
	}

	if err := m.Clients.ValidateUpload(ctx, upload); err != nil {
		return err
	}

	return nil
}
