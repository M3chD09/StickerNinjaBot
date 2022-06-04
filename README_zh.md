# StickerNinjaBot

## 描述
[StickerNinjaBot](https://t.me/StickerNinjaBot)可以帮助你将Telegram表情包转换为图片，支持png、jpg、gif格式。  
只需发送表情包或表情包集链接给机器人，它将转换为你所需求格式的图片。  
根据你的需求，单个图片或带有多个图片的 zip 文件将会被发送给你。  

## 特色
### 转换
* 支持`.webp`、`.webm`和`.tgs`格式的表情包输入。
* 支持`.png`、`.jpg`和`.gif`格式输出。
### 打包
* 发送批量表情包并在 zip 存档中接收转换后的表情包。

## 用法
* 发送命令 `/formats` 来设置你喜欢的格式。
### 单个表情包
* 将任何表情包发送到机器人并以您喜欢的格式接收转换后的表情包。
### 批量表情包
* 发送任何带有前缀`https://t.me/addstickers/`或`https://telegram.me/addstickers/`的表情包集链接，并在压缩包中接收转换后的表情包。
* 发送命令 `/newpack` 开始打包，然后将任何表情包或表情包集链接发送给机器人，发送命令 `/finish` 便可在 zip 存档中接收转换后的表情包，或者通过发送命令 `/cancel` 取消打包。

## 部署
### 安装
```bash
apt install -y ffmpeg
git clone https://github.com/M3chD09/StickerNinjaBot
cd StickerNinjaBot
go build
cp .env.example .env
```
### 配置
编辑 `.env` 文件来配置机器人：
* `BOT_TOKEN`: 机器人的 API 令牌。可以在 [Telegram BotFather](https://telegram.me/botfather) 中获取。
* `BOT_WEBHOOK`: 指定 url 并通过传出 webhook 接收传入更新，或将其留空以通过轮询接收更新。有关详细信息，请参阅 [Telegram Bot API](https://core.telegram.org/bots/api#getting-updates)。
* `BOT_DEBUG`: 设置为`true`以启用调试模式。
* `PORT`: 端口号。如果不使用 webhook，则可以留空。
* `FILESTORAGE_PATH`: 下载的表情包存放的路径。默认为 `./storage`。
* `STICKER_COUNT_LIMIT`: 打包的表情包数量限制。默认为 `100`。
* `DATABASE_TYPE`: 数据库类型。可以是 `mysql`、`pgsql` 或 `sqlite`。
* `DATABASE_URL`: 数据库连接字符串。
### 运行
```bash
./StickerNinjaBot
```

## 致谢
* [phoenixlzx/telegram-stickerimage-bot](https://github.com/phoenixlzx/telegram-stickerimage-bot)
* [Benau/tgsconverter](https://github.com/Benau/tgsconverter)
