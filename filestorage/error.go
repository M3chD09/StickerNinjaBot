package filestorage

import (
	"errors"
	"path/filepath"
	"strconv"
)

var (
	ErrConvertInputExtensionNotSupported  = errors.New("input file extension not supported")
	ErrConvertOutputExtensionNotSupported = errors.New("output file extension not supported")

	ErrDownloadStatusCode = errors.New("download status code not 200")
)

type CustomError interface {
	Error() string
	MessageID() string
	TemplateDataText() string
}

type ConvertError struct {
	Op  string
	Src string
	Dst string
	Err error
}

type DownloadError struct {
	URL        string
	StatusCode int
	Err        error
}

func NewConvertError(op, src, dst string, err error) *ConvertError {
	return &ConvertError{
		Op:  op,
		Src: src,
		Dst: dst,
		Err: err,
	}
}

func (e *ConvertError) Error() string {
	return "convert error: " + e.Op + " " + e.Src + " -> " + e.Dst + ": " + e.Err.Error()
}

func (e *ConvertError) MessageID() string {
	switch e.Err {
	case ErrConvertInputExtensionNotSupported:
		return "StickerConvertInputExtensionNotSupportedError"
	case ErrConvertOutputExtensionNotSupported:
		return "StickerConvertOutputExtensionNotSupportedError"
	default:
		return "StickerConvertOtherError"
	}
}

func (e *ConvertError) TemplateDataText() string {
	switch e.Err {
	case ErrConvertInputExtensionNotSupported:
		fallthrough
	case ErrConvertOutputExtensionNotSupported:
		return filepath.Ext(e.Src) + " -> " + filepath.Ext(e.Dst)
	default:
		return ""
	}
}

func NewDownloadError(url string, StatusCode int, err error) *DownloadError {
	return &DownloadError{
		URL:        url,
		StatusCode: StatusCode,
		Err:        err,
	}
}

func (e *DownloadError) Error() string {
	return "download error: " + e.URL + ": " + "status code: " + strconv.Itoa(e.StatusCode) + ", " + e.Err.Error()
}

func (e *DownloadError) MessageID() string {
	switch e.Err {
	case ErrDownloadStatusCode:
		return "StickerDownloadStatusCodeError"
	default:
		return "StickerDownloadOtherError"
	}
}

func (e *DownloadError) TemplateDataText() string {
	switch e.Err {
	case ErrDownloadStatusCode:
		return strconv.Itoa(e.StatusCode)
	default:
		return ""
	}
}
