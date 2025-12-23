package application

import (
	"context"
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
		func(tx *dbservice.Queries) error {
			fileId, err = tx.CreateFile(ctx, dbservice.File{
				Filename:   req.Filename,
				MimeType:   req.MimeType,
				SizeBytes:  req.SizeBytes,
				Visibility: req.Visibility,
				Bucket:     m.Cfgs.FileService.Buckets.Originals,
				ObjectKey:  objectKey,
				Status:     ct.Pending,
				Variant:    ct.Original,
			})

			if err != nil {
				return err
			}

			fmt.Printf("Creating variants %v for file %v\n", variants, fileId)

			for _, v := range variants {
				_, err := tx.CreateVariant(ctx, dbservice.File{
					Id:         fileId,
					Filename:   req.Filename,
					MimeType:   "image/webp",
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
		func(tx *dbservice.Queries) error {
			switch variant {
			case ct.Original:
				req, err = tx.GetFileById(ctx, imgId)
			default:
				req, err = tx.GetVariant(ctx, imgId, variant)
				if req.Status != ct.Complete {
					req, err = tx.GetFileById(ctx, imgId)
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
		return nil, nil, ct.ErrValidation
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
		return url, fmt.Errorf("fileId invalid, is not above 0: %w", ct.ErrValidation)
	}

	fileMeta, err := m.Queries.GetFileById(ctx, fileId)
	if err != nil {
		return url, err
	}

	if err := m.Clients.ValidateUpload(ctx, mapping.DbToModel(fileMeta)); err != nil {
		if err := m.Clients.DeleteFile(ctx, fileMeta.Bucket, fileMeta.ObjectKey); err != nil {
			return url, err
		}
		if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Failed); err != nil {
			return url, err
		}
		return url, err
	}

	if err := m.Queries.UpdateFileStatus(ctx, fileId, ct.Complete); err != nil {
		return url, err
	}

	log.Printf("Media Service: FileId %v successfully validated and marked as Complete", fileId)

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

// TODO: Make this part of ct.Ids methods
// func uniqueIds(ids ct.Ids) ct.Ids {
// 	uniq := make(map[ct.Id]struct{}, len(ids))
// 	cleaned := make([]ct.Id, 0, len(ids))
// 	for _, id := range ids {
// 		if _, ok := uniq[id]; !ok {
// 			uniq[id] = struct{}{}
// 			cleaned = append(cleaned, id)
// 		}
// 	}
// 	return ct.Ids(cleaned)
// }
