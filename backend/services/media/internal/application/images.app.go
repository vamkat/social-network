package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"social-network/services/media/internal/db/dbservice"
	"social-network/services/media/internal/mapping"
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

var (
	ErrValidation   = errors.New("validation error")
	ErrNotValidated = errors.New("file not validated")
	ErrFailed       = errors.New("file has failed validation")
	ErrInternal     = errors.New("internal error")
	ErrNotFound     = errors.New("not found")
)

// Provides a fileId and an upload url targeted on bucket Originals defined on configs.
// Creates all variant entries provided in []variants for workers to later
// create asynchronously the compressed files.
func (m *MediaService) UploadImage(ctx context.Context,
	req UploadImageReq,
	exp time.Duration,
	variants []ct.FileVariant,
) (fileId ct.Id, upUrl string, err error) {

	if err := m.validateUploadRequest(
		req,
		exp,
		variants,
	); err != nil {
		return 0, "", errors.Join(ErrValidation, fmt.Errorf("upload image: validation error %w", err))
	}

	objectKey := uuid.NewString()
	orignalsBucket := m.Cfgs.FileService.Buckets.Originals
	variantsBucket := m.Cfgs.FileService.Buckets.Variants
	var url *url.URL

	errTx := m.txRunner.RunTx(ctx,
		func(tx *dbservice.Queries) error {
			fileId, err = tx.CreateFile(ctx, dbservice.File{
				Filename:   req.Filename,
				MimeType:   req.MimeType,
				SizeBytes:  req.SizeBytes,
				Visibility: req.Visibility,
				Bucket:     orignalsBucket,
				ObjectKey:  objectKey,
				Status:     ct.Pending,
				Variant:    ct.Original,
			})

			if err != nil {
				return fmt.Errorf("upload image: creating original file db entry error for file %v: %w", req.Filename, err)
			}

			fmt.Printf("Creating variant db entries %v for file %v\n", variants, fileId)

			for _, v := range variants {
				if !v.IsValid() || v == ct.Original {
					continue
				}
				_, err := tx.CreateVariant(ctx, dbservice.File{
					Id:         fileId,
					Filename:   req.Filename,
					MimeType:   "image/webp",
					SizeBytes:  req.SizeBytes,
					Bucket:     variantsBucket,
					ObjectKey:  objectKey + "/" + v.String(),
					Visibility: req.Visibility,
					Status:     ct.Pending,
					Variant:    v,
				})
				if err != nil {
					return fmt.Errorf(
						"upload image: internal database error: failed to create variant %v for file with id %v: %w",
						v, fileId, err)
				}
			}

			url, err = m.Clients.GenerateUploadURL(ctx, orignalsBucket, objectKey, exp)
			if err != nil {
				return fmt.Errorf(
					"upload image: S3 error: failed to create upload url for file with id %v: %w",
					fileId, err)
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
		return "", ErrValidation
	}

	var req dbservice.File
	var url *url.URL

	errTx := m.txRunner.RunTx(ctx,
		func(tx *dbservice.Queries) error {
			switch variant {
			case ct.Original:
				req, err = tx.GetFileById(ctx, imgId)
				if req.Status != ct.Complete {
					return errors.Join(
						ErrNotValidated,
						fmt.Errorf("file validation status is %v", req.Status))
				}
			default:
				req, err = tx.GetVariant(ctx, imgId, variant)
				if req.Status != ct.Complete {
					req, err = tx.GetFileById(ctx, imgId)
					if req.Status != ct.Complete {
						if req.Status == ct.Failed {
							return errors.Join(
								ErrFailed,
								fmt.Errorf("file validation status is %v", req.Status))
						}
						return errors.Join(
							ErrNotValidated,
							fmt.Errorf("file validation status is %v", req.Status))
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

type FailedId struct {
	Id     ct.Id
	Status ct.UploadStatus
}

// Returns a id to download url pairs for
// an array of file ids and the prefered variant.
// Variant is common for all ids. If a variant is not present
// returns url for the original format.
// GetImages does not accept original variants in batch request
func (m *MediaService) GetImages(ctx context.Context,
	imgIds ct.Ids, variant ct.FileVariant,
) (downUrls map[ct.Id]string, failedIds []FailedId, err error) {

	fmt.Println("received ids", imgIds)
	// fmt.Println(imgIds.IsValid())
	// fmt.Println(variant.IsValid())
	// fmt.Println(variant == ct.Original)

	if !imgIds.IsValid() || !variant.IsValid() || variant == ct.Original {
		//fmt.Println("validation error", ct.ErrValidation)
		return nil, nil, ErrValidation
	}
	var na ct.Ids
	var fms []dbservice.File

	errTx := m.txRunner.RunTx(ctx, func(tx *dbservice.Queries) error {

		fms, na, err = tx.GetVariants(ctx, imgIds.Unique(), variant)
		if err != nil {
			return err
		}
		//fmt.Println("variants", fms)
		if len(na) != 0 {
			originals, err := tx.GetFiles(ctx, na)
			if err != nil {
				return err
			}
			fms = append(fms, originals...)
			//fmt.Println("fms", fms)
		}
		return nil
	})

	if errTx != nil {
		return nil, nil, err
	}
	failedIds = []FailedId{}
	downUrls = make(map[ct.Id]string, len(fms))
	for _, fm := range fms {
		if fm.Status != ct.Complete {
			failedIds = append(failedIds, FailedId{Id: fm.Id, Status: fm.Status})
			fmt.Printf("requested file %v validation status is %v", fm.Id, fm.Status)
			continue
		}
		url, err := m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
		if err != nil {
			return nil, nil, err
		}
		downUrls[fm.Id] = url.String()
		//For testing with seeds
		//downUrls[fm.Id] = fm.Filename
	}
	fmt.Println("download urls", downUrls)
	fmt.Println("failed ids", failedIds)
	return downUrls, failedIds, nil
}

// This is a call to validate an already uploaded file.
// Unvalidated files expire in 24 hours and are automatically
// deleted from file service.
func (m *MediaService) ValidateUpload(ctx context.Context,
	fileId ct.Id, returnURL bool) (url string, err error) {
	if !fileId.IsValid() {
		return url, fmt.Errorf("fileId invalid, is not above 0: %w", ErrValidation)
	}

	fileMeta, err := m.Queries.GetFileById(ctx, fileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return url, ErrNotFound
		}
		return url, err
	}

	if fileMeta.Status == ct.Failed {
		return url, ErrFailed
	}

	if fileMeta.Status != ct.Complete {
		if errOuter := m.Clients.ValidateUpload(ctx, mapping.DbToModel(fileMeta)); errOuter != nil {
			if err := m.Clients.DeleteFile(ctx, fileMeta.Bucket, fileMeta.ObjectKey); err != nil {
				return url, errors.Join(ErrFailed, errOuter, err)
			}
			if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Failed); err != nil {
				return url, errors.Join(ErrFailed, errOuter, err)
			}
			return url, errors.Join(ErrFailed, errOuter)
		}

		if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Complete); err != nil {
			return url, err
		}

		log.Printf("Media Service: FileId %v successfully validated and marked as Complete", fileId)
	}

	if returnURL {
		u, err := m.Clients.GenerateDownloadURL(ctx,
			fileMeta.Bucket,
			fileMeta.ObjectKey,
			fileMeta.Visibility.SetExp(),
		)
		if err != nil {
			log.Printf("failed to fetch url for file %v\n", fileId)
			return "", nil
		}
		url = u.String()
	}
	return url, nil
}
