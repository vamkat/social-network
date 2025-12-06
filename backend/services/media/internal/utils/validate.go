package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var (
	// Allowed MIME types
	allowedMIMEs = map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/svg+xml": true, // sanitize recommended
	}

	// Allowed file extensions
	allowedExt = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".svg":  true,
	}

	// Max upload size (e.g. 10 MB)
	maxFileSize = int64(10 * 1024 * 1024)
)

func ValidateImage(fileContent []byte, filename string) (string, error) {
	// === Size check ===
	if len(fileContent) == 0 {
		return "", errors.New("file is empty")
	}
	if int64(len(fileContent)) > maxFileSize {
		return "", fmt.Errorf("file too large (%d bytes)", len(fileContent))
	}

	// === Detect MIME type ===
	m := mimetype.Detect(fileContent)
	mime := m.String()

	if !allowedMIMEs[mime] {
		return "", fmt.Errorf("invalid MIME type: %s", mime)
	}

	// === Validate extension ===
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExt[ext] {
		return "", fmt.Errorf("invalid file extension: %s", ext)
	}

	// === Extra SVG security ===
	if mime == "image/svg+xml" {
		if err := sanitizeSVG(fileContent); err != nil {
			return "", fmt.Errorf("unsafe SVG: %v", err)
		}
	}

	return mime, nil
}

// OPTIONAL: minimal SVG sanitization
// (You can enhance this as needed)
func sanitizeSVG(data []byte) error {
	s := strings.ToLower(string(data))

	// Block embedded scripts
	if strings.Contains(s, "<script") ||
		strings.Contains(s, "onload=") ||
		strings.Contains(s, "onerror=") {
		return errors.New("SVG contains scripts or event handlers")
	}

	return nil
}
