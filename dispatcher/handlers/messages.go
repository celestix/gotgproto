package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
)

type Message struct {
	Callback      CallbackResponse
	Filters       filters.MessageFilter
	UpdateFilters filters.UpdateFilter
}

func NewMessage(filters filters.MessageFilter, response CallbackResponse) Message {
	return Message{
		Callback:      response,
		Filters:       filters,
		UpdateFilters: nil,
	}
}

func (m Message) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage == nil || u.EffectiveMessage.Out {
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
