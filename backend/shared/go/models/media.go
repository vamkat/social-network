package models

type ImageMeta struct {
	Id        int64
	Filename  string
	MimeType  string
	SizeBytes int64
	Bucket    string
	ObjectKey string
}
