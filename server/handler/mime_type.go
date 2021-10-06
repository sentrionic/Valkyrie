package handler

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// IsAllowedImageType determines if image is among types defined
// in map of allowed images
func isAllowedImageType(mimeType string) bool {
	_, exists := validImageTypes[mimeType]

	return exists
}

var validFileTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"audio/mp3":  true,
	"audio/wave": true,
}

// isAllowedFileType determines if the file is among types defined
// in map of allowed file types
func isAllowedFileType(mimeType string) bool {
	_, exists := validFileTypes[mimeType]

	return exists
}
