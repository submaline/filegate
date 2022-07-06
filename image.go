package main

import (
	"github.com/h2non/bimg"
)

func ResizeAndConvertToWebp(buf []byte, width, height int) ([]byte, error) {
	// アスペクト比保たれて、黒埋めされる。
	resized, err := bimg.NewImage(buf).Resize(width, height)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(resized).Convert(bimg.WEBP)
}

func CompressNAndConvertWebp(buf []byte, percent int) ([]byte, error) {
	webpped, err := bimg.NewImage(buf).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(webpped).Process(bimg.Options{Quality: percent})
}
