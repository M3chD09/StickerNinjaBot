package tgbot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type packingState struct {
	instance *instance
}

func (s *packingState) newPack() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PackingNewPack",
		},
	}))
}

func (s *packingState) addSticker(stickerFileID string) {
	s.instance.setState(s.instance.busy)
	defer s.instance.setState(s.instance.packing)

	s.instance.stickerFileIDs = append(s.instance.stickerFileIDs, stickerFileID)
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PackingAddSticker",
		},
		TemplateData: map[string]interface{}{
			"Count": len(s.instance.stickerFileIDs),
		},
	}))
}

func (s *packingState) addStickerSet(stickerSetName string) {
	s.instance.setState(s.instance.busy)
	defer s.instance.setState(s.instance.packing)

	stickerFileIDs := s.instance.extractStickerSet(stickerSetName)
	s.instance.stickerFileIDs = append(s.instance.stickerFileIDs, stickerFileIDs...)
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PackingAddStickerSet",
		},
		TemplateData: map[string]interface{}{
			"Count": len(s.instance.stickerFileIDs),
		},
	}))
}

func (s *packingState) finish() {
	s.instance.setState(s.instance.busy)
	defer s.instance.setState(s.instance.idle)

	if len(s.instance.stickerFileIDs) == 0 {
		s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "PackingFinishEmpty",
			},
		}))
		return
	}

	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PackingFinish",
		},
		TemplateData: map[string]interface{}{
			"Count": len(s.instance.stickerFileIDs),
		},
	}))

	s.instance.sendStickers(s.instance.stickerFileIDs)
}

func (s *packingState) cancel() {
	defer s.instance.setState(s.instance.idle)
	s.instance.stickerFileIDs = []string{}
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PackingCancel",
		},
	}))
}
