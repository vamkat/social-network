package models

import ct "social-network/shared/go/customtypes"

type FileMeta struct {
	Id        ct.Id  // db row Id
	Filename  string `json:"filename"` // the original name given by sender
	MimeType  string // content type
	SizeBytes int64
	Bucket    string // images, videos etc
	ObjectKey string // the name given to file in fileservice

	Visibility ct.FileVisibility
	Variant    ct.ImgVariant // thumb, small, medium, large
}
