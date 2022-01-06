package filters

import (
	"github.com/anonyindian/gotgproto/functions"
	"github.com/gotd/td/tg"
	"regexp"
)

// All returns true on every type of tg.Message update.
func All(m *tg.Message) bool {
	return true
}

// Chat allows the tg.Message update to process if it is from that particular chat.
func Chat(chatId int64) MessageFilter {
	return func(m *tg.Message) bool {
		return functions.GetChatIdFromPeer(m.PeerID) == chatId
	}
}

// Text returns true if tg.Message consists of text.
func Text(m *tg.Message) bool {
	return len(m.Message) > 0
}

func Regex(r_string string) (MessageFilter , error) {
	r , err := regexp.Compile(r_string)
	if err != nil {
		return func (msg *tg.Message) bool {
			return false
		},err
	}
	
	return func(msg *tg.Message) bool {
		return bool(r.MatchString(msg.Message)) 
	},nil
	

}

// Media returns true if tg.Message consists of media.
func Media(m *tg.Message) bool {
	return m.Media != nil
}

// Photo returns true if tg.Message consists of photo.
func Photo(m *tg.Message) bool {
	_, photo := m.Media.(*tg.MessageMediaPhoto)
	return photo
}

// Video returns true if tg.Message consists of video, gif etc.
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

// Animation returns true if tg.Message consists of animation.
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

// Sticker returns true if tg.Message consists of sticker.
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

// Audio returns true if tg.Message consists of audio.
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

// Edited returns true if tg.Message is an edited message.
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
