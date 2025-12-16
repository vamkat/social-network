package application

import (
	"context"
	"fmt"
	"net/url"
	"social-network/services/media/internal/db/dbservice"
	ct "social-network/shared/go/customtypes"
	"time"

	"github.com/google/uuid"
)

type UploadImageReq struct {
	Filename   string
	MimeType   string
	SizeBytes  int64
	Visibility ct.FileVisibility
}

// Provides a fileId and an upload url targeted on bucket Originals defined on configs.
// Creates all variant entries provided in []variants for workers to later
// create asynchronously the compressed files.
func (m *MediaService) UploadImage(ctx context.Context,
	req UploadImageReq,
	exp time.Duration,
	variants []ct.FileVariant,
) (fileId ct.Id, upUrl string, err error) {

	if req.Filename == "" || req.MimeType == "" ||
		req.SizeBytes < 1 || !req.Visibility.IsValid() ||
		exp < 1 {
		return 0, "", ct.ErrValidation
	}
	objectKey := uuid.NewString()
	bucket := m.Cfgs.FileService.Buckets.Originals
	var url *url.URL
	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			fileId, err = m.Queries.CreateFile(ctx, dbservice.File{
				Filename:   req.Filename,
				MimeType:   req.MimeType,
				SizeBytes:  req.SizeBytes,
				Visibility: req.Visibility,
				Bucket:     m.Cfgs.FileService.Buckets.Originals,
				ObjectKey:  objectKey,
				Status:     ct.Complete,
				Variant:    ct.Original,
			})

			if err != nil {
				return err
			}

			for _, v := range variants {
				_, err := m.Queries.CreateVariant(ctx, dbservice.File{
					Filename:   req.Filename,
					MimeType:   req.MimeType,
					SizeBytes:  req.SizeBytes,
					Bucket:     m.Cfgs.FileService.Buckets.Variants,
					ObjectKey:  objectKey + "/" + v.String(),
					Visibility: req.Visibility,
					Status:     ct.Pending,
					Variant:    v,
				})
				if err != nil {
					return fmt.Errorf(
						"internal database error: %v failed to create variant %v for file with id: %v",
						err, v, fileId)
				}
			}

			url, err = m.Clients.GenerateUploadURL(ctx, bucket, objectKey, exp)
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

// Returns an image download URL for the requested imageId and Variant.
// If the variant is not available it falls back to the original file.
func (m *MediaService) GetImage(ctx context.Context,
	imgId ct.Id, variant ct.FileVariant,
) (downUrl string, err error) {
	if !imgId.IsValid() || !variant.IsValid() {
		return "", ct.ErrValidation
	}

	var req dbservice.File
	var url *url.URL

	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {
			switch variant {
			case ct.Original:
				req, err = m.Queries.GetFileById(ctx, imgId)
			default:
				req, err = m.Queries.GetVariant(ctx, imgId, variant)
				if req.Status != ct.Complete {
					req, err = m.Queries.GetFileById(ctx, imgId)
					if req.Status != ct.Complete {
						return fmt.Errorf("file validation status is %v", req.Status)
					}
				}
			}
			if err != nil {
				return err
			}

			url, err = m.Clients.GenerateDownloadURL(ctx, req.Bucket, req.ObjectKey, req.Visibility.SetExp())
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

// Not allowing originals in batch request
func (m *MediaService) GetImages(ctx context.Context,
	imgIds ct.Ids, variant ct.FileVariant,
) (downUrls map[ct.Id]string, err error) {

	if !imgIds.IsValid() || !variant.IsValid() || variant == ct.Original {
		return nil, ct.ErrValidation
	}
	var na ct.Ids
	var fms []dbservice.File

	errTx := m.txRunner.RunTx(ctx,
		func(q dbservice.Querier) error {

			fms, na, err = m.Queries.GetVariants(ctx, uniqueIds(imgIds), variant)
			if err != nil {
				return err
			}
			if len(na) != 0 {
				originals, err := m.Queries.GetFiles(ctx, na)
				if err != nil {
					return err
				}
				fms = append(fms, originals...)
			}
			return nil
		},
	)

	if errTx != nil {
		return nil, err
	}

	downUrls = make(map[ct.Id]string, len(fms))
	for _, fm := range fms {
		if fm.Status != ct.Complete {
			fmt.Printf("requested file %v validation status is %v", fm.Id, fm.Status)
			continue
		}
		url, err := m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
		if err != nil {
			return nil, err
		}
		downUrls[fm.Id] = url.String()
	}

	return downUrls, nil
}

func (m *MediaService) ValidateUpload(ctx context.Context,
	fileId ct.Id) error {
	if !fileId.IsValid() {
		return ct.ErrValidation
	}

	fileMeta, err := m.Queries.GetFileById(ctx, fileId)
	if err != nil {
		return err
	}

	if err := m.Clients.ValidateUpload(ctx, dbToExt(fileMeta)); err != nil {
		m.Queries.UpdateFileStatus(ctx, fileId, ct.Failed)
		return err
	}

	if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Complete); err != nil {
		return err
	}

	return nil
}

func uniqueIds(ids ct.Ids) ct.Ids {
	uniq := make(map[ct.Id]struct{}, len(ids))
	cleaned := make([]ct.Id, 0, len(ids))
	for _, id := range ids {
		if _, ok := uniq[id]; !ok {
			uniq[id] = struct{}{}
			cleaned = append(cleaned, id)
		}
	}
	return ct.Ids(cleaned)
}
