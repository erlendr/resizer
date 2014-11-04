package main

import (
	"bytes"
	"github.com/erlendr/store"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
)

func main() {
	var filename = "45726c656e642d526f736a6f2e6a7067e691a4d41d8cd98f00b204e9800998ecf8427e.jpg"
	rc := store.Download(filename)

	println("Resizer - File " + filename + " downloaded")

	var buf = resizeImage(rc, 100)
	var reader = bytes.NewReader(buf.Bytes())
	store.UploadReader("thumb.jpg", reader, int64(reader.Len()))
}

func resizeImage(file io.Reader, width uint) *bytes.Buffer {
	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	m := resize.Thumbnail(width, width, img, resize.Lanczos3)

	var buf = new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {
		panic(err)
	}
	return buf
}
