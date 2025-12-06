package models

type ImageMeta struct {
	Id        int64  // db row Id
	Filename  string // the original name given by sender
	MimeType  string // content type
	SizeBytes int64
	Bucket    string // images, videos etc
	ObjectKey string // the name given to file in minIO
}
