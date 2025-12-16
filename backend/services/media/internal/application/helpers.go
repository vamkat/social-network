package application

import (
	"social-network/services/media/internal/db/dbservice"
	"social-network/shared/go/customtypes"
	"social-network/shared/go/models"
)

func extToDbFile(meta models.FileMeta) dbservice.File {
	return dbservice.File{
		Id:         meta.Id,
		Filename:   meta.Filename,
		MimeType:   meta.MimeType,
		SizeBytes:  meta.SizeBytes,
		Bucket:     meta.Bucket,
		ObjectKey:  meta.ObjectKey,
		Visibility: meta.Visibility,
		Variant:    meta.Variant,
	}
}

func extToDbFileWithStatus(meta models.FileMeta, status customtypes.UploadStatus) dbservice.File {
	f := extToDbFile(meta)
	f.Status = status
	return f
}

func dbToExt(file dbservice.File) models.FileMeta {
	return models.FileMeta{
		Id:         file.Id,
		Filename:   file.Filename,
		MimeType:   file.MimeType,
		SizeBytes:  file.SizeBytes,
		Bucket:     file.Bucket,
		ObjectKey:  file.ObjectKey,
		Visibility: file.Visibility,
		Variant:    file.Variant,
	}
}
