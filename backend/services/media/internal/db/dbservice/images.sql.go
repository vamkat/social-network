package dbservice

import (
	"context"
	"database/sql"
	ct "social-network/shared/go/customtypes"
)

func (q *Queries) CreateFile(
	ctx context.Context,
	fm File,
) (fileId ct.Id, err error) {

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
		ct.Pending,
	).Scan(&fileId)

	return fileId, err
}

func (q *Queries) GetFileById(
	ctx context.Context,
	fileId ct.Id,
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

	fm.Variant = ct.Original

	return fm, err
}

func (q *Queries) GetFiles(
	ctx context.Context,
	ids ct.Ids,
) ([]File, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	// Deduplicate IDs to avoid unnecessary DB work
	uniqueIds := make(ct.Ids, 0, len(ids))
	seen := make(map[ct.Id]struct{})
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			uniqueIds = append(uniqueIds, id)
		}
	}

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
		WHERE id = ANY($1)
	`

	rows, err := q.db.Query(ctx, query, uniqueIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var fm File
		if err := rows.Scan(
			&fm.Id,
			&fm.Filename,
			&fm.MimeType,
			&fm.SizeBytes,
			&fm.Bucket,
			&fm.ObjectKey,
			&fm.Visibility,
			&fm.Status,
		); err != nil {
			return nil, err
		}
		fm.Variant = ct.Original
		files = append(files, fm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (q *Queries) CreateVariant(
	ctx context.Context,
	fm File,
) (variantId ct.Id, err error) {

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
	fileId ct.Id,
	variant ct.FileVariant,
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

func (q *Queries) GetVariants(
	ctx context.Context,
	fileIds ct.Ids,
	variant ct.FileVariant,
) (fms []File, notComplete []ct.Id, err error) {

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
		WHERE f.id = ANY($1)
		  AND v.variant = $2
	`

	rows, err := q.db.Query(ctx, query, fileIds, variant)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	fms = make([]File, 0, len(fileIds))
	notComplete = make([]ct.Id, 0)

	for rows.Next() {
		var file File
		if err := rows.Scan(
			&file.Id,
			&file.Filename,
			&file.MimeType,
			&file.SizeBytes,
			&file.Bucket,
			&file.ObjectKey,
			&file.Visibility,
			&file.Status,
			&file.Variant,
		); err != nil {
			return nil, nil, err
		}

		if file.Status == ct.Complete {
			fms = append(fms, file)
		} else {
			notComplete = append(notComplete, file.Id)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return fms, notComplete, nil
}

func (q *Queries) UpdateVariantStatusAndSize(
	ctx context.Context,
	fileId ct.Id,
	status ct.UploadStatus,
	size int64,
) error {

	const query = `
		UPDATE file_variants
		SET 
			status = $2,
			size_bytes = $3
		WHERE file_id = $1
	`

	res, err := q.db.Exec(ctx, query, fileId, status, size)
	if err != nil {
		return err
	}

	if rows := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q *Queries) UpdateFileStatus(
	ctx context.Context,
	fileId ct.Id,
	status ct.UploadStatus,
) error {

	const query = `
		UPDATE files
		SET status = $2
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

func (q *Queries) GetPendingVariants(
	ctx context.Context) (pending []Variant, err error) {
	const query = `
		SELECT 
			fv.id,
			f.id as file_id, 
			f.filename, 
			fv.mime_type, 
			f.size_bytes, 
			fv.bucket, 
			fv.object_key, 
			f.visibility, 
			fv.variant
		FROM file_variants fv
		JOIN files f ON fv.file_id = f.id
		WHERE fv.status = 'processing' -- the variant status is changed to processing by db trigger when file status is set to complete
		  AND f.status = 'complete'
	`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fm Variant
		err := rows.Scan(
			&fm.Id,
			&fm.FileId,
			&fm.Filename,
			&fm.MimeType,
			&fm.SizeBytes,
			&fm.Bucket,
			&fm.ObjectKey,
			&fm.Visibility,
			&fm.Variant)
		if err != nil {
			return nil, err
		}
		pending = append(pending, fm)

	}

	return pending, err
}
