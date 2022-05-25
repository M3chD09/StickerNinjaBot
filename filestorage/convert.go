package filestorage

import (
	"errors"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/M3chD09/tgsconverter/libtgsconverter"
	"golang.org/x/image/webp"
)

func webp2other(webpPath, otherPath string) error {
	if filepath.Ext(webpPath) != ".webp" {
		panic("input file must be .webp")
	}

	reader, err := os.Open(webpPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	img, err := webp.Decode(reader)
	if err != nil {
		return err
	}
	writer, err := os.Create(otherPath)
	if err != nil {
		return err
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
		return errors.New("unsupported output file extension: " + ext)
	}
	if err != nil {
		return err
	}
	return nil
}

func webm2other(webmPath, otherPath string) error {
	// TODO
	return errors.New("not implemented")
}

func tgs2other(tgsPath, otherPath string) error {
	if filepath.Ext(tgsPath) != ".tgs" {
		panic("input file must be .tgs")
	}

	ext := filepath.Ext(otherPath)
	ext = ext[1:]
	if !libtgsconverter.SupportsExtension(ext) {
		return errors.New("unsupported output file extension: " + ext)
	}

	opt := libtgsconverter.NewConverterOptions()
	opt.SetExtension(ext)
	ret, err := libtgsconverter.ImportFromFile(tgsPath, opt)
	if err != nil {
		return err
	}

	if err := os.WriteFile(otherPath, ret, 0666); err != nil {
		return err
	}

	return nil
}
