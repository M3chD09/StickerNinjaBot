package filestorage

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Sticker struct {
	url      string
	filePath string
}

func NewStickerFromURL(url string) *Sticker {
	s := &Sticker{url: url}
	return s
}

func NewStickerFromFilePath(filePath string) *Sticker {
	s := &Sticker{filePath: filePath}
	return s
}

func (s *Sticker) Save(filePath string) error {
	resp, err := http.Get(s.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("http status code is not 200")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	s.filePath = filePath
	return nil
}

func (s *Sticker) FileName() string {
	if s.url != "" {
		return filepath.Base(s.url)
	} else if s.filePath != "" {
		return filepath.Base(s.filePath)
	}
	return ""
}

func (s *Sticker) Ext() string {
	return filepath.Ext(s.FileName())
}

func (s *Sticker) ReplaceExt(ext string) string {
	fileName := s.FileName()
	return fileName[:len(fileName)-len(filepath.Ext(fileName))] + "." + ext
}

func (s *Sticker) Convert(dst string) error {
	switch s.Ext() {
	case ".webp":
		return webp2other(s.filePath, dst)
	case ".webm":
		return webm2other(s.filePath, dst)
	case ".tgs":
		return tgs2other(s.filePath, dst)
	default:
		return errors.New("unsupported output file extension: " + s.Ext())
	}
}
