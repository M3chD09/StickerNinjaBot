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
		panic(err)
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

func (u *UserStorage) Zip() ([]string, error) {
	var dirs []string
	var err error
	var once sync.Once

	for _, d := range append(u.formats, "src") {
		zipDirPath := u.SubPath(d)
		if dir, err := os.ReadDir(zipDirPath); err != nil || len(dir) == 0 {
			continue
		}
		dirs = append(dirs, d)
	}

	zipFilePathList := make([]string, len(dirs))

	var wg sync.WaitGroup
	for x := range dirs {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			zipFileName := strconv.FormatInt(u.userID, 10) + "_" + dirs[a] + ".zip"
			zipFilePath := u.SubPath(zipFileName)
			zipDirPath := u.SubPath(dirs[a])
			if e := zipDir(zipDirPath, zipFilePath); e != nil {
				log.Println("Error in UserStorage zip: ", e)
				once.Do(func() { err = e })
				return
			}
			zipFilePathList[a] = zipFilePath
		}(x)
	}
	wg.Wait()
	return zipFilePathList, err
}

func (u *UserStorage) SaveSingleSticker(url string) ([]string, error) {
	var err error
	var once sync.Once

	sticker := NewStickerFromURL(url)
	filePath := filepath.Join(u.RootPath(), sticker.FileName())
	if e := sticker.Save(filePath); e != nil {
		log.Println("Error in UserStorage SaveSingleSticker: ", e)
		return nil, e
	}

	var dstList []string
	for _, f := range u.formats {
		dst := u.SubPath(sticker.ReplaceExt(f))
		if e := sticker.Convert(dst); e != nil {
			log.Println("Error in UserStorage SaveSingleSticker: ", e)
			once.Do(func() { err = e })
			continue
		}
		dstList = append(dstList, dst)
	}
	return dstList, err
}

func (u *UserStorage) SaveSticker(url string) error {
	sticker := NewStickerFromURL(url)
	u.MakeDir("src")
	filePath := filepath.Join(u.SubPath("src"), sticker.FileName())
	if err := sticker.Save(filePath); err != nil {
		log.Println("Error in UserStorage SaveSticker: ", err)
		return err
	}
	return nil
}

func (u *UserStorage) SaveStickers(urlList []string) error {
	u.MakeDir("src")
	var err error
	var once sync.Once

	var wg sync.WaitGroup
	for x := range urlList {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			sticker := NewStickerFromURL(urlList[a])
			filePath := filepath.Join(u.SubPath("src"), sticker.FileName())
			if e := sticker.Save(filePath); e != nil {
				log.Println("Error in UserStorage SaveStickers: ", e)
				once.Do(func() { err = e })
			}
		}(x)
	}
	wg.Wait()
	return err
}

func (u *UserStorage) ConvertStickers() error {
	var err error
	var once sync.Once

	for _, f := range u.formats {
		u.MakeDir(f)
	}

	var wg sync.WaitGroup
	filepath.WalkDir(u.SubPath("src"), func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sticker := NewStickerFromFilePath(p)
			for _, f := range u.formats {
				e := sticker.Convert(filepath.Join(u.SubPath(f), sticker.ReplaceExt(f)))
				if e != nil {
					log.Println("Error in UserStorage ConvertStickers: ", e)
					once.Do(func() { err = e })
				}
			}
		}(path)
		return nil
	})
	wg.Wait()

	return err
}
