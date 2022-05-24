package filestorage

import (
	"io"
	"log"
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

func (s *Sticker) Save(filePath string) {
	resp, err := http.Get(s.url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	s.filePath = filePath
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

func (s *Sticker) Convert(dst string) string {
	switch s.Ext() {
	case ".webp":
		return webp2other(s.filePath, dst)
	case ".webm":
		return webm2other(s.filePath, dst)
	case ".tgs":
		return tgs2other(s.filePath, dst)
	default:
		return ""
	}
}
