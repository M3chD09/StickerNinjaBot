package filestorage

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/M3chD09/tgsconverter/libtgsconverter"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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
		return NewConvertError("webp2other", webpPath, otherPath, err)
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
		return NewConvertError("webp2other", webpPath, otherPath, ErrConvertOutputExtensionNotSupported)
	}
	if err != nil {
		return NewConvertError("webp2other", webpPath, otherPath, err)
	}
	return nil
}

func webm2other(webmPath, otherPath string) error {
	ext := filepath.Ext(otherPath)
	if ext == ".gif" {
		return ffmpeg.Input(webmPath, ffmpeg.KwArgs{}).Output(otherPath, ffmpeg.KwArgs{}).OverWriteOutput().ErrorToStdOut().Run()
	}

	return NewConvertError("webm2other", webmPath, otherPath, ErrConvertOutputExtensionNotSupported)
}

func tgs2other(tgsPath, otherPath string) error {
	if filepath.Ext(tgsPath) != ".tgs" {
		panic("input file must be .tgs")
	}

	ext := filepath.Ext(otherPath)
	ext = ext[1:]
	if !libtgsconverter.SupportsExtension(ext) {
		return NewConvertError("tgs2other", tgsPath, otherPath, ErrConvertOutputExtensionNotSupported)
	}

	opt := libtgsconverter.NewConverterOptions()
	opt.SetExtension(ext)
	ret, err := libtgsconverter.ImportFromFile(tgsPath, opt)
	if err != nil {
		return NewConvertError("tgs2other", tgsPath, otherPath, err)
	}

	if err := os.WriteFile(otherPath, ret, 0666); err != nil {
		return err
	}

	return nil
}
