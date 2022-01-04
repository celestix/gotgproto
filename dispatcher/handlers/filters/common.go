package filters

import (
	"github.com/anonyindian/gotgproto/ext"
	"github.com/gotd/td/tg"
)

type MessageFilter func(m *tg.Message) bool
type CallbackQueryFilter func(cbq *tg.UpdateBotCallbackQuery) bool
type InlineQueryFilter func(iq *tg.UpdateInlineBotCallbackQuery) bool
type UpdateFilter func(u *ext.Update) bool

const (
	ChatTypeChannel    = "channel"
	ChatTypeSuperGroup = "supergroup"
	ChatTypeGroup      = "group"
	ChatTypeUser       = "user"
)

func Supergroup(u *ext.Update) bool {
	if c := u.GetChannel(); c != nil {
		return c.Megagroup
	}
	return false
}

func Channel(u *ext.Update) bool {
	channelType := u.GetChannel()
	chatType := u.GetChat()
	if channelType != nil && chatType == nil {
		return !channelType.Megagroup
	}
	return false
}

func Group(u *ext.Update) bool {
	return u.GetChat() != nil
}
