package ext

import "errors"

var (
	ErrPeerNotFound = errors.New("peer not found")
	ErrNotChat      = errors.New("not chat")
	ErrNotUser      = errors.New("not user")
	ErrTextEmpty    = errors.New("text was not provided")
	ErrTextInvalid  = errors.New("type of text is invalid, provide one from string and []styling.StyledTextOption")
)
