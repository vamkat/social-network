package dbservice

import (
	"context"
	"database/sql"
	"social-network/shared/go/customtypes"
)

func (q *Queries) CreateFile(
	ctx context.Context,
	fm File,
) (fileId customtypes.Id, err error) {

	const query = `
		INSERT INTO files (
			filename,
			mime_type,
			size_bytes,
			bucket,
			object_key,
			visibility,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = q.db.QueryRow(
		ctx,
		query,
		fm.Filename,
		fm.MimeType,
		fm.SizeBytes,
		fm.Bucket,
		fm.ObjectKey,
		fm.Visibility,
		fm.Status,
	).Scan(&fileId)

	return fileId, err
}

func (q *Queries) GetFileById(
	ctx context.Context,
	fileId customtypes.Id,
) (fm File, err error) {

	const query = `
		SELECT
			id,
			filename,
			mime_type,
			size_bytes,
			bucket,
			object_key,
			visibility,
			status
		FROM files
		WHERE id = $1
	`

	err = q.db.QueryRow(ctx, query, fileId).Scan(
		&fm.Id,
		&fm.Filename,
		&fm.MimeType,
		&fm.SizeBytes,
		&fm.Bucket,
		&fm.ObjectKey,
		&fm.Visibility,
		&fm.Status,
	)

	return fm, err
}

func (q *Queries) CreateVariant(
	ctx context.Context,
	fm Variant,
) (variantId customtypes.Id, err error) {

	const query = `
		INSERT INTO file_variants (
			file_id,
			variant,
			bucket,
			object_key,
			width,
			height
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err = q.db.QueryRow(
		ctx,
		query,
		fm.FileId,
		fm.Variant,
		fm.Bucket,
		fm.ObjectKey,
		fm.Width,
		fm.Height,
	).Scan(&variantId)

	return variantId, err
}

func (q *Queries) GetVariant(
	ctx context.Context,
	fileId customtypes.Id,
	variant customtypes.ImgVariant,
) (fm File, err error) {

	const query = `
		SELECT
			f.id,
			f.filename,
			f.mime_type,
			f.size_bytes,
			v.bucket,
			v.object_key,
			f.visibility,
			f.status,
			v.variant
		FROM files f
		JOIN file_variants v ON v.file_id = f.id
		WHERE f.id = $1
		  AND v.variant = $2
	`

	err = q.db.QueryRow(ctx, query, fileId, variant).Scan(
		&fm.Id,
		&fm.Filename,
		&fm.MimeType,
		&fm.SizeBytes,
		&fm.Bucket,
		&fm.ObjectKey,
		&fm.Visibility,
		&fm.Status,
		&fm.Variant,
	)

	return fm, err
}

func (q *Queries) UpdateStatus(
	ctx context.Context,
	fileId customtypes.Id,
	status customtypes.UploadStatus,
) error {

	const query = `
		UPDATE files
		SET status = $2,
		    updated_at = now()
		WHERE id = $1
	`

	res, err := q.db.Exec(ctx, query, fileId, status)
	if err != nil {
		return err
	}

	if rows := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
