package handlers

import (
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
)

type InlineQuery struct {
	Callback      CallbackResponse
	Filters       filters.InlineQueryFilter
	UpdateFilters filters.UpdateFilter
}

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
