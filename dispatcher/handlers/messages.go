package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
)

// Message handler is executed when the update consists of tg.Message with provided conditions.
type Message struct {
	Callback      CallbackResponse
	Filters       filters.MessageFilter
	UpdateFilters filters.UpdateFilter
	Outgoing      bool
}

// NewMessage creates a new Message handler bound to call its response.
func NewMessage(filters filters.MessageFilter, response CallbackResponse) Message {
	return Message{
		Callback:      response,
		Filters:       filters,
		UpdateFilters: nil,
		Outgoing:      true,
	}
}

func (m Message) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage == nil {
		return nil
	}
	if !m.Outgoing && u.EffectiveMessage.Out {
		return nil
	}
	if m.Filters != nil && !m.Filters(u.EffectiveMessage) {
		return nil
	}
	if m.UpdateFilters != nil && !m.UpdateFilters(u) {
		return nil
	}
	return m.Callback(ctx, u)
}
