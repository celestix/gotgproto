package handlers

import "github.com/anonyindian/gotgproto/ext"

type AnyUpdate struct {
	Callback CallbackResponse
}

func NewAnyUpdate(response CallbackResponse) AnyUpdate {
	return AnyUpdate{Callback: response}
}

func (au AnyUpdate) CheckUpdate(ctx *ext.Context, u *ext.Update) error {
	return au.Callback(ctx, u)
}
