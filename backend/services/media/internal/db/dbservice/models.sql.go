package dbservice

import ct "social-network/shared/go/customtypes"

type File struct {
	Id        ct.Id  `validation:"nullable"` // db row Id
	Filename  string // the original name given by sender
	MimeType  string // content type
	SizeBytes int64
	Bucket    string `validation:"nullable"` // images, videos etc
	ObjectKey string `validation:"nullable"` // the name given to file in fileservice

	Visibility ct.FileVisibility
	Status     ct.UploadStatus `validation:"nullable"` // pending, processing, complete, failed

	Variant ct.FileVariant `validation:"nullable"` // thumb, small, medium, large, original
}

type Variant struct {
	Id        ct.Id `validation:"nullable"` // db row Id
	FileId    ct.Id
	Filename  string // the original name given by sender
	MimeType  string // content type
	SizeBytes int64
	Bucket    string `validation:"nullable"` // images, videos etc
	ObjectKey string `validation:"nullable"` // the name given to file in fileservice

	SrcBucket    string // the variants origin bucket
	SrcObjectKey string // the variants origin key

	Visibility ct.FileVisibility
	Status     ct.UploadStatus `validation:"nullable"` // pending, processing, complete, failed
	Variant    ct.FileVariant  `validation:"nullable"` // thumb, small, medium, large, original
}
