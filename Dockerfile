FROM golang:latest AS builder
WORKDIR /app/
COPY . /app/
RUN go build -ldflags="-w -s" -v -o StickerNinjaBot

FROM linuxserver/ffmpeg:latest
COPY --from=builder /app/StickerNinjaBot /app/StickerNinjaBot
WORKDIR /app/
COPY ./locales /app/locales
ENTRYPOINT ["/app/StickerNinjaBot"]
