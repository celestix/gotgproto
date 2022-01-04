package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
)

// InlineQuery handler is executed when the update consists of tg.UpdateInlineBotCallbackQuery.
type InlineQuery struct {
	Callback      CallbackResponse
	Filters       filters.InlineQueryFilter
	UpdateFilters filters.UpdateFilter
}

// NewInlineQuery creates a new InlineQuery handler bound to call its response.
func NewInlineQuery(filters filters.InlineQueryFilter, response CallbackResponse) InlineQuery {
	return InlineQuery{
		Filters:       filters,
		Callback:      response,
		UpdateFilters: nil,
	}
}

func (c InlineQuery) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.InlineQuery == nil {
		return nil
	}
	if c.Filters != nil && !c.Filters(u.InlineQuery) {
		return nil
	}
	if c.UpdateFilters != nil && !c.UpdateFilters(u) {
		return nil
	}
	return c.Callback(ctx, u)
}
