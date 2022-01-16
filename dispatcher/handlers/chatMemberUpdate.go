package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
)

// ChatMemberUpdated handler is executed on all type of incoming updates.
type ChatMemberUpdated struct {
	Callback CallbackResponse
	Filters  filters.ChatMemberUpdatedFilter
}

// NewChatMemberUpdated creates a new ChatMemberUpdated handler bound to call its response.
func NewChatMemberUpdated(filters filters.ChatMemberUpdatedFilter, response CallbackResponse) ChatMemberUpdated {
	return ChatMemberUpdated{
		Callback: response,
		Filters:  filters,
	}
}

func (cm ChatMemberUpdated) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.ChatParticipant == nil && u.ChannelParticipant == nil {
		return nil
	}
	if cm.Filters != nil && !cm.Filters(u) {
		return nil
	}
	return cm.Callback(ctx, u)
}
