package filestorage

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	dirs := append(u.formats, "src")
	zipFilePathList := make([]string, len(dirs))

	var wg sync.WaitGroup
	for x := range dirs {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			zipFileName := strconv.FormatInt(u.userID, 10) + "_" + dirs[a] + ".zip"
			zipFilePath := u.SubPath(zipFileName)
			zipDirPath := u.SubPath(dirs[a])
			if err := zipDir(zipDirPath, zipFilePath); err != nil {
				log.Fatal("UserStorage zip error: ", err)
			}
			zipFilePathList[a] = zipFilePath
		}(x)
	}
	wg.Wait()
	return zipFilePathList
}

func (u *UserStorage) SaveSingleSticker(url string) []string {
	sticker := NewStickerFromURL(url)
	filePath := filepath.Join(u.RootPath(), sticker.FileName())
	if err := sticker.Save(filePath); err != nil {
		log.Fatal("UserStorage SaveSingleSticker error: ", err)
	}

	var dstList []string
	for _, f := range u.formats {
		dst := u.SubPath(sticker.ReplaceExt(f))
		err := sticker.Convert(dst)
		if err != nil {
			log.Println("UserStorage SaveSingleSticker error: ", err)
			continue
		}
		dstList = append(dstList, dst)
	}
	return dstList
}

func (u *UserStorage) SaveSticker(url string) {
	sticker := NewStickerFromURL(url)
	u.MakeDir("src")
	filePath := filepath.Join(u.SubPath("src"), sticker.FileName())
	if err := sticker.Save(filePath); err != nil {
		log.Fatal("UserStorage SaveSticker error: ", err)
	}
}

func (u *UserStorage) SaveStickers(urlList []string) {
	u.MakeDir("src")

	var wg sync.WaitGroup
	for x := range urlList {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			sticker := NewStickerFromURL(urlList[a])
			filePath := filepath.Join(u.SubPath("src"), sticker.FileName())
			if err := sticker.Save(filePath); err != nil {
				log.Fatal("UserStorage SaveStickers error: ", err)
			}
		}(x)
	}
	wg.Wait()
}

func (u *UserStorage) ConvertStickers() {
	for _, f := range u.formats {
		u.MakeDir(f)
	}

	var wg sync.WaitGroup
	filepath.WalkDir(u.SubPath("src"), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sticker := NewStickerFromFilePath(p)
			for _, f := range u.formats {
				err := sticker.Convert(filepath.Join(u.SubPath(f), sticker.ReplaceExt(f)))
				if err != nil {
					log.Println("UserStorage ConvertStickers error: ", err)
				}
			}
		}(path)
		return nil
	})
	wg.Wait()
}
