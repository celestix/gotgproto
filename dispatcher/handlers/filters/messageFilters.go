package filters

import (
	"github.com/anonyindian/gotgproto/functions"
	"github.com/gotd/td/tg"
)

func All(m *tg.Message) bool {
	return true
}

func Chat(chatId int64) MessageFilter {
	return func(m *tg.Message) bool {
		return functions.GetChatIdFromPeer(m.PeerID) == chatId
	}
}

func Text(m *tg.Message) bool {
	return len(m.Message) > 0
}

func Media(m *tg.Message) bool {
	return m.Media != nil
}

func Photo(m *tg.Message) bool {
	_, photo := m.Media.(*tg.MessageMediaPhoto)
	return photo
}

func Video(m *tg.Message) bool {
	doc := GetDocument(m)
	if doc != nil {
		for _, attr := range doc.Attributes {
			_, ok := attr.(*tg.DocumentAttributeVideo)
			if ok {
				return true
			}
		}
	}
	return false
}

func Animation(m *tg.Message) bool {
	doc := GetDocument(m)
	if doc != nil {
		for _, attr := range doc.Attributes {
			_, ok := attr.(*tg.DocumentAttributeAnimated)
			if ok {
				return true
			}
		}
	}
	return false
}

func Sticker(m *tg.Message) bool {
	doc := GetDocument(m)
	if doc != nil {
		for _, attr := range doc.Attributes {
			_, ok := attr.(*tg.DocumentAttributeSticker)
			if ok {
				return true
			}
		}
	}
	return false
}

func Audio(m *tg.Message) bool {
	doc := GetDocument(m)
	if doc != nil {
		for _, attr := range doc.Attributes {
			_, ok := attr.(*tg.DocumentAttributeAudio)
			if ok {
				return true
			}
		}
	}
	return false
}

func Edited(m *tg.Message) bool {
	return m.EditDate != 0
}

func GetDocument(m *tg.Message) *tg.Document {
	mdoc, ok := m.Media.(*tg.MessageMediaDocument)
	if !ok {
		return nil
	}
	tgdoc, ok := mdoc.Document.(*tg.Document)
	if !ok {
		return nil
	}
	return tgdoc
}
