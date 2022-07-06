package main

import (
	"github.com/h2non/bimg"
	"io"
	"os"
	"testing"
)

func TestResizeAndConvertToWebp(t *testing.T) {
	inPath := ""
	outPath := inPath + ".webp"

	fp, err := os.Open(inPath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	buf, err := io.ReadAll(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	newImg, err := ResizeAndConvertToWebp(buf, 250, 250)
	if err != nil {
		t.Fatalf("failed to resize & convert image: %v", err)
	}

	err = bimg.Write(outPath, newImg)
	if err != nil {
		t.Fatalf("faied to write image: %v", err)
	}
}
