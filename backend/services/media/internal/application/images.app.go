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
	ct "social-network/shared/go/ct"
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

	if err := m.validateUploadRequest(
		req,
		exp,
		variants,
	); err != nil {
		return 0, "", Wrap(
			ErrReqValidation,
			err,
			"upload image:",
		)
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
				return Wrap(
					ErrInternal,
					err,
					fmt.Sprintf(
						"UploadImage: creating original file db entry error for file %v",
						req.Filename,
					),
				)
			}

			for _, v := range variants {
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
					return Wrap(
						ErrInternal,
						err,
						fmt.Sprintf(
							"UploadImage: failed to create variant %v for file with id %v",
							v, fileId),
					)
				}
			}

			url, err = m.Clients.GenerateUploadURL(ctx, orignalsBucket, objectKey, exp)
			if err != nil {
				return Wrap(
					ErrInternal,
					err,
					fmt.Sprintf(
						"UploadImage: S3 error: failed to create upload url for file with id %v:",
						fileId),
				)
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
func (m *MediaService) GetImage(
	ctx context.Context,
	imgId ct.Id,
	variant ct.FileVariant,
) (string, error) {
	errMsg := fmt.Sprintf("get image err: id: %d variant: %s", imgId, variant)

	if err := ct.ValidateBatch(imgId, variant); err != nil {
		return "", Wrap(ErrReqValidation, err, errMsg)
	}

	var fm dbservice.File

	err := m.txRunner.RunTx(ctx, func(tx *dbservice.Queries) error {
		var err error

		if variant == ct.Original {
			fm, err = tx.GetFileById(ctx, imgId)
		} else {
			fm, err = tx.GetVariant(ctx, imgId, variant)
			if errors.Is(err, sql.ErrNoRows) {
				fm, err = tx.GetFileById(ctx, imgId)
			}
		}

		if err != nil {
			return mapDBError(err)
		}

		return validateFileStatus(fm)
	})

	if err != nil {
		return "", Wrap(nil, err, errMsg)
	}

	u, err := m.Clients.GenerateDownloadURL(
		ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp(),
	)
	if err != nil {
		return "", Wrap(ErrInternal, err, errMsg)
	}

	return u.String(), nil
}

type FailedId struct {
	Id     ct.Id
	Status ct.UploadStatus
}

// Returns a id to download url pairs for
// an array of file ids and the prefered variant.
// Precondition for returning a file is the variant requested to exist in the database.
// Variant is common for all ids. If a variant is present but not completed
// returns url for the original format.
// GetImages does not accept original variants in batch request
func (m *MediaService) GetImages(ctx context.Context,
	imgIds ct.Ids, variant ct.FileVariant,
) (downUrls map[ct.Id]string, failedIds []FailedId, err error) {

	errMsg := fmt.Sprintf("get images: ids: %v variant: %s", imgIds, variant)

	if err := ct.ValidateBatch(imgIds, variant); err != nil {
		return nil, nil, Wrap(ErrReqValidation, err, errMsg)
	}

	var missingVariants ct.Ids
	var fms []dbservice.File

	errTx := m.txRunner.RunTx(ctx, func(tx *dbservice.Queries) error {
		fms, missingVariants, err = tx.GetVariants(ctx, imgIds.Unique(), variant)
		if err != nil {
			return mapDBError(err)
		}

		if len(missingVariants) != 0 {
			originals, err := tx.GetFiles(ctx, missingVariants)
			if err != nil {
				return mapDBError(err)
			}
			fms = append(fms, originals...)
		}
		return nil
	})

	if errTx != nil {
		return nil, nil, err
	}

	failedIds = []FailedId{}
	downUrls = make(map[ct.Id]string, len(fms))
	for _, fm := range fms {
		if err := validateFileStatus(fm); err != nil {
			failedIds = append(failedIds, FailedId{Id: fm.Id, Status: fm.Status})
			log.Println(err.Error())
			continue
		}
		url, err := m.Clients.GenerateDownloadURL(ctx, fm.Bucket, fm.ObjectKey, fm.Visibility.SetExp())
		if err != nil {
			return nil, nil, errors.Join(
				ErrInternal,
				fmt.Errorf(
					"GetImages: GetFiles error: %w, file meta: %v",
					err,
					fm,
				),
			)
		}
		downUrls[fm.Id] = url.String()

		//For testing with seeds
		//downUrls[fm.Id] = fm.Filename
	}
	return downUrls, failedIds, nil
}

// This is a call to validate an already uploaded file.
// Unvalidated files expire in 24 hours and are automatically
// deleted from file service.
func (m *MediaService) ValidateUpload(ctx context.Context,
	fileId ct.Id, returnURL bool) (url string, err error) {

	errMsg := fmt.Sprintf("validate upload: file id: %d", fileId)

	if err := fileId.Validate(); err != nil {
		return url, Wrap(ErrReqValidation, err, errMsg)
	}

	fileMeta, err := m.Queries.GetFileById(ctx, fileId)
	if err != nil {
		return "", Wrap(nil, mapDBError(err), errMsg)
	}

	if fileMeta.Status == ct.Failed {
		return url, ErrFailed
	}

	if fileMeta.Status != ct.Complete {
		if errOuter := m.Clients.ValidateUpload(ctx, mapping.DbToModel(fileMeta)); errOuter != nil {
			if err := m.Clients.DeleteFile(ctx, fileMeta.Bucket, fileMeta.ObjectKey); err != nil {
				return url, Wrap(ErrFailed, errors.Join(errOuter, err), errMsg)
			}
			if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Failed); err != nil {
				return url, Wrap(ErrFailed, errors.Join(errOuter, err), errMsg)
			}
			return url, Wrap(ErrFailed, errors.Join(errOuter, err), errMsg)
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
