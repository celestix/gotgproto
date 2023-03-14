package gotgproto

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

func IfAuthNecessary(c *auth.Client, ctx context.Context, flow Flow) error {
	auth, err := c.Status(ctx)
	if err != nil {
		return errors.Wrap(err, "get auth status")
	}
	if auth.Authorized {
		return nil
	}
	if err := authFlow(flow, ctx, c); err != nil {
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

func authFlow(f Flow, ctx context.Context, client *auth.Client) error {
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
		hash := s.PhoneCodeHash
		code, err := f.Auth.Code(ctx, s)
		if err != nil {
			return errors.Wrap(err, "get code")
		}
		_, signInErr := client.SignIn(ctx, phone, code, hash)
		if errors.Is(signInErr, auth.ErrPasswordAuthNeeded) {
			err = signInErr
			for i := 0; err != nil && i < 3; i++ {
				if i != 0 {
					fmt.Println("The 2FA Code you just entered seems to be incorrect,")
					fmt.Println("Attempts Left:", 3-i)
					fmt.Println("Please try again.... ")
				}
				password, err1 := f.Auth.Password(ctx)
				if err1 != nil {
					return errors.Wrap(err1, "get password")
				}
				_, err = client.Password(ctx, password)
			}
			if err != nil {
				return errors.Wrap(err, "sign in with password")
			}
			// if _, err := client.Password(ctx, password); err != nil {
			// 	return errors.Wrap(err, "sign in with password")
			// }
			return nil
		}
		var signUpRequired *auth.SignUpRequired
		if errors.As(signInErr, &signUpRequired) {
			return f.handleSignUp(ctx, client, phone, hash, signUpRequired)
		}

		if signInErr != nil {
			// fmt.Println("\n\n", signInErr.Error(), "\n\n ")
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
