package validator

import (
	"net/http"
	"strings"
)

// returns true if filesize is less than maximum permitted
func MaxFileSize(filesize, maxFilesize int64) bool {
	return filesize < maxFilesize
}

// returns true if filetype is an image
func IsImage(file []byte) bool {
	fileType := http.DetectContentType(file)
	return strings.HasPrefix(fileType, "image/")
}
