package gotgproto

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/pkg/errors"
)

const (
	maxRetries = 3
)

// authFlow handles the authentication flow for Telegram clients
type authFlow struct {
	client      *auth.Client
	conversator AuthConversator
	flow        auth.Flow
	phone       string
	timeout     time.Time
}

// newAuthFlow creates a new authentication flow instance
func newAuthFlow(client *auth.Client, conversator AuthConversator, phone string, sendOpts auth.SendCodeOptions) *authFlow {
	return &authFlow{
		client:      client,
		conversator: conversator,
		flow: auth.NewFlow(
			termAuth{
				phone:       phone,
				client:      client,
				conversator: conversator,
			},
			sendOpts,
		),
		phone: phone,
	}
}

// Execute runs the authentication flow
func (f *authFlow) Execute(ctx context.Context) error {
	if f.flow.Auth == nil {
		return errors.New("no UserAuthenticator provided")
	}

	sentCode, err := f.sendVerificationCode(ctx)
	if err != nil {
		fmt.Println("sendVerificationCode erR", err)
		return err
	}

	return f.handleSentCode(ctx, sentCode)
}

// sendVerificationCode handles sending the verification code with retries
func (f *authFlow) sendVerificationCode(ctx context.Context) (tg.AuthSentCodeClass, error) {
	SendAuthStatus(f.conversator, AuthStatusPhoneAsked)

	var (
		sentCode tg.AuthSentCodeClass
		phone    string
		err      error
	)

	for i := 0; i < maxRetries; i++ {
		if !f.timeout.IsZero() && time.Now().Before(f.timeout) {
			time.Sleep(time.Until(f.timeout) + time.Second)
		}

		phone, err = f.getPhoneNumber(ctx, i)
		if err != nil {
			return nil, fmt.Errorf("get phone number: %w", err)
		}

		sentCode, err = f.client.SendCode(ctx, phone, f.flow.Options)

		if timeout, ok := tgerr.AsFloodWait(err); ok {
			f.timeout = time.Now().Add(timeout)
			SendAuthStatusFloodWait(f.conversator, f.timeout)
			continue
		}

		if !tgerr.Is(err, "PHONE_NUMBER_INVALID") {
			break
		}
	}

	if err != nil {
		SendAuthStatus(f.conversator, AuthStatusPhoneFailed)
		return nil, fmt.Errorf("send code: %w", err)
	}

	return sentCode, nil
}

// getPhoneNumber handles phone number input with retries
func (f *authFlow) getPhoneNumber(ctx context.Context, attempt int) (string, error) {
	if attempt == 0 {
		return f.flow.Auth.Phone(ctx)
	}

	SendAuthStatusWithRetrials(f.conversator, AuthStatusPhoneRetrial, maxRetries-attempt)

	return f.conversator.AskPhoneNumber()
}

// handleSignIn manages the sign-in process with verification code
func (f *authFlow) handleSignIn(ctx context.Context, phone, hash string) error {
	for i := 0; i < maxRetries; i++ {
		code, err := f.getVerificationCode(ctx, i)
		if err != nil {
			return err
		}

		_, signInErr := f.client.SignIn(ctx, phone, code, hash)
		if signInErr == nil {
			SendAuthStatus(f.conversator, AuthStatusSuccess)
			return nil
		}

		if errors.Is(signInErr, auth.ErrPasswordAuthNeeded) {
			return f.handlePasswordAuth(ctx)
		}

		var signUpRequired *auth.SignUpRequired
		if errors.As(signInErr, &signUpRequired) {
			return f.handleSignUp(ctx, phone, hash, signUpRequired)
		}

		if !tgerr.Is(signInErr, "PHONE_CODE_INVALID") {
			SendAuthStatus(f.conversator, AuthStatusPhoneCodeFailed)
			return errors.Wrap(signInErr, "sign in failed")
		}
	}

	SendAuthStatus(f.conversator, AuthStatusPhoneCodeFailed)
	return errors.New("max verification code attempts exceeded")
}

// getVerificationCode handles verification code input with retries
func (f *authFlow) getVerificationCode(ctx context.Context, attempt int) (string, error) {
	if attempt == 0 {
		SendAuthStatus(f.conversator, AuthStatusPhoneCodeAsked)
		return f.flow.Auth.Code(ctx, nil) // Note: This might need adaptation based on your needs
	}

	SendAuthStatusWithRetrials(f.conversator, AuthStatusPhoneCodeRetrial, maxRetries-attempt)
	return f.conversator.AskCode()
}

// handlePasswordAuth manages 2FA password authentication
func (f *authFlow) handlePasswordAuth(ctx context.Context) error {
	SendAuthStatus(f.conversator, AuthStatusPasswordAsked)
	var err error

	for i := 0; i < maxRetries; i++ {
		password, err := f.getPassword(ctx, i)
		if err != nil {
			return err
		}

		_, err = f.client.Password(ctx, password)
		if err != auth.ErrPasswordInvalid {
			break
		}
	}

	if err != nil {
		SendAuthStatus(f.conversator, AuthStatusPasswordFailed)
		return errors.Wrap(err, "password authentication failed")
	}

	SendAuthStatus(f.conversator, AuthStatusSuccess)
	return nil
}

// getPassword handles password input with retries
func (f *authFlow) getPassword(ctx context.Context, attempt int) (string, error) {
	if attempt == 0 {
		return f.flow.Auth.Password(ctx)
	}

	SendAuthStatusWithRetrials(f.conversator, AuthStatusPasswordRetrial, maxRetries-attempt)
	return f.conversator.AskPassword()
}

// handleSignUp manages the sign-up process for new users
func (f *authFlow) handleSignUp(ctx context.Context, phone, hash string, s *auth.SignUpRequired) error {
	if err := f.flow.Auth.AcceptTermsOfService(ctx, s.TermsOfService); err != nil {
		return errors.Wrap(err, "failed to confirm TOS")
	}

	info, err := f.flow.Auth.SignUp(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get sign up info")
	}

	if _, err := f.client.SignUp(ctx, auth.SignUp{
		PhoneNumber:   phone,
		PhoneCodeHash: hash,
		FirstName:     info.FirstName,
		LastName:      info.LastName,
	}); err != nil {
		return errors.Wrap(err, "failed to sign up")
	}

	return nil
}

// handleSentCode processes the response from sending verification code
func (f *authFlow) handleSentCode(ctx context.Context, sentCode tg.AuthSentCodeClass) error {
	if sentCode == nil {
		return errors.New("sentCode is nil, something went wrong")
	}

	switch s := sentCode.(type) {
	case *tg.AuthSentCode:
		return f.handleSignIn(ctx, f.phone, s.PhoneCodeHash)

	case *tg.AuthSentCodeSuccess:
		return f.handleAuthorizationSuccess(ctx, s)

	default:
		return fmt.Errorf("unexpected sent code type: %T", sentCode)
	}
}

// handleAuthorizationSuccess processes successful authorization responses
func (f *authFlow) handleAuthorizationSuccess(ctx context.Context, s *tg.AuthSentCodeSuccess) error {
	switch a := s.Authorization.(type) {
	case *tg.AuthAuthorization:
		SendAuthStatus(f.conversator, AuthStatusSuccess)
		return nil

	case *tg.AuthAuthorizationSignUpRequired:
		return f.handleSignUp(ctx, f.phone, "", &auth.SignUpRequired{
			TermsOfService: a.TermsOfService,
		})

	default:
		return fmt.Errorf("unexpected authorization type: %T", a)
	}
}
