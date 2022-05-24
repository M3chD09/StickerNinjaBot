package filestorage

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/Benau/tgsconverter/libtgsconverter"
	"golang.org/x/image/webp"
)

func webp2other(webpPath, otherPath string) string {
	if filepath.Ext(webpPath) != ".webp" {
		return ""
	}

	reader, err := os.Open(webpPath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	img, err := webp.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	writer, err := os.Create(otherPath)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	ext := filepath.Ext(otherPath)
	switch ext {
	case ".png":
		err = png.Encode(writer, img)
	case ".jpg":
		err = jpeg.Encode(writer, img, nil)
	case ".gif":
		err = gif.Encode(writer, img, nil)
	default:
		log.Println("unsupported file extension: " + ext)
		return ""
	}
	if err != nil {
		log.Fatal(err)
	}
	return otherPath
}

func webm2other(webmPath, otherPath string) string {
	// TODO
	return ""
}

func tgs2other(tgsPath, otherPath string) string {
	if filepath.Ext(tgsPath) != ".tgs" {
		return ""
	}

	ext := filepath.Ext(otherPath)
	ext = ext[1:]
	if !libtgsconverter.SupportsExtension(ext) {
		log.Println("unsupported file extension: " + ext)
		return ""
	}

	opt := libtgsconverter.NewConverterOptions()
	opt.SetExtension(ext)
	ret, err := libtgsconverter.ImportFromFile(tgsPath, opt)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(otherPath, ret, 0666)

	return otherPath
}
