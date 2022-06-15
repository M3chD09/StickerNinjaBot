# StickerNinjaBot

[![Go](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/go.yml/badge.svg)](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/go.yml)
[![CodeQL](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/codeql-analysis.yml)
[![Heroku](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/heroku.yml/badge.svg)](https://github.com/M3chD09/StickerNinjaBot/actions/workflows/heroku.yml)

[中文文档](README_zh.md)

## Description
[StickerNinjaBot](https://t.me/StickerNinjaBot) can help you convert telegram stickers to images in png, jpg and gif format.  
Just send the sticker or sticker set link to the bot and it will convert the image to your preferred format.  
Single image or zip file with multiple images will be sent to you, depending on your needs.  

## Features
### Convertion
* Support `.webp`, `.webm` and `.tgs` format stickers input.
* Support `.png`, `.jpg` and `.gif` format output.
### Packing
* Send bulk stickers and receive converted stickers in a zip archive.

## Usage
* Send command `/formats` to set your preferred format.
### Single sticker
* Send any stickers to the bot and receive converted stickers in your preferred format.
### Multiple stickers
* Send any sticker set link with the prefix `https://t.me/addstickers/` or `https://telegram.me/addstickers/` and receive converted stickers in a zip archive.
* Send command `/newpack` to start packing, then send any stickers or sticker set link to the bot, and receive converted stickers in a zip archive until the command `/finish` is sent, or cancel the packing by sending the command `/cancel`.

## Deploying
### Deploy on [Heroku](https://heroku.com)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)
### Installation
```bash
apt install -y ffmpeg
git clone https://github.com/M3chD09/StickerNinjaBot
cd StickerNinjaBot
go build
cp .env.example .env
```
### Configuration
Edit the `.env` file to configure the bot:  
* `BOT_TOKEN`: Your telegram bot token. Get it from [Telegram Botfather](https://telegram.me/botfather).
* `BOT_WEBHOOK`: Specify a url and receive incoming updates via an outgoing webhook, or leave it blank to receive updates via polling. For more information, see [Telegram Bot API](https://core.telegram.org/bots/api#getting-updates).
* `BOT_DEBUG`: Set to `true` to enable debug mode.
* `PORT`: The port to listen on. Leave it blank if you don't want to use a webhook.
* `FILESTORAGE_PATH`: The path to the directory where the downloaded stickers will be stored. Default is `./storage`.
* `STICKER_COUNT_LIMIT`: The maximum number of stickers that can be sent in a single pack. Default is `100`.
* `CACHE_TICK`: The interval at which the user cache is refreshed. Default is `10s`.
* `CACHE_EXPIRE`: The time after which the user cache is expired. Default is `15m`.
* `DATABASE_TYPE`: The type of database to use. Currently `mysql`, `pgsql`, and `sqlite` are supported.
* `DATABASE_URL`: The url of the database to connect to.
### Running
```bash
./StickerNinjaBot
```

## Credits
* [phoenixlzx/telegram-stickerimage-bot](https://github.com/phoenixlzx/telegram-stickerimage-bot)
* [Benau/tgsconverter](https://github.com/Benau/tgsconverter)
