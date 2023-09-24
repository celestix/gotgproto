package handlers

import (
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
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
	msg := u.EffectiveMessage
	if msg == nil {
		return nil
	}
	if !m.Outgoing && msg.Out {
		return nil
	}
	if m.Filters != nil && !m.Filters(msg) {
		return nil
	}
	if m.UpdateFilters != nil && !m.UpdateFilters(u) {
		return nil
	}
	return m.Callback(ctx, u)
}
