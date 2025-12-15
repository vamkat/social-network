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
			visibility
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
			visibility
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
	)

	return fm, err
}

func (q *Queries) CreateVariant(
	ctx context.Context,
	fm File,
) (variantId customtypes.Id, err error) {

	const query = `
		INSERT INTO file_variants (
			file_id,
			mime_type,
			size_bytes,
			variant,
			bucket,
			object_key,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = q.db.QueryRow(
		ctx,
		query,
		fm.Id, // fileId
		fm.MimeType,
		fm.SizeBytes,
		fm.Variant,
		fm.Bucket,
		fm.ObjectKey,
		fm.Status,
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
			v.mime_type,
			v.size_bytes,
			v.bucket,
			v.object_key,
			f.visibility,
			v.status,
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

func (q *Queries) UpdateVariantStatus(
	ctx context.Context,
	fileId customtypes.Id,
	variant customtypes.ImgVariant,
	status customtypes.UploadStatus,
) error {

	const query = `
		UPDATE file_variants
		SET status = $2,
		    updated_at = now()
		WHERE file_id = $1
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
