package filters

import (
	"github.com/anonyindian/gotgproto/functions"
	"github.com/gotd/td/tg"
	"regexp"
)

type messageFilters struct{}

// All returns true on every type of tg.Message update.
func (*messageFilters) All(_ *tg.Message) bool {
	return true
}

// Chat allows the tg.Message update to process if it is from that particular chat.
func (*messageFilters) Chat(chatId int64) MessageFilter {
	return func(m *tg.Message) bool {
		return functions.GetChatIdFromPeer(m.PeerID) == chatId
	}
}

// Text returns true if tg.Message consists of text.
func (*messageFilters) Text(m *tg.Message) bool {
	return len(m.Message) > 0
}

// Regex returns true if the Message field of tg.Message matches the regex filter
func (*messageFilters) Regex(rString string) (MessageFilter, error) {
	r, err := regexp.Compile(rString)
	if err != nil {
		return nil, err
	}
	return func(m *tg.Message) bool {
		return r.MatchString(m.Message)
	}, nil
}

// Media returns true if tg.Message consists of media.
func (*messageFilters) Media(m *tg.Message) bool {
	return m.Media != nil
}

// Photo returns true if tg.Message consists of photo.
func (*messageFilters) Photo(m *tg.Message) bool {
	_, photo := m.Media.(*tg.MessageMediaPhoto)
	return photo
}

// Video returns true if tg.Message consists of video, gif etc.
func (*messageFilters) Video(m *tg.Message) bool {
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

// Animation returns true if tg.Message consists of animation.
func (*messageFilters) Animation(m *tg.Message) bool {
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

// Sticker returns true if tg.Message consists of sticker.
func (*messageFilters) Sticker(m *tg.Message) bool {
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

// Audio returns true if tg.Message consists of audio.
func (*messageFilters) Audio(m *tg.Message) bool {
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

// Edited returns true if tg.Message is an edited message.
func (*messageFilters) Edited(m *tg.Message) bool {
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
