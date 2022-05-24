package filestorage

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var storageRootPath string

type UserStorage struct {
	userID  int64
	formats []string
}

func Config(path string) string {
	if path == "" {
		storageRootPath = "./storage"
	} else {
		storageRootPath = path
	}

	err := os.MkdirAll(storageRootPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return storageRootPath
}

func NewUserStorage(userID int64, formats []string) *UserStorage {
	userStorage := &UserStorage{
		userID:  userID,
		formats: formats,
	}
	userStorage.MakeDir("")
	return userStorage
}

func (u *UserStorage) RootPath() string {
	return filepath.Join(storageRootPath, strconv.FormatInt(u.userID, 10))
}

func (u *UserStorage) SubPath(sub string) string {
	return filepath.Join(u.RootPath(), sub)
}

func (u *UserStorage) MakeDir(sub string) {
	err := os.MkdirAll(u.SubPath(sub), 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func (u *UserStorage) Remove(sub string) {
	err := os.RemoveAll(u.SubPath(sub))
	if err != nil {
		log.Fatal(err)
	}
}

func (u *UserStorage) Zip() []string {
	var zipFilePathList []string
	for _, f := range append(u.formats, "src") {
		zipFileName := strconv.FormatInt(u.userID, 10) + "_" + f + ".zip"
		zipFilePath := u.SubPath(zipFileName)
		zipDirPath := u.SubPath(f)
		zipDir(zipDirPath, zipFilePath)
		zipFilePathList = append(zipFilePathList, zipFilePath)
	}
	return zipFilePathList
}

func (u *UserStorage) SaveSingleSticker(url string) []string {
	sticker := NewStickerFromURL(url)
	filePath := filepath.Join(u.RootPath(), sticker.FileName())
	sticker.Save(filePath)

	var dstList []string
	for _, f := range u.formats {
		dst := u.SubPath(sticker.ReplaceExt(f))
		if sticker.Convert(dst) != "" {
			dstList = append(dstList, dst)
		}
	}
	return dstList
}

func (u *UserStorage) SaveSticker(url string) {
	sticker := NewStickerFromURL(url)
	u.MakeDir("src")
	filePath := filepath.Join(u.SubPath("src"), sticker.FileName())
	sticker.Save(filePath)
}

func (u *UserStorage) SaveStickers(urlList []string) {
	for _, url := range urlList {
		u.SaveSticker(url)
	}
}

func (u *UserStorage) ConvertStickers() {
	for _, f := range u.formats {
		u.MakeDir(f)
	}

	filepath.Walk(u.SubPath("src"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, f := range u.formats {
			sticker := NewStickerFromFilePath(path)
			sticker.Convert(filepath.Join(u.SubPath(f), sticker.ReplaceExt(f)))
		}
		return nil
	})
}
