package filters

import (
	"github.com/anonyindian/gotgproto/ext"
	"github.com/gotd/td/tg"
)

var (
	Message             = messageFilters{}
	CallbackQuery       = callbackQueryFilters{}
	InlineQuery         = inlineQuery{}
	PendingJoinRequests = pendingJoinRequests{}
)

type (
	UpdateFilter              func(u *ext.Update) bool
	MessageFilter             func(m *tg.Message) bool
	CallbackQueryFilter       func(cbq *tg.UpdateBotCallbackQuery) bool
	InlineQueryFilter         func(iq *tg.UpdateBotInlineQuery) bool
	PendingJoinRequestsFilter func(cjr *tg.UpdatePendingJoinRequests) bool
)

// Supergroup returns true if the update is from a supergroup.
func Supergroup(u *ext.Update) bool {
	if c := u.GetChannel(); c != nil {
		return c.Megagroup
	}
	return false
}

// Channel returns true if the update is from a channel.
func Channel(u *ext.Update) bool {
	channelType := u.GetChannel()
	chatType := u.GetChat()
	if channelType != nil && chatType == nil {
		return !channelType.Megagroup
	}
	return false
}

// Group returns true if the update is from a normal group.
func Group(u *ext.Update) bool {
	return u.GetChat() != nil
}
