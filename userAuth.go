package gotgproto

import (
	"context"
	"errors"
	"strings"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (noSignUp) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (noSignUp) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	client      *auth.Client
	conversator AuthConversator
	noSignUp

	phone string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	if a.phone != "" {
		return a.phone, nil
	}
	phone, err := a.conversator.AskPhoneNumber()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(phone), nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	pass, err := a.conversator.AskPassword()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(pass), nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	code, err := a.conversator.AskCode()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
