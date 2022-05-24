package tgbot

type state interface {
	newPack()
	addSticker(stickerFileID string)
	addStickerSet(stickerSetName string)
	finish()
	cancel()
}
