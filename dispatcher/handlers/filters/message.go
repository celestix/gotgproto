package filters

import (
	"regexp"

	"github.com/amupxm/gotgproto/functions"
	"github.com/amupxm/gotgproto/types"
	"github.com/gotd/td/tg"
)

type messageFilters struct{}

// All returns true on every type of types.Message update.
func (*messageFilters) All(_ *types.Message) bool {
	return true
}

type ChatType int

const (
	ChatTypeUser ChatType = iota
	ChatTypeChat
	ChatTypeChannel
)

func (*messageFilters) ChatType(chatType ChatType) MessageFilter {
	return func(m *types.Message) bool {
		chatPeer := m.PeerID
		switch chatType {
		case ChatTypeUser:
			_, ok := chatPeer.(*tg.PeerUser)
			return ok
		case ChatTypeChat:
			_, ok := chatPeer.(*tg.PeerChat)
			return ok
		case ChatTypeChannel:
			_, ok := chatPeer.(*tg.PeerChannel)
			return ok
		}
		return false
	}
}

// Chat allows the types.Message update to process if it is from that particular chat.
func (*messageFilters) Chat(chatId int64) MessageFilter {
	return func(m *types.Message) bool {
		return functions.GetChatIdFromPeer(m.PeerID) == chatId
	}
}

// Text returns true if types.Message consists of text.
func (*messageFilters) Text(m *types.Message) bool {
	return m.Text != ""
}

// Regex returns true if the Message field of types.Message matches the regex filter
func (*messageFilters) Regex(rString string) (MessageFilter, error) {
	r, err := regexp.Compile(rString)
	if err != nil {
		return nil, err
	}
	return func(m *types.Message) bool {
		return r.MatchString(m.Text)
	}, nil
}

// Media returns true if types.Message consists of media.
func (*messageFilters) Media(m *types.Message) bool {
	return m.Media != nil
}

// Photo returns true if types.Message consists of photo.
func (*messageFilters) Photo(m *types.Message) bool {
	_, photo := m.Media.(*tg.MessageMediaPhoto)
	return photo
}

// Video returns true if types.Message consists of video, gif etc.
func (*messageFilters) Video(m *types.Message) bool {
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

// Animation returns true if types.Message consists of animation.
func (*messageFilters) Animation(m *types.Message) bool {
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

// Sticker returns true if types.Message consists of sticker.
func (*messageFilters) Sticker(m *types.Message) bool {
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

// Audio returns true if types.Message consists of audio.
func (*messageFilters) Audio(m *types.Message) bool {
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

// Edited returns true if types.Message is an edited message.
func (*messageFilters) Edited(m *types.Message) bool {
	return m.EditDate != 0
}

func GetDocument(m *types.Message) *tg.Document {
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
