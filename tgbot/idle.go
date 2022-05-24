package tgbot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/M3chD09/StickerNinjaBot/filestorage"
)

type idleState struct {
	instance *instance
}

func (s *idleState) newPack() {
	if len(s.instance.formats) == 0 {
		s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatNoneErr",
			},
		}))
		return
	}

	s.instance.setState(s.instance.packing)
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "IdleNewPack",
		},
	}))
}

func (s *idleState) addSticker(stickerFileID string) {
	s.instance.setState(s.instance.busy)
	defer s.instance.setState(s.instance.idle)

	if len(s.instance.formats) == 0 {
		s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatNoneErr",
			},
		}))
		return
	}

	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "IdleAddSticker",
		},
	}))

	us := filestorage.NewUserStorage(s.instance.userID, s.instance.formats)
	defer us.Remove("")

	urlList := s.instance.fetchStickers([]string{stickerFileID})
	filePathList := us.SaveSingleSticker(urlList[0])

	if len(filePathList) == 0 {
		s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "StickerConvertNotSupport",
			},
		}))
		return
	}

	for _, filePath := range filePathList {
		s.instance.sendFileMessage(filePath)
	}
}

func (s *idleState) addStickerSet(stickerSetName string) {
	s.instance.setState(s.instance.busy)
	defer s.instance.setState(s.instance.idle)

	if len(s.instance.formats) == 0 {
		s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "FormatNoneErr",
			},
		}))
		return
	}

	if len(stickerSetName) == 0 {
		return
	}

	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "IdleAddStickerSet",
		},
	}))

	stickerFileIDs := s.instance.extractStickerSet(stickerSetName)
	if s.instance.isStickerCountTooMany(len(stickerFileIDs)) {
		return
	}

	s.instance.sendStickers(stickerFileIDs)
}

func (s *idleState) finish() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "IdleFinish",
		},
	}))
}

func (s *idleState) cancel() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "IdleCancel",
		},
	}))
}
