FROM golang:1.27 AS builder
WORKDIR /app/
COPY . /app/
RUN go build -ldflags="-w -s" -v -o StickerNinjaBot

FROM linuxserver/ffmpeg:9.0-cli-ls79
COPY --from=builder /app/StickerNinjaBot /app/StickerNinjaBot
WORKDIR /app/
COPY ./locales /app/locales
ENTRYPOINT ["/app/StickerNinjaBot"]
