package errors

import "errors"

var (
	ErrClientAlreadyRunning = errors.New("client is already running")
	ErrClientNotInitialized = errors.New("client is not initialized")
	ErrSessionUnauthorized  = errors.New("session is unauthorized")
	ErrNoOptions            = errors.New("no options provided")
)

var (
	ErrPeerNotFound     = errors.New("peer not found")
	ErrNotChat          = errors.New("not chat")
	ErrNotChannel       = errors.New("not channel")
	ErrNotUser          = errors.New("not user")
	ErrTextEmpty        = errors.New("text was not provided")
	ErrTextInvalid      = errors.New("type of text is invalid, provide one from string and []styling.StyledTextOption")
	ErrMessageNotExist  = errors.New("message not exist")
	ErrReplyNotMessage  = errors.New("reply header is not a message")
	ErrUnknownTypeMedia = errors.New("unknown type media")
)
