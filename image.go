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
	webpImg, err := bimg.NewImage(buf).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(webpImg).Process(bimg.Options{Quality: percent})
}

func ThumbnailWebp(buf []byte) ([]byte, error) {
	// アスペクト比保たれて、黒埋めされる。
	resized, err := bimg.NewImage(buf).Resize(640, 480)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(resized).Process(bimg.Options{Quality: 80})
}
