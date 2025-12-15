package dbservice

import ct "social-network/shared/go/customtypes"

type File struct {
	Id        ct.Id  // db row Id
	Filename  string // the original name given by sender
	MimeType  string // content type
	SizeBytes int64
	Bucket    string // images, videos etc
	ObjectKey string // the name given to file in fileservice

	Visibility ct.FileVisibility
	Status     ct.UploadStatus // pending, processing, complete, failed

	Variant ct.ImgVariant // thumb, small, medium, large, original
}
