package userdb

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	ID         uint
	TelegramId int64 `gorm:"uniqueIndex"`
	Language   string
	FormatCode uint8
}

const (
	FormatsNone uint8 = 0
	FormatsPNG  uint8 = 1
	FormatsJPG  uint8 = 2
	FormatsGIF  uint8 = 4
)

func DBConfig(dbType string, dbURL string) {
	var err error
	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	case "pgsql":
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	default:
		panic("Unknown database type. Supported types: mysql, pgsql, sqlite")
	}
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
}

func DBSaveLanguage(id int64, language string) {
	user := User{}
	db.Where(User{TelegramId: id}).First(&user)
	if user == (User{}) {
		user = User{
			TelegramId: id,
			Language:   language,
			FormatCode: FormatsPNG,
		}
		db.Create(&user)
	} else {
		user.Language = language
		db.Save(&user)
	}
}

func DBGetLanguage(id int64) string {
	user := User{}
	db.Where(User{TelegramId: id}).First(&user)
	if user == (User{}) {
		return ""
	}

	return user.Language
}

func DBSaveFormatCode(id int64, formatCode uint8) {
	user := User{}
	db.Where(User{TelegramId: id}).First(&user)
	if user == (User{}) {
		user = User{
			TelegramId: id,
			Language:   "en",
			FormatCode: formatCode,
		}
		db.Create(&user)
	} else {
		user.FormatCode = user.FormatCode ^ formatCode
		db.Save(&user)
	}
}

func DBGetFormats(id int64) []string {
	user := User{}
	db.Where(User{TelegramId: id}).First(&user)
	if user == (User{}) {
		return []string{}
	}

	return DBParseFormatCode(user.FormatCode)
}

func DBParseFormatCode(formatCode uint8) []string {
	var formats []string
	if formatCode&FormatsPNG != 0 {
		formats = append(formats, "png")
	}
	if formatCode&FormatsJPG != 0 {
		formats = append(formats, "jpg")
	}
	if formatCode&FormatsGIF != 0 {
		formats = append(formats, "gif")
	}
	return formats
}
