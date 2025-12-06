package application

import (
	"context"
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/db/sqlc"
	"social-network/services/media/internal/utils"

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

type ImageMeta struct {
	Id        string
	MimeType  string
	SizeBytes int64
	Bucket    string
	ObjectKey string
}

func (m MediaService) SaveImage(ctx context.Context, file []byte, filename string) (ImageMeta, error) {
	contentType, err := utils.ValidateImage(file, filename)
	if err != nil {
		return ImageMeta{}, err
	}
	info, err := m.Clients.UploadToMinIO(ctx, file, filename, "images", contentType)
	if err != nil {
		return ImageMeta{}, err
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
		return ImageMeta{}, err
	}
	return ImageMeta{
		Id:        row[0].ID.String(),
		MimeType:  contentType,
		SizeBytes: int64(len(file)),
		ObjectKey: info.Key,
		Bucket:    info.Bucket,
	}, nil
}
