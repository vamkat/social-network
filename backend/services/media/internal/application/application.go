package application

import (
	"context"
	"fmt"
	"io"
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/db/sqlc"
	"social-network/services/media/internal/utils"
	"social-network/shared/go/customtypes"
	"social-network/shared/go/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type MediaService struct {
	Pool     *pgxpool.Pool
	Clients  *client.Clients
	Queries  sqlc.Querier
	txRunner TxRunner
}

func NewMediaService(pool *pgxpool.Pool, clients *client.Clients, queries sqlc.Querier) *MediaService {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := queries.(*sqlc.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}
	return &MediaService{
		Pool:     pool,
		Clients:  clients,
		Queries:  queries,
		txRunner: txRunner,
	}
}

func (m *MediaService) SaveImage(ctx context.Context, file []byte, filename string) (models.FileMeta, error) {
	contentType, err := utils.ValidateImage(file, filename)
	if err != nil {
		return models.FileMeta{}, err
	}
	info, err := m.Clients.UploadToMinIO(ctx, file, filename, "images", contentType)
	if err != nil {
		return models.FileMeta{}, err
	}
	row, err := m.Queries.SaveImageMetadata(ctx,
		sqlc.SaveImageMetadataParams{
			OriginalName: filename,
			Bucket:       info.Bucket,
			ObjectKey:    info.Key,
			MimeType:     contentType,
			SizeBytes:    info.Size,
		},
	)
	if err != nil {
		return models.FileMeta{}, err
	}
	return models.FileMeta{
		Id:        row[0].ID,
		Filename:  filename,
		MimeType:  contentType,
		SizeBytes: int64(len(file)),
		ObjectKey: info.Key,
		Bucket:    info.Bucket,
	}, nil
}

func (m *MediaService) RetriveImageById(ctx context.Context, imageId customtypes.Id) (reader io.ReadCloser, meta models.FileMeta, err error) {
	// Call db
	if !imageId.IsValid() {
		return nil, meta, fmt.Errorf("invalid image id: %v", imageId)
	}

	var imageMeta sqlc.Image
	resp, err := m.Queries.GetImageById(ctx, imageId.Int64())
	if err != nil {
		return reader, meta, err
	}

	imageMeta = resp[0]
	info := models.FileMeta{
		Id:        imageMeta.ID,
		Filename:  imageMeta.OriginalName,
		MimeType:  imageMeta.MimeType,
		SizeBytes: imageMeta.SizeBytes,
		Bucket:    imageMeta.Bucket,
		ObjectKey: imageMeta.ObjectKey,
	}

	obj, err := m.Clients.GetFromMiniIo(ctx, info)
	if err != nil {
		return reader, meta, err
	}
	return io.ReadCloser(obj), meta, nil
}
