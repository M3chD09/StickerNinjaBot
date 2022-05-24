package filestorage

import (
	"os"
	"testing"
)

func TestUserStorage(t *testing.T) {
	storageRootPath := Config("")
	userStorage := NewUserStorage(1, []string{"png"})
	if userStorage.RootPath() != "storage/1" {
		t.Errorf("UserStorage.RootPath() error %v", userStorage.RootPath())
	}
	if _, err := os.Stat(userStorage.RootPath()); os.IsNotExist(err) {
		t.Error("UserStorage is not exist error")
	}
	userStorage.MakeDir("src")
	if _, err := os.Stat(userStorage.SubPath("src")); os.IsNotExist(err) {
		t.Error("UserStorage.SubPath is not exist error")
	}
	userStorage.Remove("src")
	if _, err := os.Stat(userStorage.SubPath("src")); !os.IsNotExist(err) {
		t.Error("UserStorage.SubPath is exist error")
	}
	userStorage.Remove("")
	if _, err := os.Stat(userStorage.RootPath()); !os.IsNotExist(err) {
		t.Error("UserStorage is exist error")
	}
	os.RemoveAll(storageRootPath)
}
