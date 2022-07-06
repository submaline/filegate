package main

import "net/http"

func IsAllowedFileType(data []byte) (string, bool) {
	mime := http.DetectContentType(data)
	switch mime {
	case "image/png", "image/jpeg", "image/gif":
		// ok
	default:
		return mime, false
	}

	return mime, true
}
