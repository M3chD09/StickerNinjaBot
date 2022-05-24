package tgbot

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/M3chD09/StickerNinjaBot/filestorage"
	"github.com/M3chD09/StickerNinjaBot/userdb"
)

var stickerCountLimit = 100
var userInstanceCache = userdb.NewCache[int64](time.Second*10, true)
var bundle *i18n.Bundle
var langKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("English", "en"),
		tgbotapi.NewInlineKeyboardButtonData("简体中文", "zh-hans"),
	),
)
var formatsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("PNG", strconv.FormatUint(uint64(userdb.FormatsPNG), 10)),
		tgbotapi.NewInlineKeyboardButtonData("JPG", strconv.FormatUint(uint64(userdb.FormatsJPG), 10)),
		tgbotapi.NewInlineKeyboardButtonData("GIF", strconv.FormatUint(uint64(userdb.FormatsGIF), 10)),
	),
)

type instance struct {
	idle    state
	packing state
	busy    state

	currentState state

	userID         int64
	formats        []string
	bot            *tgbotapi.BotAPI
	stickerFileIDs []string
	localizer      *i18n.Localizer
}

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("locales/en.toml")
	bundle.MustLoadMessageFile("locales/zh-hans.toml")
}

func Config(count string) {
	if count != "" {
		if c, err := strconv.Atoi(count); err == nil && c > 0 {
			stickerCountLimit = c
		}
	}
}

func GetInstance(userID int64, bot *tgbotapi.BotAPI) *instance {
	if i, ok := userInstanceCache.Get(userID); ok {
		return i.(*instance)
	}

	i := &instance{
		userID: userID,
		bot:    bot,
	}
	if lang := userdb.DBGetLanguage(userID); lang != "" {
		i.localizer = i18n.NewLocalizer(bundle, lang)
	} else {
		userdb.DBSaveLanguage(userID, "en")
		i.localizer = i18n.NewLocalizer(bundle, "en")
	}
	i.formats = userdb.DBGetFormats(userID)

	i.idle = &idleState{instance: i}
	i.packing = &packingState{instance: i}
	i.busy = &busyState{instance: i}
	i.currentState = i.idle
	userInstanceCache.Set(userID, i, time.Hour)

	return i
}

func (i *instance) Help() {
	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Help",
		},
	}))
}

func (i *instance) Lang() {
	msg := tgbotapi.NewMessage(i.userID, "Please select language")
	msg.ReplyMarkup = langKeyboard
	if _, err := i.bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func (i *instance) LangApply(lang string) {
	go userdb.DBSaveLanguage(i.userID, lang)
	i.localizer = i18n.NewLocalizer(bundle, lang)

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "LangApply",
		},
	}))
}

func (i *instance) Formats() {
	text := ""
	if len(i.formats) == 0 {
		text = i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatNone",
			},
		})
	} else {
		text = i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "Formats",
			},
			TemplateData: map[string]interface{}{
				"Formats": strings.Join(i.formats, ", "),
			},
		})
	}
	msg := tgbotapi.NewMessage(i.userID, text)
	msg.ReplyMarkup = formatsKeyboard
	if _, err := i.bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func (i *instance) FormatsApply(formatCode uint8) {
	userdb.DBSaveFormatCode(i.userID, formatCode)
	i.formats = userdb.DBGetFormats(i.userID)

	text := ""
	if len(i.formats) == 0 {
		text = i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatNone",
			},
		})
	} else {
		text = i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatsApply",
			},
			TemplateData: map[string]interface{}{
				"Formats": strings.Join(i.formats, ", "),
			},
		})
	}
	i.sendTextMessage(text)
}

func (i *instance) NewPack() {
	i.currentState.newPack()
}

func (i *instance) AddSticker(stickerFileID string) {
	i.currentState.addSticker(stickerFileID)
}

func (i *instance) AddStickerSet(stickerSetName string) {
	i.currentState.addStickerSet(stickerSetName)
}

func (i *instance) Finish() {
	i.currentState.finish()
}

func (i *instance) Cancel() {
	i.currentState.cancel()
}

func (i *instance) setState(s state) {
	i.currentState = s
}

func (i *instance) extractStickerSet(stickerSetName string) []string {
	stickerSet, err := i.bot.GetStickerSet(tgbotapi.GetStickerSetConfig{Name: stickerSetName})
	if err != nil {
		return []string{}
	}

	stickerFileIDs := make([]string, len(stickerSet.Stickers))
	var wg sync.WaitGroup
	for x := range stickerSet.Stickers {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			stickerFileIDs[a] = stickerSet.Stickers[a].FileID
		}(x)
	}
	wg.Wait()

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "StickerSetExtract",
		},
		PluralCount: len(stickerFileIDs),
		TemplateData: map[string]interface{}{
			"Count": len(stickerFileIDs),
		},
	}))
	return stickerFileIDs
}

func (i *instance) fetchStickers(stickerFileIDs []string) []string {
	urlList := make([]string, len(stickerFileIDs))
	var wg sync.WaitGroup
	for x := range stickerFileIDs {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			url, err := i.bot.GetFileDirectURL(stickerFileIDs[a])
			if err != nil {
				log.Fatal(err)
			}
			urlList[a] = url
		}(x)
	}
	wg.Wait()
	return urlList
}

func (i *instance) isStickerCountTooMany(count int) bool {
	if count <= stickerCountLimit {
		return false
	}

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "StickerCountTooMany",
		},
		TemplateData: map[string]interface{}{
			"Count": stickerCountLimit,
		},
	}))
	return true
}

func (i *instance) sendStickers(stickerFileIDs []string) {
	us := filestorage.NewUserStorage(i.userID, i.formats)
	defer us.Remove("")

	urlList := i.fetchStickers(stickerFileIDs)
	i.stickerFileIDs = []string{}
	us.SaveStickers(urlList)

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "StickersDownload",
		},
	}))

	us.ConvertStickers()

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "StickersConvert",
		},
	}))

	filePathList := us.Zip()

	i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "StickersZip",
		},
	}))

	if len(filePathList) == 1 {
		i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "StickerConvertNotSupport",
			},
		}))
	}

	i.sendMultiFileMessage(filePathList)
}

func (i *instance) sendTextMessage(text string) {
	msg := tgbotapi.NewMessage(i.userID, text)
	_, err := i.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *instance) sendFileMessage(filePath string) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileName := filepath.Base(filePath)

	if len(b) > (1<<20)*50 {
		i.sendTextMessage(i.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "StickersSizeTooLarge",
			},
			TemplateData: map[string]interface{}{
				"FileName": fileName,
			},
		}))
		return
	}

	msg := tgbotapi.NewDocument(i.userID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: b,
	})
	_, err = i.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *instance) sendMultiFileMessage(filePathList []string) {
	var wg sync.WaitGroup
	for x := range filePathList {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			i.sendFileMessage(filePathList[a])
		}(x)
	}
	wg.Wait()
}
