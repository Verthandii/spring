package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// MimeFromIncipit returns the mime type of an image file from its first few
// bytes or the empty string if the file does not look like a known file type
func MimeFromIncipit(incipit []byte) bool {
	_, _, err := image.Decode(bytes.NewReader(incipit))
	return err == nil
}
