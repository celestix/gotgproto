package handlers

import (
	"github.com/celestix/gotgproto/ext"
)

// AnyUpdate handler is executed on all type of incoming updates.
type AnyUpdate struct {
	Callback CallbackResponse
}

// NewAnyUpdate creates a new AnyUpdate handler bound to call its response.
func NewAnyUpdate(response CallbackResponse) AnyUpdate {
	return AnyUpdate{Callback: response}
}

func (au AnyUpdate) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	return au.Callback(ctx, u)
}
