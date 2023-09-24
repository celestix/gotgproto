package handlers

import (
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
)

// CallbackQuery handler is executed when the update consists of tg.UpdateBotCallbackQuery.
type CallbackQuery struct {
	Filters       filters.CallbackQueryFilter
	Callback      CallbackResponse
	UpdateFilters filters.UpdateFilter
}

// NewCallbackQuery creates a new CallbackQuery handler bound to call its response.
func NewCallbackQuery(filters filters.CallbackQueryFilter, response CallbackResponse) CallbackQuery {
	return CallbackQuery{
		Filters:       filters,
		Callback:      response,
		UpdateFilters: nil,
	}
}

func (c CallbackQuery) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	if u.CallbackQuery == nil {
		return nil
	}
	if c.Filters != nil && !c.Filters(u.CallbackQuery) {
		return nil
	}
	if c.UpdateFilters != nil && !c.UpdateFilters(u) {
		return nil
	}
	return c.Callback(ctx, u)
}
