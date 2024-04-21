package gotgproto

import (
	"context"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/pkg/errors"
)

func ifAuthNecessary(ctx context.Context, c *auth.Client, conversator AuthConversator, phone string, sendOpts auth.SendCodeOptions) error {
	auth, err := c.Status(ctx)
	if err != nil {
		return errors.Wrap(err, "get auth status")
	}
	if auth.Authorized {
		return nil
	}
	if err := authFlow(ctx, c, conversator, phone, sendOpts); err != nil {
		return errors.Wrap(err, "auth flow")
	}
	return nil

}

type Flow auth.Flow

func (f Flow) handleSignUp(ctx context.Context, client auth.FlowClient, phone, hash string, s *auth.SignUpRequired) error {
	if err := f.Auth.AcceptTermsOfService(ctx, s.TermsOfService); err != nil {
		return errors.Wrap(err, "confirm TOS")
	}
	info, err := f.Auth.SignUp(ctx)
	if err != nil {
		return errors.Wrap(err, "sign up info not provided")
	}
	if _, err := client.SignUp(ctx, auth.SignUp{
		PhoneNumber:   phone,
		PhoneCodeHash: hash,
		FirstName:     info.FirstName,
		LastName:      info.LastName,
	}); err != nil {
		return errors.Wrap(err, "sign up")
	}
	return nil
}

func authFlow(ctx context.Context, client *auth.Client, conversator AuthConversator, phone string, sendOpts auth.SendCodeOptions) error {
	f := Flow(auth.NewFlow(
		termAuth{
			phone:       phone,
			client:      client,
			conversator: conversator,
		},
		auth.SendCodeOptions{},
	))
	if f.Auth == nil {
		return errors.New("no UserAuthenticator provided")
	}

	phone, err := f.Auth.Phone(ctx)
	if err != nil {
		return errors.Wrap(err, "get phone")
	}

	sentCode, err := client.SendCode(ctx, phone, f.Options)
	if err != nil {
		return errors.Wrap(err, "send code")
	}

	switch s := sentCode.(type) {
	case *tg.AuthSentCode:
		var (
			code string
			hash = s.PhoneCodeHash
		)

		var signInErr error
		for i := 0; i < 3; i++ {
			if code, err = f.Auth.Code(ctx, s); err != nil {
				return errors.Wrap(err, "get code")
			}

			_, signInErr = client.SignIn(ctx, phone, code, hash)
			if !tgerr.Is(signInErr, "PHONE_CODE_INVALID") {
				break
			}
		}

		if errors.Is(signInErr, auth.ErrPasswordAuthNeeded) {
			err = signInErr
			for i := 0; err != nil && i < 3; i++ {
				var password string
				var err1 error
				if i == 0 {
					password, err1 = f.Auth.Password(ctx)
				} else {
					password, err1 = conversator.RetryPassword(3 - i)
				}
				if err1 != nil {
					return errors.Wrap(err1, "get password")
				}
				_, err = client.Password(ctx, password)
			}
			if err != nil {
				return errors.Wrap(err, "sign in with password")
			}
			return nil
		}
		var signUpRequired *auth.SignUpRequired
		if errors.As(signInErr, &signUpRequired) {
			return f.handleSignUp(ctx, client, phone, hash, signUpRequired)
		}

		if signInErr != nil {
			return errors.Wrap(signInErr, "sign in")
		}
	case *tg.AuthSentCodeSuccess:
		switch a := s.Authorization.(type) {
		case *tg.AuthAuthorization:
			// Looks that we are already authorized.
			return nil
		case *tg.AuthAuthorizationSignUpRequired:
			if err := f.handleSignUp(ctx, client, phone, "", &auth.SignUpRequired{
				TermsOfService: a.TermsOfService,
			}); err != nil {
				// TODO: not sure that blank hash will work here
				return errors.Wrap(err, "sign up after auth sent code success")
			}
			return nil
		default:
			return errors.Errorf("unexpected authorization type: %T", a)
		}
	default:
		return errors.Errorf("unexpected sent code type: %T", sentCode)
	}

	return nil
}
