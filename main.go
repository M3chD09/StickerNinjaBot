package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/M3chD09/StickerNinjaBot/filestorage"
	"github.com/M3chD09/StickerNinjaBot/tgbot"
	"github.com/M3chD09/StickerNinjaBot/userdb"
)

func main() {
	godotenv.Load()

	tgbot.Config(os.Getenv("STICKER_COUNT_LIMIT"))
	filestorage.Config(os.Getenv("FILESTORAGE_PATH"))
	userdb.DBConfig(os.Getenv("DATABASE_TYPE"), os.Getenv("DATABASE_URL"))

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = os.Getenv("BOT_DEBUG") == "true"

	log.Printf("Authorized on account %s", bot.Self.UserName)

	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Show help message",
		},
		tgbotapi.BotCommand{
			Command:     "lang",
			Description: "Change language",
		},
		tgbotapi.BotCommand{
			Command:     "formats",
			Description: "Set preferred formats",
		},
		tgbotapi.BotCommand{
			Command:     "newpack",
			Description: "Create a new pack",
		},
		tgbotapi.BotCommand{
			Command:     "finish",
			Description: "Finish the pack",
		},
		tgbotapi.BotCommand{
			Command:     "cancel",
			Description: "Cancel the pack",
		},
	)
	if _, err := bot.Request(setCommands); err != nil {
		log.Fatal("Unable to set commands", err)
	}

	var updates tgbotapi.UpdatesChannel
	if os.Getenv("BOT_WEBHOOK") != "" {
		secretPath := make([]byte, 16)
		rand.Read(secretPath)
		secretPath = []byte(base64.StdEncoding.EncodeToString(secretPath))

		wh, _ := tgbotapi.NewWebhook(os.Getenv("BOT_WEBHOOK") + string(secretPath))
		_, err = bot.Request(wh)
		if err != nil {
			log.Fatal(err)
		}

		info, err := bot.GetWebhookInfo()
		if err != nil {
			log.Fatal(err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}
		updates = bot.ListenForWebhook("/" + string(secretPath))
		go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	} else {
		wh := tgbotapi.DeleteWebhookConfig{DropPendingUpdates: false}
		_, err = bot.Request(wh)
		if err != nil {
			log.Fatal(err)
		}
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates = bot.GetUpdatesChan(u)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Fatal(err)
			}

			instance := tgbot.GetInstance(update.CallbackQuery.Message.Chat.ID, bot)
			if f, err := strconv.ParseUint(update.CallbackQuery.Data, 10, 64); err == nil {
				instance.FormatsApply(uint8(f))
			} else {
				instance.LangApply(update.CallbackQuery.Data)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		instance := tgbot.GetInstance(update.Message.Chat.ID, bot)

		if update.Message.Sticker != nil {
			go instance.AddSticker(update.Message.Sticker.FileID)
			continue
		}

		if !update.Message.IsCommand() {
			if update.Message.Text == "" {
				continue
			}
			if strings.HasPrefix(update.Message.Text, "https://t.me/addstickers/") ||
				strings.HasPrefix(update.Message.Text, "https://telegram.me/addstickers/") {
				stickerSetName := filepath.Base(update.Message.Text)
				go instance.AddStickerSet(stickerSetName)
			} else {
				go instance.Help()
			}
			continue
		}

		switch update.Message.Command() {
		case "help":
			go instance.Help()
		case "lang":
			go instance.Lang()
		case "formats":
			go instance.Formats()
		case "newpack":
			go instance.NewPack()
		case "finish":
			go instance.Finish()
		case "cancel":
			go instance.Cancel()
		default:
			go instance.Help()
		}
	}
}
