package handlers

import (
	"github.com/anonyindian/gotgproto/ext"
)

type CallbackResponse func(*ext.Context, *ext.Update) error
