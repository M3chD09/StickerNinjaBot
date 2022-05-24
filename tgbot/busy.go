package tgbot

import "github.com/nicksnyder/go-i18n/v2/i18n"

type busyState struct {
	instance *instance
}

func (s *busyState) newPack() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Busy",
		},
	}))
}

func (s *busyState) addSticker(stickerFileID string) {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Busy",
		},
	}))
}

func (s *busyState) addStickerSet(stickerSetName string) {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Busy",
		},
	}))
}

func (s *busyState) finish() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Busy",
		},
	}))
}

func (s *busyState) cancel() {
	s.instance.sendTextMessage(s.instance.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Busy",
		},
	}))
}
