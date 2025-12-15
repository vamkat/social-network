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
	Status     ct.UploadStatus // pending, complete, failed
	Variant    ct.ImgVariant   // thumb, small, medium, large
}

type Variant struct {
	Id        ct.Id
	FileId    ct.Id
	Variant   ct.ImgVariant
	Bucket    string // images, videos etc
	ObjectKey string // the name given to file in fileservice

	Width  int
	Height int
}
